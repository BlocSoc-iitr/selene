package evm
type EvmBuilder[EXT interface{}, DB Database] struct {
	context Context[EXT, DB]
	handler Handler[Context[EXT, DB], EXT, DB]
	phantom struct{}
}
func NewDefaultEvmBuilder[EXT interface{}]() *EvmBuilder[EXT, *EmptyDB] {
    handlerCfg := NewHandlerCfg(LATEST)
    return &EvmBuilder[EXT, *EmptyDB]{
        context: DefaultContext[EXT](),
        handler: Handler[Context[EXT, *EmptyDB], EXT, *EmptyDB]{Cfg: handlerCfg},
        phantom: struct{}{},
    }
}
func WithNewDB[EXT interface{}, DB1, DB2 Database](
    eb *EvmBuilder[EXT, DB1],
    db DB2,
) *EvmBuilder[EXT, DB2] {
    return &EvmBuilder[EXT, DB2]{
        context: NewContext[EXT, DB2](
            WithNewEvmDB(eb.context.Evm, db),
            eb.context.External.(EXT),
        ),
        handler: Handler[Context[EXT, DB2], EXT, DB2]{Cfg: eb.handler.Cfg},
        phantom: struct{}{},
    }
}
func (eb *EvmBuilder[EXT, DB]) WithEnv(env *Env) *EvmBuilder[EXT, DB] {
	eb.context.Evm.Inner.Env = env
	return eb
}

func (eb *EvmBuilder[EXT, DB]) Build() Evm[EXT,DB] {
	return NewEvm(eb.context, eb.handler)
}

func WithContextWithHandlerCfg[EXT interface{},DB1,DB2 Database](eb *EvmBuilder[EXT, DB1],contextWithHandlerCfg ContextWithHandlerCfg[EXT, DB2]) *EvmBuilder[EXT, DB2] {
	return &EvmBuilder[EXT, DB2]{
		context: contextWithHandlerCfg.Context,
		handler: Handler[Context[EXT, DB2], EXT, DB2]{Cfg: contextWithHandlerCfg.Cfg},
		phantom: struct{}{},
	}
}
//Not in use
func (eb EvmBuilder[EXT, DB]) WithDB(db DB) *EvmBuilder[EXT, DB] {
	return &EvmBuilder[EXT, DB]{
		context: NewContext[EXT, DB](eb.context.Evm.WithDB(db), eb.context.External.(EXT)), //Doubt
		handler: Handler[Context[EXT, DB], EXT, DB]{Cfg: eb.handler.Cfg},                                              //Doubt
		phantom: struct{}{},
	}
}
//Not in use
type PhantomData[T any] struct{}
//Not in use
/*
func WithContextWithHandlerCfg[EXT interface{}, DB Database](eb *EvmBuilder[EXT, Database], contextWithHandlerCfg ContextWithHandlerCfg[EXT, Database]) *EvmBuilder[EXT, Database] {
    return &EvmBuilder[EXT, Database]{
        context: contextWithHandlerCfg.Context,
        handler: Handler[Context[EXT, Database], EXT, Database]{Cfg: contextWithHandlerCfg.Cfg},
        phantom: struct{}{},
    }
}*/
/*
func (eb EvmBuilder[EXT, DB]) WithContextWithHandlerCfg(contextWithHandlerCfg ContextWithHandlerCfg[EXT, DB]) *EvmBuilder[EXT, DB] {
	return &EvmBuilder[EXT, DB]{
		context: contextWithHandlerCfg.Context,
		handler: Handler[Context[EXT, DB], EXT, DB]{Cfg: contextWithHandlerCfg.Cfg},
		phantom: struct{}{},
	}
}*/


// Tuple implementation in golang:

// LogData represents the data structure of a log entry.


