package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {

	rand.Seed(int64(time.Now().Second()))

	//Obtener el String a encondear en el ADN, pasarlo a un ARRAY
	reader := bufio.NewReader(os.Stdin)
	mensaje, _ := reader.ReadString('\n')
	mensaje = mensaje[:len(mensaje)-1]
	mensaje = strings.ToUpper(mensaje)
	fmt.Println("Mensaje:", mensaje, len(mensaje))

	//Generar un número random entre 1000 y 10000.
	result := generaCrisp()

	for i := 0; i < len(mensaje); i++ {
		letra := string(mensaje[i])
		result = result + letra
		result = result + generaCrisp()
	}

	fmt.Println("DNA: ", result)

	d1 := []byte(result)
	err := ioutil.WriteFile("dna.txt", d1, 0644)
	if err != nil {
		fmt.Println("Error grabando file")
	}

}

var letras string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generaCrisp() string {

	tamanioCrisp := rand.Intn(30000) + 1000
	fmt.Println("Random: ", tamanioCrisp)

	//Generar un string random con X letras, de ese tamaño

	var crisp string
	for i := 0; i < tamanioCrisp; i++ {
		posLetra := rand.Intn(len(letras))
		letra := string(letras[posLetra])
		crisp = crisp + letra
	}

	var invertedCrisp string
	for i := len(crisp) - 1; i >= 0; i-- {
		letra := string(crisp[i])
		invertedCrisp = invertedCrisp + letra
	}

	return crisp + invertedCrisp

}
