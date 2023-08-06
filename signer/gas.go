package signer

import (
	"bytes"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ququzone/verifying-paymaster-service/types"
)

func CalcCallDataCost(op *types.UserOperation) float64 {
	cost := float64(0)
	for _, b := range op.Pack() {
		if b == byte(0) {
			cost += 4
		} else {
			cost += 16
		}
	}

	return cost
}

func CalcPerUserOpCost(op *types.UserOperation) float64 {
	opLen := math.Floor(float64(len(op.Pack())+31) / 32)
	cost := (25 * opLen) + 22874

	return cost
}

func CalcPreVerificationGas(op *types.UserOperation) (*big.Int, error) {
	// Sanitize fields to reduce as much variability due to length and zero bytes
	data, err := op.ToMap()
	if err != nil {
		return nil, err
	}
	data["preVerificationGas"] = hexutil.EncodeBig(big.NewInt(100000))
	data["verificationGasLimit"] = hexutil.EncodeBig(big.NewInt(1000000))
	data["callGasLimit"] = hexutil.EncodeBig(big.NewInt(1000000))
	data["signature"] = hexutil.Encode(bytes.Repeat([]byte{1}, len(op.Signature)))
	tmp, err := types.NewUserOperation(data)
	if err != nil {
		return nil, err
	}

	// Calculate the additional gas for adding this userOp to a batch.
	batchOv := (21000 / 1) + CalcCallDataCost(tmp)

	// The total PVG is the sum of the batch overhead and the overhead for this userOp's validation and
	// execution.
	pvg := batchOv + CalcPerUserOpCost(tmp)
	static := big.NewInt(int64(math.Round(pvg)))

	return static, nil
}
