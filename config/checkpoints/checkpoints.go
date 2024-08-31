package checkpoints

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/BlocSoc-iitr/selene/config"
	"github.com/avast/retry-go"
	"gopkg.in/yaml.v2"
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
func get(url string) (*http.Response, error) {
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
func deserialize_slot(data []byte) (*uint64, error) {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}

	if len(s) > 1 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}

	value, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return nil, err
	}

	return &value, nil

}

// create a new CheckpointFallback object
// @return CheckpointFallback
func (ch CheckpointFallback) new() CheckpointFallback {
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
func (ch CheckpointFallback) build() (CheckpointFallback, error) {
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

		serviceList := serviceListRaw.([]interface{})
		for _, service := range serviceList {
			serviceMap := service.(map[interface{}]interface{})
			endpoint := serviceMap["endpoint"].(string)
			name := serviceMap["name"].(string)
			state := serviceMap["state"].(bool)
			verification := serviceMap["verification"].(bool)
			contacts := serviceMap["contacts"].(*yaml.MapSlice)
			notes := serviceMap["notes"].(*yaml.MapSlice)
			health := serviceMap["health"].(map[interface{}]interface{})
			healthResult := health["result"].(bool)
			healthDate := health["date"].(string)

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
func (ch CheckpointFallback) fetch_latest_checkpoint(network config.Network) byte256 {
	services := ch.get_healthy_fallback_services(network)
	checkpoint, error := ch.fetch_latest_checkpoint_from_services(services)
	if error != nil {
		return byte256{}
	}
	return checkpoint

}

// fetch the latest checkpoint from the given endpoint
func (ch CheckpointFallback) query_service(endpoint string) (*RawSlotResponse, error) {
	constructed_url := ch.construct_url(endpoint)
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
func (ch CheckpointFallback) fetch_latest_checkpoint_from_services(services []CheckpointFallbackService) (byte256, error) {
	var (
		slots      []Slot
		wg         sync.WaitGroup
		slotChan   = make(chan Slot)
		errorsChan = make(chan error)
	)

	for _, service := range services {
		wg.Add(1)
		go func(service CheckpointFallbackService) {
			defer wg.Done()
			raw, err := ch.query_service(service.Endpoint)
			if err != nil {
				errorsChan <- fmt.Errorf("failed to fetch checkpoint from service %s: %w", service.Endpoint, err)
				return
			}

			if len(raw.Data.Slots) > 0 {
				for _, slot := range raw.Data.Slots {
					if slot.Block_root != nil {
						slotChan <- slot
						return
					}
				}
			}
		}(service)
	}

	wg.Wait()
	close(slotChan)
	close(errorsChan)

	var allErrors error
	for err := range errorsChan {
		if allErrors == nil {
			allErrors = err
		} else {
			allErrors = fmt.Errorf("%v; %v", allErrors, err)
		}
	}

	if allErrors != nil {
		return byte256{}, allErrors
	}

	for slot := range slotChan {
		slots = append(slots, slot)
	}

	if len(slots) == 0 {
		return byte256{}, fmt.Errorf("failed to find max epoch from checkpoint slots")
	}

	maxEpochSlot := slots[0]
	for _, slot := range slots {
		if slot.Epoch > maxEpochSlot.Epoch {
			maxEpochSlot = slot
		}
	}
	maxEpoch := maxEpochSlot.Epoch

	var maxEpochSlots []Slot
	for _, slot := range slots {
		if slot.Epoch == maxEpoch {
			maxEpochSlots = append(maxEpochSlots, slot)
		}
	}

	checkpoints := make(map[byte256]int)
	for _, slot := range maxEpochSlots {
		if slot.Block_root != nil {
			checkpoints[*slot.Block_root]++
		}
	}

	var mostCommon byte256
	maxCount := 0
	for blockRoot, count := range checkpoints {
		if count > maxCount {
			mostCommon = blockRoot
			maxCount = count
		}
	}

	if maxCount == 0 {
		return byte256{}, fmt.Errorf("no checkpoint found")
	}

	return mostCommon, nil
}

func (ch CheckpointFallback) fetch_latest_checkpoint_from_api(url string) (byte256, error) {
	constructed_url := ch.construct_url(url)
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

func (ch CheckpointFallback) construct_url(endpoint string) string {
	return fmt.Sprintf("%s/checkpointz/v1/beacon/slots", endpoint)
}

func (ch CheckpointFallback) get_all_fallback_endpoints(network config.Network) []string {
	var endpoints []string
	for _, service := range ch.Services[network] {
		endpoints = append(endpoints, service.Endpoint)
	}

	return endpoints

}

func (ch CheckpointFallback) get_healthy_fallback_endpoints(network config.Network) []string {
	var healthyEndpoints []string
	for _, service := range ch.Services[network] {
		if service.Health_from_fallback.Result {
			healthyEndpoints = append(healthyEndpoints, service.Endpoint)
		}
	}
	return healthyEndpoints

}
func (ch CheckpointFallback) get_healthy_fallback_services(network config.Network) []CheckpointFallbackService {
	var healthyServices []CheckpointFallbackService
	for _, service := range ch.Services[network] {
		if service.Health_from_fallback.Result {
			healthyServices = append(healthyServices, service)
		}
	}
	return healthyServices

}
func (ch CheckpointFallback) get_fallback_services(network config.Network) []CheckpointFallbackService {
	return ch.Services[network]
}
