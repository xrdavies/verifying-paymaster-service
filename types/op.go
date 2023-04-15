package types

import (
	"encoding/hex"
	"errors"
	"math/big"
	"reflect"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/go-playground/validator"
	"github.com/mitchellh/mapstructure"
)

var (
	UserOpPrimitives = []abi.ArgumentMarshaling{
		{Name: "sender", InternalType: "Sender", Type: "address"},
		{Name: "nonce", InternalType: "Nonce", Type: "uint256"},
		{Name: "initCode", InternalType: "InitCode", Type: "bytes"},
		{Name: "callData", InternalType: "CallData", Type: "bytes"},
		{Name: "callGasLimit", InternalType: "CallGasLimit", Type: "uint256"},
		{Name: "verificationGasLimit", InternalType: "VerificationGasLimit", Type: "uint256"},
		{Name: "preVerificationGas", InternalType: "PreVerificationGas", Type: "uint256"},
		{Name: "maxFeePerGas", InternalType: "MaxFeePerGas", Type: "uint256"},
		{Name: "maxPriorityFeePerGas", InternalType: "MaxPriorityFeePerGas", Type: "uint256"},
		{Name: "paymasterAndData", InternalType: "PaymasterAndData", Type: "bytes"},
		{Name: "signature", InternalType: "Signature", Type: "bytes"},
	}

	UserOpType, _ = abi.NewType("tuple", "op", UserOpPrimitives)
	UserOpArr, _  = abi.NewType("tuple[]", "ops", UserOpPrimitives)

	validate = validator.New()
	onlyOnce = sync.Once{}
)

func getAbiArgs() abi.Arguments {
	return abi.Arguments{
		{Name: "UserOp", Type: UserOpType},
	}
}

type UserOperation struct {
	Sender               common.Address `json:"sender"               mapstructure:"sender"               validate:"required"`
	Nonce                *big.Int       `json:"nonce"                mapstructure:"nonce"                validate:"required"`
	InitCode             []byte         `json:"initCode"             mapstructure:"initCode"             validate:"required"`
	CallData             []byte         `json:"callData"             mapstructure:"callData"             validate:"required"`
	CallGasLimit         *big.Int       `json:"callGasLimit"         mapstructure:"callGasLimit"         validate:"required"`
	VerificationGasLimit *big.Int       `json:"verificationGasLimit" mapstructure:"verificationGasLimit" validate:"required"`
	PreVerificationGas   *big.Int       `json:"preVerificationGas"   mapstructure:"preVerificationGas"   validate:"required"`
	MaxFeePerGas         *big.Int       `json:"maxFeePerGas"         mapstructure:"maxFeePerGas"         validate:"required"`
	MaxPriorityFeePerGas *big.Int       `json:"maxPriorityFeePerGas" mapstructure:"maxPriorityFeePerGas" validate:"required"`
	PaymasterAndData     []byte         `json:"paymasterAndData"     mapstructure:"paymasterAndData"     validate:"required"`
	Signature            []byte         `json:"signature"            mapstructure:"signature"`
}

func (op *UserOperation) GetFactory() common.Address {
	if len(op.InitCode) < common.AddressLength {
		return common.HexToAddress("0x")
	}

	return common.BytesToAddress(op.InitCode[:common.AddressLength])
}

func (op *UserOperation) Pack() []byte {
	args := getAbiArgs()
	packed, _ := args.Pack(&struct {
		Sender               common.Address
		Nonce                *big.Int
		InitCode             []byte
		CallData             []byte
		CallGasLimit         *big.Int
		VerificationGasLimit *big.Int
		PreVerificationGas   *big.Int
		MaxFeePerGas         *big.Int
		MaxPriorityFeePerGas *big.Int
		PaymasterAndData     []byte
		Signature            []byte
	}{
		op.Sender,
		op.Nonce,
		op.InitCode,
		op.CallData,
		op.CallGasLimit,
		op.VerificationGasLimit,
		op.PreVerificationGas,
		op.MaxFeePerGas,
		op.MaxPriorityFeePerGas,
		op.PaymasterAndData,
		op.Signature,
	})

	enc := hexutil.Encode(packed)
	enc = "0x" + enc[66:]
	return (hexutil.MustDecode(enc))
}
func exactFieldMatch(mapKey, fieldName string) bool {
	return mapKey == fieldName
}

func decodeOpTypes(
	f reflect.Kind,
	t reflect.Kind,
	data interface{}) (interface{}, error) {
	// String to common.Address conversion
	if f == reflect.String && t == reflect.Array {
		return common.HexToAddress(data.(string)), nil
	}

	// String to big.Int conversion
	if f == reflect.String && t == reflect.Struct {
		n := new(big.Int)
		n, ok := n.SetString(data.(string), 0)
		if !ok {
			return nil, errors.New("bigInt conversion failed")
		}
		return n, nil
	}

	// Float64 to big.Int conversion
	if f == reflect.Float64 && t == reflect.Struct {
		n, ok := data.(float64)
		if !ok {
			return nil, errors.New("bigInt conversion failed")
		}
		return big.NewInt(int64(n)), nil
	}

	// String to []byte conversion
	if f == reflect.String && t == reflect.Slice {
		byteStr := data.(string)
		if len(byteStr) < 2 || byteStr[:2] != "0x" {
			return nil, errors.New("not byte string")
		}

		b, err := hex.DecodeString(byteStr[2:])
		if err != nil {
			return nil, err
		}
		return b, nil
	}

	return data, nil
}

func validateAddressType(field reflect.Value) interface{} {
	value, ok := field.Interface().(common.Address)
	if !ok || value == common.HexToAddress("0x") {
		return nil
	}

	return field
}

func validateBigIntType(field reflect.Value) interface{} {
	value, ok := field.Interface().(big.Int)
	if !ok || value.Cmp(big.NewInt(0)) == -1 {
		return nil
	}

	return field
}

func NewUserOperation(data map[string]any) (*UserOperation, error) {
	var op UserOperation

	// Convert map to struct
	config := &mapstructure.DecoderConfig{
		DecodeHook: decodeOpTypes,
		Result:     &op,
		ErrorUnset: false,
		MatchName:  exactFieldMatch,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return nil, err
	}
	if err := decoder.Decode(data); err != nil {
		return nil, err
	}

	// Validate struct
	onlyOnce.Do(func() {
		validate.RegisterCustomTypeFunc(validateAddressType, common.Address{})
		validate.RegisterCustomTypeFunc(validateBigIntType, big.Int{})
	})
	err = validate.Struct(op)
	if err != nil {
		return nil, err
	}

	return &op, nil
}
