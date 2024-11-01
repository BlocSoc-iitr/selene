package execution

import (
	seleneCommon "github.com/BlocSoc-iitr/selene/common"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
)

// Changes have been made as HttpRpc was not satisfying ExecutionRpc
type ExecutionRpc interface {
	New(rpc *string) (ExecutionRpc, error)
	GetProof(address *seleneCommon.Address, slots *[]common.Hash, block uint64) (EIP1186ProofResponse, error)
	CreateAccessList(opts CallOpts, block seleneCommon.BlockTag) (types.AccessList, error)
	GetCode(address *seleneCommon.Address, block uint64) ([]byte, error)
	SendRawTransaction(bytes *[]byte) (common.Hash, error)
	GetTransactionReceipt(tx_hash *common.Hash) (types.Receipt, error)
	GetTransaction(tx_hash *common.Hash) (seleneCommon.Transaction, error)
	GetLogs(filter *ethereum.FilterQuery) ([]types.Log, error)
	GetFilterChanges(filer_id *uint256.Int) ([]types.Log, error)
	UninstallFilter(filter_id *uint256.Int) (bool, error)
	GetNewFilter(filter *ethereum.FilterQuery) (uint256.Int, error)
	GetNewBlockFilter() (uint256.Int, error)
	GetNewPendingTransactionFilter() (uint256.Int, error)
	ChainId() (uint64, error)
	GetFeeHistory(block_count uint64, last_block uint64, reward_percentiles *[]float64) (FeeHistory, error)
}
