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

void imprimePersona(Persona *per) {
    fprintf(stderr, "Edad: %d, Teles: %d\n", per->edad, per->televisores);
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
import "fmt"
import "unsafe"
import "net/http"
import "strconv"

func webserver() {
	http.Handle("/", http.HandlerFunc(processRequest))
	http.ListenAndServe("0.0.0.0:8080", nil)
}

func processRequest(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Random desde webserver: " + strconv.Itoa(Random())))
}

func main() {

	go webserver()

	fmt.Println("Hola Fer")
	fmt.Println("Crea Semilla")
	Seed(23)
	fmt.Println("Random: ", Random())
	
	persona := GoPersona{25, 2}	
	C.imprimePersona((*C.Persona)(unsafe.Pointer(&persona)))
	
	animal := GoAnimal{}
	animal.nombre = C.CString("Gato")
	animal.patas = 4

	var ciudad GoCiudad = GoCiudad(C.pass_struct_animal((*C.Animal)(unsafe.Pointer(&animal))))
	fmt.Println("Ciudad: ", C.GoString(ciudad.nombre))
	C.free(unsafe.Pointer(ciudad.nombre))

}

//GoPersona Tipo de Dato que linkea a la estructura de C
type GoPersona _Ctype_Persona
//GoAnimal Tipo de Dato que linkea a la estructura de C
type GoAnimal _Ctype_Animal
//GoCiudad Tipo de Dato que linkea a la estructura de C
type GoCiudad _Ctype_struct_ciuda


//Seed crea una semilla
func Seed(semilla int) {
	C.srandom(C.uint(semilla))
}

//Random Devuelve un número aleatorio desde C
func Random() int {
	return int(C.random())
}
