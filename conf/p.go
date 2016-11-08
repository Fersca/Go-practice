package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {

	fmt.Println("Hola conf!!")

	numero, _ := strconv.Atoi(os.Args[1])

	timeout := time.NewTimer(2 * time.Second)

	finalizados := make(chan int)
	cancelar := make(chan int)

	for i := 1; i <= numero; i++ {
		go facto(i, finalizados, cancelar)
	}

	contador := 0

	for {
		select {
		case value := <-finalizados:
			contador++
			if contador == numero {
				fmt.Println("Fin: ", value)
				return
			}

		case <-timeout.C:
			fmt.Println("Time out!!!")
			close(cancelar)
		}
	}

}

func facto(numero int, finalizados chan<- int, cancelar <-chan int) {

	resultado := numero

	for i := numero - 1; i > 1; i-- {
		resultado = resultado * i

		select {
		case <-cancelar:
			fmt.Println("Cancelado, ", numero)
			finalizados <- numero
			return
		default:
		}

		time.Sleep(100 * time.Millisecond)
	}

	fmt.Printf("fact(%v)=%v\n", numero, resultado)

	finalizados <- numero

}
