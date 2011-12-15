package main

import (
	"fmt"
	"math"
)

//Solo se puede acceder a las cosas que exporta el package, y eso es lo que está en mayúsculas
func main() {
	fmt.Println("Happy", math.Pi, "Day")
	fmt.Println("Número: ",add(4,5))
	fmt.Println("Número: ",mult(4,5))
	a,b := swap("4","5")
	fmt.Println("Al revez: ",a,b)
	fmt.Println(split(17))
	variables()
	forSample()
	pruebaStruct()
	pruebaMaps()
}

//Las funciones pueden ir en cualquier lugar, no hace falta ponerlas antes
//hay que definir el tipo de dato, luego de los parámetros
func add(x int, y int) int {
	return x+y
}

//restar, pero especificando solo 1 tipo de dato, porque son iguales, el último hay que ponerlo y lo repite
func mult(x, y int) int {
	return x*y
}

//Una funció puede devolver cualquier cantidad de resultados: 
func swap(x, y string) (string, string){
	return y,x
}

//si pones return sin especificar, lo que devuelve son las dos variables que se había dicho 
//que tenía la función con el valor que tengan en ese momento
func split(sum int) (x, y int) {
	x = sum /2
	y = sum /2
	return
}

//si se pone "var" luego se pueden poner muchas variables y al final el tipo de dato
func variables(){
	var x,y,t,e int = 1,2,3,4 		//esto es un inicializador, le asigna los números
	var f,h,j = false, true, "noooo" 	//si se ponen inicializadores, no hace falta poner el tipo de dato, lo indiere
	a,b,c := 1,2,3 				//si estoy dentro de una función puedo poner := en lugar del var
	const mundo = "世界" 			//si se pone const ... = se declara una constante
	const PI = 3.14

	fmt.Println(x,y,t,e)
	fmt.Println(f,h,j)
	fmt.Println(a,b,c)
	fmt.Println(mundo, PI)
}

//lo único que existe es el for, no hay que poner () pero si {}
func forSample(){

	cant := 0
	for i:=0;i<10;i++ {
		cant++
	}
	fmt.Println("cant:",cant)

	//se puede dejar vacío la primera y última instrucción del for, simplemente no hace nada
	sum := 1
	for ; sum < 1000; { //no hace falta ponerlos quedaría así:  for sum < 1000 { --> como un while, para que haga un loop forever: for ;; { o for {
		sum += sum
	}
	fmt.Println(sum)

	//lo mismo aplica para el if, sin () y con {}
	if cant<10 {
		fmt.Println("si")
	}

	//Se le pueden poner precondiciones al if, así como al for, que sólo son de scope local al if:

	if fer:=4;fer<cant {
		fmt.Println("fer es menor a cant")
	} else {
		fmt.Println(fer) //la variable fer está disponible en el if
	}
}

//Método que prueba el manejo de structs
func pruebaStruct(){

	type vertice struct {
		X int
		Y int
	}
	
	a:= vertice{1, 2}
	fmt.Println("vertice: ",a)

	a.Y=4 //se accede a la struct con un "."
	fmt.Println("Y: ",a.Y)

	b:= &a //asigna a "b" un puntero a "a"
	fmt.Println("b Y: ",b.Y)

	c:= vertice{} 		//implicitamente inicializa las variables a 0
	d:= vertice{Y:1} 	//si le pongo el nombre asigna a la variable
	e:= &vertice{Y:1} 	//es de tipo puntero a vertice
 	fmt.Println(c)
	fmt.Println(d)
	fmt.Println(e)
	fmt.Println(e.Y) 	//resuelve sola la indirección de punteros

	//la expresión new(T) donde "T" es un struct, crea un puntero a una struct nueva inicializada vacía ej: vertice{}
	var nuevo *vertice = new(vertice)
	nuevo2 := new(vertice) //más fácil
	
	fmt.Println(nuevo)
	fmt.Println(nuevo2)
	fmt.Println(nuevo2.X)
	
}

//Método que muestra el manejo de mapas
func pruebaMaps(){
	
	type vertice struct {
		X int
		Y int
	}

	//declaramos la variable del tipo map
	var coordenadas map[string]vertice //crea un variable map llamada "coordenadas" donde las keys son strings y los values son del tipo vertice

	//se crea el mapa de tipo coordenada
	coordenadas = make(map[string]vertice) //con el make se crea un mapa, es lo mismo que el new

	coordenadas["cero absoluto"] = vertice{0,0}
	fmt.Println("cero: ",coordenadas["cero absoluto"])

	//mapas literales
	var mapa = map[string]vertice{
					"cero":vertice{3,4},
					"uno":vertice{4,5},
					"dos":vertice{5,6},
       				}

	fmt.Println("mapa literal: ",mapa)

	//pruebo de acceder directo, anda joya!	
	fmt.Println("mapa uno: ",mapa["cero"].X)

	//se puede iniciar sin especificar que el tipo de elemento es vertice si ya fue definido así en el mapa
	mapa = map[string]vertice{
				"cero":{3,4},
				"uno":{4,5},
				"dos":{5,6},
       				}

	fmt.Println("mapa literal: ",mapa)

}

/* 
NOTAS:
------

Todos los tipos de dato son estos:

bool
string
int  int8  int16  int32  int64
uint uint8 uint16 uint32 uint64 uintptr
float32 float64
complex64 complex128

*/

