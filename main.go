package main

import (
	"bufio"
	"fmt"
	"os"
	"github.com/XLex0/XCipherGo/utils"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Ingrese la ruta:")
	path, _ := reader.ReadString('\n')

	path = path[:len(path)-1]

	err := utils.BinToFile(path)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Conversión completada.")
	}
}
