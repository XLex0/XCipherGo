package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/XLex0/XCipherGo/utils"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("1. BinToFile")
	fmt.Println("2. FileToBin")
	fmt.Printf("Seleccione una opción: ")
	opcion, _ := reader.ReadString('\n')
	opcion = strings.TrimSpace(opcion)

	fmt.Printf("Ingrese la ruta: ")
	path, _ := reader.ReadString('\n')
	path = strings.TrimSpace(path)

	var err error
	if opcion == "1" {
		err = utils.BinToFile(path)
	} else if opcion == "2" {
		err = utils.FileToBin(path)
	} else {
		fmt.Println("Opción no válida")
		return
	}

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Conversión completada.")
	}
}

/*Para ingresar la ruta en la consola, quitar las comillas
Ejemplo de ruta: C:\Users\Emilio\OneDrive\Escritorio\FAMILIA GENERAL\Quinga Quishpe Emilio Josue.pdf
*/
