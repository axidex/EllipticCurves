package cypher

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"io"
	"math/big"
)

var (
	ErrKeyDataTooLong = fmt.Errorf("ecies: can't supply requested key data")
	ErrSharedTooLong  = fmt.Errorf("ecies: shared secret is too long")
	ErrInvalidMessage = fmt.Errorf("ecies: invalid message")
)

// PublicKey is a representation of an elliptic curve public key.
type PublicKey struct {
	X *big.Int
	Y *big.Int
	elliptic.Curve
	Params *ECIESParams
}

// PrivateKey is a representation of an elliptic curve private key.
type PrivateKey struct {
	PublicKey
	D *big.Int
}

func MaxSharedKeyLength(pub *PublicKey) int {
	return (pub.Curve.Params().BitSize + 7) / 8
}

// GenerateKey Generate an elliptic curve public / private keypair. If params is nil,
// the recommended default parameters for the key will be chosen.
func GenerateKey(rand io.Reader, curve elliptic.Curve, params *ECIESParams) (prv *PrivateKey, err error) {
	pb, err := ecdsa.GenerateKey(curve, rand)
	if err != nil {
		return
	}

	prv = new(PrivateKey)
	prv.PublicKey.X = pb.X
	prv.PublicKey.Y = pb.Y
	prv.PublicKey.Curve = curve
	prv.D = pb.D
	if params == nil {
		params = ParamsFromCurve(curve)
	}
	prv.PublicKey.Params = params
	return
}

// GenerateShared ECDH key agreement method used to establish secret keys for encryption.
func (prv *PrivateKey) GenerateShared(pub *PublicKey, skLen, macLen int) (sk []byte, err error) {
	if prv.PublicKey.Curve != pub.Curve {
		return nil, ErrInvalidCurve
	}
	if skLen+macLen > MaxSharedKeyLength(pub) {
		return nil, ErrSharedKeyTooBig
	}
	x, _ := pub.Curve.ScalarMult(pub.X, pub.Y, prv.D.Bytes())
	if x == nil {
		return nil, ErrSharedKeyIsPointAtInfinity
	}

	sk = make([]byte, skLen+macLen)
	skBytes := x.Bytes()
	copy(sk[len(sk)-len(skBytes):], skBytes)
	return sk, nil
}

// ExportECDSA Export an ECIES private key as an ECDSA private key.
func (prv *PrivateKey) ExportECDSA() *ecdsa.PrivateKey {
	pub := &prv.PublicKey
	pubECDSA := pub.ExportECDSA()
	return &ecdsa.PrivateKey{PublicKey: *pubECDSA, D: prv.D}
}

// ExportECDSA Export an ECIES public key as an ECDSA public key.
func (pub *PublicKey) ExportECDSA() *ecdsa.PublicKey {
	return &ecdsa.PublicKey{Curve: pub.Curve, X: pub.X, Y: pub.Y}
}
