package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

func main() {

	frec := frecuencia()

	dat, err := ioutil.ReadFile("chino.txt")
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

	f, err := os.Create("chino-result.txt")
	check(err)
	defer f.Close()

	for caracter, cantidad := range caracteres {
		if frec[caracter] != 0 {
			frase := caracter + "	" + strconv.Itoa(cantidad) + "	" + strconv.Itoa(frec[caracter]) + "\n"
			f.WriteString(frase)
		}
	}

	f.Sync()

	fmt.Println("fin")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func frecuencia() map[string]int {
	file, err := os.Open("caracteres.txt")
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	contador := 1
	frecuencia := make(map[string]int)
	for scanner.Scan() {
		linea := scanner.Text()
		pos := contador
		caracter := ""
		for i, c := range linea {
			caracter = fmt.Sprintf("%c", c)
			if i == (len(strconv.Itoa(contador)) + 1) {
				break
			}
		}
		frecuencia[caracter] = pos
		contador++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return frecuencia
}
