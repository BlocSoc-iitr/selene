//go:build !optimism
// +build !optimism

package evm

import (
	"encoding/json"
	"fmt"
	"math/big"

)
type TxEnv struct {
    // Caller aka Author aka transaction signer
    Caller Address `json:"caller"`
    // The gas limit of the transaction
    GasLimit uint64 `json:"gas_limit"`
    // The gas price of the transaction
    GasPrice *big.Int `json:"gas_price"`
    // The destination of the transaction
    TransactTo TxKind `json:"transact_to"`
    // The value sent to TransactTo
    Value *big.Int `json:"value"`
    // The data of the transaction
    Data Bytes `json:"data"`
    // The nonce of the transaction
    // If nil, nonce validation against the account's nonce is skipped
    Nonce *uint64 `json:"nonce,omitempty"`
    // The chain ID of the transaction
    // If nil, no checks are performed (EIP-155)
    ChainID *uint64 `json:"chain_id,omitempty"`
    // List of addresses and storage keys that the transaction plans to access (EIP-2930)
    AccessList []AccessListItem `json:"access_list"`
    // The priority fee per gas (EIP-1559)
    GasPriorityFee *big.Int `json:"gas_priority_fee,omitempty"`
    // The list of blob versioned hashes (EIP-4844)
    BlobHashes []B256 `json:"blob_hashes"`
    // The max fee per blob gas (EIP-4844)
    MaxFeePerBlobGas *big.Int `json:"max_fee_per_blob_gas,omitempty"`
    // List of authorizations for EOA account code (EIP-7702)
    AuthorizationList *AuthorizationList `json:"authorization_list,omitempty"`
}
const isOptimismEnabled = false

// SpecId represents the specification IDs and their activation block.
type SpecId uint8
const (
	FRONTIER        SpecId = 0   // Frontier               0
	FRONTIER_THAWING SpecId = 1   // Frontier Thawing       200000
	HOMESTEAD       SpecId = 2   // Homestead              1150000
	DAO_FORK        SpecId = 3   // DAO Fork               1920000
	TANGERINE       SpecId = 4   // Tangerine Whistle      2463000
	SPURIOUS_DRAGON  SpecId = 5   // Spurious Dragon        2675000
	BYZANTIUM       SpecId = 6   // Byzantium              4370000
	CONSTANTINOPLE  SpecId = 7   // Constantinople         7280000 is overwritten with PETERSBURG
	PETERSBURG      SpecId = 8   // Petersburg             7280000
	ISTANBUL        SpecId = 9   // Istanbul               9069000
	MUIR_GLACIER    SpecId = 10  // Muir Glacier           9200000
	BERLIN          SpecId = 11  // Berlin                 12244000
	LONDON          SpecId = 12  // London                 12965000
	ARROW_GLACIER   SpecId = 13  // Arrow Glacier          13773000
	GRAY_GLACIER    SpecId = 14  // Gray Glacier           15050000
	MERGE           SpecId = 15  // Paris/Merge            15537394 (TTD: 58750000000000000000000)
	SHANGHAI        SpecId = 16  // Shanghai               17034870 (Timestamp: 1681338455)
	CANCUN          SpecId = 17  // Cancun                 19426587 (Timestamp: 1710338135)
	PRAGUE          SpecId = 18  // Prague                 TBD
	PRAGUE_EOF      SpecId = 19  // Prague+EOF             TBD
	LATEST          SpecId = 255  // LATEST = u8::MAX
)

// String method to convert SpecId to string
func (s SpecId) String() string {
	switch s {
	case FRONTIER:
		return "FRONTIER"
	case FRONTIER_THAWING:
		return "FRONTIER_THAWING"
	case HOMESTEAD:
		return "HOMESTEAD"
	case DAO_FORK:
		return "DAO_FORK"
	case TANGERINE:
		return "TANGERINE"
	case SPURIOUS_DRAGON:
		return "SPURIOUS_DRAGON"
	case BYZANTIUM:
		return "BYZANTIUM"
	case CONSTANTINOPLE:
		return "CONSTANTINOPLE"
	case PETERSBURG:
		return "PETERSBURG"
	case ISTANBUL:
		return "ISTANBUL"
	case MUIR_GLACIER:
		return "MUIR_GLACIER"
	case BERLIN:
		return "BERLIN"
	case LONDON:
		return "LONDON"
	case ARROW_GLACIER:
		return "ARROW_GLACIER"
	case GRAY_GLACIER:
		return "GRAY_GLACIER"
	case MERGE:
		return "MERGE"
	case SHANGHAI:
		return "SHANGHAI"
	case CANCUN:
		return "CANCUN"
	case PRAGUE:
		return "PRAGUE"
	case PRAGUE_EOF:
		return "PRAGUE_EOF"
	default:
		return "UNKNOWN"
	}
}
func DefaultSpecId() SpecId {
	return LATEST
}

type (
	FrontierSpec       struct{ BaseSpec }
	HomesteadSpec      struct{ BaseSpec }
	TangerineSpec      struct{ BaseSpec }
	SpuriousDragonSpec struct{ BaseSpec }
	ByzantiumSpec      struct{ BaseSpec }
	PetersburgSpec     struct{ BaseSpec }
	IstanbulSpec       struct{ BaseSpec }
	BerlinSpec         struct{ BaseSpec }
	LondonSpec         struct{ BaseSpec }
	MergeSpec          struct{ BaseSpec }
	ShanghaiSpec       struct{ BaseSpec }
	CancunSpec         struct{ BaseSpec }
	PragueSpec         struct{ BaseSpec }
	PragueEofSpec      struct{ BaseSpec }
	LatestSpec         struct{ BaseSpec }
)
func SpecToGeneric(specID SpecId) Spec {
	switch specID {
	case FRONTIER, FRONTIER_THAWING:
		return FrontierSpec{NewSpec(FRONTIER)}
	case HOMESTEAD, DAO_FORK:
		return HomesteadSpec{NewSpec(HOMESTEAD)}
	case TANGERINE:
		return TangerineSpec{NewSpec(TANGERINE)}
	case SPURIOUS_DRAGON:
		return SpuriousDragonSpec{NewSpec(SPURIOUS_DRAGON)}
	case BYZANTIUM:
		return ByzantiumSpec{NewSpec(BYZANTIUM)}
	case PETERSBURG, CONSTANTINOPLE:
		return PetersburgSpec{NewSpec(PETERSBURG)}
	case ISTANBUL, MUIR_GLACIER:
		return IstanbulSpec{NewSpec(ISTANBUL)}
	case BERLIN:
		return BerlinSpec{NewSpec(BERLIN)}
	case LONDON, ARROW_GLACIER, GRAY_GLACIER:
		return LondonSpec{NewSpec(LONDON)}
	case MERGE:
		return MergeSpec{NewSpec(MERGE)}
	case SHANGHAI:
		return ShanghaiSpec{NewSpec(SHANGHAI)}
	case CANCUN:
		return CancunSpec{NewSpec(CANCUN)}
	case PRAGUE:
		return PragueSpec{NewSpec(PRAGUE)}
	case PRAGUE_EOF:
		return PragueEofSpec{NewSpec(PRAGUE_EOF)}
	default:
		return LatestSpec{NewSpec(LATEST)}
	}
}


// MarshalJSON implements the json.Marshaler interface for SpecId.
// MarshalJSON for serialization
func (s SpecId) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// UnmarshalJSON for deserialization
func (s *SpecId) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return err
	}

	switch name {
	case "FRONTIER":
		*s = FRONTIER
	case "FRONTIER_THAWING":
		*s = FRONTIER_THAWING
	case "HOMESTEAD":
		*s = HOMESTEAD
	case "DAO_FORK":
		*s = DAO_FORK
	case "TANGERINE":
		*s = TANGERINE
	case "SPURIOUS_DRAGON":
		*s = SPURIOUS_DRAGON
	case "BYZANTIUM":
		*s = BYZANTIUM
	case "CONSTANTINOPLE":
		*s = CONSTANTINOPLE
	case "PETERSBURG":
		*s = PETERSBURG
	case "ISTANBUL":
		*s = ISTANBUL
	case "MUIR_GLACIER":
		*s = MUIR_GLACIER
	case "BERLIN":
		*s = BERLIN
	case "LONDON":
		*s = LONDON
	case "ARROW_GLACIER":
		*s = ARROW_GLACIER
	case "GRAY_GLACIER":
		*s = GRAY_GLACIER
	case "MERGE":
		*s = MERGE
	case "SHANGHAI":
		*s = SHANGHAI
	case "CANCUN":
		*s = CANCUN
	case "PRAGUE":
		*s = PRAGUE
	case "PRAGUE_EOF":
		*s = PRAGUE_EOF
	default:
		return fmt.Errorf("unknown SpecId: %s", name)
	}
	return nil
}

// func SpecToGeneric(specID SpecId) Spec {
// 	switch specID {
// 	case FRONTIER, FRONTIER_THAWING:
// 		return FrontierSpec{NewSpec(FRONTIER)}
// 	case HOMESTEAD, DAO_FORK:
// 		return HomesteadSpec{NewSpec(HOMESTEAD)}
// 	case TANGERINE:
// 		return TangerineSpec{NewSpec(TANGERINE)}
// 	case SPURIOUS_DRAGON:
// 		return SpuriousDragonSpec{NewSpec(SPURIOUS_DRAGON)}
// 	case BYZANTIUM:
// 		return ByzantiumSpec{NewSpec(BYZANTIUM)}
// 	case PETERSBURG, CONSTANTINOPLE:
// 		return PetersburgSpec{NewSpec(PETERSBURG)}
// 	case ISTANBUL, MUIR_GLACIER:
// 		return IstanbulSpec{NewSpec(ISTANBUL)}
// 	case BERLIN:
// 		return BerlinSpec{NewSpec(BERLIN)}
// 	case LONDON, ARROW_GLACIER, GRAY_GLACIER:
// 		return LondonSpec{NewSpec(LONDON)}
// 	case MERGE:
// 		return MergeSpec{NewSpec(MERGE)}
// 	case SHANGHAI:
// 		return ShanghaiSpec{NewSpec(SHANGHAI)}
// 	case CANCUN:
// 		return CancunSpec{NewSpec(CANCUN)}
// 	case PRAGUE:
// 		return PragueSpec{NewSpec(PRAGUE)}
// 	case PRAGUE_EOF:
// 		return PragueEofSpec{NewSpec(PRAGUE_EOF)}
// 	default:
// 		return LatestSpec{NewSpec(LATEST)}
// 	}
// }
/*
	func specToGeneric[H Host, EXT any, DB Database](
		specID SpecId,
		h *EvmHandler[H, EXT, DB],
		isOptimism bool,
	) (*EvmHandler[H, EXT, DB], error) {
		switch specID {
		case FRONTIER, FRONTIER_THAWING:
			return createSpecHandler[H, EXT, DB](h, "frontier", isOptimism)
		case HOMESTEAD, DAO_FORK:
			return createSpecHandler[H, EXT, DB](h, "homestead", isOptimism)
		case TANGERINE:
			return createSpecHandler[H, EXT, DB](h, "tangerine", isOptimism)
		case SPURIOUS_DRAGON:
			return createSpecHandler[H, EXT, DB](h, "spurious_dragon", isOptimism)
		case BYZANTIUM:
			return createSpecHandler[H, EXT, DB](h, "byzantium", isOptimism)
		case PETERSBURG, CONSTANTINOPLE:
			return createSpecHandler[H, EXT, DB](h, "petersburg", isOptimism)
		case ISTANBUL, MUIR_GLACIER:
			return createSpecHandler[H, EXT, DB](h, "istanbul", isOptimism)
		case BERLIN:
			return createSpecHandler[H, EXT, DB](h, "berlin", isOptimism)
		case LONDON, ARROW_GLACIER, GRAY_GLACIER:
			return createSpecHandler[H, EXT, DB](h, "london", isOptimism)
		case MERGE:
			return createSpecHandler[H, EXT, DB](h, "merge", isOptimism)
		case SHANGHAI:
			return createSpecHandler[H, EXT, DB](h, "shanghai", isOptimism)
		case CANCUN:
			return createSpecHandler[H, EXT, DB](h, "cancun", isOptimism)
		case PRAGUE:
			return createSpecHandler[H, EXT, DB](h, "prague", isOptimism)
		case PRAGUE_EOF:
			return createSpecHandler[H, EXT, DB](h, "prague_eof", isOptimism)
		case LATEST:
			return createSpecHandler[H, EXT, DB](h, "latest", isOptimism)
		}
	
		// Optimism-specific specs
		/*
		if isOptimism {
			switch specID {
			case BEDROCK:
				return createSpecHandler[H, EXT, DB](h, "bedrock", true)
			case REGOLITH:
				return createSpecHandler[H, EXT, DB](h, "regolith", true)
			case CANYON:
				return createSpecHandler[H, EXT, DB](h, "canyon", true)
			case ECOTONE:
				return createSpecHandler[H, EXT, DB](h, "ecotone", true)
			case FJORD:
				return createSpecHandler[H, EXT, DB](h, "fjord", true)
			}
		}*/
	
		//return nil, fmt.Errorf("unsupported spec ID: %d", specID)
	//}*/

/*
// Spec interface defining the behavior of Ethereum specs.
type Spec interface {
	GetSpecID() SpecId
}

// CreateSpec creates a new specification struct based on the provided SpecId.
func CreateSpec(specId SpecId) Spec {
	switch specId {
	case FRONTIER:
		return &FrontierSpec{}
	case FRONTIER_THAWING:
		// No changes for EVM spec
		return nil
	case HOMESTEAD:
		return &HomesteadSpec{}
	case DAO_FORK:
		// No changes for EVM spec
		return nil
	case TANGERINE:
		return &TangerineSpec{}
	case SPURIOUS_DRAGON:
		return &SpuriousDragonSpec{}
	case BYZANTIUM:
		return &ByzantiumSpec{}
	case PETERSBURG:
		return &PetersburgSpec{}
	case ISTANBUL:
		return &IstanbulSpec{}
	case MUIR_GLACIER:
		// No changes for EVM spec
		return nil
	case BERLIN:
		return &BerlinSpec{}
	case LONDON:
		return &LondonSpec{}
	case ARROW_GLACIER:
		// No changes for EVM spec
		return nil
	case GRAY_GLACIER:
		// No changes for EVM spec
		return nil
	case MERGE:
		return &MergeSpec{}
	case SHANGHAI:
		return &ShanghaiSpec{}
	case CANCUN:
		return &CancunSpec{}
	case PRAGUE:
		return &PragueSpec{}
	case PRAGUE_EOF:
		return &PragueEofSpec{}
	case LATEST:
		return &LatestSpec{}
	default:
		return nil
	}
}

// Spec structures
type FrontierSpec struct{}
func (s *FrontierSpec) GetSpecID() SpecId { return FRONTIER }

type HomesteadSpec struct{}
func (s *HomesteadSpec) GetSpecID() SpecId { return HOMESTEAD }

type TangerineSpec struct{}
func (s *TangerineSpec) GetSpecID() SpecId { return TANGERINE }

type SpuriousDragonSpec struct{}
func (s *SpuriousDragonSpec) GetSpecID() SpecId { return SPURIOUS_DRAGON }

type ByzantiumSpec struct{}
func (s *ByzantiumSpec) GetSpecID() SpecId { return BYZANTIUM }

type PetersburgSpec struct{}
func (s *PetersburgSpec) GetSpecID() SpecId { return PETERSBURG }

type IstanbulSpec struct{}
func (s *IstanbulSpec) GetSpecID() SpecId { return ISTANBUL }

type BerlinSpec struct{}
func (s *BerlinSpec) GetSpecID() SpecId { return BERLIN }

type LondonSpec struct{}
func (s *LondonSpec) GetSpecID() SpecId { return LONDON }

type MergeSpec struct{}
func (s *MergeSpec) GetSpecID() SpecId { return MERGE }

type ShanghaiSpec struct{}
func (s *ShanghaiSpec) GetSpecID() SpecId { return SHANGHAI }

type CancunSpec struct{}
func (s *CancunSpec) GetSpecID() SpecId { return CANCUN }

type PragueSpec struct{}
func (s *PragueSpec) GetSpecID() SpecId { return PRAGUE }

type PragueEofSpec struct{}
func (s *PragueEofSpec) GetSpecID() SpecId { return PRAGUE_EOF }

type LatestSpec struct{}
func (s *LatestSpec) GetSpecID() SpecId { return LATEST }
*/
