package consensus

// uses config mmodule
type Database interface{
    new()
    save_checkpoint()
    load_checkpoint()
}

//NOTE: parameters are not included
type FileDB struct{Database}

func (f FileDB)new() FileDB{}
func (f FileDB)save_checkpoint() FileDB{}
func (f FileDB)load_checkpoint() FileDB{}



type ConfigDB struct{Database}
func (cdb ConfigDB)new() ConfigDB{}
func (cdb ConfigDB)save_checkpoint() ConfigDB{}
func (cdb ConfigDB)load_checkpoint() ConfigDB{}