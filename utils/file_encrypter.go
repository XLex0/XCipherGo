package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
	"io/ioutil"
	"encoding/base64"
	"github.com/h2non/filetype"
)

/*
Función para procesar la clave
*/

func padKey(password string) []byte {
	key := []byte(password)
	length := len(key)
	switch {
	case length <= 16:
		return append(key, make([]byte, 16-length)...)
	case length <= 24:
		return append(key, make([]byte, 24-length)...)
	case length <= 32:
		return append(key, make([]byte, 32-length)...)
	default:
		return key[:32]
	}
}

func EncryptFile(inputPath, outputPath string, pass string) error {
	key := padKey(pass)  // Correcta declaración de 'key'

	plainData, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	ciphertext := gcm.Seal(nonce, nonce, plainData, nil)
	return os.WriteFile(outputPath, ciphertext, 0644)
}

// Desencripta un archivo usando AES-GCM 
func DecryptFile(inputPath, outputPath string, pass string) error {
	key := padKey(pass)  // Correcta declaración de 'key'

	cipherData, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonceSize := gcm.NonceSize()
	if len(cipherData) < nonceSize {
		return errors.New("datos encriptados inválidos o incompletos")
	}

	nonce, ciphertext := cipherData[:nonceSize], cipherData[nonceSize:]
	plainData, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return errors.New("clave incorrecta o archivo corrupto")
	}

	return os.WriteFile(outputPath, plainData, 0644)
}

/*
Este convierte el binario a un formato legible para Windows [extensión]
A Linux le da igual XD
input: /ruta/archivo.bin [encrypted]
output: nil
*/

func BinToFile(pathFile string) error {
	binData, err := ioutil.ReadFile(pathFile)
	if err != nil {
		return fmt.Errorf("error al leer archivo: %v", err)
	}

	kind, err := filetype.Match(binData)
	if err != nil {
		return fmt.Errorf("Error en conversión")
	}

	if kind == filetype.Unknown {
		kind.Extension = ".unknown"
	}

	outputFileName := pathFile + "." + kind.Extension

	err = ioutil.WriteFile(outputFileName, binData, 0644)
	if err != nil {
		return fmt.Errorf("error al guardar el archivo restaurado: %v", err)
	}

	fmt.Printf("Archivo restaurado guardado como: %s\n", outputFileName)
	return nil
}

/*
Convierte cualquier archivo a binario
input: /ruta/archivo.ext
output: /ruta/archivo.bin
*/
func FileToBin(pathFile string) error {
	data, err := ioutil.ReadFile(pathFile)
	if err != nil {
		return fmt.Errorf("error al leer archivo: %v", err)
	}

	outputFileName := pathFile + ".bin"
	err = ioutil.WriteFile(outputFileName, data, 0644)
	if err != nil {
		return fmt.Errorf("error al guardar el archivo binario: %v", err)
	}

	fmt.Printf("Archivo convertido a binario guardado como: %s\n", outputFileName)
	return nil
}

func FileToBase64(pathFile string) (string, error) {
	data, err := os.ReadFile(pathFile)
	if err != nil {
		return "", fmt.Errorf("error al leer archivo: %v", err)
	}

	encoded := base64.StdEncoding.EncodeToString(data)
	return encoded, nil
}

func Base64ToFile(base64Str, outputPath string) error {
	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return fmt.Errorf("error al decodificar base64: %v", err)
	}

	err = os.WriteFile(outputPath, data, 0644)
	if err != nil {
		return fmt.Errorf("error al escribir archivo: %v", err)
	}

	return nil
}

func CaesarEncrypt(base64Str string, shift int) string {
	shift = shift % 256
	encrypted := make([]byte, len(base64Str))

	for i, b := range []byte(base64Str) {
		encrypted[i] = byte((int(b) + shift) % 256)
	}

	return string(encrypted)
}

func CaesarDecrypt(encryptedStr string, shift int) string {
	shift = shift % 256
	decrypted := make([]byte, len(encryptedStr))

	for i, b := range []byte(encryptedStr) {
		decrypted[i] = byte((int(b) - shift + 256) % 256)
	}

	return string(decrypted)
}
