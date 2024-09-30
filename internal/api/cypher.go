package api

import (
	"crypto/rand"
	"github.com/axidex/elliptic/internal/cypher"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary Encrypt data
// @Description Encrypt the provided text using the given public key
// @Tags encryption
// @Accept application/json
// @Produce text/plain
// @Param payload body EncryptRequest true "Payload"
// @Success 200 {string} string "Encrypted data"
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/cypher/elliptic/encrypt [post]
func (app *App) encrypt(c *gin.Context) {
	var req EncryptRequest

	// Попытка привязки данных из JSON тела
	if err := c.ShouldBindJSON(&req); err != nil {
		app.logger.Warnf("Invalid input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Проверка текста
	if len(req.Text) == 0 {
		app.logger.Warnf("Provided text is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "text is empty"})
		return
	}

	// 	Public string `json:"public"  form:"public"`
	app.logger.Infof("Got task encryption")
	app.logger.Infof("Creating public key from user input")
	key, err := cypher.ImportPublicPEM([]byte(req.PEMKey))
	if err != nil {
		app.logger.Infof("Not valid key: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "provide valid key"})
		return
	}

	encryptedText, err := cypher.Encrypt(rand.Reader, key, []byte(req.Text), nil, nil)
	if err != nil {
		app.logger.Infof("Encryption error %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "encryption error"})
		return
	}

	c.String(http.StatusOK, encryptedText)
}

// @Summary Decrypt data
// @Description Decrypt the provided text using the given public key
// @Tags encryption
// @Accept application/json
// @Produce text/plain
// @Param payload body EncryptRequest true "Payload"
// @Success 200 {string} string "Decrypted data"
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/cypher/elliptic/decrypt [post]
func (app *App) decrypt(c *gin.Context) {
	var req EncryptRequest

	// Попытка привязки данных из JSON тела
	if err := c.ShouldBindJSON(&req); err != nil {
		app.logger.Warnf("Invalid input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Проверка текста
	if len(req.Text) == 0 {
		app.logger.Warnf("Provided text is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "text is empty"})
		return
	}

	app.logger.Infof("Got task decryption")

	key, err := cypher.ImportPrivatePEM([]byte(req.PEMKey))
	if err != nil {
		app.logger.Infof("Not valid key: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "provide valid key"})
		return
	}

	//app.logger.Infof("Decrypting data %s", encryptedBytes)

	decryptText, err := key.Decrypt(rand.Reader, req.Text, nil, nil)
	if err != nil {
		app.logger.Infof("Decryption error %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "decryption error"})
		return
	}

	c.String(http.StatusOK, string(decryptText))
}

// @Summary Generate a public key
// @Description Generate a public key using the elliptic curve algorithm
// @Tags keys
// @Accept json
// @Produce json
// @Success 200 {object} Keys
// @Router /api/cypher/elliptic/keys [get]
func (app *App) generateKey(c *gin.Context) {
	keys, err := cypher.GenerateKey(rand.Reader, cypher.DefaultCurve, nil)
	if err != nil {
		app.logger.Errorf("GenerateKey err: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "generating keys error"})
		return
	}

	private, err := cypher.ExportPrivatePEM(keys)
	if err != nil {
		app.logger.Errorf("Encoding err: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "encoding keys error"})
		return
	}

	public, err := cypher.ExportPublicPEM(&keys.PublicKey)
	if err != nil {
		app.logger.Errorf("Encoding err: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "encoding keys error"})
		return
	}

	app.logger.Info("private key: ", keys.D, "\tpublic key: ", keys.PublicKey)

	c.JSON(http.StatusOK, Keys{
		Private: string(private),
		Public:  string(public),
	})
}
