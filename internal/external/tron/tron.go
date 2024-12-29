package tron

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/OwodDEV/crypto-service/internal/config"
	"github.com/OwodDEV/crypto-service/internal/models"
	"github.com/OwodDEV/crypto-service/pkg/utils"

	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

const (
	usdtAddress        = "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"
	usdtBalanceMethod  = "70a08231"
	usdtTransferMethod = "a9059cbb"
	usdtDecimals       = 6
)

type Tron struct {
	Config *config.Config
	client *client.GrpcClient
}

func NewTronService(cfg *config.Config) (s *Tron, err error) {
	s = &Tron{
		Config: cfg,
	}
	return
}

func (s *Tron) Connect() (err error) {
	slog.Info("initializing Tron external service connection...")
	s.client = client.NewGrpcClient(s.Config.External.Tron.RPCEndpoint)
	err = s.client.Start(grpc.WithInsecure())
	if err != nil {
		slog.Error("failed to connect to Tron", slog.Any("error", err))
		return
	}

	return nil
}

func (s *Tron) Shutdown() {
	slog.Info("shutting down Tron external service...")
	s.client.Stop()
}

func (s *Tron) GetBalance(ctx context.Context, addr, token string) (balance string, err error) {
	logger := slog.With(
		slog.String("request_id", ctx.Value("request_id").(string)),
		slog.String("func", "external.Tron.GetBalance()"),
		slog.String("address", addr),
		slog.String("token", token),
	)

	var tokenAddress string
	var tokenBalanceMethod string
	var tokenDecimals int
	switch token {
	case "USDT":
		tokenAddress = usdtAddress
		tokenBalanceMethod = usdtBalanceMethod
		tokenDecimals = usdtDecimals
	default:
		err = errors.New("unknown token")
		logger.Warn(err.Error())
		return
	}

	addrBytes21, err := address.Base58ToAddress(addr)
	if err != nil {
		logger.Error("failed to convert address to 21 bytes format", slog.Any("error", err))
		return
	}
	addrBytes20 := addrBytes21.Bytes()[1:] // remove first byte of version
	addrPadded := make([]byte, 32)
	copy(addrPadded[12:], addrBytes20)

	// invoke
	data := tokenBalanceMethod + hex.EncodeToString(addrPadded)
	callResult, err := s.client.TRC20Call(addr, tokenAddress, data, true, 0)
	if err != nil {
		logger.Error("failed to invoke contract with balanceOf method", slog.Any("error", err))
		return
	}

	// parse result
	if len(callResult.ConstantResult) == 0 {
		err = errors.New("balanceOf method has no result")
		logger.Warn(err.Error())
		return
	}

	balanceBytes := callResult.ConstantResult[0]
	balanceRaw := new(big.Int).SetBytes(balanceBytes)
	balance = utils.FormatCurrency(balanceRaw, tokenDecimals)
	return
}

func (s *Tron) GetTransaction(ctx context.Context, hash, token string) (result models.Transaction, err error) {
	logger := slog.With(
		slog.String("request_id", ctx.Value("request_id").(string)),
		slog.String("func", "external.Tron.GetTransaction()"),
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

	trx, err := s.client.GetTransactionByID(hash)
	if err != nil {
		logger.Error("failed to get transaction by hash", slog.Any("error", err))
		return
	}

	trxContract := trx.GetRawData().GetContract()
	if len(trxContract) == 0 {
		err = errors.New("failed to get contract")
		logger.Error(err.Error())
		return
	}

	parameter := trxContract[0].GetParameter()
	if parameter == nil {
		err = errors.New("failed to get parameter of contract")
		logger.Error(err.Error())
		return
	}

	scData := core.TriggerSmartContract{}
	err = proto.Unmarshal(parameter.GetValue(), &scData)
	if err != nil {
		logger.Error("failed to unmarshal smartcontract data", slog.Any("error", err))
		return
	}

	trxContractAddress := common.EncodeCheck(scData.ContractAddress)
	if trxContractAddress != tokenAddress {
		err = errors.New("the transaction does not involve in requested token transfers")
		logger.Warn(err.Error())
		return
	}

	trxInput := scData.GetData()
	if len(trxInput) != 4+32+32 { // 4 bytes for signature, 2 params
		err = errors.New("unknown input data")
		logger.Warn(err.Error())
		return
	}

	trxMethodSignature := hex.EncodeToString(trxInput[:4])
	fmt.Println("Method Signature:", trxMethodSignature)
	if trxMethodSignature != tokenTransferMethod {
		err = errors.New("not a transfer method")
		logger.Warn(err.Error())
		return
	}

	trxFrom := common.EncodeCheck(scData.OwnerAddress)

	trxParams := trxInput[4:]
	trxToBytes := append([]byte{0x41}, trxParams[12:32]...)
	trxTo := common.EncodeCheck(trxToBytes)

	trxAmountBytes := trxParams[32:]
	trxAmountRaw := new(big.Int).SetBytes(trxAmountBytes)
	trxAmount := utils.FormatCurrency(trxAmountRaw, tokenDecimals)

	result = models.Transaction{
		Hash:   hash,
		From:   trxFrom,
		To:     trxTo,
		Amount: trxAmount,
	}
	return
}
