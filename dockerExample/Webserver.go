package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/bradfitz/gomemcache/memcache"
)

var mc *memcache.Client

func handler(w http.ResponseWriter, r *http.Request) {

	it, err := mc.Get("visita")

	quantity := 0

	if err == nil && it != nil {
		xs := string(it.Value)
		x, _ := strconv.Atoi(xs)
		quantity = x + 1
	} else {
		fmt.Println("Error", err)
	}

	quantity++
	fmt.Fprintf(w, "Visita %s!", strconv.Itoa(quantity))
	fmt.Println("Visita: ", strconv.Itoa(quantity))

	mc.Set(&memcache.Item{Key: "visita", Value: []byte(strconv.Itoa(quantity))})

}

func main() {
	fmt.Println("Iniciando Webserver....")

	//Esta IP es la IP del container, también se puede poner la de la máquina virtual y funciona 
	//(que es más práctico porque es siempre la misma para todos los containers)
	mc = memcache.New("memcached1:11211")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
