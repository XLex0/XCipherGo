package main

import (
	"net/http"
	"os"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/XLex0/XCipherGo/utils"
	"log"
	"fmt"
	"time"
)

func main() {
	// Inicializamos Gin
	r := gin.Default()

	// Rutas para cada operación de cifrado/desencriptado
	r.POST("/encrypt", encryptHandler)      // Para cifrar archivos
	r.POST("/decrypt", decryptHandler)      // Para descifrar archivos
	r.POST("/cesarenc", cesarEncryptHandler) // Para encriptar con César
	r.POST("/cesardec", cesarDecryptHandler) // Para desencriptar con César

	// Iniciamos el servidor en localhost:8080
	r.Run(":8080")
}

// Handler para cifrar archivos
func encryptHandler(c *gin.Context) {
	// Obtener el archivo subido
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "Archivo no proporcionado")
		log.Println("Error al obtener el archivo:", err)
		return
	}

	// Obtener la contraseña
	pass := c.PostForm("pass")
	if pass == "" {
		c.String(http.StatusBadRequest, "Contraseña no proporcionada")
		log.Println("Contraseña no proporcionada")
		return
	}

	// Crear un nombre único para el archivo temporal
	inputPath := fmt.Sprintf("./temp_in_%d", time.Now().UnixNano()) // Usamos la hora actual en nanosegundos como parte del nombre
	outputPath := fmt.Sprintf("./temp_out_%d", time.Now().UnixNano())

	log.Println("Guardando archivo en:", inputPath)
	if err := c.SaveUploadedFile(file, inputPath); err != nil {
		c.String(http.StatusInternalServerError, "Error al guardar el archivo")
		log.Println("Error al guardar el archivo:", err)
		return
	}
	log.Println("Archivo guardado con éxito:", inputPath)

	// Realizar la encriptación
	if err := utils.EncryptFile(inputPath, outputPath, pass); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		log.Println("Error al encriptar el archivo:", err)
		return
	}

	// Enviar el archivo encriptado como respuesta
	c.FileAttachment(outputPath, "encrypted_"+file.Filename)
	log.Println("Archivo encriptado y enviado:", outputPath)

	// Eliminar archivos temporales después de enviarlos
	os.Remove(inputPath)
	os.Remove(outputPath)
}

// Handler para descifrar archivos
func decryptHandler(c *gin.Context) {
	// Obtener el archivo subido
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "Archivo no proporcionado")
		log.Println("Error al obtener el archivo:", err)
		return
	}

	// Obtener la contraseña
	pass := c.PostForm("pass")
	if pass == "" {
		c.String(http.StatusBadRequest, "Contraseña no proporcionada")
		log.Println("Contraseña no proporcionada")
		return
	}

	// Crear un nombre único para el archivo temporal
	inputPath := fmt.Sprintf("./temp_in_%d", time.Now().UnixNano())
	outputPath := fmt.Sprintf("./temp_out_%d", time.Now().UnixNano())

	log.Println("Guardando archivo en:", inputPath)
	if err := c.SaveUploadedFile(file, inputPath); err != nil {
		c.String(http.StatusInternalServerError, "Error al guardar el archivo")
		log.Println("Error al guardar el archivo:", err)
		return
	}
	log.Println("Archivo guardado con éxito:", inputPath)

	// Realizar la desencriptación
	if err := utils.DecryptFile(inputPath, outputPath, pass); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		log.Println("Error al desencriptar el archivo:", err)
		return
	}

	// Enviar el archivo desencriptado como respuesta
	c.FileAttachment(outputPath, "decrypted_"+file.Filename)
	log.Println("Archivo desencriptado y enviado:", outputPath)

	// Eliminar archivos temporales después de enviarlos
	os.Remove(inputPath)
	os.Remove(outputPath)
}

// Handler para encriptar con cifrado César
func cesarEncryptHandler(c *gin.Context) {
	// Obtener el archivo subido
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "Archivo no proporcionado")
		log.Println("Error al obtener el archivo:", err)
		return
	}

	// Obtener el desplazamiento para César
	shiftStr := c.PostForm("shift")
	shift, _ := strconv.Atoi(shiftStr) // Ya es obligatorio en el frontend

	// Crear un nombre único para el archivo temporal
	inputPath := fmt.Sprintf("./temp_in_%d", time.Now().UnixNano())
	log.Println("Guardando archivo en:", inputPath)
	if err := c.SaveUploadedFile(file, inputPath); err != nil {
		c.String(http.StatusInternalServerError, "Error al guardar el archivo")
		log.Println("Error al guardar el archivo:", err)
		return
	}
	log.Println("Archivo guardado con éxito:", inputPath)

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
	os.Remove(inputPath)
}

// Handler para desencriptar con cifrado César
func cesarDecryptHandler(c *gin.Context) {
	// Obtener el archivo subido
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "Archivo no proporcionado")
		log.Println("Error al obtener el archivo:", err)
		return
	}

	// Obtener el desplazamiento para César
	shiftStr := c.PostForm("shift")
	shift, _ := strconv.Atoi(shiftStr) // Ya es obligatorio en el frontend

	// Crear un nombre único para el archivo temporal
	inputPath := fmt.Sprintf("./temp_in_%d", time.Now().UnixNano())
	log.Println("Guardando archivo en:", inputPath)
	if err := c.SaveUploadedFile(file, inputPath); err != nil {
		c.String(http.StatusInternalServerError, "Error al guardar el archivo")
		log.Println("Error al guardar el archivo:", err)
		return
	}
	log.Println("Archivo guardado con éxito:", inputPath)

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
	outputPath := fmt.Sprintf("./temp_out_%d", time.Now().UnixNano())
	err = utils.Base64ToFile(base64Str, outputPath)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		log.Println("Error al convertir Base64 a archivo:", err)
		return
	}

	// Enviar archivo desencriptado como respuesta
	c.FileAttachment(outputPath, "cesar_decrypted")
	log.Println("Archivo desencriptado y enviado:", outputPath)

	// Eliminar archivos temporales después de enviarlos
	os.Remove(inputPath)
	os.Remove(outputPath)
}
