package evm
import "bytes"
type Evm[EXT interface{}, DB Database] struct {
	Context Context[EXT, DB]
	Handler Handler[Context[EXT, DB], EXT, DB]
}
func NewEvm[EXT any, DB Database](context Context[EXT, DB], handler Handler[Context[EXT, DB], EXT, DB]) Evm[EXT, DB] {
	context.Evm.Inner.JournaledState.SetSpecId(handler.Cfg.specID)
	return Evm[EXT, DB]{
		Context: context,
		Handler: handler,
	}
}
func (e Evm[EXT, DB]) specId() SpecId {
	return e.Handler.Cfg.specID
}
func (e Evm[EXT, DB]) IntoContextWithHandlerCfg() ContextWithHandlerCfg[EXT, DB] {
    return ContextWithHandlerCfg[EXT, DB]{
        Context: e.Context,
        Cfg:     e.Handler.Cfg,
    }
}
var EOF_MAGIC_BYTES = []byte{0xef, 0x00}
type EVMResult[D any] EVMResultGeneric[ResultAndState, D]
func (e *Evm[EXT, DB]) Transact() EVMResult[DatabaseError] {
	initialGasSpend, err := e.preverifyTransactionInner()
	if len(err.Message) != 0 {
		e.clear()
		return EVMResult[DatabaseError]{Err: err}
	}
	output := e.TransactPreverifiedInner(initialGasSpend)
	output2, err := e.Handler.PostExecution.End(&e.Context, EVMResultGeneric[ResultAndState, EvmError](output))
	e.clear()
	return EVMResult[DatabaseError]{Err: err, Value: output2}
}
func (e *Evm[EXT, DB]) preverifyTransactionInner() (uint64, EvmError) {
	err := e.Handler.Validation.Env(e.Context.Evm.Inner.Env) //?
	if err != nil {
		return 0, EvmError{Message: err.Error()}
	}
	initialGasSpend, err := e.Handler.Validation.InitialTxGas(e.Context.Evm.Inner.Env)
	err = e.Handler.Validation.TxAgainstState(&e.Context)
	return initialGasSpend, EvmError{}
}
func (e *Evm[EXT, DB]) clear() {
	e.Handler.PostExecution.Clear(&e.Context)
}
func (e *Evm[EXT, DB]) TransactPreverifiedInner(initialGasSpend uint64) EVMResult[DatabaseError] {
	specId := e.Handler.Cfg.specID
	ctx := &e.Context
	preExec := e.Handler.PreExecution
	err := preExec.LoadAccounts(ctx)
	if err != nil {
		return EVMResult[DatabaseError]{Err: err}
	}
	precompiles := preExec.LoadPrecompilesFunction()
	ctx.Evm.SetPrecompiles(precompiles)
	err_caller := preExec.DeductCaller(ctx)	//?
	if err_caller.Err != nil {
		return EVMResult[DatabaseError](err_caller)
	}
	gasLimit := ctx.Evm.Inner.Env.Tx.GasLimit - initialGasSpend
	exec := e.Handler.Execution
	var firstFrameOrResult FrameOrResult
		if ctx.Evm.Inner.Env.Tx.TransactTo.Type == Call2 {
			result, err := exec.Call(ctx, (&CallInputs{}).NewBoxed(&ctx.Evm.Inner.Env.Tx, gasLimit))
			if len(err.Message) != 0 {
				return EVMResult[DatabaseError]{Err: err}
			}
			firstFrameOrResult = result
		} else {
			if specId.IsEnabledIn(PRAGUE_EOF) && bytes.Equal(ctx.Evm.Inner.Env.Tx.Data[:2], EOF_MAGIC_BYTES) {
				result, err := exec.EOFCreate(ctx, EOFCreateInputs{}.NewTx(&ctx.Evm.Inner.Env.Tx, gasLimit))
				if len(err.Message) != 0 {
					return EVMResult[DatabaseError]{Err: err}
				}
				firstFrameOrResult = result
			} else {
				result, err := exec.Create(ctx, (&CreateInputs{}).NewBoxed(&ctx.Evm.Inner.Env.Tx, gasLimit))
				if len(err.Message) != 0 {
					return EVMResult[DatabaseError]{Err: err}
				}
				firstFrameOrResult = result
			}
		}
	var result FrameResult
	if firstFrameOrResult.Type == Frame_FrameOrResult {
		result_new, err_new := e.RunTheLoop(firstFrameOrResult.Frame)
		if len(err_new.Message) != 0 {
			return EVMResult[DatabaseError]{Err: err_new}
		}
		result = result_new
	} else {
		result = firstFrameOrResult.Result
	}
	ctx = &e.Context
	err = e.Handler.Execution.LastFrameReturn(ctx, &result)
	postExec := e.Handler.PostExecution
	postExec.ReimburseCaller(ctx, result.Gas())	//!
	postExec.RewardBeneficiary(ctx, result.Gas())	//!
	res, err := postExec.Output(ctx, result)
	return EVMResult[DatabaseError]{Value: res, Err: err}
}
func (e *Evm[EXT, DB]) RunTheLoop(firstFrame Frame) (FrameResult, EvmError) {
	callStack := make([]Frame, 0, 1025)
	callStack = append(callStack, firstFrame)
	
	var sharedMemory SharedMemory
	if e.Context.Evm.Inner.Env.Cfg.MemoryLimit > 0 {
		sharedMemory = SharedMemory{}.NewWithMemoryLimit(e.Context.Evm.Inner.Env.Cfg.MemoryLimit)
	} else {
		sharedMemory = SharedMemory{}.New()
	}
	sharedMemory.NewContext()

	stackFrame := &callStack[len(callStack) - 1]

	for {
		nextAction, err := e.Handler.ExecuteFrame(stackFrame,&sharedMemory,&e.Context)
		if err != nil {
			return FrameResult{}, EvmError{Message: err.Error()}
		}
		// Take error from the context, if any
		if err := e.Context.Evm.Inner.TakeError(); len(err.Message) != 0 {
			return FrameResult{}, err
		}
		exec := &e.Handler.Execution
		var frameOrResult FrameOrResult 
		{
			if nextAction.actionType == ActionTypeCall {
				res, err := exec.Call(&e.Context, nextAction.callInputs)
				if len(err.Message) != 0 {
					return FrameResult{}, err
				}
				frameOrResult = res
			} else if nextAction.actionType == ActionTypeCreate {
				res, err := exec.Create(&e.Context, nextAction.createInputs)
				if len(err.Message) != 0 {
					return FrameResult{},err
				}
				frameOrResult = res
			} else if nextAction.actionType == ActionTypeEOFCreate {
				res, err := exec.EOFCreate(&e.Context, *nextAction.eofCreateInputs)
				if len(err.Message) != 0 {
					return FrameResult{},err
				}
				frameOrResult = res
			} else if nextAction.actionType == ActionTypeReturn {
				sharedMemory.FreeContext()
				returnedFrame := callStack[len(callStack) - 1]
				callStack = callStack[:len(callStack) - 1]
				ctx := &e.Context

				var result FrameResult
				frameOrResult2 := FrameOrResult{
					Type: Result_FrameOrResult,
					Result: result,
				}
				switch returnedFrame.Type {
				case Call_Frame:
					result.ResultType = Call
					res, err := exec.CallReturn(ctx, returnedFrame.Call, *nextAction.result)
					if len(err.Message) != 0 {
						return FrameResult{}, err
					}
					result.Call = &res
				case Create_Frame:
					result.ResultType = Create
					res, err := exec.CreateReturn(ctx, returnedFrame.Create, *nextAction.result)
					if len(err.Message) != 0 {
						return FrameResult{}, err
					}
					result.Create = &res
				case EOFCreate_Frame:
					result.ResultType = EOFCreate
					res, err :=  exec.EOFCreateReturn(ctx, returnedFrame.EOFCreate, *nextAction.result)
					if len(err.Message) != 0 {
						return FrameResult{}, err
					}
					result.Create = &res
				}
				frameOrResult = frameOrResult2
			} else {
				return FrameResult{}, EvmError{
					Message: "InterpreterAction::None is not expected",
				}
			}
		}
		switch frameOrResult.Type {
		case Frame_FrameOrResult:
			sharedMemory.NewContext()
			callStack = append(callStack, frameOrResult.Frame)
			stackFrame = &callStack[len(callStack) - 1]
		case Result_FrameOrResult:
			topFrame := &callStack[len(callStack) - 1]
			if topFrame == nil {
				return frameOrResult.Result, EvmError{}
			}
			stackFrame = topFrame
			ctx := &e.Context
			switch frameOrResult.Result.ResultType {
			case Call:
				exec.InsertCallOutcome(ctx, stackFrame, &sharedMemory, *frameOrResult.Result.Call)
			case Create:
				exec.InsertCreateOutcome(ctx, stackFrame, *frameOrResult.Result.Create)
			case EOFCreate:
				exec.InsertCreateOutcome(ctx, stackFrame, *frameOrResult.Result.Create)
			}
		}
	}
}