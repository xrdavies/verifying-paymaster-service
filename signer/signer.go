package signer

import (
	"crypto/ecdsa"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ququzone/verifying-paymaster-service/config"
	"github.com/ququzone/verifying-paymaster-service/container"
	"github.com/ququzone/verifying-paymaster-service/contracts"
	"github.com/ququzone/verifying-paymaster-service/logger"
	"github.com/ququzone/verifying-paymaster-service/types"
	"github.com/ququzone/verifying-paymaster-service/utils"
)

var (
	// one day
	validTimeDelay = new(big.Int).SetInt64(86400)
	uint48Ty, _    = abi.NewType("uint256", "uint48", []abi.ArgumentMarshaling{})
	timeRangeABI   = abi.Arguments{
		{Name: "validUntil", Type: uint48Ty},
		{Name: "validAfter", Type: uint48Ty},
	}
	emptySignature = make([]byte, 65)
)

type Signer struct {
	Container  container.Container
	Contract   common.Address
	Paymaster  *contracts.VerifyingPaymaster
	PrivateKey *ecdsa.PrivateKey
}

func NewSigner(con container.Container) (*Signer, error) {
	conf := config.Config()
	keyData, err := os.ReadFile(conf.Keystore)
	if err != nil {
		return nil, err
	}
	keystore, err := keystore.DecryptKey(keyData, conf.Passphrase)
	if err != nil {
		return nil, err
	}
	logger.S().Infof("VerifyingPaymaster contract: %s", conf.Contract)
	logger.S().Infof("VerifyingPaymaster signer: %s", keystore.Address.String())

	rpc, err := ethclient.Dial(conf.RPC)
	if err != nil {
		return nil, err
	}

	contract := common.HexToAddress(conf.Contract)
	paymaster, err := contracts.NewVerifyingPaymaster(contract, rpc)
	if err != nil {
		return nil, err
	}

	return &Signer{
		Container:  con,
		Contract:   contract,
		Paymaster:  paymaster,
		PrivateKey: keystore.PrivateKey,
	}, nil
}

func (s *Signer) Eth_signVerifyingPaymaster(op map[string]any) (string, error) {
	userOp, err := types.NewUserOperation(op)
	if err != nil {
		return "", err
	}
	// TODO: verify op rules:
	//  1. normal gas
	//  2. only for create

	validAfter := new(big.Int).SetInt64(time.Now().Unix())
	validUntil := new(big.Int).Add(validAfter, validTimeDelay)
	timeRangeData, err := timeRangeABI.Pack(validUntil, validAfter)
	if err != nil {
		return "", err
	}
	userOp.PaymasterAndData = append(append(s.Contract.Bytes(), timeRangeData...), emptySignature...)
	userOp.Signature = []byte{}

	hash, err := s.Paymaster.GetHash(nil, contracts.UserOperation{
		Sender:               userOp.Sender,
		Nonce:                userOp.Nonce,
		InitCode:             userOp.InitCode,
		CallData:             userOp.CallData,
		CallGasLimit:         userOp.CallGasLimit,
		VerificationGasLimit: userOp.VerificationGasLimit,
		PreVerificationGas:   userOp.PreVerificationGas,
		MaxFeePerGas:         userOp.MaxFeePerGas,
		MaxPriorityFeePerGas: userOp.MaxPriorityFeePerGas,
		PaymasterAndData:     userOp.PaymasterAndData,
		Signature:            userOp.Signature,
	}, validUntil, validAfter)
	if err != nil {
		return "", err
	}
	signature, err := utils.SignMessage(s.PrivateKey, hash[:])
	if err != nil {
		return "", err
	}

	return hexutil.Encode(append(append(s.Contract.Bytes(), timeRangeData...), signature...)), nil
}
