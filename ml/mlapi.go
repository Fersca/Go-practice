package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

//Cliente http para relizar las conexiones contra la API
var client *http.Client

func main() {

	//crea un transport para la conxion
	tr := &http.Transport{
		DisableCompression:  false,
		DisableKeepAlives:   false,
		MaxIdleConnsPerHost: 100,
	}

	//crea un cliente http para conectarse contra la api
	client = &http.Client{Transport: tr}

	//Pide una palabra a buscar desde la linea de comandos
	fmt.Print("Buscar Producto: ")
	var producto string
	fmt.Scanln(&producto)
	fmt.Println(producto)

	f, err := os.OpenFile("/Users/fscasserra/code/go/src/github.com/Fersca/Go-practice/ml/data.txt", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	page := 0
	cant := 0
	for ; page < 39; page++ {

		//busca en el listado de mercadolibre los productos
		body, _, _, _ := get("https://api.mercadolibre.com/sites/MLA/search?q=" + producto + "&offset=" + strconv.Itoa(page*50))

		//Create the Json element
		d := json.NewDecoder(strings.NewReader(body))
		d.UseNumber()
		var fi interface{}
		err3 := d.Decode(&fi)

		if err3 != nil {
			panic(err3)
		}

		mapa := fi.(map[string]interface{})

		resultados := mapa["results"].([]interface{})

		//fmt.Println(resultados)

		for _, v := range resultados {

			var km string
			var year string
			var trans string
			var version string

			valor := v.(map[string]interface{})

			if valor["condition"] == "used" {

				attributes := valor["attributes"].([]interface{})

				for _, v1 := range attributes {
					att := v1.(map[string]interface{})
					if att["id"] == "MLA1744-KMTS" {
						km = att["value_name"].(string)
					}
					if att["id"] == "MLA1744-YEAR" {
						year = att["value_name"].(string)
					}
					if att["id"] == "MLA1744-TRANS" {
						trans = att["value_name"].(string)
					}
					if att["id"] == "MLA6628-VERS" {
						version = att["value_name"].(string)
					}

				}

				fmt.Println(valor["title"], ";", valor["price"], ";", km, ";", year, ";", trans, ";", version)

				precio := valor["price"].(json.Number).String()

				//por cada uno arma una linea
				text := valor["title"].(string) + ";" + precio + ";" + km + ";" + year + ";" + trans + ";" + version + "\n"

				//guarda la linea en un archivo
				if _, err = f.WriteString(text); err != nil {
					panic(err)
				}

				cant++

			}
		}
	}
	fmt.Println("cantidad:", cant)
}

func get(urlString string) (string, int, map[string][]string, error) {

	//crea el request a la api
	req, err := http.NewRequest("GET", urlString, nil)

	if err != nil {
		println("Error in http.NewRequest:", err)
	}

	//Setea los headers para llamar a las apis
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Connection", "Keep-Alive")

	//realiza el request mediante el cliente
	resp, err2 := client.Do(req)
	if err2 != nil {
		return "", 0, nil, err2
	}

	code := resp.StatusCode
	//leo el body de la respuesta
	body, err3 := ioutil.ReadAll(resp.Body)
	if err3 != nil {
		return "", code, nil, err2
	}

	//difiere el cierr del body
	defer func() {
		//chequeo si no es nil porque en caso de error al abrirlo falla esta respuesta.
		if resp != nil {
			resp.Body.Close()
		}
	}()

	return string(body), code, resp.Header, nil

}
