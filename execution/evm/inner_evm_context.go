package evm
type InnerEvmContext[DB Database] struct {
    Env                    *Env
    JournaledState        JournaledState
    DB                     DB
    Error                  error
    ValidAuthorizations    []Address
    L1BlockInfo           *L1BlockInfo // For optimism feature
}
func NewInnerEvmContext[DB Database](db DB) InnerEvmContext[DB] {
	SpecId:=DefaultSpecId()
	return InnerEvmContext[DB]{
		Env: &Env{},
		JournaledState: NewJournalState(SpecId,make(map[Address]struct{})),//to be changed
		DB: db,
		Error: nil,
		ValidAuthorizations: nil,
		L1BlockInfo: nil,
	}
}
func WithNewInnerDB[DB1, DB2 Database](
    inner InnerEvmContext[DB1],
    db DB2,
) InnerEvmContext[DB2] {
    return InnerEvmContext[DB2]{
        Env:                 inner.Env,
        JournaledState:     inner.JournaledState,
        DB:                 db,
        Error:              inner.Error,
        ValidAuthorizations: nil,
        L1BlockInfo:        inner.L1BlockInfo,
    }
}
func (i *InnerEvmContext[DB]) TakeError() EvmError {
    currentError := i.Error
    i.Error = nil
    return EvmError{Message: currentError.Error()}
}
type L1BlockInfo struct {
	L1BaseFee           U256
	L1FeeOverhead       *U256
	L1BaseFeeScalar     U256
	L1BlobBaseFee       *U256
	L1BlobBaseFeeScalar *U256
	EmptyScalars        bool
}
//Not being used in the code
func (inner InnerEvmContext[DB]) WithDB(db DB) InnerEvmContext[DB] {
	return InnerEvmContext[DB]{
		Env:                 inner.Env,
		JournaledState:      inner.JournaledState,
		DB:                  db,
		Error:               inner.Error,
		ValidAuthorizations: nil,
		L1BlockInfo:         inner.L1BlockInfo,
	}
}
//Not being used in the code
func (inner InnerEvmContext[DB]) specId() SpecId {
	return inner.JournaledState.Spec
}

