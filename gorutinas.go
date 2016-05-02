package main

import (
	"fmt"
	"time"
    "sync"
)

func main() {

    //Unbuffered channel
	//ejemplo1()
    
    //Buffered channel
	//ejemplo2()
    
    //Cancel event
	//ejemplo3()
    
    //Wait Group
    ejemplo4()

}

func ejemplo1() {
	//sync channel
	canal := make(chan int)

	contador := 0
	go printSecond(&contador, canal)
	go printSecond(&contador, canal)
	go printSecond(&contador, canal)

	for i := 0; i < 3; i++ {
		<-canal
		time.Sleep(1 * time.Second)
	}

}

func printSecond(contador *int, canal chan int) {

	canal <- 1
	*contador = *contador + 1
	now := time.Now()
	fmt.Println("Ahora: ", now.Second(), ", contador: ", *contador)

}

func ejemplo2() {
	//Buffered channels
	canalCliente := make(chan int, 5)

	for clientes := 1; clientes < 20; clientes++ {
		go atenderCliente(canalCliente, clientes)
	}

	for i := 1; i < 10; i++ {
		time.Sleep(1000 * time.Millisecond)
		<-canalCliente
		<-canalCliente
	}

	time.Sleep(2 * time.Second)

}

func atenderCliente(canalCliente chan<- int, cliente int) {

	canalCliente <- 1
	fmt.Println("Atendiendo cliente :", cliente)
	time.Sleep(1000 * time.Millisecond)
	//<- canalCliente
}

func ejemplo3() {

	//Buffered channels
	canalCliente := make(chan int, 5)
	canalCancelar := make(chan int)

	for clientes := 1; clientes < 20; clientes++ {
		go atenderCliente2(canalCliente, canalCancelar, clientes)
	}

	go cancelar(canalCancelar,canalCliente)

    for range canalCliente{
        time.Sleep(1000 * time.Millisecond)        
    }

	time.Sleep(2 * time.Second)

}

func atenderCliente2(canalCliente chan<- int, canalCancelar <-chan int, cliente int) {

	select {
	case canalCliente <- 1:
		fmt.Println("Atendiendo cliente :", cliente)
		time.Sleep(1000 * time.Millisecond)

	case <-canalCancelar:
		fmt.Println("Cancelado")
		return
	}
}

func cancelar(canalCancelar,canalCliente chan int) {

	time.Sleep(2 * time.Second)
    fmt.Println("Cancelando....")
	canalCancelar <- 1
	time.Sleep(1 * time.Second)
    fmt.Println("Cancelando TODOS!!")        
    close(canalCancelar)
    time.Sleep(1 * time.Second)
    close(canalCliente)
}

func ejemplo4(){
    
    var done sync.WaitGroup
    
    for i:=0;i<10;i++{                
        done.Add(1)        
        go func(){
            time.Sleep(1*time.Second)
            fmt.Println("Hola ...")
            done.Done()            
        }()        
    }
    
    done.Wait()
    fmt.Println("Terminaron...")    
    time.Sleep(2 * time.Second)
    fmt.Println("Fin...")
    
}
