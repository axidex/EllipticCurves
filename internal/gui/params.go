package gui

import (
	"crypto/elliptic"
	"github.com/axidex/elliptic/internal/cypher"
)

type CurveName string

const (
	curveP256 = "P256"
	curveP384 = "P384"
	curveP521 = "P521"
)

var CurveNames = []string{
	curveP256,
	curveP384,
	curveP521,
}

func GetNameByCurve(curve elliptic.Curve) CurveName {
	switch curve {
	case elliptic.P384():
		return curveP384
	case elliptic.P521():
		return curveP521
	default:
		return curveP256

	}
}

func (name CurveName) GetCurveByName() elliptic.Curve {
	switch name {
	case curveP384:
		return elliptic.P384()
	case curveP521:
		return elliptic.P521()
	default:
		return elliptic.P256()
	}
}

type ParamName string

const (
	paramAes128Sha256 = "AES128 HMAC-SHA-256-16"
	paramAes256Sha256 = "AES256 HMAC-SHA-256-32"
	paramAes256Sha384 = "AES256 HMAC-SHA-384-48"
	paramAes256Sha512 = "AES256 HMAC-SHA-512-64"
)

var ParamNames = []string{
	paramAes128Sha256,
	paramAes256Sha256,
	paramAes256Sha384,
	paramAes256Sha512,
}

func GetNameByParam(param *cypher.ECIESParams) ParamName {
	switch param {
	case cypher.EciesAes256Sha256:
		return paramAes256Sha256
	case cypher.EciesAes256Sha384:
		return paramAes256Sha384
	case cypher.EciesAes256Sha512:
		return paramAes256Sha512
	default:
		return paramAes128Sha256
	}
}

func (name ParamName) GetParamByName() *cypher.ECIESParams {
	switch name {
	case paramAes256Sha256:
		return cypher.EciesAes256Sha256
	case paramAes256Sha384:
		return cypher.EciesAes256Sha384
	case paramAes256Sha512:
		return cypher.EciesAes256Sha512
	default:
		return cypher.EciesAes128Sha256
	}
}
