package evm
type PreExecutionHandler[EXT any, DB Database] struct {
	LoadPrecompiles LoadPrecompilesHandle[DB]
	LoadAccounts    LoadAccountsHandle[EXT, DB]
	DeductCaller    DeductCallerHandle[EXT, DB]
}
type ValidationHandler[EXT any, DB Database] struct {
	InitialTxGas   ValidateInitialTxGasHandle[DB]
	TxAgainstState ValidateTxEnvAgainstState[EXT, DB]
	Env            ValidateEnvHandle[DB]
}
type ValidateEnvHandle[DB Database] func(env *Env) error
type ValidateTxEnvAgainstState[EXT any, DB Database] func(ctx *Context[EXT, DB]) error
type ValidateInitialTxGasHandle[DB Database] func(env *Env) (uint64, error)
func (p *PreExecutionHandler[EXT, DB]) LoadPrecompilesFunction() ContextPrecompiles[DB] {
	return p.LoadPrecompiles()
}
