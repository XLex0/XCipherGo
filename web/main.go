package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/XLex0/XCipherGo/utils"
	"github.com/gin-gonic/gin"
)

// Función para guardar archivo y generar un nombre único
func saveFile(c *gin.Context) (string, string, error) {
	file, err := c.FormFile("file")
	if err != nil {
		return "", "", fmt.Errorf("archivo no proporcionado: %v", err)
	}

	inputPath := fmt.Sprintf("./temp_in_%d", time.Now().UnixNano())
	if err := c.SaveUploadedFile(file, inputPath); err != nil {
		return "", "", fmt.Errorf("error al guardar el archivo: %v", err)
	}

	log.Println("Archivo guardado con éxito:", inputPath)
	outputPath := fmt.Sprintf("./temp_out_%d", time.Now().UnixNano())
	return inputPath, outputPath, nil
}

// Función para eliminar archivos temporales
func removeTempFiles(inputPath, outputPath string) {
	os.Remove(inputPath)
	os.Remove(outputPath)
}

// Handler para cifrar archivos
func encryptHandler(c *gin.Context) {
	inputPath, outputPath, err := saveFile(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		log.Println(err)
		return
	}

	pass := c.PostForm("pass")
	if pass == "" {
		c.String(http.StatusBadRequest, "Contraseña no proporcionada")
		log.Println("Contraseña no proporcionada")
		return
	}

	// Encriptar el archivo
	if err := utils.EncryptFile(inputPath, outputPath, pass); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		log.Println("Error al encriptar el archivo:", err)
		return
	}

	// Enviar el archivo encriptado como respuesta
	c.FileAttachment(outputPath, "encrypted_"+c.DefaultPostForm("file", ""))
	log.Println("Archivo encriptado y enviado:", outputPath)

	// Eliminar archivos temporales después de enviarlos
	removeTempFiles(inputPath, outputPath)
}

// Handler para descifrar archivos
func decryptHandler(c *gin.Context) {
	inputPath, outputPath, err := saveFile(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		log.Println(err)
		return
	}

	pass := c.PostForm("pass")
	if pass == "" {
		c.String(http.StatusBadRequest, "Contraseña no proporcionada")
		log.Println("Contraseña no proporcionada")
		return
	}

	// Desencriptar el archivo
	if err := utils.DecryptFile(inputPath, outputPath, pass); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		log.Println("Error al desencriptar el archivo:", err)
		return
	}

	// Usamos la función BinToFile para restaurar el archivo y darle la extensión correcta
	if err := utils.BinToFile(outputPath); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		log.Println("Error al restaurar el archivo:", err)
		return
	}

	// Enviar el archivo desencriptado con la extensión correcta
	c.FileAttachment(outputPath, "decrypted_"+c.DefaultPostForm("file", ""))
	log.Println("Archivo desencriptado y enviado:", outputPath)

	// Eliminar archivos temporales después de enviarlos
	removeTempFiles(inputPath, outputPath)
}

// Handler para encriptar con cifrado César
func cesarEncryptHandler(c *gin.Context) {
	inputPath, outputPath, err := saveFile(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		log.Println(err)
		return
	}

	// Obtener el desplazamiento para César
	shiftStr := c.PostForm("shift")
	shift, _ := strconv.Atoi(shiftStr)

	// Convertir el archivo a Base64
	base64Str, err := utils.FileToBase64(inputPath)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		log.Println("Error al convertir el archivo a Base64:", err)
		return
	}

	// Realizar encriptación César
	cipher := utils.CaesarEncrypt(base64Str, shift)
	log.Println("Encriptación César realizada con éxito")

	// Enviar archivo cifrado como respuesta
	c.Header("Content-Disposition", "attachment; filename=cesar_encrypted.txt")
	c.Data(http.StatusOK, "text/plain", []byte(cipher))
	log.Println("Archivo cifrado enviado")

	// Eliminar archivo temporal después de enviarlo
	removeTempFiles(inputPath, outputPath)
}

// Handler para desencriptar con cifrado César
func cesarDecryptHandler(c *gin.Context) {
	inputPath, outputPath, err := saveFile(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		log.Println(err)
		return
	}

	// Obtener el desplazamiento para César
	shiftStr := c.PostForm("shift")
	shift, _ := strconv.Atoi(shiftStr)

	// Leer archivo cifrado
	data, err := os.ReadFile(inputPath)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		log.Println("Error al leer el archivo cifrado:", err)
		return
	}

	// Desencriptar con César
	base64Str := utils.CaesarDecrypt(string(data), shift)

	// Convertir Base64 a archivo y guardarlo
	err = utils.Base64ToFile(base64Str, outputPath)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		log.Println("Error al convertir Base64 a archivo:", err)
		return
	}

	// Usamos la función BinToFile para restaurar el archivo y darle la extensión correcta
	if err := utils.BinToFile(outputPath); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		log.Println("Error al restaurar el archivo:", err)
		return
	}

	// Enviar archivo desencriptado con la extensión correcta
	c.FileAttachment(outputPath, "cesar_decrypted")
	log.Println("Archivo desencriptado y enviado:", outputPath)

	// Eliminar archivos temporales después de enviarlos
	removeTempFiles(inputPath, outputPath)
}

func main() {
	// Inicializamos Gin
	r := gin.Default()

	r.Static("/", "./static")

	r.POST("/encrypt", encryptHandler)       // Para cifrar archivos
	r.POST("/decrypt", decryptHandler)       // Para descifrar archivos
	r.POST("/cesarenc", cesarEncryptHandler) // Para encriptar con César
	r.POST("/cesardec", cesarDecryptHandler) // Para desencriptar con César

	// Iniciamos el servidor en localhost:8080
	r.Run(":8080")
}
