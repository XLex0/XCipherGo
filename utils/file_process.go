package utils

import (
	"fmt"
	"io/ioutil"
	"github.com/h2non/filetype"
)
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

	outputFileName := pathFile +"."+ kind.Extension

	err = ioutil.WriteFile(outputFileName, binData, 0644)
	if err != nil {
		return fmt.Errorf("error al guardar el archivo restaurado: %v", err)
	}

	fmt.Printf("Archivo restaurado guardado como: %s\n", outputFileName)
	return nil
}
