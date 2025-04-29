package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"github.com/XLex0/XCipherGo/utils"
)

func main() {
	// Subcomandos
	encryptCmd := flag.NewFlagSet("encrypt", flag.ExitOnError)
	decryptCmd := flag.NewFlagSet("decrypt", flag.ExitOnError)
	cesarEncCmd := flag.NewFlagSet("cesarenc", flag.ExitOnError)
	cesarDecCmd := flag.NewFlagSet("cesardec", flag.ExitOnError)

	var (
		input    string
		output   string
		password string
		shift    int
	)

	// Flags para encrypt y decrypt
	encryptCmd.StringVar(&input, "in", "", "Archivo de entrada")
	encryptCmd.StringVar(&output, "out", "", "Archivo de salida")
	encryptCmd.StringVar(&password, "pass", "", "Contraseña")

	decryptCmd.StringVar(&input, "in", "", "Archivo cifrado")
	decryptCmd.StringVar(&output, "out", "", "Archivo de salida")
	decryptCmd.StringVar(&password, "pass", "", "Contraseña")

	// Flags para César
	cesarEncCmd.StringVar(&input, "in", "", "Archivo a cifrar (cifrado César sobre base64)")
	cesarEncCmd.IntVar(&shift, "shift", 3, "Desplazamiento César")
	cesarEncCmd.StringVar(&output, "out", "", "Archivo de salida")

	cesarDecCmd.StringVar(&input, "in", "", "Archivo cifrado (César sobre base64)")
	cesarDecCmd.IntVar(&shift, "shift", 3, "Desplazamiento César")
	cesarDecCmd.StringVar(&output, "out", "", "Archivo restaurado")

	if len(os.Args) < 2 {
		fmt.Println("Uso: encrypt | decrypt | cesarenc | cesardec")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "encrypt":
		encryptCmd.Parse(os.Args[2:])
		check(input, output, password)
		checkErr(utils.EncryptFile(input, output, password))

	case "decrypt":
		decryptCmd.Parse(os.Args[2:])
		check(input, output, password)
		checkErr(utils.DecryptFile(input, output, password))

	case "cesarenc":
		cesarEncCmd.Parse(os.Args[2:])
		check(input, output, "")
		base64Str, err := utils.FileToBase64(input)
		checkErr(err)
		cipher := utils.CaesarEncrypt(base64Str, shift)
		checkErr(os.WriteFile(output, []byte(cipher), 0644))

	case "cesardec":
		cesarDecCmd.Parse(os.Args[2:])
		check(input, output, "")
		data, err := os.ReadFile(input)
		checkErr(err)
		base64Str := utils.CaesarDecrypt(string(data), shift)
		checkErr(utils.Base64ToFile(base64Str, output))

	default:
		fmt.Println("Comando no reconocido. Usa: encrypt | decrypt | cesarenc | cesardec")
		os.Exit(1)
	}
}

func check(in, out, pass string) {
	if in == "" {
		log.Fatal("Falta -in")
	}
	if out == "" {
		log.Fatal("Falta -out")
	}
	if (os.Args[1] == "encrypt" || os.Args[1] == "decrypt") && pass == "" {
		log.Fatal("Falta -pass")
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
