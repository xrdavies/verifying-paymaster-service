package signer

type Signer struct {
}

func NewSigner() (*Signer, error) {
	return &Signer{}, nil
}

func (s *Signer) Eth_signVerifyingPaymaster() (string, error) {
	return "hello", nil
}
