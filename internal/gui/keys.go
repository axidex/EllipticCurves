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

	eciesInfoString := "Algorithm: %s\n" +
		"BlockSize: %d\tKeyLength: %d\n" +
		"X: %d\n" +
		"Y: %d\n" +
		"D: %s"

	paramName := GetNameByParam(keys.Params)
	eciesInfo := fmt.Sprintf(
		eciesInfoString,
		paramName,
		keys.PublicKey.Params.KeyLen, keys.Params.BlockSize,
		keys.X, keys.Y, keys.D,
	)

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

	eciesInfoString := "Algorithm: %s\n" +
		"BlockSize: %d\tKeyLength: %d\n" +
		"X: %d\n" +
		"Y: %d\n" +
		"D: 0"

	paramName := GetNameByParam(keys.Params)
	eciesInfo := fmt.Sprintf(
		eciesInfoString,
		paramName,
		keys.Params.KeyLen, keys.Params.BlockSize,
		keys.X, keys.Y,
	)

	app.curveInfoEntry.SetText(curveInfo)
	app.eciesInfo.SetText(eciesInfo)
}
