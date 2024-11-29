//go:build optimism
// +build optimism

package evm

const isOptimismEnabled = true
import (
	"encoding/json"
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
    // Optimism fields (only included when build tag is set)
    Optimism OptimismFields `json:"optimism,omitempty"`
}
func TryFromUint8(specID uint8) (SpecId, bool) {
	if specID > uint8(PRAGUE_EOF) && specID != uint8(LATEST) {
		return 0, false
	}
	return SpecId(specID), true
}
type SpecId uint8

const (
    FRONTIER       SpecId = 0   // Frontier
    FRONTIER_THAWING SpecId = 1   // Frontier Thawing
    HOMESTEAD      SpecId = 2   // Homestead
    DAO_FORK       SpecId = 3   // DAO Fork
    TANGERINE      SpecId = 4   // Tangerine Whistle
    SPURIOUS_DRAGON SpecId = 5   // Spurious Dragon
    BYZANTIUM      SpecId = 6   // Byzantium
    CONSTANTINOPLE SpecId = 7   // Constantinople
    PETERSBURG     SpecId = 8   // Petersburg
    ISTANBUL       SpecId = 9   // Istanbul
    MUIR_GLACIER   SpecId = 10  // Muir Glacier
    BERLIN         SpecId = 11  // Berlin
    LONDON         SpecId = 12  // London
    ARROW_GLACIER  SpecId = 13  // Arrow Glacier
    GRAY_GLACIER   SpecId = 14  // Gray Glacier
    MERGE          SpecId = 15  // Paris/Merge
    BEDROCK        SpecId = 16  // Bedrock
    REGOLITH       SpecId = 17  // Regolith
    SHANGHAI       SpecId = 18  // Shanghai
    CANYON         SpecId = 19  // Canyon
    CANCUN         SpecId = 20  // Cancun
    ECOTONE        SpecId = 21  // Ecotone
    FJORD          SpecId = 22  // Fjord
    PRAGUE         SpecId = 23  // Prague
    PRAGUE_EOF     SpecId = 24  // Prague+EOF
    LATEST         SpecId = 255 // LATEST = u8::MAX
)

// Default value for SpecId
func DefaultSpecId() SpecId {
	return LATEST
}

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
    case BEDROCK:
        return "BEDROCK"
    case REGOLITH:
        return "REGOLITH"
    case SHANGHAI:
        return "SHANGHAI"
    case CANYON:
        return "CANYON"
    case CANCUN:
        return "CANCUN"
    case ECOTONE:
        return "ECOTONE"
    case FJORD:
        return "FJORD"
    case PRAGUE:
        return "PRAGUE"
    case PRAGUE_EOF:
        return "PRAGUE_EOF"
    default:
        return "UNKNOWN"
    }
}
type (
    FrontierSpec       struct{ BaseSpec }
    FrontierThawingSpec struct{ BaseSpec }
    HomesteadSpec      struct{ BaseSpec }
    DaoForkSpec        struct{ BaseSpec }
    TangerineSpec      struct{ BaseSpec }
    SpuriousDragonSpec struct{ BaseSpec }
    ByzantiumSpec      struct{ BaseSpec }
    ConstantinopleSpec struct{ BaseSpec }
    PetersburgSpec     struct{ BaseSpec }
    IstanbulSpec       struct{ BaseSpec }
    MuirGlacierSpec    struct{ BaseSpec }
    BerlinSpec         struct{ BaseSpec }
    LondonSpec         struct{ BaseSpec }
    ArrowGlacierSpec   struct{ BaseSpec }
    GrayGlacierSpec    struct{ BaseSpec }
    MergeSpec          struct{ BaseSpec }
    BedrockSpec        struct{ BaseSpec }
    RegolithSpec       struct{ BaseSpec }
    ShanghaiSpec       struct{ BaseSpec }
    CanyonSpec         struct{ BaseSpec }
    CancunSpec         struct{ BaseSpec }
    EcotoneSpec        struct{ BaseSpec }
    FjordSpec          struct{ BaseSpec }
    PragueSpec         struct{ BaseSpec }
    PragueEofSpec      struct{ BaseSpec }
    LatestSpec         struct{ BaseSpec }
)
// MarshalJSON for serialization
func (s SpecId) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}



func SpecToGeneric(specID SpecId) interface{} {
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
    case BEDROCK:
        return BedrockSpec{NewSpec(BEDROCK)}
    case REGOLITH:
        return RegolithSpec{NewSpec(REGOLITH)}
    case CANYON:
        return CanyonSpec{NewSpec(CANYON)}
    case ECOTONE:
        return EcotoneSpec{NewSpec(ECOTONE)}
    case FJORD:
        return FjordSpec{NewSpec(FJORD)}
    default:
        return LatestSpec{NewSpec(LATEST)}
    }
}

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
    case "BEDROCK":
        *s = BEDROCK
    case "REGOLITH":
        *s = REGOLITH
    case "CANYON":
        *s = CANYON
    case "ECOTONE":
        *s = ECOTONE
    case "FJORD":
        *s = FJORD
    default:
        return fmt.Errorf("unknown SpecId: %s", name)
    }
    return nil
}