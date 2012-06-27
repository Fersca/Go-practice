// Copyright 2012 Fernando Scasserra twitter: @fersca. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"container/list"
	"net/http"
	"runtime"
)

//Create the list to support the LRU List
var lista list.List

//Create the map to store the key-elements
var mapa map[string]*list.Element

//Channes to sync the List and map modifications
var lisChan chan int
var mapChan chan int

//Print information
const enablePrint bool = false 

/*
 * Init the system variables
 */
func init(){

	//Welcome Message
	fmt.Println("Starting Fercacher HTTP Key-Value Server")

	//Set the thread quantity based on the number of CPU's
	coreNum := runtime.NumCPU()
	fmt.Println("Core numbers: ",coreNum)
	runtime.GOMAXPROCS(coreNum)

	//Create a new doble-linked list to act as LRU
	lista  = *list.New()

	//Create a new Map to search for elements
	mapa = make(map[string]*list.Element)

	//Create the channels
	lisChan = make(chan int,1)
	mapChan = make(chan int,1)

	fmt.Println("Ready.")
}

/*
 * Create the server
 */
func main() {
	
	//Create the webserver
	http.Handle("/", http.HandlerFunc(processRequest))
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		fmt.Printf("Fercacher ListenAndServe Error",err)
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
	//PrintInformation
	printRequest(req)

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
			if enablePrint {fmt.Println("Not Supported: ", req.Method)}
			 //Method Not Allowed
			w.WriteHeader(405)
	}

}

/*
 * Print the request information 
 */
func printRequest(req *http.Request){

	//Print request information
	if enablePrint {
		fmt.Println("-------------------")
		fmt.Println("Method: ",req.Method)
		fmt.Println("URL: ",req.URL)
		fmt.Println("Headers: ",req.Header)
	}
}

/*
 * Save the sent value to the map and the 
 */
func createElement(clave string, valor string){

	//Add the value to the list and get a pointer to the node	
	lisChan <- 1 
	elemento := lista.PushFront(valor)
	<- lisChan	

	mapChan <- 1 
	//Save the node in the map
	mapa[clave] = elemento
	<- mapChan

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

	//Move the element to the front of the LRU-List using a goru
	go func(){
		//Move the element
		lisChan <- 1 
		lista.MoveToFront(elemento)
		<- lisChan
		if enablePrint {fmt.Println("LRU Updated")}
	}()

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

	//Remove the element from cache and list in a separated gorutine
	go func(){
		
		//Delete the element in the LRU List 
		lisChan <- 1 
		lista.Remove(elemento)
		<- lisChan 

		//Delete the key in the map
		mapChan <- 1
		delete(mapa, clave)
		<- mapChan	

		//Print message
		if enablePrint {fmt.Println("Delete successfull: ",clave)}
	}()

	return true

}


