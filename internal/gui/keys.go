package gui

import (
	"fmt"
	"github.com/axidex/elliptic/internal/cypher"
)

func (app *AppGui) setPrivateKeyInfo(keys *cypher.PrivateKey) {
	curveInfoString := "y² = x³ - 3x + b\n" +
		"B: %d\n" +
		"X: %d\n" +
		"Y: %d\n" +
		"Curve: %s\n" +
		"BitSize: %d"

	curveInfo := fmt.Sprintf(curveInfoString, keys.Curve.Params().B, keys.Curve.Params().Gx, keys.Curve.Params().Gy, keys.Curve.Params().Name, keys.Curve.Params().BitSize)

	eciesInfoString := "CryptoAlgorithm: %s BlockSize: %d\n" +
		"HashAlgorithm: %s\n" +
		"X: %d\n" +
		"Y: %d\n" +
		"D: %s\n" +
		"KeyLength: %d"
	eciesInfo := fmt.Sprintf(eciesInfoString, "AES", keys.Params.BlockSize, "SHA256", keys.X, keys.Y, keys.D, keys.PublicKey.Params.KeyLen)

	app.curveInfoEntry.SetText(curveInfo)
	app.eciesInfo.SetText(eciesInfo)
}

func (app *AppGui) setPublicKeyInfo(keys *cypher.PublicKey) {
	curveInfoString := "y² = x³ - 3x + b\n" +
		"B: %d\n" +
		"X: %d\n" +
		"Y: %d\n" +
		"Curve: %s\n" +
		"BitSize: %d"

	curveInfo := fmt.Sprintf(curveInfoString, keys.Curve.Params().B, keys.Curve.Params().Gx, keys.Curve.Params().Gy, keys.Curve.Params().Name, keys.Curve.Params().BitSize)

	eciesInfoString := "CryptoAlgorithm: %s BlockSize: %d\n" +
		"HashAlgorithm: %s\n" +
		"X: %d\n" +
		"Y: %d\n" +
		"KeyLength: %d"
	eciesInfo := fmt.Sprintf(eciesInfoString, "AES", keys.Params.BlockSize, "SHA256", keys.X, keys.Y, keys.Params.KeyLen)

	app.curveInfoEntry.SetText(curveInfo)
	app.eciesInfo.SetText(eciesInfo)
}
