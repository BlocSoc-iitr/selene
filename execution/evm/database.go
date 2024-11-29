package evm
import(
	"errors"
	"math/big"
)
	type Database interface {
		Basic(address Address) (AccountInfo, error)
		BlockHash(number uint64) (B256, error)
		Storage(address Address, index U256) (U256, error)
		CodeByHash(codeHash B256) (Bytecode, error)	
	}
	type EmptyDB struct{}
	func (db *EmptyDB) Basic(address Address) (AccountInfo, error) {
		return AccountInfo{}, errors.New("EmptyDB: Basic not implemented")
	}
	func (db *EmptyDB) BlockHash(number uint64) (B256, error) {
		return B256{}, errors.New("EmptyDB: BlockHash not implemented")
	}
	func (db *EmptyDB) Storage(address Address, index U256) (U256, error) {
		return big.NewInt(0), errors.New("EmptyDB: Storage not implemented")
	}
	func (db *EmptyDB) CodeByHash(codeHash B256) (Bytecode, error) {
		return Bytecode{}, errors.New("EmptyDB: CodeByHash not implemented")
	}
	func NewEmptyDB() *EmptyDB {
		return &EmptyDB{}
	}
type DatabaseError struct {
	msg string
	err error
}
func (e *DatabaseError) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return e.msg
}