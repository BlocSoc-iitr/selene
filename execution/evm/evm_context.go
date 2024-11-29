package evm
type EvmContext[DB Database] struct {
	Inner InnerEvmContext[DB]
	Precompiles ContextPrecompiles[DB]
}
func NewEvmContext[DB Database](db DB) EvmContext[DB] {
	return EvmContext[DB]{
		Inner: NewInnerEvmContext[DB](db),
		Precompiles: DefaultContextPrecompiles[DB](),
	}
}
func WithNewEvmDB[DB1, DB2 Database](
    ec EvmContext[DB1],
    db DB2,
) EvmContext[DB2] {
    return EvmContext[DB2]{
        Inner:	   WithNewInnerDB(ec.Inner,db),
        Precompiles: DefaultContextPrecompiles[DB2](),
    }
}
