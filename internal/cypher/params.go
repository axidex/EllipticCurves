package cypher

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
)

// Standard ECIES parameters:
// * ECIES using AES128 and HMAC-SHA-256-16
// * ECIES using AES256 and HMAC-SHA-256-32
// * ECIES using AES256 and HMAC-SHA-384-48
// * ECIES using AES256 and HMAC-SHA-512-64

var (
	EciesAes128Sha256 = &ECIESParams{
		Hash:      sha256.New,
		hashAlgo:  crypto.SHA256,
		Cipher:    aes.NewCipher,
		BlockSize: aes.BlockSize,
		KeyLen:    16,
		MacLen:    16,
	}

	EciesAes256Sha256 = &ECIESParams{
		Hash:      sha256.New,
		hashAlgo:  crypto.SHA256,
		Cipher:    aes.NewCipher,
		BlockSize: aes.BlockSize,
		KeyLen:    32,
		MacLen:    16,
	}

	EciesAes256Sha384 = &ECIESParams{
		Hash:      sha512.New384,
		hashAlgo:  crypto.SHA384,
		Cipher:    aes.NewCipher,
		BlockSize: aes.BlockSize,
		KeyLen:    32,
		MacLen:    16,
	}

	EciesAes256Sha512 = &ECIESParams{
		Hash:      sha512.New,
		hashAlgo:  crypto.SHA512,
		Cipher:    aes.NewCipher,
		BlockSize: aes.BlockSize,
		KeyLen:    32,
		MacLen:    32,
	}
)

var paramsFromCurve = map[elliptic.Curve]*ECIESParams{
	elliptic.P256(): EciesAes128Sha256,
	elliptic.P384(): EciesAes256Sha384,
	elliptic.P521(): EciesAes256Sha512,
}

// ParamsFromCurve selects parameters optimal for the selected elliptic curve.
// Only the curves P256, P384, and P512 are supported.
func ParamsFromCurve(curve elliptic.Curve) (params *ECIESParams) {
	//return paramsFromCurve[curve]
	switch curve {
	case elliptic.P256():
		return EciesAes128Sha256
	case elliptic.P384():
		return EciesAes256Sha384
	case elliptic.P521():
		return EciesAes256Sha512
	default:
		return nil
	}
}

// DefaultCurve The default curve for this package is the NIST P256 curve, which
// provides security equivalent to AES-128.
var DefaultCurve = elliptic.P256()

var (
	ErrUnsupportedECDHAlgorithm   = fmt.Errorf("ecies: unsupported ECDH algorithm")
	ErrUnsupportedECIESParameters = fmt.Errorf("ecies: unsupported ECIES parameters")
)

type ECIESParams struct {
	Hash      func() hash.Hash // hash function
	hashAlgo  crypto.Hash
	Cipher    func([]byte) (cipher.Block, error) // symmetric cipher
	BlockSize int                                // block size of symmetric cipher
	KeyLen    int                                // length of symmetric key
	MacLen    int                                // HMAC Length
}
