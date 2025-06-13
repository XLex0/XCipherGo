package utils

import (
	"os"
	"testing"

	"github.com/XLex0/XCipherGo/utils"
)

func TestEncryptDecryptFile(t *testing.T) {
	input := "test_input.txt"
	encrypted := "test_input_encrypted.bin"
	decrypted := "test_input_decrypted.txt"
	content := []byte("Contenido secreto de prueba")
	password := "clave123"

	err := os.WriteFile(input, content, 0644)
	if err != nil {
		t.Fatalf("Error al crear archivo de entrada: %v", err)
	}

	defer os.Remove(input)
	defer os.Remove(encrypted)
	defer os.Remove(decrypted)

	if err := utils.EncryptFile(input, encrypted, password); err != nil {
		t.Errorf("Falló EncryptFile: %v", err)
	}

	if err := utils.DecryptFile(encrypted, decrypted, password); err != nil {
		t.Errorf("Falló DecryptFile: %v", err)
	}

	result, err := os.ReadFile(decrypted)
	if err != nil {
		t.Fatalf("No se pudo leer el archivo desencriptado: %v", err)
	}

	if string(result) != string(content) {
		t.Errorf("Contenido no coincide. Esperado %q, obtenido %q", string(content), string(result))
	}
}

func TestBase64Conversion(t *testing.T) {
	original := "base64_test.txt"
	decoded := "base64_decoded.txt"
	content := []byte("Texto para codificar en base64")

	err := os.WriteFile(original, content, 0644)
	if err != nil {
		t.Fatalf("Error al crear archivo de prueba base64: %v", err)
	}

	defer os.Remove(original)
	defer os.Remove(decoded)

	encoded, err := utils.FileToBase64(original)
	if err != nil {
		t.Fatalf("Falló FileToBase64: %v", err)
	}

	if err := utils.Base64ToFile(encoded, decoded); err != nil {
		t.Fatalf("Falló Base64ToFile: %v", err)
	}

	data, err := os.ReadFile(decoded)
	if err != nil {
		t.Fatalf("No se pudo leer el archivo restaurado: %v", err)
	}

	if string(data) != string(content) {
		t.Errorf("Contenido base64 no coincide. Esperado %q, obtenido %q", content, data)
	}
}

func TestCaesarEncryptDecrypt(t *testing.T) {
	original := "Texto en base64 de prueba"
	shift := 5

	encrypted := utils.CaesarEncrypt(original, shift)
	decrypted := utils.CaesarDecrypt(encrypted, shift)

	if decrypted != original {
		t.Errorf("Error en CaesarEncrypt/CaesarDecrypt. Esperado %q, obtenido %q", original, decrypted)
	}
}
