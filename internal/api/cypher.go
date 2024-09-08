package api

import (
	"crypto/rand"
	"github.com/axidex/elliptic/internal/cypher"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

// @Summary Encrypt data
// @Description Encrypt the provided text using the given public key
// @Tags encryption
// @Accept multipart/form-data
// @Produce octet-stream
// @Param text formData string true "Text"
// @Param pemKey formData file true "PEM public key file"
// @Success 200 {string} string "Encrypted data as binary"
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/cypher/elliptic/encrypt [post]
func (app *App) encrypt(c *gin.Context) {
	text := c.Request.FormValue("text")
	if len(text) == 0 {
		app.logger.Warnf("Provided text is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "text is empty"})
		return
	}

	pemFile, _, err := c.Request.FormFile("pemKey")
	if err != nil {
		app.logger.Infof("Error retrieving PEM key: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "PEM key is required"})
		return
	}
	defer pemFile.Close()

	pemKey, err := io.ReadAll(pemFile)
	if err != nil {
		app.logger.Infof("Error reading PEM key: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading PEM key"})
		return
	}

	// 	Public string `json:"public"  form:"public"`
	app.logger.Infof("Got task encryption")
	app.logger.Infof("Creating public key from user input")
	key, err := cypher.ImportPublicPEM(pemKey)
	if err != nil {
		app.logger.Infof("Not valid key: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "provide valid key"})
		return
	}

	encryptedText, err := cypher.Encrypt(rand.Reader, key, []byte(text), nil, nil)
	if err != nil {
		app.logger.Infof("Encryption error %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "encryption error"})
		return
	}

	c.Data(http.StatusOK, "application/octet-stream", encryptedText)
}

// decrypt handles decryption of data using a provided PEM key.
// @Summary Decrypt encrypted data using a PEM key
// @Description Decrypts the provided encrypted data using the provided PEM key
// @Tags encryption
// @Accept multipart/form-data
// @Produce json
// @Param pemKey formData file true "PEM private key file"
// @Param encryptedData formData file true "Encrypted data file"
// @Success 200 {object} ResultData "Successfully decrypted data"
// @Failure 400 {object} map[string]any "Bad request"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/cypher/elliptic/decrypt [post]
func (app *App) decrypt(c *gin.Context) {
	// Retrieve the PEM key from the form
	pemFile, _, err := c.Request.FormFile("pemKey")
	if err != nil {
		app.logger.Infof("Error retrieving PEM key: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "PEM key is required"})
		return
	}
	defer pemFile.Close()

	pemKey, err := io.ReadAll(pemFile)
	if err != nil {
		app.logger.Infof("Error reading PEM key: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading PEM key"})
		return
	}

	// Retrieve the encrypted data from the form
	encryptedFile, _, err := c.Request.FormFile("encryptedData")
	if err != nil {
		app.logger.Infof("Error retrieving encrypted data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Encrypted data is required"})
		return
	}
	defer encryptedFile.Close()

	encryptedBytes, err := io.ReadAll(encryptedFile)
	if err != nil {
		app.logger.Infof("Error reading encrypted data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading encrypted data"})
		return
	}

	app.logger.Infof("Got task decryption")

	key, err := cypher.ImportPrivatePEM(pemKey)
	if err != nil {
		app.logger.Infof("Not valid key: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "provide valid key"})
		return
	}

	//app.logger.Infof("Decrypting data %s", encryptedBytes)

	decryptText, err := key.Decrypt(rand.Reader, encryptedBytes, nil, nil)
	if err != nil {
		app.logger.Infof("Decryption error %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "decryption error"})
		return
	}

	c.JSON(http.StatusOK, ResultData{string(decryptText)})
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
