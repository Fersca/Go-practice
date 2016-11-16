package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

var paises = make(map[string]pais)

func main() {

	fmt.Println("Hola Conf")

	cargarDatos()

	if len(os.Args) < 2 {
		fmt.Println("creando webserver ...")
		crearWebserver()
	} else {
		codigo := os.Args[1]

		p := paises[codigo]

		p.darInfo()

	}

}

func crearWebserver() {

	http.HandleFunc("/", processReq)
	http.ListenAndServe(":8080", nil)

}

func processReq(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	codigo := r.Form.Get("id")

	p := paises[codigo]

	j, _ := json.Marshal(p)

	fmt.Fprintf(w, string(j))

}

func cargarDatos() {

	mexico := pais{"Mexico", 120000000, "Español"}
	brasil := pais{"Brasil", 250000000, "Portugues"}

	paises["MEX"] = mexico
	paises["BRA"] = brasil

}

type pais struct {
	Nombre    string `json:"nombre"`
	Poblacion int    `json:"pobl"`
	Idioma    string `json:"idioma"`
}

func (p *pais) darInfo() {
	fmt.Println("El pais, " + p.Nombre + " tiene " + strconv.Itoa(p.Poblacion) + " habitantes")
}
