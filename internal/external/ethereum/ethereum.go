package ethereum

import (
	"context"
	"encoding/hex"
	"errors"
	"log/slog"
	"math/big"
	"strings"

	"github.com/OwodDEV/crypto-service/internal/config"
	"github.com/OwodDEV/crypto-service/internal/models"
	"github.com/OwodDEV/crypto-service/pkg/utils"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	usdtAddress        = "0xdac17f958d2ee523a2206206994597c13d831ec7"
	usdtTransferMethod = "a9059cbb"
	usdtDecimals       = 6
)

type Ethereum struct {
	Config    *config.Config
	client    *ethclient.Client
	parsedABI abi.ABI
}

func NewEthereumService(cfg *config.Config) (s *Ethereum, err error) {
	logger := slog.With(
		slog.String("func", "external.ethereum.NewEthereumService()"),
	)

	s = &Ethereum{
		Config: cfg,
	}

	s.parsedABI, err = abi.JSON(strings.NewReader(`
		[
		  {
			"constant": true,
			"inputs": [
			  {"name": "account", "type": "address"}
			],
			"name": "balanceOf",
			"outputs": [
			  {"name": "", "type": "uint256"}
			],
			"payable": false,
			"stateMutability": "view",
			"type": "function"
		  },
		  {
			"inputs": [
			  {"name": "recipient", "type": "address"},
			  {"name": "amount", "type": "uint256"}
			],
			"name": "transfer",
			"outputs": [],
			"payable": false,
			"stateMutability": "nonpayable",
			"type": "function"
		  }
		]
	`))
	if err != nil {
		logger.Error("failed to parse ABI", slog.Any("error", err))
		return
	}
	return
}

func (s *Ethereum) Connect() (err error) {
	slog.Info("initializing Ethereum external service connection...")
	s.client, err = ethclient.Dial(s.Config.External.Ethereum.RPCEndpoint)
	if err != nil {
		slog.Error("failed to connect to Ethereum", slog.Any("error", err))
		return
	}
	return nil
}

func (s *Ethereum) Shutdown() {
	slog.Info("shutting down Ethereum external service...")
	s.client.Close()
}

func (s *Ethereum) GetBalance(ctx context.Context, address, token string) (balance string, err error) {
	logger := slog.With(
		slog.String("request_id", ctx.Value("request_id").(string)),
		slog.String("func", "external.Ethereum.GetBalance()"),
		slog.String("address", address),
		slog.String("token", token),
	)

	var tokenAddress string
	var tokenDecimals int
	switch token {
	case "USDT":
		tokenAddress = usdtAddress
		tokenDecimals = usdtDecimals
	default:
		err = errors.New("unknown token")
		logger.Warn(err.Error())
		return
	}

	data, err := s.parsedABI.Pack("balanceOf", common.HexToAddress(address))
	if err != nil {
		logger.Error("failed to pack data for balanceOf method", slog.Any("error", err))
		return
	}

	// invoke
	tokenAddressCommon := common.HexToAddress(tokenAddress)
	msg := ethereum.CallMsg{
		To:   &tokenAddressCommon,
		Data: data,
	}
	callResult, err := s.client.CallContract(ctx, msg, nil)
	if err != nil {
		logger.Error("failed to invoke contract with balanceOf method", slog.Any("error", err))
		return
	}

	// parse result
	rawBalance := new(big.Int)
	err = s.parsedABI.UnpackIntoInterface(&rawBalance, "balanceOf", callResult)
	if err != nil {
		logger.Error("failed to unpack result of balanceOf method", slog.Any("error", err))
		return
	}
	balance = utils.FormatCurrency(rawBalance, tokenDecimals)
	return
}

func (s *Ethereum) GetTransaction(ctx context.Context, hash, token string) (result models.Transaction, err error) {
	logger := slog.With(
		slog.String("request_id", ctx.Value("request_id").(string)),
		slog.String("func", "external.Ethereum.GetTransaction()"),
		slog.String("hash", hash),
		slog.String("token", token),
	)

	var tokenAddress string
	var tokenTransferMethod string
	var tokenDecimals int
	switch token {
	case "USDT":
		tokenAddress = usdtAddress
		tokenTransferMethod = usdtTransferMethod
		tokenDecimals = usdtDecimals
	default:
		err = errors.New("unknown token")
		logger.Warn(err.Error())
		return
	}

	// invoke
	trx, _, err := s.client.TransactionByHash(ctx, common.HexToHash(hash))
	if err != nil {
		logger.Error("failed to get transaction by hash", slog.Any("error", err))
		return
	}

	// parse result
	trxContractAddress := strings.ToLower(trx.To().String())
	if trxContractAddress != tokenAddress {
		err = errors.New("the transaction does not involve in requested token transfers")
		logger.Warn(err.Error())
		return
	}

	trxInput := trx.Data()
	trxMethodSignature := hex.EncodeToString(trxInput[:4]) // Первые 4 байта — метод.
	if trxMethodSignature != tokenTransferMethod {
		err = errors.New("not a transfer method")
		logger.Warn(err.Error())
		return
	}

	dataMap := make(map[string]interface{})
	err = s.parsedABI.Methods["transfer"].Inputs.UnpackIntoMap(dataMap, trxInput[4:])
	if err != nil {
		logger.Error("failed to unpack transfer data of transaction", slog.Any("error", err))
		return
	}

	trxFromCommon, err := types.Sender(types.NewLondonSigner(trx.ChainId()), trx)
	if err != nil {
		logger.Error("not able to retrieve sender", slog.Any("error", err))
		return
	}
	trxFrom := trxFromCommon.Hex()

	trxToCommon := dataMap["recipient"].(common.Address)
	trxTo := strings.ToLower(trxToCommon.Hex())

	trxAmountRaw := dataMap["amount"].(*big.Int)
	trxAmount := utils.FormatCurrency(trxAmountRaw, tokenDecimals)

	result = models.Transaction{
		Hash:   hash,
		From:   trxFrom,
		To:     trxTo,
		Amount: trxAmount,
	}
	return
}
