package cypher

import (
	"crypto/cipher"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"hash"
	"io"
)

var (
	ErrImport                     = fmt.Errorf("ecies: failed to import key")
	ErrInvalidCurve               = fmt.Errorf("ecies: invalid elliptic curve")
	ErrInvalidParams              = fmt.Errorf("ecies: invalid ECIES parameters")
	ErrInvalidPublicKey           = fmt.Errorf("ecies: invalid public key")
	ErrSharedKeyIsPointAtInfinity = fmt.Errorf("ecies: shared key is point at infinity")
	ErrSharedKeyTooBig            = fmt.Errorf("ecies: shared key params are too big")
)

// messageTag computes the MAC of a message (called the tag) as per
// SEC 1, 3.5.
func messageTag(hash func() hash.Hash, km, msg, shared []byte) []byte {
	if shared == nil {
		shared = make([]byte, 0)
	}
	mac := hmac.New(hash, km)
	mac.Write(msg)
	tag := mac.Sum(nil)
	return tag
}

// Generate an initialisation vector for CTR mode.
func generateIV(params *ECIESParams, rand io.Reader) (iv []byte, err error) {
	iv = make([]byte, params.BlockSize)
	_, err = io.ReadFull(rand, iv)
	return
}

// symEncrypt carries out CTR encryption using the block cipher specified in the
// parameters.
func symEncrypt(rand io.Reader, params *ECIESParams, key, m []byte) (ct []byte, err error) {
	c, err := params.Cipher(key)
	if err != nil {
		return
	}

	iv, err := generateIV(params, rand)
	if err != nil {
		return
	}
	ctr := cipher.NewCTR(c, iv)

	ct = make([]byte, len(m)+params.BlockSize)
	copy(ct, iv)
	ctr.XORKeyStream(ct[params.BlockSize:], m)
	return
}

// symDecrypt carries out CTR decryption using the block cipher specified in
// the parameters
func symDecrypt(rand io.Reader, params *ECIESParams, key, ct []byte) (m []byte, err error) {
	c, err := params.Cipher(key)
	if err != nil {
		return
	}

	ctr := cipher.NewCTR(c, ct[:params.BlockSize])

	m = make([]byte, len(ct)-params.BlockSize)
	ctr.XORKeyStream(m, ct[params.BlockSize:])
	return
}

func Encrypt(rand io.Reader, pub *PublicKey, m, s1, s2 []byte) (ctBase64 string, err error) {
	params := pub.Params
	if params == nil {
		fmt.Println("Params is nil! Generating params from curve")
		if params = ParamsFromCurve(pub.Curve); params == nil {
			err = ErrUnsupportedECIESParameters
			return
		}
	}
	R, err := GenerateKey(rand, pub.Curve, params)
	if err != nil {
		return
	}

	hash := params.Hash()
	z, err := R.GenerateShared(pub, params.KeyLen, params.MacLen)
	if err != nil {
		return
	}
	K, err := concatKDF(hash, z, s1, params.KeyLen+params.MacLen)
	if err != nil {
		return
	}
	Ke := K[:params.KeyLen]
	Km := K[params.MacLen:]
	hash.Write(Km)
	Km = hash.Sum(nil)
	hash.Reset()

	em, err := symEncrypt(rand, params, Ke, m)
	if err != nil || len(em) <= params.BlockSize {
		return
	}

	d := messageTag(params.Hash, Km, em, s2)

	Rb := elliptic.Marshal(pub.Curve, R.PublicKey.X, R.PublicKey.Y)
	ct := make([]byte, len(Rb)+len(em)+len(d))
	copy(ct, Rb)
	copy(ct[len(Rb):], em)
	copy(ct[len(Rb)+len(em):], d)
	ctBase64 = base64.StdEncoding.EncodeToString(ct)
	return
}

// Decrypt decrypts an ECIES ciphertext.
func (prv *PrivateKey) Decrypt(rand io.Reader, ct string, s1, s2 []byte) (m []byte, err error) {
	c, err := base64.StdEncoding.DecodeString(ct)
	if c == nil || len(c) == 0 || err != nil {
		err = ErrInvalidMessage
		return
	}
	params := prv.PublicKey.Params
	if params == nil {
		if params = ParamsFromCurve(prv.PublicKey.Curve); params == nil {
			err = ErrUnsupportedECIESParameters
			return
		}
	}
	curHash := params.Hash()

	var (
		rLen   int
		hLen   int = curHash.Size()
		mStart int
		mEnd   int
	)

	switch c[0] {
	case 2, 3, 4:
		rLen = (prv.PublicKey.Curve.Params().BitSize + 7) / 4
		if len(c) < (rLen + hLen + 1) {
			err = ErrInvalidMessage
			return
		}
	default:
		err = ErrInvalidPublicKey
		return
	}

	mStart = rLen
	mEnd = len(c) - hLen

	R := new(PublicKey)
	R.Curve = prv.PublicKey.Curve
	R.X, R.Y = elliptic.Unmarshal(R.Curve, c[:rLen])
	if R.X == nil {
		err = ErrInvalidPublicKey
		return
	}

	z, err := prv.GenerateShared(R, params.KeyLen, params.MacLen)
	if err != nil {
		return
	}

	K, err := concatKDF(curHash, z, s1, params.KeyLen+params.MacLen)
	if err != nil {
		return
	}

	Ke := K[:params.KeyLen]
	Km := K[params.MacLen:]
	curHash.Write(Km)
	Km = curHash.Sum(nil)
	curHash.Reset()

	d := messageTag(params.Hash, Km, c[mStart:mEnd], s2)
	if subtle.ConstantTimeCompare(c[mEnd:], d) != 1 {
		err = ErrInvalidMessage
		return
	}

	m, err = symDecrypt(rand, params, Ke, c[mStart:mEnd])
	return
}
