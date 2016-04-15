package main

/*
#include <stdlib.h>
#include <stdio.h>

typedef struct {
    int a;
    int b;
} Foo;

void pass_struct(Foo *in) {
    fprintf(stderr, "[%d, %d]\n", in->a, in->b);
}

*/
import "C"

import (
	"fmt"
    "reflect"
    "unsafe"
)

func main() {

	fmt.Println("Random, number from C: ", Random())    
    fmt.Println("Random con Seed, number from C: ", Seed(23))
    fmt.Println("Pass struct: ")
    
    foo := Gofoo{25, 26}
    C.pass_struct((*C.Foo)(unsafe.Pointer(&foo)))
       
}

/*
type Gofoo struct {
    A int32
    B int32
}
*/
//De esta forma se tiene referenciada a la misma estructura y evitamos errores
type Gofoo _Ctype_Foo

func Random() int {
	return int(C.random())
}

func Seed(i int) int {
    fmt.Println("Hola ", i)
    value := C.srandom(C.uint(i))
    tipe := reflect.TypeOf(value)
    fmt.Println("value type: ", tipe)
    return int(C.random())
}
