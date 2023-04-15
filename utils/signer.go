package utils

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)

func SignMessage(privateKey *ecdsa.PrivateKey, message []byte) ([]byte, error) {
	data := []byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(message)))
	data = append(data, message...)
	sha := sha3.NewLegacyKeccak256()
	sha.Write(data)
	hash := sha.Sum(nil)

	signature, err := crypto.Sign(hash, privateKey)
	if err != nil {
		return nil, err
	}
	signature[crypto.RecoveryIDOffset] += 27

	return signature, nil
}
