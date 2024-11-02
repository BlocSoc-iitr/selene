package evm
type PostExecutionHandler[EXT any, DB Database] struct {
	ReimburseCaller   ReimburseCallerHandle[EXT, DB]
	RewardBeneficiary RewardBeneficiaryHandle[EXT, DB]
	Output            OutputHandle[EXT, DB]
	End               EndHandle[EXT, DB]
	Clear             ClearHandle[EXT, DB]
}
type ReimburseCallerHandle[EXT any, DB Database] func(ctx *Context[EXT, DB], gas *Gas) EVMResultGeneric[struct{}, any]
type RewardBeneficiaryHandle[EXT any, DB Database] ReimburseCallerHandle[EXT, DB]
type OutputHandle[EXT any, DB Database] func(ctx *Context[EXT, DB], frameResult FrameResult) (ResultAndState, EvmError)
type EndHandle[EXT any, DB Database] func(ctx *Context[EXT, DB], result EVMResultGeneric[ResultAndState, EvmError]) (ResultAndState, EvmError)
type ClearHandle[EXT any, DB Database] func(ctx *Context[EXT, DB])
type LoadPrecompilesHandle[DB Database] func() ContextPrecompiles[DB]
type LoadAccountsHandle[EXT any, DB Database] func(ctx *Context[EXT, DB]) error
type DeductCallerHandle[EXT any, DB Database] func(ctx *Context[EXT, DB]) EVMResultGeneric[ ResultAndState, DatabaseError]
type EVMResultGeneric[T any, DBError any] struct {
	Value T
	Err   error
}
