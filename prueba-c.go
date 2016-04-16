package main

/*
#include <stdlib.h>
#include <stdio.h>

typedef struct {
    int edad;
    int televisores;
} Persona;

typedef struct {
    char* nombre;
    int patas;
} Animal;

typedef struct ciuda {
    char* nombre;
} Ciudad;

void pass_struct(Persona *per) {
    fprintf(stderr, "[%d, %d]\n", per->edad, per->televisores);
};

struct ciuda pass_struct_animal(Animal *ani) {
    fprintf(stderr, "Animal: %s, patas: %d\n", ani->nombre, ani->patas);    
    struct ciuda ciu;
    printf("Nombre de la ciudad: ");
    char* aver = (char*)malloc(50);
    scanf("%s", aver);    
    ciu.nombre = aver;
    return ciu;
}

*/
import "C"

import (
	"fmt"
    "reflect"
    "unsafe"
	"net/http"
	"strconv"        
)

func webserver() {    
	http.Handle("/", http.HandlerFunc(processRequest))
	http.ListenAndServe("0.0.0.0:8080", nil)
}

func processRequest(w http.ResponseWriter, req *http.Request) {    
    w.Write([]byte("Random desde webserver: "+strconv.Itoa(Random())))
}

func main() {

    go webserver()
    
	fmt.Println("Número random: ", Random())    
    fmt.Println("Número random (semilla): ", Seed(23))
    
    persona := GoPersona{25, 2}
    animal := GoAnimal{}
    animal.nombre = C.CString("Gato")
    animal.patas = 4
    
    C.pass_struct((*C.Persona)(unsafe.Pointer(&persona)))
    
    var ciudad GoCiudad = GoCiudad(C.pass_struct_animal((*C.Animal)(unsafe.Pointer(&animal))))
    fmt.Println("Ciudad: ", C.GoString(ciudad.nombre))
    C.free(unsafe.Pointer(ciudad.nombre))
}

/*
type Gofoo struct {
    A int32
    B int32
}
*/

//De esta forma se tiene referenciada a la misma estructura y evitamos errores
type GoPersona _Ctype_Persona
type GoAnimal _Ctype_Animal
type GoCiudad _Ctype_struct_ciuda

func Random() int {
	return int(C.random())
}

func Seed(semilla int) int {
    fmt.Println("Semilla ", semilla)
    valor := C.srandom(C.uint(semilla))
    tipo := reflect.TypeOf(valor)
    fmt.Println("Tipo: ", tipo)
    return int(C.random())
}
