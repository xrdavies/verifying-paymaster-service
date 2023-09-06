package api

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/ququzone/verifying-paymaster-service/contracts"
	"github.com/ququzone/verifying-paymaster-service/types"
	"github.com/ququzone/verifying-paymaster-service/utils"
)

type RPCError struct {
	code    int
	message string
	data    any
}

// New returns a new custom RPCError.
func NewRPCError(code int, message string, data any) error {
	return &RPCError{code, message, data}
}

// Code returns the message field of the JSON-RPC error object.
func (e *RPCError) Error() string {
	return e.message
}

// Data returns the data field of the JSON-RPC error object.
func (e *RPCError) Data() any {
	return e.data
}

type ExecutionResultRevert struct {
	PreOpGas      *big.Int
	Paid          *big.Int
	ValidAfter    *big.Int
	ValidUntil    *big.Int
	TargetSuccess bool
	TargetResult  []byte
}

func executionResult() abi.Error {
	uint256, _ := abi.NewType("uint256", "", nil)
	uint48, _ := abi.NewType("uint48", "", nil)
	boolean, _ := abi.NewType("bool", "", nil)
	bytes, _ := abi.NewType("bytes", "", nil)
	return abi.NewError("ExecutionResult", abi.Arguments{
		{Name: "preOpGas", Type: uint256},
		{Name: "paid", Type: uint256},
		{Name: "validAfter", Type: uint48},
		{Name: "validUntil", Type: uint48},
		{Name: "targetSuccess", Type: boolean},
		{Name: "targetResult", Type: bytes},
	})
}

func NewExecutionResult(err error) (*ExecutionResultRevert, error) {
	rpcErr, ok := err.(rpc.DataError)
	if !ok {
		return nil, errors.New("executionResult: cannot assert type: error is not of type rpc.DataError")
	}

	data, ok := rpcErr.ErrorData().(string)
	if !ok {
		return nil, errors.New("executionResult: cannot assert type: data is not of type string")
	}

	sim := executionResult()
	revert, err := sim.Unpack(common.Hex2Bytes(data[2:]))
	if err != nil {
		return nil, fmt.Errorf("executionResult: %s", err)
	}

	args, ok := revert.([]any)
	if !ok {
		return nil, errors.New("executionResult: cannot assert type: args is not of type []any")
	}
	if len(args) != 6 {
		return nil, fmt.Errorf("executionResult: invalid args length: expected 6, got %d", len(args))
	}

	return &ExecutionResultRevert{
		PreOpGas:      args[0].(*big.Int),
		Paid:          args[1].(*big.Int),
		ValidAfter:    args[2].(*big.Int),
		ValidUntil:    args[3].(*big.Int),
		TargetSuccess: args[4].(bool),
		TargetResult:  args[5].([]byte),
	}, nil
}

type FailedOpRevert struct {
	OpIndex int
	Reason  string
}

func failedOp() abi.Error {
	opIndex, _ := abi.NewType("uint256", "uint256", nil)
	reason, _ := abi.NewType("string", "string", nil)
	return abi.NewError("FailedOp", abi.Arguments{
		{Name: "opIndex", Type: opIndex},
		{Name: "reason", Type: reason},
	})
}

func NewFailedOp(err error) (*FailedOpRevert, error) {
	rpcErr, ok := err.(rpc.DataError)
	if !ok {
		return nil, fmt.Errorf(
			"failedOp: cannot assert type: error is not of type rpc.DataError, err: %s",
			err,
		)
	}

	data, ok := rpcErr.ErrorData().(string)
	if !ok {
		return nil, fmt.Errorf(
			"failedOp: cannot assert type: data is not of type string, err: %s, data: %s",
			rpcErr.Error(),
			rpcErr.ErrorData(),
		)
	}

	failedOp := failedOp()
	revert, err := failedOp.Unpack(common.Hex2Bytes(data[2:]))
	if err != nil {
		return nil, fmt.Errorf("failedOp: %s", err)
	}

	args, ok := revert.([]any)
	if !ok {
		return nil, errors.New("failedOp: cannot assert type: args is not of type []any")
	}
	if len(args) != 2 {
		return nil, fmt.Errorf("failedOp: invalid args length: expected 2, got %d", len(args))
	}

	opIndex, ok := args[0].(*big.Int)
	if !ok {
		return nil, errors.New("failedOp: cannot assert type: opIndex is not of type *big.Int")
	}

	reason, ok := args[1].(string)
	if !ok {
		return nil, errors.New("failedOp: cannot assert type: reason is not of type string")
	}

	return &FailedOpRevert{
		OpIndex: int(opIndex.Int64()),
		Reason:  reason,
	}, nil
}

func estimate(
	client *ethclient.Client,
	key *ecdsa.PrivateKey,
	paymasterAddr common.Address,
	paymaster *contracts.VerifyingPaymaster,
	entryPoint common.Address,
	op *types.UserOperation,
) (preVerificationGas *big.Int, verificationGas *big.Int, callGas *big.Int, err error) {
	defaultGas := big.NewInt(1000000)
	op.CallGasLimit = defaultGas
	op.VerificationGasLimit = defaultGas
	validAfter := new(big.Int).SetInt64(time.Now().Unix())
	validUntil := new(big.Int).Add(validAfter, validTimeDelay)
	timeRangeData, err := timeRangeABI.Pack(validUntil, validAfter)
	if err != nil {
		return nil, nil, nil, err
	}
	op.PaymasterAndData = append(append(paymasterAddr.Bytes(), timeRangeData...), emptySignature...)
	preSignature := op.Signature
	op.Signature = []byte{}

	hash, err := paymaster.GetHash(nil, contracts.UserOperation{
		Sender:               op.Sender,
		Nonce:                op.Nonce,
		InitCode:             op.InitCode,
		CallData:             op.CallData,
		CallGasLimit:         op.CallGasLimit,
		VerificationGasLimit: op.VerificationGasLimit,
		PreVerificationGas:   op.PreVerificationGas,
		MaxFeePerGas:         op.MaxFeePerGas,
		MaxPriorityFeePerGas: op.MaxPriorityFeePerGas,
		PaymasterAndData:     op.PaymasterAndData,
		Signature:            op.Signature,
	}, validUntil, validAfter)
	if err != nil {
		return nil, nil, nil, err
	}
	signature, err := utils.SignMessage(key, hash[:])
	if err != nil {
		return nil, nil, nil, err
	}

	op.PaymasterAndData = append(append(paymasterAddr.Bytes(), timeRangeData...), signature...)
	op.Signature = preSignature
	parsedABI, err := abi.JSON(strings.NewReader(contracts.EntryPointABI))
	if err != nil {
		return nil, nil, nil, err
	}
	input, err := parsedABI.Pack("simulateHandleOp", op, common.BigToAddress(big.NewInt(0)), []byte{})

	if err != nil {
		return nil, nil, nil, err
	}
	data, err := client.CallContract(
		context.Background(),
		ethereum.CallMsg{
			From: common.BigToAddress(common.Big0),
			To:   &entryPoint,
			Data: input,
		},
		nil,
	)
	if err != nil {
		return nil, nil, nil, err
	}
	err = &revertError{
		reason: hexutil.Encode(data),
	}

	sim, simErr := NewExecutionResult(err)
	if simErr != nil {
		fo, foErr := NewFailedOp(err)
		if foErr != nil {
			return nil, nil, nil, fmt.Errorf("%s, %s", simErr, foErr)
		}
		return nil, nil, nil, NewRPCError(-32500, fo.Reason, fo)
	}

	code, err := client.CodeAt(context.Background(), op.Sender, nil)
	if err != nil {
		return nil, nil, nil, err
	}
	var est uint64 = 100000
	if len(code) > 0 {
		est, err = client.EstimateGas(context.Background(), ethereum.CallMsg{
			From: entryPoint,
			To:   &op.Sender,
			Data: op.CallData,
		})
		if err != nil {
			return nil, nil, nil, err
		}
	}

	pvg, err := CalcPreVerificationGas(op)
	if err != nil {
		return nil, nil, nil, err
	}

	return pvg, sim.PreOpGas, big.NewInt(int64(est)), err
}
