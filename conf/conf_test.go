package main

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func Test_API(t *testing.T) {

	go main()

	resp, _ := http.Get("http://localhost:8080/pipi")

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if string(body) != "Hola Amigo pipi" {
		t.Fatalf("Error en respuesta de la API")
	}

}
