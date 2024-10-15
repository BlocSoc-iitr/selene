package checkpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/BlocSoc-iitr/selene/config"
	"github.com/avast/retry-go"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// / The location where the list of checkpoint services are stored.
const CHECKPOINT_SYNC_SERVICES_LIST = "https://raw.githubusercontent.com/ethpandaops/checkpoint-sync-health-checks/master/_data/endpoints.yaml"

type byte256 [32]byte

type StartEndTime struct {
	/// An ISO 8601 formatted UTC timestamp.
	Start_time string
	/// An ISO 8601 formatted UTC timestamp.
	End_time string
}

// /Struct to health check checkpoint sync service
type Health struct {
	/// returns true if node is healthy and false otherwise
	Result bool
	/// An ISO 8601 formatted UTC timestamp.
	Date string
}
type CheckpointFallbackService struct {
	/// The endpoint of the service.
	Endpoint string
	///The checkpoint sync service name
	Name string
	///True if the checkpoint sync service is avalible to send requests
	State bool
	///True if the checkpoint sync service is verified
	Verification bool
	///The contacts of the checkpoint sync service maintainers
	Contacts *yaml.MapSlice `yaml:"contacts,omitempty"` // Option type in Go becomes a pointer
	///service notes
	Notes *yaml.MapSlice `yaml:"notes,omitempty"` // Using MapSlice for generic YAML structure
	///Health check of the checkpoint sync service
	Health_from_fallback *Health
}

// / The CheckpointFallback manages checkpoint fallback services.
type CheckpointFallback struct {
	///services map
	Services map[config.Network][]CheckpointFallbackService
	///networks list
	///available network - [SEPOLIA, MAINNET, GOERLI]
	Networks []config.Network
}
type RawSlotResponse struct {
	Data RawSlotResponseData
}
type RawSlotResponseData struct {
	Slots []Slot
}
type Slot struct {
	Slot       uint64
	Block_root *byte256
	State_root *byte256
	Epoch      uint64
	Time       StartEndTime
}

// @param url: string - url to fetch
// @return *http.Response, error
// Fetches the response from the given url
func Get(url string) (*http.Response, error) {
	var resp *http.Response

	// Retry with default settings
	err := retry.Do(
		func() error {
			var err error
			resp, err = http.Get(url)
			if err != nil {
				return err
			}

			// Check if response status is not OK and trigger retry if so
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("received non-200 response: %d", resp.StatusCode)
			}

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return resp, nil

}

// @param data: []byte - data to deserialize
// @return *uint64, error
// Deserializes the given data to uint64
func DeserializeSlot(input []byte) (*uint64, error) {
	var value interface{}
	if err := json.Unmarshal(input, &value); err != nil {
		return nil, err
	}

	switch v := value.(type) {
	case string:
		// Try to parse a string as a number
		num, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid string format: %s", v)
		}
		return &num, nil
	case float64:
		num := uint64(v)
		return &num, nil
	default:
		return nil, fmt.Errorf("unexpected type: %T", v)
	}
}

// create a new CheckpointFallback object
// @return CheckpointFallback
func (ch CheckpointFallback) New() CheckpointFallback {
	return CheckpointFallback{
		Services: make(map[config.Network][]CheckpointFallbackService),
		Networks: []config.Network{
			config.SEPOLIA,
			config.MAINNET,
			config.GOERLI,
		},
	}

}

// build the CheckpointFallback object from the fetch data from get request

// / Build the checkpoint fallback service from the community-maintained list by [ethPandaOps](https://github.com/ethpandaops).
// /
// / The list is defined in [ethPandaOps/checkpoint-fallback-service](https://github.com/ethpandaops/checkpoint-sync-health-checks/blob/master/_data/endpoints.yaml).
func (ch CheckpointFallback) Build() (CheckpointFallback, error) {
	resp, err := http.Get(CHECKPOINT_SYNC_SERVICES_LIST)
	if err != nil {
		return ch, fmt.Errorf("failed to fetch services list: %w", err)
	}
	defer resp.Body.Close()

	yamlData, err := io.ReadAll(resp.Body)
	if err != nil {
		return ch, fmt.Errorf("failed to read response body: %w", err)
	}

	data := make(map[string]interface{})
	err = yaml.Unmarshal(yamlData, &data)
	if err != nil {
		return ch, fmt.Errorf("failed to parse YAML: %w", err)
	}

	for _, network := range ch.Networks {
		networkFormatted := strings.ToLower(string(network))
		serviceListRaw, ok := data[networkFormatted]
		if !ok {
			return ch, fmt.Errorf("no services found for network %s", network)
		}

		serviceList, ok := serviceListRaw.([]interface{})
		if !ok {
			return ch, fmt.Errorf("expected a list for services in network %s", network)
		}

		for _, serviceRaw := range serviceList {
			serviceMap, ok := serviceRaw.(map[interface{}]interface{})
			if !ok {
				return ch, fmt.Errorf("expected a map for service in network %s", network)
			}

			endpoint, _ := serviceMap["endpoint"].(string)       // Handle potential nil
			name, _ := serviceMap["name"].(string)               // Handle potential nil
			state, _ := serviceMap["state"].(bool)               // Handle potential nil
			verification, _ := serviceMap["verification"].(bool) // Handle potential nil

			// Check contacts and notes
			var contacts *yaml.MapSlice
			if c, ok := serviceMap["contacts"].(*yaml.MapSlice); ok {
				contacts = c
			}

			var notes *yaml.MapSlice
			if n, ok := serviceMap["notes"].(*yaml.MapSlice); ok {
				notes = n
			}

			healthRaw, ok := serviceMap["health"].(map[interface{}]interface{})
			if !ok {
				return ch, fmt.Errorf("expected a map for health in service %s", name)
			}
			healthResult, _ := healthRaw["result"].(bool) // Handle potential nil
			healthDate, _ := healthRaw["date"].(string)   // Handle potential nil

			ch.Services[network] = append(ch.Services[network], CheckpointFallbackService{
				Endpoint:     endpoint,
				Name:         name,
				State:        state,
				Verification: verification,
				Contacts:     contacts,
				Notes:        notes,
				Health_from_fallback: &Health{
					Result: healthResult,
					Date:   healthDate,
				},
			})
		}
	}

	return ch, nil
}

// fetch the latest checkpoint from the given network
func (ch CheckpointFallback) FetchLatestCheckpoint(network config.Network) byte256 {
	services := ch.GetHealthyFallbackServices(network)
	checkpoint, error := ch.FetchLatestCheckpointFromServices(services)
	if error != nil {
		return byte256{}
	}
	return checkpoint

}

// fetch the latest checkpoint from the given endpoint
func (ch CheckpointFallback) QueryService(endpoint string) (*RawSlotResponse, error) {
	constructed_url := ch.ConstructUrl(endpoint)
	resp, err := http.Get(constructed_url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var raw RawSlotResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	return &raw, nil
}

// fetch the latest checkpoint from the given services
func (ch CheckpointFallback) FetchLatestCheckpointFromServices(services []CheckpointFallbackService) (byte256, error) {
	var (
		slots      []Slot
		wg         sync.WaitGroup
		slotChan   = make(chan Slot, len(services)) // Buffered channel
		errorsChan = make(chan error, len(services))
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, service := range services {
		wg.Add(1)
		go func(service CheckpointFallbackService) {
			defer wg.Done()
			raw, err := ch.QueryService(service.Endpoint)
			if err != nil {
				errorsChan <- fmt.Errorf("failed to fetch checkpoint from service %s: %w", service.Endpoint, err)
				return
			}
			if len(raw.Data.Slots) > 0 {
				slotChan <- raw.Data.Slots[0] // Send the first valid slot
			}
		}(service)
	}

	go func() {
		wg.Wait()
		close(slotChan)
		close(errorsChan)
	}()

	for {
		select {
		case slot, ok := <-slotChan:
			if !ok {
				// Channel closed, all slots processed
				if len(slots) == 0 {
					return byte256{}, fmt.Errorf("failed to find max epoch from checkpoint slots")
				}
				return processSlots(slots)
			}
			slots = append(slots, slot)
		case err := <-errorsChan:
			if err != nil {
				log.Printf("Error fetching checkpoint: %v", err) // Log only if the error is not nil.
			}
		case <-ctx.Done():
			if len(slots) == 0 {
				return byte256{}, ctx.Err()
			}
			return processSlots(slots)
		}
	}
}

func processSlots(slots []Slot) (byte256, error) {
	maxEpochSlot := slots[0]
	for _, slot := range slots {
		if slot.Epoch > maxEpochSlot.Epoch {
			maxEpochSlot = slot
		}
	}

	if maxEpochSlot.Block_root == nil {
		return byte256{}, fmt.Errorf("no valid block root found")
	}

	return *maxEpochSlot.Block_root, nil
}

func (ch CheckpointFallback) FetchLatestCheckpointFromApi(url string) (byte256, error) {
	constructed_url := ch.ConstructUrl(url)
	resp, err := http.Get(constructed_url)

	if err != nil {
		return byte256{}, fmt.Errorf("failed to fetch checkpoint from API: %w", err)
	}
	defer resp.Body.Close()

	var raw RawSlotResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return byte256{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(raw.Data.Slots) == 0 {
		return byte256{}, fmt.Errorf("no slots found in response")
	}

	slot := raw.Data.Slots[0]
	if slot.Block_root == nil {
		return byte256{}, fmt.Errorf("no block root found in response")
	}

	return *slot.Block_root, nil

}

func (ch CheckpointFallback) ConstructUrl(endpoint string) string {
	return fmt.Sprintf("%s/checkpointz/v1/beacon/slots", endpoint)
}

func (ch CheckpointFallback) GetAllFallbackEndpoints(network config.Network) []string {
	var endpoints []string
	for _, service := range ch.Services[network] {
		endpoints = append(endpoints, service.Endpoint)
	}

	return endpoints

}

func (ch CheckpointFallback) GetHealthyFallbackEndpoints(network config.Network) []string {
	var healthyEndpoints []string
	for _, service := range ch.Services[network] {
		if service.Health_from_fallback.Result {
			healthyEndpoints = append(healthyEndpoints, service.Endpoint)
		}
	}
	return healthyEndpoints

}
func (ch CheckpointFallback) GetHealthyFallbackServices(network config.Network) []CheckpointFallbackService {
	var healthyServices []CheckpointFallbackService
	for _, service := range ch.Services[network] {
		if service.Health_from_fallback.Result {
			healthyServices = append(healthyServices, service)
		}
	}
	return healthyServices

}
func (ch CheckpointFallback) GetFallbackServices(network config.Network) []CheckpointFallbackService {
	return ch.Services[network]
}
