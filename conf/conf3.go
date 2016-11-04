package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {

	numero, _ := strconv.Atoi(os.Args[1])
	fmt.Println("Inicia")

	finalizados := make(chan int)
	timer := time.NewTimer(3 * time.Second)

	for i := 1; i <= numero; i++ {
		go facto(finalizados, i)
	}
	fmt.Println("Creados..")

	cant := 0

	for {
		select {
		case value := <-finalizados:
			cant++
			if cant == numero {
				fmt.Println("terminaron todos! valor: ", value)
				return
			}

		case <-timer.C:
			fmt.Println("time OUT!!")
			close(finalizados)
			//return
		}
	}

}

func facto(finalizados chan<- int, numero int) {

	resultado := numero

	for i := numero - 1; i > 1; i-- {
		resultado = resultado * i
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Printf("facto(%v) = %v\n", numero, resultado)

	finalizados <- numero
}
