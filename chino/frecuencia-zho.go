package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func frecuencia() map[string]int {
	file, err := os.Open("/Users/fscasserra/code/go/src/github.com/Fersca/Go-practice/caracteres.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	contador := 1
	frecuencia := make(map[string]int)
	for scanner.Scan() {
		linea := scanner.Text()
		pos := contador
		caracter := ""
		for i, c := range linea {
			s := fmt.Sprintf("%c", c)
			caracter = s
			if i == (len(strconv.Itoa(i)) + 1) {
				break
			}
		}
		fmt.Println(caracter, " - ", pos)
		frecuencia[caracter] = pos
		contador++
		if contador == 100 {
			return frecuencia
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return frecuencia
}
