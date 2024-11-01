package execution

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	seleneCommon "github.com/BlocSoc-iitr/selene/common"
	"github.com/BlocSoc-iitr/selene/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

func MakeNewRpc(t *testing.T) ExecutionRpc {
	rpcUrl := "https://eth-mainnet.g.alchemy.com/v2/j28GcevSYukh-GvSeBOYcwHOfIggF1Gt"

	var httpRpc ExecutionRpc
	httpRpc, err := (&HttpRpc{}).New(&rpcUrl)

	if err != nil {
		t.Errorf("Error in creating new rpc: %v", err)
	}

	return httpRpc
}

func TestGetProof(t *testing.T) {
	rpc := MakeNewRpc(t)

	addressBytes, err := utils.Hex_str_to_bytes("0xB856af30B938B6f52e5BfF365675F358CD52F91B")
	if err != nil {
		t.Errorf("Error in decoding address string:, %v", err)
	}

	var address seleneCommon.Address = seleneCommon.Address{Addr: [20]byte(addressBytes)}
	slots := []common.Hash{}
	var block uint64 = 14900001

	proof, err := rpc.GetProof(&address, &slots, block)
	// fmt.Printf("Proof: %v", proof)
	if err != nil {
		t.Errorf("Error in fetching proof: %v", err)
	}

	var accountProof EIP1186ProofResponse
	proofString := `{
    "address": "0xb856af30b938b6f52e5bff365675f358cd52f91b",
    "accountProof": [
      "0xf90211a021162657aa1e0af5eef47130ffc3362cb675ebbccfc99ee38ef9196144623507a073dec98f4943e2ab00f5751c23db67b65009bb3cb178d33f5aa93f0c08d583dda0d85b4e33773aaab742db880f8b64ea71f348c6eccb0a56854571bbd3db267f24a0bdcca489de03a49f109c1a2e7d3bd4e644d72de38b7b26dca2f8d3f112110c6fa05c7e8fdff6de07c4cb9ca6bea487a6e5be04af538c25480ce30761901b17e4bfa0d9891f4870e745509cfe17a31568f870b367a36329c892f1b2a37bf59e547183a0af08f747d2ea66efa5bcd03729a95f56297ef9b1e8533ac0d3c7546ebefd2418a0a107595919d4b102afaa0d9b91d9f554f83f0ad61a1e04487e5091543eb81db8a0a0725da6da3b62f88fc573a3fd0dd9dea9cba1750786021da836fd95b7295636a0fd7a768700af3caadaf52a08a23ab0b71ca52830f2b88b1a6b23a52f9ee05507a059434ae837706d7d317e4f7d03cd91f94ed0465fa8b99eaf18ca363bb318c7b3a09e9b831a5f59b781efd5dae8bea30bfd81b9fd5ea231d6b7e82da495c95dd35da0e72d02a01ed9bc928d94cad59ae0695f45120b7fbdbce43a2239a7e5bc81f731a0184bfb9a4051cbaa79917183d004c8d574d7ed5becaf9614c650ed40e8d123d9a0fa4797dc4a35af07f1cd6955318e3ff59578d4df32fd2174ed35f6c4db3471f9a0fec098d1fee8e975b5e78e19003699cf7cd746f47d03692d8e11c5fd58ba92a680",
      "0xf90211a07fc5351578eb6ab7618a31e18c87b2b8b2703c682f2d4c1d01aaa8b53343036ea0e8871ae1828c54b9c9bbf7530890a2fe4e160fb62f72c740c7e79a756e07dbf3a04dd116a7d37146cd0ec730172fa97e84b1f13e687d56118e2d666a02a31a629fa08949d66b81ba98e5ca453ba1faf95c8476873d4c32ff6c9a2558b772c51c5768a028db2de6d80f3a06861d3acc082e3a6bb4a6948980a8e5527bd354a2da037779a09b01ba0fe0193c511161448c602bb9fff88b87ab0ded3255606a15f8bca9d348a0c1c1c6a89f2fdbee0840ff309b5cecd9764b5b5815b385576e75e235d1f04656a04e827215bb9511b3a288e33bb418132940a4d42d589b8db0f796ec917e8f9373a099398993d1d6fdd15d6082be370e4d2cc5d9870923d22770aaec9418f4b675d7a00cd1db5e131341b472af1bdf9a1bf1f1ca82bc5b280c8a50a20bcfff1ab0bdd4a09bbcc86c94be1aabf5c5ceced29f462f59103aa6dafe0fc60172bb2c549a8dbaa0902df0ba9eed7e8a6ebff2d06de8bcec5785bb98cba7606f7f40648408157ef4a0ba9dfd07c453e54504d41b7a44ea42e8220767d1e2a0e6e91ae8d5677ac70e50a0f02f2a5e26d7848f0e5a07de68cbbbd24253d545afb74aac81b35a70b6323f1ca0218b955deca7177f8f58c2da188611b333e5c7ef9212000f64ca92cd5bb6e5a0a049cd750f59e2d6f411d7b611b21b17c8eefe637ca01e1566c53f412308b34c6280",
      "0xf90211a05303302919681c3ad0a56c607c9530ed910f44515f6b40c9633d1900bbbc7e0fa0459fc49e57f39ca6471b1c6905ede7eaa6d7186c8249485cc28338ba18c540cba0825307726d1b7c9d74973d37c12e8b035bf0334836e48ec3f2ff18bf2232dabea0a67ef68daba820c7d6343d1b147b73430ce5c5915a27581cfd12946c2307dc49a003c9b0f0b784de7d72f3b5d5fea87e30dc5fc6f93a0c218240f79a2c74b0f8e2a05a38ddf70df168305b8ba38e8ee54dfadc3f7d81335ec849cb143a10d9738a91a058f0692b5cb07a1c8c800fcf8a70c6e6189a5d72f24ca0040423cf450df1da44a0890dbc62e7429fcca3f1507ad2cd4799c0a7aab25db37ccad434ae23ae889807a075be60d2f635292e69dbc600600605cb8eaf75e96425fd3f2347a5c644db43b9a07b65ba06ee9d2b5dab0a9acc1b8b521cb42f91566de9c174636e817c3d990265a0de65bc6092e28b0cc1ed454fcc07ce26df21bb05efe0a4b4472ff63278e28b95a08077cd7de83428d376ff7588b32df903d2764af7d41deb9c873e9ae889823cd3a0af2f63837dc01e2efb9e40815761015a0d740c2d2549593eefd291a05d40b55aa0c3214baa8d347bd5703926b6fe3ee2f607d0debc0fd73352696dc10f4cbc517da01756cf85b4785bda4a9867a453f8ca1948f015bd329b84083f14d313bddafb80a00dac89194bc1f28d3971b9ca9d1e16a49c6383557187d7bbeb899626d60bfb1980",
      "0xf90211a0bf50b49fae6cfe8b7671e3fa0c163aed76f6457720a2b7c18f567b3c02194c29a070dc71bb7e399e5ae66958261108c84b75e8aacc8d255ce28cd7c9029358872ba0d2ae86d376e65eee52338ad4a1951deb9312f2c161fdf5cfc3e36d5a07ee4239a0f2029dea5033d0e788191ba25fc25bb0570bdbbaf321dfdb076f6695c649a07ca066074b59980560ecdd8ebc96eeb93f50dc1e92983659ea4a6a61a4cff0f474cba01ad85159ddc98609ea628cd17897fe08b0d9a7bb07a2087d92a673e063039aaba04921580f8766f8f156546abd8f0e44af250b34e7323f35c40fdc078223822344a034e07b24a1c17f5dcff27b766099c206fbbb6e549d3f4c02fd8db0241061482aa0c852267182c35e2e5014ab6d656672e9446aaf79c6248d103870d55ee36368b1a00aed203f7e2684942a64f05306e57d64fd44eded94e2ce95e462be93adff640da02cea88d74264c91c546de3822b6169a427559781a774511409864d70a834706ea01d542f8a9b69674e58a5bb89fafb5e79cbe3607732455b09c2a996df48e48837a04c3bbb4f47041018455347567a4e3af472fafe179871f667c3d26038f5dabacfa03c4f12f7cdd35126ce5452aa8322bc8b497eb06c5c41741b590d40645a8fd14ea01334e9a4160b44b622e9523cb587d8ba4795bbf9ad3dd0aa1a2b7f5c6a5cbe94a08feae3d50602063d65763185633aa6e23bb47eca9a39982a4863a7cc6d3586ff80",
      "0xf90211a03fc22103871f30d114942583d24adfab1ce2e651ff6705a05964318fde7c425aa01c86fd2e9d2a823db33bf4089ce1af41332b4e3069b31bb70a67861944d71688a08a90ae88b4479d21135517195f50df20fc29dbe495e09440b0fa5797fc0352c8a02a195c4a89ab6322d8daa221124274d711a9435587406addfb289f9360c0b1caa02f7ed0113a1b72febc7ceb7d9193baeaa093aebea76eea4729821f53d29b302ea07d5cbbeffd22fb0f9d510576cac47a604b2121ef8b08588eaefc46a07844515ca01c0c09d203e342fab9f80835f3aa7bb7e94cbf94d3a18b21ce905e75b690673fa07c310f931f12d1651dbb9bdffbe5e0a16db981dccfb3f4a838592e2347b1b187a091b24e00d37034ed70e0c6653f8616363efbea43be86d52aabf6a9c5d5049d3aa017cc8ab2e63508691dee64ddb2dbe5352fe531f55728368cab3d8634450730dca0d8ee92d688eabcdf28af1830d7217fcd1431c0b0e3311039c422c5f45f9d525da0abe4323ef90fb783ffd6bf29d240818b0a2477d2ad4577fce57642a5ba476957a0eee3fbc510b1a6d8da176b9eab2035837769988b216fcef67f6a215e5e261d5ea082bc27591d8b0408739712c2f0e623e3d296d12afd6b7356dae237a315a6ce3ba0634affb8f9744fed774851772cb5ce495c50212961a64262d915632c2bede721a0345017d846f3be29dca1f56a734886c26cb49d3bcbbae5f7356e55e71d84147d80",
      "0xf90211a0a152d08043b3865248d8ae9c4594d6e09079f61ee4ad9afca98d3befb3a9307ca01313ea4df7ec991f5ba1b2175408765ea67111857ae25cafe21986810d633353a05cb00f30a4b6749cb8d01dda2ad665a2c570b7c6959b364c46da75fc2aebfa14a0a493d42fa40d6c8fb23090f20cebe9f1cde8a47d9594e666096c3f76219cfb34a06d2149e05ba1a31bcc352fd3a79fd42d95383f14a312862fa5f4b7b7bdb63254a0f40094679bf0599fdeecae8c800a423ad2499b67f9546d4085d5b7a351561072a05036fc625ed5a13d143fd0984e99969d2d48f962baec80b3f0e78323c8e864ffa0c0887db54ab0d4309ed5f563448bfc78a4a88c3a5473fcbd9bd263c8fcca4b9fa04287b193a315cba13a49482b4d83c068cf08b622593111e416b7a3b815e3595aa051b82224b54dc4050703157cdd6c74c618872c798badce7192fe1fd534814d5da017ac02273956b25dc2429750156e0a45bf461437bec84b784d19a5b964dd6882a05be9c25f80a6f34e9eb526d6c3c89e3ec2c5dd769d2d915d835208cac3f56d36a0ea7ae8e74baeae9d6307371f85cce46ba7ad46b5c2b3616a3573a26e1260bc31a04446678b6bf75075ae0261c179d88e49fae1d9482c214ec8c693239b583a6b18a09ec91d47f671c343cea224def0afcd62a57a408ac0e36b79cf29a4495ff9055ca0c3da71d14030daf8fb9157ec84f077df97dff395c1fcf8f04361e02aa1af36ff80",
      "0xf8b1808080a038f5e1b2680d95eaf7f225b996fc482f60cabcebeae26f883a4b58e1e9c7bbeda0c220c5d76d85ad38d8f82d0d6d6f48db3c23ae3657d1ac3ca6e2d98b4e48bfde8080a06fe32c7c1f8f80ebe5128f11a7af3a5bc47dbc6ac0af705069532b1cbefd6792a0015eb7d24835c910fc5f906627968a1a9e810cb164afd634ca683b6bb34f0241808080808080a03ca1d0ead152d38c16bda91dac49e3f0c9a0afeaa67598c1fce2506f5f03162680",
      "0xf86d9d3c3738deb88e49108e7a5bd83c14ad65b5ba598e2932551dc9b9ad1879b84df84b10874ef05b2fe9d8c8a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421a0c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"
    ],
    "balance": "0x4ef05b2fe9d8c8",
    "codeHash": "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
    "nonce": "0x10",
    "storageHash": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
    "storageProof": []
  }`
	err = json.Unmarshal([]byte(proofString), &accountProof)
	assert.NoError(t, err, "Found unexpected error")
	assert.Equal(t, accountProof, proof, "Proof didn't match expected value")
}

// func TestCreateAccessList(t *testing.T) {
// 	rpc := MakeNewRpc(t)

// 	fromAddress := common.HexToAddress("0x8885d8FFefb29AC94FCe584014266A6fE8437356")
// 	toAddress := common.HexToAddress("0xb5C3dd56b4Dc6108b21b0Dac4A20898A464b8c41")
// 	inputData := common.Hex2Bytes("0x")
// 	// data, err := utils.Hex_str_to_bytes("0x608060806080608155")

// 	opts := CallOpts{
// 		From: &fromAddress,
// 		To: &toAddress,
// 		Data: inputData,
// 	}
// 	block := seleneCommon.BlockTag{
// 		Finalized: true,
// 		Number: 20985649,
// 	}

// 	accessList, err := rpc.CreateAccessList(opts, block)
// 	fmt.Println("Access list:", accessList)
// 	assert.NoError(t, err, "Found Error")

// }

func TestGetCode(t *testing.T) {
	rpc := MakeNewRpc(t)

	addressBytes, err := utils.Hex_str_to_bytes("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")
	if err != nil {
		t.Errorf("Error in decoding address string:, %v", err)
	}

	var address seleneCommon.Address = seleneCommon.Address{Addr: [20]byte(addressBytes)}
	blockNumber := 20983632

	code, _ := rpc.GetCode(&address, uint64(blockNumber))

	codeString := "0x6060604052600436106100af576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806306fdde03146100b9578063095ea7b31461014757806318160ddd146101a157806323b872dd146101ca5780632e1a7d4d14610243578063313ce5671461026657806370a082311461029557806395d89b41146102e2578063a9059cbb14610370578063d0e30db0146103ca578063dd62ed3e146103d4575b6100b7610440565b005b34156100c457600080fd5b6100cc6104dd565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561010c5780820151818401526020810190506100f1565b50505050905090810190601f1680156101395780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b341561015257600080fd5b610187600480803573ffffffffffffffffffffffffffffffffffffffff1690602001909190803590602001909190505061057b565b604051808215151515815260200191505060405180910390f35b34156101ac57600080fd5b6101b461066d565b6040518082815260200191505060405180910390f35b34156101d557600080fd5b610229600480803573ffffffffffffffffffffffffffffffffffffffff1690602001909190803573ffffffffffffffffffffffffffffffffffffffff1690602001909190803590602001909190505061068c565b604051808215151515815260200191505060405180910390f35b341561024e57600080fd5b61026460048080359060200190919050506109d9565b005b341561027157600080fd5b610279610b05565b604051808260ff1660ff16815260200191505060405180910390f35b34156102a057600080fd5b6102cc600480803573ffffffffffffffffffffffffffffffffffffffff16906020019091905050610b18565b6040518082815260200191505060405180910390f35b34156102ed57600080fd5b6102f5610b30565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561033557808201518184015260208101905061031a565b50505050905090810190601f1680156103625780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b341561037b57600080fd5b6103b0600480803573ffffffffffffffffffffffffffffffffffffffff16906020019091908035906020019091905050610bce565b604051808215151515815260200191505060405180910390f35b6103d2610440565b005b34156103df57600080fd5b61042a600480803573ffffffffffffffffffffffffffffffffffffffff1690602001909190803573ffffffffffffffffffffffffffffffffffffffff16906020019091905050610be3565b6040518082815260200191505060405180910390f35b34600360003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825401925050819055503373ffffffffffffffffffffffffffffffffffffffff167fe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c346040518082815260200191505060405180910390a2565b60008054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156105735780601f1061054857610100808354040283529160200191610573565b820191906000526020600020905b81548152906001019060200180831161055657829003601f168201915b505050505081565b600081600460003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925846040518082815260200191505060405180910390a36001905092915050565b60003073ffffffffffffffffffffffffffffffffffffffff1631905090565b600081600360008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054101515156106dc57600080fd5b3373ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff16141580156107b457507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600460008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205414155b156108cf5781600460008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020541015151561084457600080fd5b81600460008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825403925050819055505b81600360008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254039250508190555081600360008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825401925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040518082815260200191505060405180910390a3600190509392505050565b80600360003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205410151515610a2757600080fd5b80600360003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825403925050819055503373ffffffffffffffffffffffffffffffffffffffff166108fc829081150290604051600060405180830381858888f193505050501515610ab457600080fd5b3373ffffffffffffffffffffffffffffffffffffffff167f7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65826040518082815260200191505060405180910390a250565b600260009054906101000a900460ff1681565b60036020528060005260406000206000915090505481565b60018054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610bc65780601f10610b9b57610100808354040283529160200191610bc6565b820191906000526020600020905b815481529060010190602001808311610ba957829003601f168201915b505050505081565b6000610bdb33848461068c565b905092915050565b60046020528160005260406000206020528060005260406000206000915091505054815600a165627a7a72305820deb4c2ccab3c2fdca32ab3f46728389c2fe2c165d5fafa07661e4e004f6c344a0029"
	codeBytes, err := utils.Hex_str_to_bytes(codeString)

	// fmt.Println("Code:", code)
	assert.NoError(t, err, "Found Error")
	assert.Equal(t, codeBytes, code, "Code didn't match expected value")
}

// ** Not possible to test as to send actual transaction, ETH is needed on mainnet
// func TestSendRawTransaction(t *testing.T) {
// 	rpc := MakeNewRpc(t)
// // Raw Transaction Data
// 	data, _ := utils.Hex_str_to_bytes("")
// 	// transaction := seleneCommon.Transaction{
// 	// 	From: "",
// 	// 	Nonce: ,
// 	// }
// 	// data, err := json.Marshal(transaction)
// 	hash, err := rpc.SendRawTransaction(&data)

// 	fmt.Println("Hash: ", hash)
// 	assert.NoError(t, err, "Found Error")
// }

func TestGetTransactionReceipt(t *testing.T) {
	rpc := MakeNewRpc(t)
	txHash, _ := utils.Hex_str_to_bytes("0x4bc11033063e445e038e52e72266f5054845d3879704d0cf38bedeb86c924cec")

	_, err := rpc.GetTransactionReceipt((*common.Hash)(txHash))
	// expectedReceipt := types.Receipt{}

	assert.NoError(t, err, "Found Error")
	// assert.Equal(t, expectedReceipt ,txnReceipt, "Receipt didn't match")
}

func TestGetTransaction(t *testing.T) {
	rpc := MakeNewRpc(t)
	txHash, _ := utils.Hex_str_to_bytes("0x4bc11033063e445e038e52e72266f5054845d3879704d0cf38bedeb86c924cec")

	_, err := rpc.GetTransaction((*common.Hash)(txHash))
	// expectedTxn := seleneCommon.Transaction{}

	assert.NoError(t, err, "Found Error")
	// assert.Equal(t, expectedTxn ,txn, "Receipt didn't match")
}

func TestGetLogs(t *testing.T) {
	rpc := MakeNewRpc(t)
	address := common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")
	blockHash := common.HexToHash("0x945b6b3568c39736fe1595234f0529efc8d64930a5343196aa8ef799b02dd609")

	filterQuery := ethereum.FilterQuery{
		Addresses: []common.Address{address},
		BlockHash: &blockHash,
	}

	_, err := rpc.GetLogs(&filterQuery)
	// expectedLogs := []types.Log{}
	fmt.Println("")
	assert.NoError(t, err, "Found Error")
	// assert.Equal(t, expectedLogs, logs, "Logs didn't match")
}

func TestGetChainId(t *testing.T) {
	rpc := MakeNewRpc(t)

	chainId, err := rpc.ChainId()

	assert.NoError(t, err, "Found Error")
	assert.Equal(t, chainId, uint64(1), "Expected chain Id to be 1")
}

func TestGetNewFilter(t *testing.T) {
	rpc := MakeNewRpc(t)
	address := common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")
	blockHash := common.HexToHash("0x945b6b3568c39736fe1595234f0529efc8d64930a5343196aa8ef799b02dd609")

	filterQuery := ethereum.FilterQuery{
		Addresses: []common.Address{address},
		BlockHash: &blockHash,
		FromBlock: big.NewInt(200000000),
	}

	_, err := rpc.GetNewFilter(&filterQuery)

	assert.NoError(t, err, "Found Error")
	// assert.Equal(t, *uint256.MustFromHex("0x35ef5ae2a28c50b0a6a8fd0903c99e7f"), filterId, "Filter Id doesn't match")
}

func TestGetFilterChanges(t *testing.T) {
	rpc := MakeNewRpc(t)

	address := common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")
	blockHash := common.HexToHash("0x945b6b3568c39736fe1595234f0529efc8d64930a5343196aa8ef799b02dd609")

	filterQuery := ethereum.FilterQuery{
		Addresses: []common.Address{address},
		BlockHash: &blockHash,
		FromBlock: big.NewInt(200000000),
	}

	filterId, _ := rpc.GetNewFilter(&filterQuery)

	logs, err := rpc.GetFilterChanges(&filterId)
	expectedLogs := []types.Log{}
	assert.NoError(t, err, "Found Error")
	assert.Equal(t, expectedLogs, logs, "Logs didn't match")
}

func TestUninstallFilter(t *testing.T) {
	rpc := MakeNewRpc(t)

	address := common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")
	blockHash := common.HexToHash("0x945b6b3568c39736fe1595234f0529efc8d64930a5343196aa8ef799b02dd609")

	filterQuery := ethereum.FilterQuery{
		Addresses: []common.Address{address},
		BlockHash: &blockHash,
		FromBlock: big.NewInt(200000000),
	}

	filterId, _ := rpc.GetNewFilter(&filterQuery)

	success, err := rpc.UninstallFilter(&filterId)

	assert.NoError(t, err, "Found Error")
	assert.Equal(t, true, success, "Success didn't match")

	//Check whether the filter still exists
	_, err = rpc.GetFilterChanges(&filterId)
	assert.Error(t, err, "Found Error")
}

func TestGetNewBlockFilter(t *testing.T) {
	rpc := MakeNewRpc(t)

	_, err := rpc.GetNewBlockFilter()

	assert.NoError(t, err, "Found Error")
	// assert.Equal(t, 0, filterId, "Filter ID not equal")
	// _, err = rpc.GetFilterChanges(&filterId)
	// assert.NoError(t, err, "Found Error")
}

func TestGetNewPendingTransactionFilter(t *testing.T) {
	rpc := MakeNewRpc(t)

	_, err := rpc.GetNewPendingTransactionFilter()

	assert.NoError(t, err, "Found Error")
}

func TestGetFeeHistory(t *testing.T) {
	rpc := MakeNewRpc(t)
	blockCount := uint64(5)
	lastBlock := uint64(20000000)
	rewardPercentiles := []float64{10, 40, 90}
	_, err := rpc.GetFeeHistory(blockCount, lastBlock, &rewardPercentiles)

	assert.NoError(t, err, "Found Error")
	// assert.Equal(t, FeeHistory{}, feeHistory, "Fee History does not match")
}
