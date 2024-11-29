package evm
type Handler[H Host, EXT any, DB Database] struct {
	Cfg              HandlerCfg
	InstructionTable InstructionTables[Context[EXT, DB]]
	Registers        []HandleRegisters[H, EXT, DB]
	Validation       ValidationHandler[EXT, DB]
	PreExecution     PreExecutionHandler[EXT, DB]
	PostExecution    PostExecutionHandler[EXT, DB]
	Execution        ExecutionHandler[EXT, DB]
}
func NewEvmHandler[H Host, EXT any, DB Database](cfg HandlerCfg) (*EvmHandler[H, EXT, DB], error) {
	h := &EvmHandler[H, EXT, DB]{Cfg: cfg}
	if h.Cfg.isOptimism {
		return h.optimismWithSpec(cfg.specID)
	}
	return h.mainnetWithSpec(cfg.specID)
}
func (h *EvmHandler[H, EXT, DB]) optimismWithSpec(specid SpecId) (*EvmHandler[H, EXT, DB], error) {
	return &EvmHandler[H, EXT, DB]{
        Cfg: HandlerCfg{
            specID: specid,
            isOptimism: true,
        },
        InstructionTable: NewPlainInstructionTable(InstructionTable[Context[EXT, DB]]{}),
        Registers: []HandleRegisters[H, EXT, DB]{},
        Validation: ValidationHandler[EXT, DB]{},
        PreExecution: PreExecutionHandler[EXT, DB]{},
        PostExecution: PostExecutionHandler[EXT, DB]{},
        Execution: ExecutionHandler[EXT, DB]{},
    }, nil
}
func (h *EvmHandler[H, EXT, DB]) mainnetWithSpec(specid SpecId) (*EvmHandler[H, EXT, DB], error) {
	return &EvmHandler[H, EXT, DB]{
        Cfg: HandlerCfg{
            specID: specid,
            isOptimism: false,
        },
        InstructionTable: NewPlainInstructionTable(InstructionTable[Context[EXT, DB]]{}),
        Registers: []HandleRegisters[H, EXT, DB]{},
        Validation: ValidationHandler[EXT, DB]{},
        PreExecution: PreExecutionHandler[EXT, DB]{},
        PostExecution: PostExecutionHandler[EXT, DB]{},
        Execution: ExecutionHandler[EXT, DB]{},
    }, nil
}
type HandlerCfg struct {
	specID     SpecId
	isOptimism bool
}
func NewHandlerCfg(specID SpecId) HandlerCfg {
	return HandlerCfg{
		specID:     specID,
		isOptimism: getDefaultOptimismSetting(),
	}
}
//Note there might be error in specToGeneric needs to be well tested

/*
Enabling run time flexibility
# Run with optimism disabled (default)
go run main.go

# Run with optimism enabled
OPTIMISM_ENABLED=true go run main.go*/
// SpecId is the enumeration of various Ethereum hard fork specifications.

// PhantomData is a Go struct that simulates the behavior of Rust's PhantomData.
func (s SpecId) Enabled(a SpecId, b SpecId) bool {
	return uint8(a) >= uint8(b)
}
func (h *Handler[H, EXT, DB]) ExecuteFrame(frame *Frame, sharedMemory *SharedMemory, context *Context[EXT, DB]) (InterpreterAction, error) {
	return h.Execution.ExecuteFrame(frame, sharedMemory, &h.InstructionTable, context)
}
func (h *EvmHandler[H, EXT, DB]) specId() SpecId {
	return h.Cfg.specID
}
func (h *EvmHandler[H, EXT, DB]) IsOptimism() bool {
	return isOptimismEnabled && h.Cfg.isOptimism
}

// EvmHandler is a type alias for Handler with specific type parameters.
type EvmHandler[H Host, EXT any, DB Database] Handler[H, EXT, DB]

// HandleRegister defines a function type that accepts a pointer to EvmHandler.
type HandleRegister[H Host, EXT any, DB Database] func(handler *EvmHandler[H, EXT, DB])

// HandleRegisterBox defines a function type as a boxed register.
type HandleRegisterBox[H Host, EXT any, DB Database] func(handler *EvmHandler[H, EXT, DB])

// HandleRegisters is an interface that defines the `Register` method.
type HandleRegisters[H Host, EXT any, DB Database] interface {
	Register(handler *EvmHandler[H, EXT, DB])
}

// PlainRegister struct implements HandleRegisters with a plain function register.
type PlainRegister[H Host, EXT any, DB Database] struct {
	RegisterFn HandleRegister[H, EXT, DB]
}

// BoxRegister struct implements HandleRegisters with a boxed function register.
type BoxRegister[H Host, EXT any, DB Database] struct {
	RegisterFn HandleRegisterBox[H, EXT, DB]
}

// Register method for PlainRegister, to satisfy HandleRegisters interface.
func (p PlainRegister[H, EXT, DB]) Register(handler *EvmHandler[H, EXT, DB]) {
	p.RegisterFn(handler)
}

// Register method for BoxRegister, to satisfy HandleRegisters interface.
func (b BoxRegister[H, EXT, DB]) Register(handler *EvmHandler[H, EXT, DB]) {
	b.RegisterFn(handler)
}
/*
func (s SpecId) IsEnabledIn(spec SpecId) bool {
	return s.Enabled(s, spec)
}
*/
// String method to provide a string representation for each CallScheme variant.
// func (cs CallScheme) String() string {
// 	switch cs {
// 	case Call_Scheme:
// 		return "Call"
// 	case CallCode:
// 		return "CallCode"
// 	case DelegateCall:
// 		return "DelegateCall"
// 	case StaticCall:
// 		return "StaticCall"
// 	case ExtCall:
// 		return "ExtCall"
// 	case ExtStaticCall:
// 		return "ExtStaticCall"
// 	case ExtDelegateCall:
// 		return "ExtDelegateCall"
// 	default:
// 		return "Unknown"
// 	}
// }

/*
func (s Spec) SpecID() SpecId {
    return s.SPEC_ID
}
*/
// IsEnabled checks if a specific SpecId is enabled.
/*func (s Spec) IsEnabled(specId SpecId) bool {
    return SpecIdEnabled(s.SpecID(), specId)
}
func SpecIdEnabled(our, other SpecId) bool {
    return our >= other
}*/
/*
func NewHandler[H Host, EXT any, DB Database](cfg HandlerCfg) Handler[H, EXT, DB] {
    return createHandlerWithConfig[H, EXT, DB](cfg)
}*/
/*
func createHandlerWithConfig[H Host, EXT any, DB Database](cfg HandlerCfg) Handler[H, EXT, DB] {
    spec := GetSpecforID[H, EXT, DB](cfg.specID)

    if cfg.isOptimism {
        return createOptimismHandler(spec)
    }
    return createMainnetHandler(spec)
}
*/

/*
func (h *EvmHandler[H,EXT, DB])optimism() *EvmHandler[H,EXT, DB]{
    handler:=
}*/
//Doubt in newplaininstructiontable
/*
func (h *EvmHandler[H,EXT, DB])mainnet(spec Spec) *EvmHandler[H,EXT, DB]{
    specID := spec.SpecID()

    return &EvmHandler[H, EXT, DB]{
        Cfg: NewHandlerCfg(specID),
        InstructionTable: NewPlainInstructionTable(InstructionTable[H]{}),
        Registers: []HandleRegister[H, EXT, DB]{},


    }
}*/
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
    }

    return nil, fmt.Errorf("unsupported spec ID: %d", specID)
}

*/
/*
// Concrete spec implementations
type FrontierSpec struct{}
type HomesteadSpec struct{}
type TangerineSpec struct{}
type SpuriousDragonSpec struct{}
type ByzantiumSpec struct{}
type PetersburgSpec struct{}
type IstanbulSpec struct{}
type BerlinSpec struct{}
type LondonSpec struct{}
type MergeSpec struct{}
type ShanghaiSpec struct{}
type CancunSpec struct{}
type PragueSpec struct{}
type PragueEofSpec struct{}
type LatestSpec struct{}

// Optimism-specific specs
type BedrockSpec struct{}
type RegolithSpec struct{}
type CanyonSpec struct{}
type EcotoneSpec struct{}
type FjordSpec struct{}

func (s *FrontierSpec) Name() string         { return "frontier" }
func (s *FrontierSpec) Version() uint        { return 1 }
func (s *HomesteadSpec) Name() string        { return "homestead" }
func (s *HomesteadSpec) Version() uint       { return 1 }
func (s *TangerineSpec) Name() string        { return "tangerine" }
func (s *TangerineSpec) Version() uint       { return 1 }
func (s *SpuriousDragonSpec) Name() string   { return "spurious_dragon" }
func (s *SpuriousDragonSpec) Version() uint  { return 1 }
func (s *ByzantiumSpec) Name() string        { return "byzantium" }
func (s *ByzantiumSpec) Version() uint       { return 1 }
func (s *PetersburgSpec) Name() string       { return "petersburg" }
func (s *PetersburgSpec) Version() uint      { return 1 }
func (s *IstanbulSpec) Name() string         { return "istanbul" }
func (s *IstanbulSpec) Version() uint        { return 1 }
func (s *BerlinSpec) Name() string           { return "berlin" }
func (s *BerlinSpec) Version() uint          { return 1 }
func (s *LondonSpec) Name() string           { return "london" }
func (s *LondonSpec) Version() uint          { return 1 }
func (s *MergeSpec) Name() string            { return "merge" }
func (s *MergeSpec) Version() uint           { return 1 }
func (s *ShanghaiSpec) Name() string         { return "shanghai" }
func (s *ShanghaiSpec) Version() uint        { return 1 }
func (s *CancunSpec) Name() string           { return "cancun" }
func (s *CancunSpec) Version() uint          { return 1 }
func (s *PragueSpec) Name() string           { return "prague" }
func (s *PragueSpec) Version() uint          { return 1 }
func (s *PragueEofSpec) Name() string        { return "prague_eof" }
func (s *PragueEofSpec) Version() uint       { return 1 }
func (s *LatestSpec) Name() string           { return "latest" }
func (s *LatestSpec) Version() uint          { return 1 }
func (s *BedrockSpec) Name() string          { return "bedrock" }
func (s *BedrockSpec) Version() uint         { return 1 }
func (s *RegolithSpec) Name() string         { return "regolith" }
func (s *RegolithSpec) Version() uint        { return 1 }
func (s *CanyonSpec) Name() string           { return "canyon" }
func (s *CanyonSpec) Version() uint          { return 1 }
func (s *EcotoneSpec) Name() string          { return "ecotone" }
func (s *EcotoneSpec) Version() uint         { return 1 }
func (s *FjordSpec) Name() string            { return "fjord" }
func (s *FjordSpec) Version() uint           { return 1 }
*/
// createSpecHandler now creates the appropriate spec implementation
/*
func createSpecHandler[H Host, EXT any, DB Database](
    h *EvmHandler[H, EXT, DB],
    specName string,
    isOptimism bool,
) (*EvmHandler[H, EXT, DB], error) {
    newHandler := *h

    var spec Spec
    switch specName {
    case "frontier":
        spec = &FrontierSpec{}
    case "homestead":
        spec = &HomesteadSpec{}
    case "tangerine":
        spec = &TangerineSpec{}
    case "spurious_dragon":
        spec = &SpuriousDragonSpec{}
    case "byzantium":
        spec = &ByzantiumSpec{}
    case "petersburg":
        spec = &PetersburgSpec{}
    case "istanbul":
        spec = &IstanbulSpec{}
    case "berlin":
        spec = &BerlinSpec{}
    case "london":
        spec = &LondonSpec{}
    case "merge":
        spec = &MergeSpec{}
    case "shanghai":
        spec = &ShanghaiSpec{}
    case "cancun":
        spec = &CancunSpec{}
    case "prague":
        spec = &PragueSpec{}
    case "prague_eof":
        spec = &PragueEofSpec{}
    case "latest":
        spec = &LatestSpec{}
    /*
    case "bedrock":
        spec = &BedrockSpec{}
    case "regolith":
        spec = &RegolithSpec{}
    case "canyon":
        spec = &CanyonSpec{}
    case "ecotone":
        spec = &EcotoneSpec{}
    case "fjord":
        spec = &FjordSpec{}
*/
/*
    default:
        return nil, fmt.Errorf("unknown spec: %s", specName)
    }

    // Add the spec to the handler
    newHandler = spec
    return &newHandler, nil
}
*/
/*
func withSpec[EXT any, DB Database, S Spec](h *EVMHandler[EXT, DB]) (*EVMHandler[EXT, DB], error) {
    var spec S
    rules := spec.Rules()

    // Create a new handler with the specific spec configuration
    newHandler := &EVMHandler[EXT, DB]{
        cfg: h.cfg,
        spec: spec,
        rules: *rules,
    }

    return newHandler, nil
}*/

/*
type HandleRegister[EXT any, DB Database] func(handler *EvmHandler[EXT, DB])
type EvmHandler[EXT any, DB Database] Handler[EXT, DB]
type HandleRegisterBox[EXT any, DB Database] interface {
    Call(handler *EvmHandler[EXT, DB])
}
type HandleRegisters[EXT any, DB Database] struct {
    Plain HandleRegister[EXT, DB]      // Plain function register
    Box   HandleRegisterBox[EXT, DB]   // Boxed function register
}*/
