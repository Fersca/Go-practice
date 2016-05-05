package main

import (
	"fmt"
	"time"
)

func main() {

	contador := 0

	canal := make(chan int, 5)
	canalHola := make(chan int, 5)
	canalCancel := make(chan int, 5)

	go mandarHola(canalCancel, canalHola, canal)

	for i := 0; i < 15; i++ {
		go procesar(&contador, canal, canalHola, canalCancel)
	}

	for range canal {
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("fin")
}

func procesar(contador *int, canal, canalHola, canalCancel chan int) {
	for {
		select {
		case canal <- 1:
			*contador++
			fmt.Println("Hola", *contador)
			time.Sleep(1 * time.Second)
		case <-canalHola:
			fmt.Println("LALALLAL 2")
		case e := <-canalCancel:
			fmt.Println("CANCELADA ", e)
			return
		}
	}
}

func mandarHola(canalCancelar, canalHola chan int, canal chan int) {

	cont := 0
	for {
		time.Sleep(2 * time.Second)
		canalHola <- 1
		if cont == 2 {
			close(canalCancelar)
			time.Sleep(500 * time.Millisecond)
			close(canal)
		}
		cont++
	}
}
