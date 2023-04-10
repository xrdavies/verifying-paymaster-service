package signer

type Signer struct {
}

func NewSigner() (*Signer, error) {
	return &Signer{}, nil
}

func (s *Signer) Eth_signVerifyingPaymaster(op map[string]any) (string, error) {
	return "hello", nil
}
