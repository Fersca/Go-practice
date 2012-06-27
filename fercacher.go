package main

import (
	"fmt"
	"container/list"
	"net/http"
)

//Create the list to support the LRU List
var lista list.List

//Create the map to store the key-elements
var mapa map[string]*list.Element

/*
 * Init the system variables
 */
func init(){
	//Create a new doble-linked list to act as LRU
	lista  = *list.New()

	//Create a new Map to search for elements
	mapa = make(map[string]*list.Element)
}

/*
 * Create the server
 */
func main() {
	
	//Process the http commands
	fmt.Println("Starting Fercacher HTTP Key-Value Server ... ")

	//Create the webserver
	http.Handle("/", http.HandlerFunc(processRequest))
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		fmt.Printf("ListenAndServe Error",err)
	}
}

/*
 * Process the commands recived from internet
 */
func processRequest(w http.ResponseWriter, req *http.Request){
	//Get the headers map	
	headerMap := w.Header()
	//Add the new headers
	headerMap.Add("System","Fercacher 1.0")

	//Print request information
	fmt.Println("-------------------")
	fmt.Println("Method: ",req.Method)
	fmt.Println("URL: ",req.URL)
	fmt.Println("Headers: ",req.Header)
	
	//Performs action based on the request Method
	switch req.Method {

		case "GET":
			//Get the vale from the cache
			value := getElement(req.URL.Path[1:])
			
			if value==nil {
				//Return a not-found				
				w.WriteHeader(404)
			} else {
				//Write the response to the client
				w.Write([]byte(value.Value.(string)))				
			}

		case "PUT":
			fallthrough
		case "POST":
			//Create the array to hold the body
			var p []byte = make([]byte,req.ContentLength)
			
			//Reads the body content 
			req.Body.Read(p)
			
			//Save the element in the cache
			createElement(req.URL.Path[1:], string(p))
			
			//Response the 201 - created to the client
			w.WriteHeader(201)

		case "DELETE":
			//Get the vale from the cache
			result := deleteElement(req.URL.Path[1:])
			
			if result==false {
				//Return a not-found				
				w.WriteHeader(404)
			} else {
				//Return a Ok
				w.WriteHeader(200)
			}

		default:
			fmt.Println("Not Supported: ", req.Method)
			 //Method Not Allowed
			w.WriteHeader(405)
	}

}

/*
 * Save the sent value to the map and the 
 */
func createElement(clave string, valor string){

	//Add the value to the list and get a pointer to the node	
	elemento := lista.PushFront(valor)

	//Save the node in the map
	mapa[clave] = elemento

}

/*
 * Get the element from the Map and push the element to the first position of the LRU-List 
*/ 
func getElement(clave string) *list.Element {

	//Get the element from the map
	elemento := mapa[clave]	

	//checks if the element exists in the cache
	if elemento==nil {
		return nil
	} 

	//Move the element to the front of the LRU-List
	lista.MoveToFront(elemento)

	//Return the element
	return elemento

}

/*
 * Delete the key in the cache
 */
func deleteElement(clave string) bool {

	//Get the element from the map
	elemento := mapa[clave]	

	//checks if the element exists in the cache
	if elemento==nil {
		return false
	} 

	//Delete the element in the LRU List
	lista.Remove(elemento)

	//Delete the key in the map
	delete(mapa, clave)

	return true

}


