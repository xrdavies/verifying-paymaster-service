package signer

import (
	"crypto/ecdsa"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/ququzone/verifying-paymaster-service/config"
	"github.com/ququzone/verifying-paymaster-service/logger"
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
	logger.S().Infof("VerifyingPaymaster contract: %s", conf.Contract)
	logger.S().Infof("VerifyingPaymaster signer: %s", keystore.Address.String())
	return &Signer{
		PrivateKey: keystore.PrivateKey,
	}, nil
}

func (s *Signer) Eth_signVerifyingPaymaster(op map[string]any) (string, error) {
	return "hello", nil
}
