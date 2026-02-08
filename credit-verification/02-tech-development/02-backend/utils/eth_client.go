package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// EthClient 以太坊客户端实例
var EthClient *ethclient.Client

// CreditContractABI CreditContract的ABI
var CreditContractABI abi.ABI

// 仅保留 CreditContract 实例
var CreditContractInstance *bind.BoundContract

// InitEthClient 初始化以太坊客户端和合约实例（仅CreditContract）
func InitEthClient() {
	// 1. 连接以太坊节点
	rpcUrl := GlobalConfig.Ethereum.RpcUrl
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		log.Fatalf("连接以太坊节点失败: %v", err)
	}
	EthClient = client
	log.Println("以太坊节点连接成功")

	// 2. 加载 CreditContract ABI
	creditABIFile := "./contract/abi/credit_contract.json"
	creditABIJson, err := ioutil.ReadFile(creditABIFile)
	if err != nil {
		log.Fatalf("读取CreditContract ABI失败: %v", err)
	}
	// 直接解析为ABI（极简方式）
	err = json.Unmarshal(creditABIJson, &CreditContractABI)
	if err != nil {
		log.Fatalf("解析CreditContract ABI失败: %v", err)
	}

	// 3. 实例化 CreditContract
	creditContractAddr := common.HexToAddress(GlobalConfig.Ethereum.CreditContractAddr)
	CreditContractInstance = bind.NewBoundContract(creditContractAddr, CreditContractABI, client, client, client)

	log.Println("合约实例化成功（仅CreditContract）")
}

// GetTransactOpts 保留原有逻辑（无需修改）
func GetTransactOpts() (*bind.TransactOpts, error) {
	privateKeyStr := GlobalConfig.Ethereum.PrivateKey
	if strings.HasPrefix(privateKeyStr, "0x") {
		privateKeyStr = privateKeyStr[2:]
	}
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return nil, fmt.Errorf("解析私钥失败: %v", err)
	}

	fromAddr := crypto.PubkeyToAddress(privateKey.PublicKey)
	nonce, err := EthClient.PendingNonceAt(context.Background(), fromAddr)
	if err != nil {
		return nil, fmt.Errorf("获取Nonce失败: %v", err)
	}

	chainID, err := EthClient.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("获取链ID失败: %v", err)
	}

	transactOpts, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("创建交易选项失败: %v", err)
	}

	transactOpts.Nonce = big.NewInt(int64(nonce))
	transactOpts.From = fromAddr
	transactOpts.GasLimit = uint64(300000)
	transactOpts.GasPrice = big.NewInt(1000000000)

	return transactOpts, nil
}
