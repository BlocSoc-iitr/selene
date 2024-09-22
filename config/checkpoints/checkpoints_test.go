package checkpoints
import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/BlocSoc-iitr/selene/config"
)
type CustomTransport struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}
func (t *CustomTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.DoFunc(req)
}
func TestCheckpointFallbackNew(t *testing.T) {
	cf := CheckpointFallback{}.New()

	if len(cf.Services) != 0 {
		t.Errorf("Expected empty Services map, got %v", cf.Services)
	}

	expectedNetworks := []config.Network{config.SEPOLIA, config.MAINNET, config.GOERLI}
	if !equalNetworks(cf.Networks, expectedNetworks) {
		t.Errorf("Expected Networks %v, got %v", expectedNetworks, cf.Networks)
	}
}

func TestDeserializeSlot(t *testing.T) {
	tests := []struct {
		input    []byte
		expected uint64
		hasError bool
	}{
		{[]byte(`"12345"`), 12345, false},
		{[]byte(`12345`), 12345, false},
		{[]byte(`"abc"`), 0, true},
		{[]byte(`{}`), 0, true},
	}

	for _, test := range tests {
		result, err := DeserializeSlot(test.input)
		if test.hasError {
			if err == nil {
				t.Errorf("Expected error for input %s, got nil", string(test.input))
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for input %s: %v", string(test.input), err)
			}
			if result == nil || *result != test.expected {
				t.Errorf("Expected %d, got %v for input %s", test.expected, result, string(test.input))
			}
		}
	}
}
func TestGetHealthyFallbackEndpoints(t *testing.T) {
	cf := CheckpointFallback{
		Services: map[config.Network][]CheckpointFallbackService{
			config.MAINNET: {
				{Endpoint: "http://healthy1.com", Health_from_fallback: &Health{Result: true}},
				{Endpoint: "http://unhealthy.com", Health_from_fallback: &Health{Result: false}},
				{Endpoint: "http://healthy2.com", Health_from_fallback: &Health{Result: true}},
			},
		},
	}

	healthyEndpoints := cf.GetHealthyFallbackEndpoints(config.MAINNET)
	expected := []string{"http://healthy1.com", "http://healthy2.com"}

	if !equalStringSlices(healthyEndpoints, expected) {
		t.Errorf("Expected healthy endpoints %v, got %v", expected, healthyEndpoints)
	}
}

func TestBuild(t *testing.T) {
	yamlData := `
sepolia:
  - endpoint: "https://sepolia1.example.com"
    name: "Sepolia 1"
    state: true
    verification: true
    health:
      result: true
mainnet:
  - endpoint: "https://mainnet1.example.com"
    name: "Mainnet 1"
    state: true
    verification: true
    health:
      result: true
goerli:
  - endpoint: "https://goerli1.example.com"
    name: "Goerli 1"
    state: true
    verification: true
    health:
      result: true
`

	client := &http.Client{
		Transport: &CustomTransport{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewBufferString(yamlData)),
					Header:     make(http.Header),
				}, nil
			},
		},
	}

	cf := CheckpointFallback{}.New()
	originalClient := http.DefaultClient
	http.DefaultClient = client
	defer func() { http.DefaultClient = originalClient }()

	builtCf, err := cf.Build()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(builtCf.Services[config.SEPOLIA]) != 1 {
		t.Errorf("Expected 1 service for Sepolia, got %d", len(builtCf.Services[config.SEPOLIA]))
	}

	if len(builtCf.Services[config.MAINNET]) != 1 {
		t.Errorf("Expected 1 service for Mainnet, got %d", len(builtCf.Services[config.MAINNET]))
	}

	sepoliaService := builtCf.Services[config.SEPOLIA][0]
	if sepoliaService.Endpoint != "https://sepolia1.example.com" {
		t.Errorf("Expected endpoint https://sepolia1.example.com, got %s", sepoliaService.Endpoint)
	}
}

func TestFetchLatestCheckpointFromServices(t *testing.T) {
    server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if err := json.NewEncoder(w).Encode(RawSlotResponse{
            Data: RawSlotResponseData{
                Slots: []Slot{
                    {Epoch: 100, Block_root: &byte256{1}},
                },
            },
        }); err != nil {
            t.Fatalf("Failed to encode response: %v", err)
        }
    }))
    defer server1.Close()

    server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if err := json.NewEncoder(w).Encode(RawSlotResponse{
            Data: RawSlotResponseData{
                Slots: []Slot{
                    {Epoch: 101, Block_root: &byte256{2}},
                },
            },
        }); err != nil {
            t.Fatalf("Failed to encode response: %v", err)
        }
    }))
    defer server2.Close()

    cf := CheckpointFallback{}
    services := []CheckpointFallbackService{
        {Endpoint: server1.URL},
        {Endpoint: server2.URL},
    }

    checkpoint, err := cf.FetchLatestCheckpointFromServices(services)
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }

    expected := byte256{2}
    if checkpoint != expected {
        t.Errorf("Expected checkpoint %v, got %v", expected, checkpoint)
    }
}
func TestQueryService(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if err := json.NewEncoder(w).Encode(RawSlotResponse{
            Data: RawSlotResponseData{
                Slots: []Slot{
                    {Epoch: 100, Block_root: &byte256{1}},
                },
            },
        }); err != nil {
            t.Fatalf("Failed to encode response: %v", err)
        }
    }))
    defer server.Close()

    cf := CheckpointFallback{}
    response, err := cf.QueryService(server.URL)
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }

    if len(response.Data.Slots) != 1 {
        t.Errorf("Expected 1 slot, got %d", len(response.Data.Slots))
    }

    if response.Data.Slots[0].Epoch != 100 {
        t.Errorf("Expected epoch 100, got %d", response.Data.Slots[0].Epoch)
    }
}

func TestFetchLatestCheckpointFromApi(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if err := json.NewEncoder(w).Encode(RawSlotResponse{
            Data: RawSlotResponseData{
                Slots: []Slot{
                    {Epoch: 100, Block_root: &byte256{1}},
                },
            },
        }); err != nil {
            t.Fatalf("Failed to encode response: %v", err)
        }
    }))
    defer server.Close()

    cf := CheckpointFallback{}
    checkpoint, err := cf.FetchLatestCheckpointFromApi(server.URL)
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }

    expected := byte256{1}
    if checkpoint != expected {
        t.Errorf("Expected checkpoint %v, got %v", expected, checkpoint)
    }
}

func TestConstructUrl(t *testing.T) {
	cf := CheckpointFallback{}
	endpoint := "https://example.com"
	expected := "https://example.com/checkpointz/v1/beacon/slots"
	result := cf.ConstructUrl(endpoint)
	if result != expected {
		t.Errorf("Expected URL %s, got %s", expected, result)
	}
}

func TestGetAllFallbackEndpoints(t *testing.T) {
	cf := CheckpointFallback{
		Services: map[config.Network][]CheckpointFallbackService{
			config.MAINNET: {
				{Endpoint: "http://endpoint1.com"},
				{Endpoint: "http://endpoint2.com"},
			},
		},
	}

	endpoints := cf.GetAllFallbackEndpoints(config.MAINNET)
	expected := []string{"http://endpoint1.com", "http://endpoint2.com"}

	if !equalStringSlices(endpoints, expected) {
		t.Errorf("Expected endpoints %v, got %v", expected, endpoints)
	}
}

func TestGetHealthyFallbackServices(t *testing.T) {
	cf := CheckpointFallback{
		Services: map[config.Network][]CheckpointFallbackService{
			config.MAINNET: {
				{Endpoint: "http://healthy1.com", Health_from_fallback: &Health{Result: true}},
				{Endpoint: "http://unhealthy.com", Health_from_fallback: &Health{Result: false}},
				{Endpoint: "http://healthy2.com", Health_from_fallback: &Health{Result: true}},
			},
		},
	}

	healthyServices := cf.GetHealthyFallbackServices(config.MAINNET)
	if len(healthyServices) != 2 {
		t.Errorf("Expected 2 healthy services, got %d", len(healthyServices))
	}

	for _, service := range healthyServices {
		if !service.Health_from_fallback.Result {
			t.Errorf("Expected healthy service, got unhealthy for endpoint %s", service.Endpoint)
		}
	}
}

func equalNetworks(a, b []config.Network) bool {
	if len(a) != len(b) {
		return false
	}
	networkMap := make(map[config.Network]struct{})
	for _, network := range a {
		networkMap[network] = struct{}{}
	}
	for _, network := range b {
		if _, exists := networkMap[network]; !exists {
			return false
		}
	}
	return true
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	stringMap := make(map[string]struct{})
	for _, str := range a {
		stringMap[str] = struct{}{}
	}
	for _, str := range b {
		if _, exists := stringMap[str]; !exists {
			return false
		}
	}
	return true
}
