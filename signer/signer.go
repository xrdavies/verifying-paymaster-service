package signer

import (
	"crypto/ecdsa"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/ququzone/verifying-paymaster-service/config"
)

type Signer struct {
	PrivateKey *ecdsa.PrivateKey
}

func NewSigner() (*Signer, error) {
	conf := config.Config()
	keyData, err := os.ReadFile(conf.Keystore)
	if err != nil {
		return nil, err
	}
	keystore, err := keystore.DecryptKey(keyData, conf.Passphrase)
	if err != nil {
		return nil, err
	}
	return &Signer{
		PrivateKey: keystore.PrivateKey,
	}, nil
}

func (s *Signer) Eth_signVerifyingPaymaster(op map[string]any) (string, error) {
	return "hello", nil
}
