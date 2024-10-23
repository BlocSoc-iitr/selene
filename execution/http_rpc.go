package execution

import (
	seleneCommon "github.com/BlocSoc-iitr/selene/common"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/holiman/uint256"
	"strconv"
)

type HttpRpc struct {
	url      string
	provider *rpc.Client
}

func New(rpcUrl string) (*HttpRpc, error) {
	client, err := rpc.Dial(rpcUrl)
	if err != nil {
		return nil, err
	}

	return &HttpRpc{
		url:      rpcUrl,
		provider: client,
	}, nil
}

func (h *HttpRpc) GetProof(address *common.Address, slots *[]common.Hash, block uint64) (EIP1186ProofResponse, error) {
	resultChan := make(chan struct {
		proof EIP1186ProofResponse
		err   error
	})

	go func() {
		var proof EIP1186ProofResponse
		err := h.provider.Call(&proof, "eth_getProof", address, slots, toBlockNumArg(block))
		resultChan <- struct {
			proof EIP1186ProofResponse
			err   error
		}{proof, err}
		close(resultChan)
	}()
	result := <-resultChan
	if result.err != nil {
		return EIP1186ProofResponse{}, result.err
	}
	return result.proof, nil
}

func (h *HttpRpc) CreateAccessList(opts CallOpts, block seleneCommon.BlockTag) (types.AccessList, error) {
	resultChan := make(chan struct {
		accessList types.AccessList
		err        error
	})

	go func() {
		var accessList types.AccessList
		err := h.provider.Call(&accessList, "eth_createAccessList", opts, block)
		resultChan <- struct {
			accessList types.AccessList
			err        error
		}{accessList, err}
		close(resultChan)
	}()

	result := <-resultChan
	if result.err != nil {
		return nil, result.err
	}
	return result.accessList, nil
}

func (h *HttpRpc) GetCode(address *common.Address, block uint64) (hexutil.Bytes, error) {
	resultChan := make(chan struct {
		code hexutil.Bytes
		err  error
	})

	go func() {
		var code hexutil.Bytes
		err := h.provider.Call(&code, "eth_getCode", address, toBlockNumArg(block))
		resultChan <- struct {
			code hexutil.Bytes
			err  error
		}{code, err}
		close(resultChan)
	}()

	result := <-resultChan
	if result.err != nil {
		return nil, result.err
	}
	return result.code, nil
}

func (h *HttpRpc) SendRawTransaction(data *[]byte) (common.Hash, error) {
	resultChan := make(chan struct {
		txHash common.Hash
		err    error
	})

	go func() {
		var txHash common.Hash
		err := h.provider.Call(&txHash, "eth_sendRawTransaction", hexutil.Bytes(*data))
		resultChan <- struct {
			txHash common.Hash
			err    error
		}{txHash, err}
		close(resultChan)
	}()

	result := <-resultChan
	if result.err != nil {
		return common.Hash{}, result.err
	}
	return result.txHash, nil
}

func (h *HttpRpc) GetTransactionReceipt(txHash *common.Hash) (types.Receipt, error) {
	resultChan := make(chan struct {
		receipt types.Receipt
		err     error
	})

	go func() {
		var receipt types.Receipt
		err := h.provider.Call(&receipt, "eth_getTransactionReceipt", txHash)
		resultChan <- struct {
			receipt types.Receipt
			err     error
		}{receipt, err}
		close(resultChan)
	}()

	result := <-resultChan
	if result.err != nil {
		return types.Receipt{}, result.err
	}
	return result.receipt, nil
}

func (h *HttpRpc) GetTransaction(txHash *common.Hash) (seleneCommon.Transaction, error) {
	resultChan := make(chan struct {
		tx  seleneCommon.Transaction
		err error
	})

	go func() {
		var tx seleneCommon.Transaction
		err := h.provider.Call(&tx, "eth_getTransactionByHash", txHash)
		resultChan <- struct {
			tx  seleneCommon.Transaction
			err error
		}{tx, err}
		close(resultChan)
	}()

	result := <-resultChan
	if result.err != nil {
		return seleneCommon.Transaction{}, result.err
	}
	return result.tx, nil
}

func (h *HttpRpc) GetLogs(filter *ethereum.FilterQuery) ([]types.Log, error) {
	resultChan := make(chan struct {
		logs []types.Log
		err  error
	})

	go func() {
		var logs []types.Log
		err := h.provider.Call(&logs, "eth_getLogs", toFilterArg(*filter))
		resultChan <- struct {
			logs []types.Log
			err  error
		}{logs, err}
		close(resultChan)
	}()

	result := <-resultChan
	if result.err != nil {
		return nil, result.err
	}
	return result.logs, nil
}

func (h *HttpRpc) GetFilterChanges(filterID *uint256.Int) ([]types.Log, error) {
	resultChan := make(chan struct {
		logs []types.Log
		err  error
	})

	go func() {
		var logs []types.Log
		err := h.provider.Call(&logs, "eth_getFilterChanges", filterID)
		resultChan <- struct {
			logs []types.Log
			err  error
		}{logs, err}
		close(resultChan)
	}()

	result := <-resultChan
	if result.err != nil {
		return nil, result.err
	}
	return result.logs, nil
}

func (h *HttpRpc) UninstallFilter(filterID *uint256.Int) (bool, error) {
	resultChan := make(chan struct {
		result bool
		err    error
	})

	go func() {
		var result bool
		err := h.provider.Call(&result, "eth_uninstallFilter", filterID)
		resultChan <- struct {
			result bool
			err    error
		}{result, err}
		close(resultChan)
	}()

	result := <-resultChan
	if result.err != nil {
		return false, result.err
	}
	return result.result, nil
}

func (h *HttpRpc) GetNewFilter(filter *ethereum.FilterQuery) (uint256.Int, error) {
	resultChan := make(chan struct {
		filterID uint256.Int
		err      error
	})

	go func() {
		var filterID uint256.Int
		err := h.provider.Call(&filterID, "eth_newFilter", toFilterArg(*filter))
		resultChan <- struct {
			filterID uint256.Int
			err      error
		}{filterID, err}
		close(resultChan)
	}()

	result := <-resultChan
	if result.err != nil {
		return uint256.Int{}, result.err
	}
	return result.filterID, nil
}

func (h *HttpRpc) GetNewBlockFilter() (uint256.Int, error) {
	resultChan := make(chan struct {
		filterID uint256.Int
		err      error
	})

	go func() {
		var filterID uint256.Int
		err := h.provider.Call(&filterID, "eth_newBlockFilter")
		resultChan <- struct {
			filterID uint256.Int
			err      error
		}{filterID, err}
		close(resultChan)
	}()

	result := <-resultChan
	if result.err != nil {
		return uint256.Int{}, result.err
	}
	return result.filterID, nil
}

func (h *HttpRpc) GetNewPendingTransactionFilter() (uint256.Int, error) {
	resultChan := make(chan struct {
		filterID uint256.Int
		err      error
	})

	go func() {
		var filterID uint256.Int
		err := h.provider.Call(&filterID, "eth_newPendingTransactionFilter")
		resultChan <- struct {
			filterID uint256.Int
			err      error
		}{filterID, err}
		close(resultChan)
	}()

	result := <-resultChan
	if result.err != nil {
		return uint256.Int{}, result.err
	}
	return result.filterID, nil
}

func (h *HttpRpc) ChainId() (uint64, error) {
	resultChan := make(chan struct {
		chainID uint64
		err     error
	})

	go func() {
		var chainID uint64
		err := h.provider.Call(&chainID, "eth_chainId")
		resultChan <- struct {
			chainID uint64
			err     error
		}{chainID, err}
		close(resultChan)
	}()

	result := <-resultChan
	if result.err != nil {
		return 0, result.err
	}
	return result.chainID, nil
}

func (h *HttpRpc) GetFeeHistory(blockCount uint64, lastBlock uint64, rewardPercentiles *[]float64) (ethereum.FeeHistory, error) {
	resultChan := make(chan struct {
		feeHistory ethereum.FeeHistory
		err        error
	})

	go func() {
		var feeHistory ethereum.FeeHistory
		err := h.provider.Call(&feeHistory, "eth_feeHistory", blockCount, lastBlock, rewardPercentiles)
		resultChan <- struct {
			feeHistory ethereum.FeeHistory
			err        error
		}{feeHistory, err}
		close(resultChan)
	}()

	result := <-resultChan
	if result.err != nil {
		return ethereum.FeeHistory{}, result.err
	}
	return result.feeHistory, nil
}

func toBlockNumArg(number uint64) string {
	if number == 0 {
		return "latest"
	}
	return "0x" + strconv.FormatUint(number, 16)
}

func toFilterArg(q ethereum.FilterQuery) map[string]interface{} {
	arg := make(map[string]interface{})
	if len(q.Addresses) > 0 {
		arg["address"] = q.Addresses
	}
	if len(q.Topics) > 0 {
		arg["topics"] = q.Topics
	}
	if q.FromBlock != nil {
		arg["fromBlock"] = q.FromBlock.String()
	}
	if q.ToBlock != nil {
		arg["toBlock"] = q.ToBlock.String()
	}
	return arg
}
