package main

import (
	"fmt"
	"math"
	"time"
	"os"
	"net/http"
	"image"
	"runtime"
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
	pruebaSlices()
	functions()
	ranges()
	pruebaSwitch()
	methods()
	interfaces()
	errores()
	imagenes()
	//webservers()
	goroutines()
	channels()
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

	//modificar algo en el mapa es como simpre
	mapa["cero"] = vertice{2,3} //al parecer no hay que poner new ... raro, pensé que había que crear el objeto, 
	mapa["fer"] = *new(vertice) //al parecer así funciona, pone un puntero al elemento en el mapa
	element := mapa["cero"]     //así se busca un elemento
	fmt.Println("elemento: ",element)	 //se imprime lo más bien
	fmt.Println("elemento: ",mapa["fer"])	 //veo que se imprime el vacío que lo cree con el puntero
	delete(mapa,"cero")			 //borra un elemento del array

	//esto es loco, si pongo dos variables a la primera le 
	//asigna el elemento a la segundo true o false si está o no
	element,ok := mapa["uno"]		 
	fmt.Println("El valor :",element, "Presente: ",ok) //imprimo si está o no, debería
	element,ok = mapa["cero"]	//notese que no es := sinó = solo porque ya las declaré		 
	fmt.Println("El valor :",mapa["cero"], "Presente: ",ok) //debería poner false

}

func pruebaSlices(){

	//un slice apunta a un array de valores y también incluye un tamaño
	//[]T es un slice del tipo de elemento T
	s := []int{1,2,3,4,5,6}
	fmt.Println("slice: ",s)
	
	//imprimo el slice, uso el tamaño
	for i:=0;i<len(s);i++ {
		fmt.Println("Elemento: ",i,s[i])
	}

	//un slice puede ser recreado, para hacer esto hay que accederlo con slice[x:y], devuelve el rango X a Y-1
	fmt.Println("Slice: [1-4]: ",s[1:4])
	//si no se pone el inicio o el fin toma los límites
	fmt.Println("Slice: [:-4]: ",s[:4])
	fmt.Println("Slice: [1-:]: ",s[1:])
	
	//para crear un slice se puede usar el make, ahi se le puede definir el tamaño y la capacidad
	s2 := make([]int, 5)	//crea un slice de tamaño 5, pero con capacidad 5
	fmt.Println("tamaño: ",len(s2)," capacidad: ",cap(s2)," valores: ", s2)
	s3 := make([]int,3,4)	//crea un slice con tamaño 3 pero con capacidad 4
	fmt.Println("tamaño: ",len(s3)," capacidad: ",cap(s3)," valores: ", s3)
	//este último puso 3 elementos en el array, no entiendo entonces que es la capacidad, voy a probar de poner un elemento en [3]
	//s3[3] = 1	
	//no funcionó con lo cual no entiendo para que es la capacidad
	var z []int
	fmt.Println("tamaño: ",len(z)," capacidad: ",cap(z)," valores: ", z)
	if z == nil {
		fmt.Println("es nil!")
	}	

	//pasa saber más de slides ver este link: http://golang.org/doc/articles/slices_usage_and_internals.html
	//estaría bueno ver que diferencia exactamente hay entre length y capacity
}

func functions(){

	//las funciones puede definirse como valores:
	hipotenusa:= func(x,y float64) float64 {
			return math.Sqrt(x*x + y*y)
			}
	
	fmt.Println(hipotenusa(3,4))

	//al parecer se puede devolver una funcion desde una función... muy bueno, y obvio ya que es como un valor:
	
	//ver la función adder() declarada después de esta función.
	
	//se crean dos variables del tipo adder() entonces lo que contienen es una función
	pos, neg := adder(), adder()
	
	//las puedo usar después
	for i:=0;i<10;i++ {
		fmt.Println(pos(i))			
		fmt.Println(neg(-2*i))			
	}

	//como se puede ver después de ejecutar el ejemplo, cada función creada mantiene su propia variable interna.
}

//la función adder devuelve una función
func adder() func(int) int {
	sum :=0
	return func(x int) int {
		sum +=x
		return sum
	}
}

func ranges(){
	
	//un range sirve para iterar sobre un slice o un map
	
	numeros := []int{1,2,3,4,5,6,7,8}

	//imprime primero la posicion y luego el elemento en esa posicion, está bueno.
	for i, v:= range numeros {
		fmt.Println("Numero: ",i,v)
	}
	
	//no hace falta usar key y value, solo se puede usar key, pero para usar value hay que poner _
	otros := make([]int,10)
	
	//así se usa solo el i
	for i:= range otros {
		otros[i] = i*2
	}

	//así se usa solo el v
	for _,v:=range otros {
		fmt.Println("Numero: ",v)
	}
}

func pruebaSwitch(){

	//el switch es como siempre

	today := time.Now().Weekday()

	//cuando encuentra una que machea, corta salvo que se ponga fallthrough al final
	switch time.Saturday {

		case today+0:
			fmt.Println("Hoy")
		case today+1:
			fmt.Println("Mañana")
		case today+2:
			fmt.Println("Pasado Mañana")
		default:
			fmt.Println("Falta mucho")
	}

	//se pueden poner switch sin condiciones, es lo mismo que poner switch true, es una manera mas limpia de hacer muchos if-else

	t := time.Now()
	switch {
		case t.Hour() < 12:
		    fmt.Println("Buenos días!")
		case t.Hour() < 17:
		    fmt.Println("buenas tardes.")
		default:
		    fmt.Println("buenas noches.")
	}
}

func methods(){
	//go no tiene clases pero se pueden definir metodos en los structs

	//para crear un método lo que se hace es crear una struct y después atachar a esa struct un metodo, así:
	
	//Ver abajo de esto la creación del struct y el agregado del método

	//se crea un vertex y se llama a su método
	v := &Vertex{3, 4}
	fmt.Println(v.Abs())

	//de hecho se puede agregar un método a todo lo que esté definido en el package, no solo a los structs
	//este es un ejemplo donde una variable tiene un método.
	f:= MyFloat(-math.Sqrt(2))
	fmt.Println("numero: ",f.Abs())
	//entonces se define un type del tipo float y se le pone un método.

	//como pudo verse, un método recibió un puntero y otro un valor, ambos funcionan pero como siempre
	//el del puntero modifica otra estructura, no una copia de lo que le pasaron.
}

type MyFloat float64

func (a MyFloat) Abs() float64{
	//como se puede ver, el parámetro se puede usar dentro
	if a<0 {
		return float64(-a)
	} 
	return float64(a)
}

//se crea la struct
type Vertex struct {
	X,Y float64
}


//se crea la función Abs() y antes del nombre se le pone el receptor de la función "vertex"
func (v *Vertex) Abs() float64{
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func interfaces(){
	//las interfaces son lo mismo que en otro lenguaje, un contrato que dice que se van a implementar una serie de métodos
	//entonces si se crea una variable de ese tipo de interface, se sabe que todo lo que se le asigne implementa esos métodos.
	
	//definimos una interface afuera "Abser"
	//se crea una variable del tipo Abser, con lo cual tiene el metodo Abs()
	var a Abser
	
	f:= MyFloat(3)
	v:= Vertex{3,4}
	a = &v //funciona porque hay un método Abs() que recibe un *Vertex
	a = f //funciona porque MyFloat tiene un método que es Abs()
	//a = v //Vertex solo no implementa un metodo Abs() que reviva un v

	/*
	Devuelve esto:
	cannot use v (type Vertex) as type Abser in assignment:
	Vertex does not implement Abser (Abs method requires pointer receiver)
	*/

	fmt.Println(a.Abs())

	//Lo interesante es que una interface se puede definir como que implementa otras interfaces
	//por ejemplo, vea lo siguiente.
	var w Writer //se crea una variable del tipo writer que es una interface

	// os.Stdout implementa Writer
	w = os.Stdout

	fmt.Fprintf(w, "hello, writer\n")	
}

//se crea la interface
type Abser interface{
	Abs() float64
}

type Writer interface{
	Write(b []byte) (n int, err error)
}

func errores(){

	//al parecer un error es una interface que define el método Error() que devuelve un string
	//entonces hay muchos packages que sabiendo esto, esperan ejecutar el metodo Error()

	r:=run() //se crea una variable r la cual es un error. Entonces apunta a una estructura que implementa Error()
	//Si se llama a primeln va a imprimir el error porque sabe que es un error
	fmt.Println(r) 
	//lo que quiero explicar es que printLn cuando recibe una estructura que implementa error, se da cuenta de que se 
	//quiere imprimir un error entonces ejecuta el metodo Error() de la misma.

}

//Esta es la estructura
type MyError struct {
	When time.Time
	What string
}

//tiene el método de error (con lo cual implementa error):
func (e *MyError) Error() string {
	return fmt.Sprintf("Fer at %v, %s",
		e.When, e.What)
}

//Esta es la función pedorra que crea el objeto MyError.
func run() error {
	return &MyError{
		time.Now(),
		"it didn't work",
	}
}

func webservers(){
	//el package de HTTP maneja pedidos https usando cualquier value que implemente Handler.

	//por ejemplo, creamos una struct vacía y le agregamos una función llamada ServeHTTP
	//después se la podemos pasar como handler al http server, porque la misma implementa los metodos de un handler
	var h Hello
	http.ListenAndServe("localhost:4000",h)
	//como puede verse levanta un webserver y responde con h, que es el handler llamado Hello.

	//La interface es así:
	/*
	type Handler interface {
		ServeHTTP(w ResponseWriter,
			  r *Request)
	}
	*/

}


type Hello struct{}
func (h Hello) ServeHTTP(w http.ResponseWriter,r *http.Request) {
	fmt.Println("Request recibido")
	fmt.Fprint(w, "Hello!")
}

func imagenes(){
	
	/*

	La interface de la imagen es así:
	type Image interface {
		ColorModel() color.Model
		Bounds() Rectangle
		At(x, y int) color.Color
	}

	O sea, tiene un método que devuelve el modelo de color, otro que devuelve el color en una posicion x,y y otro que dice los bordes
	*/

	m := image.NewRGBA(image.Rect(0, 0, 100, 100)) //Se crea una imagen
	fmt.Println(m.Bounds())	//se le pide los bordes
	fmt.Println(m.At(0, 0).RGBA()) //se imprime el clor en la posicion 0,0, a su vez, se ejecuta el método RGBA de ese color.

}

func goroutines(){

	//cuando se pone delante la palabra "go" la función que se ejecuta después se ejecuta en un gorutine
	go say("hola en gorutine 0")
	go say("hola en gorutine 1")
	go say("hola en gorutine 2")
	go say("hola en gorutine 3")
	say("hola")

}

func say(s string){
	for i:=0; i<5; i++{		
		runtime.Gosched()
		time.Sleep(time.Second*1)
		fmt.Println(s)
	}
}


func channels(){

	//los channels son conductos en donde se envía y obtiene información.
	//se tienen que crear con make
	//y se envía información con la flecha chan c int, c <- 5, y se obtiene así a int, a <- c
	//por default el envío y recepción de información se bloquea hasta que la otra parte esté lista, esto permite sincronizar.
	c := make(chan int) //crea un canal de ints
	go suma(2,3,c)
	go suma2(4,5,c)
	x := <- c 
	fmt.Println("termina 1")	
	y := <- c 
	fmt.Println("termina 2")	
	fmt.Println("result: ",x+y) 
	//queda claro en este ejemplo que se bloquea cuando quiere recibir y aun no hay datos en "c"

	//ahora voy a tratar de poner dos datos al mismo tiempo en el channel y ver si realmente se bloquea al poner
	go suma3(4,5,c) //agrega e imprime cuando agrega
	go suma4(6,7,c)	//agrega e imprime cuando agrega, pero no debería poder porque ya hay un dato en el channel
	time.Sleep(time.Second*1)
	x = <- c
	time.Sleep(time.Second*1)
	x = <- c
	time.Sleep(time.Second*1)
	fmt.Println("fin")	
	//realmente se bloquean al poner datos en "c", porque el "poniendo x" lo imprime luego de que saco datos del channel

	//se puede hacer tambien un channel con un buffer, entonces solo se bloquea si se llena el buffer o si queda vacío, ta bueno, a probarlo....
	c2 := make(chan int,3) //crea un canal de ints con 3 valores de buffer	
	go suma3(4,5,c2) 
	go suma4(6,7,c2) 
	go suma3(4,5,c2) 
	go suma4(6,7,c2) //debería bloquearse acá.....
	time.Sleep(time.Second*2) //espero 2 segundos.
	x = <- c2 //saco un dato del channel, con lo cual el último gorutine debería poder poner datos
	fmt.Println("fin buffer")
	// :D excelente, funcionó perfecto!	
	
	//un channel se puede cerrar, lo tiene que hacer el sender, al cerrarlo se le dice que no se van a recibir más datos.
	//esto sirve porque se puede hacer un loop sobre un channel y el loop se queda esperando pero si se cierra el channel termina.
	c3 := make(chan int,100)
	go meteDatos(c3)
	for i:= range c3 {
		fmt.Println("Datos: ",i)
	}
	fmt.Println("fin cerrando channel")
	//anduvo perfecto!!! cuando se cierra el channel termina el range.
	
}

func suma(a int, b int, c chan int){
	result := a+b
	time.Sleep(time.Second*1)
	c <- result
}
func suma2(a int, b int, c chan int){
	result := a+b
	time.Sleep(time.Second*2)
	c <- result
}
func suma3(a int, b int, c chan int){
	result := a+b
	c <- result
	fmt.Println("Poniendo 3")	
}
func suma4(a int, b int, c chan int){
	result := a+b
	c <- result
	fmt.Println("Poniendo 4")	
}

func meteDatos(c chan int){
	for i:=0;i<5;i++{
		c <- i
		time.Sleep(time.Second*1)
	}
	close(c) //acá se cierra el canal, entonces el range termina.
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

