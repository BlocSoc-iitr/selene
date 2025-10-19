package evm

import (
	"errors"
	"math/big"
	"reflect"
	"testing"
)

func TestNewEmptyDB(t *testing.T) {
    db := NewEmptyDB()
    if db == nil {
        t.Error("NewEmptyDB() returned nil")
    }
}

func TestEmptyDBBasic(t *testing.T) {
    db := NewEmptyDB()
    
    testCases := []struct {
        name    string
        address Address
        wantErr string
    }{
        {
            name:    "basic query",
            address: Address{}, // zero address
            wantErr: "EmptyDB: Basic not implemented",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            info, err := db.Basic(tc.address)
            
            // Check if error matches expected
            if err == nil || err.Error() != tc.wantErr {
                t.Errorf("Basic() error = %v, want %v", err, tc.wantErr)
            }
            
            // Verify empty AccountInfo is returned
            if !reflect.DeepEqual(info, AccountInfo{}) {
                t.Errorf("Basic() info = %v, want empty AccountInfo", info)
            }
        })
    }
}

func TestEmptyDBBlockHash(t *testing.T) {
    db := NewEmptyDB()
    
    testCases := []struct {
        name    string
        number  uint64
        wantErr string
    }{
        {
            name:    "block hash query",
            number:  1,
            wantErr: "EmptyDB: BlockHash not implemented",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            hash, err := db.BlockHash(tc.number)
            
            // Check if error matches expected
            if err == nil || err.Error() != tc.wantErr {
                t.Errorf("BlockHash() error = %v, want %v", err, tc.wantErr)
            }
            
            // Verify empty B256 is returned
            if !reflect.DeepEqual(hash, B256{}) {
                t.Errorf("BlockHash() hash = %v, want empty B256", hash)
            }
        })
    }
}

func TestEmptyDBStorage(t *testing.T) {
    db := NewEmptyDB()
    
    testCases := []struct {
        name    string
        address Address
        index   U256
        wantErr string
    }{
        {
            name:    "storage query",
            address: Address{},
            index:   big.NewInt(0),
            wantErr: "EmptyDB: Storage not implemented",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            value, err := db.Storage(tc.address, tc.index)
            
            // Check if error matches expected
            if err == nil || err.Error() != tc.wantErr {
                t.Errorf("Storage() error = %v, want %v", err, tc.wantErr)
            }
            
            // Verify zero value is returned
            if value.Cmp(big.NewInt(0)) != 0 {
                t.Errorf("Storage() value = %v, want 0", value)
            }
        })
    }
}

func TestEmptyDBCodeByHash(t *testing.T) {
    db := NewEmptyDB()
    
    testCases := []struct {
        name     string
        codeHash B256
        wantErr  string
    }{
        {
            name:     "code query",
            codeHash: B256{},
            wantErr:  "EmptyDB: CodeByHash not implemented",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            code, err := db.CodeByHash(tc.codeHash)
            
            // Check if error matches expected
            if err == nil || err.Error() != tc.wantErr {
                t.Errorf("CodeByHash() error = %v, want %v", err, tc.wantErr)
            }
            
            // Verify empty Bytecode is returned
            if !reflect.DeepEqual(code, Bytecode{}) {
                t.Errorf("CodeByHash() code = %v, want empty Bytecode", code)
            }
        })
    }
}

func TestDatabaseError(t *testing.T) {
    testCases := []struct {
        name    string
        err     *DatabaseError
        wantErr string
    }{
        {
            name: "with underlying error",
            err: &DatabaseError{
                msg: "test message",
                err: errors.New("underlying error"),
            },
            wantErr: "underlying error",
        },
        {
            name: "without underlying error",
            err: &DatabaseError{
                msg: "test message",
                err: nil,
            },
            wantErr: "test message",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            if err := tc.err.Error(); err != tc.wantErr {
                t.Errorf("DatabaseError.Error() = %v, want %v", err, tc.wantErr)
            }
        })
    }
}