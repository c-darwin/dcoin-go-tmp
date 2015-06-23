package controllers
import (
	"utils"
	"log"
	"fmt"
	"errors"
	"crypto/rand"
	"crypto/rsa"
	"encoding/pem"
	"crypto/x509"
	"encoding/base64"
	"strings"
	"encoding/json"
	"crypto/aes"
	"crypto/cipher"
	"io"
)

func genKeys() (string, string) {
	privatekey, _ := rsa.GenerateKey(rand.Reader, 2048)
	var pemkey = &pem.Block{Type : "RSA PRIVATE KEY", Bytes : x509.MarshalPKCS1PrivateKey(privatekey)}
	PrivBytes0 := pem.EncodeToMemory(&pem.Block{Type:  "RSA PRIVATE KEY", Bytes: pemkey.Bytes})

	PubASN1, _ := x509.MarshalPKIXPublicKey(&privatekey.PublicKey)
	pubBytes := pem.EncodeToMemory(&pem.Block{Type:  "RSA PUBLIC KEY", Bytes: PubASN1})
	s := strings.Replace(string(pubBytes),"-----BEGIN RSA PUBLIC KEY-----","",-1)
	s = strings.Replace(s,"-----END RSA PUBLIC KEY-----","",-1)
	sDec, _ := base64.StdEncoding.DecodeString(s)

	return string(PrivBytes0), fmt.Sprintf("%x", sDec)
}

func encrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

func (c *Controller) GenerateNewPrimaryKey() (string, error) {

	if c.SessRestricted!=0 {
		return "", errors.New("Permission denied")
	}

	c.r.ParseForm()
	password := c.r.FormValue("password")

	priv, pub := genKeys()
	if len(password) > 0 {
		encKey, err := encrypt(utils.Md5(password), []byte(priv))
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		priv = string(encKey)
	}
	json, err := json.Marshal(map[string]string{"private_key": priv, "public_key": pub, "password_hash": string(utils.DSha256(password))})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	log.Println(json)
	return string(json), nil
}

