package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

func main() {

	//leer el archivo
	dat, err := ioutil.ReadFile("/home/fersca/code/go/chino.txt")
	check(err)

	caracteres := make(map[string]int)
	dat8 := string(dat)
	for _, c := range dat8 {
		cant, exists := caracteres[string(c)]
		if exists {
			caracteres[string(c)] = cant + 1
		} else {
			caracteres[string(c)] = 1
		}
	}

	f, err := os.Create("/home/fersca/code/go/chino-result.txt")
	check(err)
	defer f.Close()

	for caracter, cantidad := range caracteres {
		frase := caracter + "	" + strconv.Itoa(cantidad) + "\n"
		f.WriteString(frase)
	}

	f.Sync()
	fmt.Println("fin")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
