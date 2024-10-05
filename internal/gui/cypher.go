package gui

import (
	"crypto/rand"
	"errors"
	"fmt"
	"fyne.io/fyne/v2/dialog"
	"github.com/axidex/elliptic/internal/cypher"
)

var (
	ErrGeneratingKeys = errors.New("error generating keys")
	ErrImportingKeys  = errors.New("error importing keys")
	ErrEncrypt        = errors.New("error encrypt")
	ErrDecrypt        = errors.New("error decrypt")
)

func (app *AppGui) generateKeys() {
	app.logger.Infof("Generating keys")

	if app.publicKeyEntry == nil || app.privateKeyEntry == nil {
		app.logger.Errorf("Public or private key is nil")
		dialog.ShowError(ErrGeneratingKeys, app.w)
		return
	}

	// Генерируем случайные ключи
	keys, err := cypher.GenerateKey(rand.Reader, cypher.DefaultCurve, nil)
	if err != nil {
		app.logger.Errorf("GenerateKey err: %s", err)
		dialog.ShowError(ErrGeneratingKeys, app.w)
		return
	}

	curveInfoString := "y² = x³ - 3x + b\n" +
		"B: %d\n" +
		"X: %d\n" +
		"Y: %d\n" +
		"Curve: %s\n" +
		"BitSize: %d"

	curveInfo := fmt.Sprintf(curveInfoString, keys.Curve.Params().B, keys.Curve.Params().Gx, keys.Curve.Params().Gy, keys.Curve.Params().Name, keys.Curve.Params().BitSize)

	app.curveInfoEntry.SetText(curveInfo)

	private, err := cypher.ExportPrivatePEM(keys)
	if err != nil {
		app.logger.Errorf("Encoding err: %s", err)
		dialog.ShowError(ErrGeneratingKeys, app.w)
		return
	}

	public, err := cypher.ExportPublicPEM(&keys.PublicKey)
	if err != nil {
		app.logger.Errorf("Encoding err: %s", err)
		dialog.ShowError(ErrGeneratingKeys, app.w)
		return
	}

	// Обновляем текст в окнах с ключами
	app.privateKeyEntry.SetText(string(private))
	app.publicKeyEntry.SetText(string(public))
}

func (app *AppGui) encryptData() {
	text := app.openText.Text

	pemKey := []byte(app.publicKeyEntry.Text)

	app.logger.Infof("Got task encryption")
	key, err := cypher.ImportPublicPEM(pemKey)
	if err != nil {
		app.logger.Infof("Not valid key: %v", err)
		dialog.ShowError(ErrImportingKeys, app.w)
		return
	}

	encryptedText, err := cypher.Encrypt(rand.Reader, key, []byte(text), nil, nil)
	if err != nil {
		app.logger.Errorf("Encryption error %v", err)
		dialog.ShowError(ErrEncrypt, app.w)
		return
	}

	app.closedText.SetText(encryptedText)
}

func (app *AppGui) decryptData() {
	encryptedBytes := app.closedText.Text

	pemKey := []byte(app.privateKeyEntry.Text)

	app.logger.Infof("Got task decryption")

	key, err := cypher.ImportPrivatePEM(pemKey)
	if err != nil {
		app.logger.Infof("Not valid key: %v", err)
		dialog.ShowError(ErrImportingKeys, app.w)
		return
	}

	//app.logger.Infof("Decrypting data %s", encryptedBytes)

	decryptText, err := key.Decrypt(rand.Reader, encryptedBytes, nil, nil)
	if err != nil {
		app.logger.Infof("Decryption error %v", err)
		dialog.ShowError(ErrDecrypt, app.w)
		return
	}

	app.openText.SetText(string(decryptText))
}
