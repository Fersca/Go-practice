package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func Test_APICall(t *testing.T) {

	go cargarDatos()
	go crearWebserver()
	time.Sleep(300 * time.Millisecond)

	resp, err := http.Get("http://localhost:8080?id=MEX")

	if err != nil {
		t.Fatalf("Erro obteniendo recurso")
	}

	data, err2 := ioutil.ReadAll(resp.Body)

	if err2 != nil {
		t.Fatalf("Error leyendo el body")
	}

	var f interface{}

	err3 := json.Unmarshal(data, &f)

	if err3 != nil {
		t.Fatalf("error en unmarshalear")
	}

	p := f.(map[string]interface{})

	if p["nombre"] != "Mexico" {
		t.Fatalf("Erro obteniendo mexico")
	}

}
