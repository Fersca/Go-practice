/* Copyright 2012 Fernando Scasserra twitter: @fersca. All rights reserved.
 *
 * Fercacher is a HTTP cache system that performs in constant time.
 * It keeps a MAP to store the object internally, and a Double Linked list to purge the LRU elements.
 * 
 * To Store element, do a PUT/POST call to: /element_ID and the body of the message will be stored.
 * To get an element do a GET /element_ID, you will receive the previous stored message.
 * To delete a key do a DELETE /element_ID.
 *
 * LRU updates are done in backgrounds gorutines.
 * LRU and MAP modifications are performed through channels in order to keep them synchronized.
 * Bytes stored are counted in order to limit the amount of memory used by the application.
 */

package main

import (
	"fmt"
	"container/list"
	"net/http"
	"runtime"
	"encoding/json"
)

//Create the list to support the LRU List
var lista list.List

//Create the map to store the key-elements
var mapa map[string]*list.Element

//Max byte in memory (Key + Data), today set to 100KB
const maxMemBytes int64 = 1048576
var memBytes int64 = 0

//Channes to sync the List, map
var lisChan chan int
var mapChan chan int

//Print information
const enablePrint bool = false

//Struct to hold the value and the key in the LRU
type node struct {
	K string
	V map[string]interface{}
}

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

			if req.URL.Path[1:]=="search" {
				key := req.FormValue("key")
				value := req.FormValue("value")
				result := search(key, value)
				b, err := json.Marshal(result)				
				if err!=nil {				
					fmt.Println(b)
				}
				w.Write([]byte(b))
				return
			} 
			if req.URL.Path[1:]=="elements" {
				b, err := json.Marshal(len(mapa))				
				if err!=nil {				
					fmt.Println(b)
				}
				w.Write([]byte(b))
				return
			} 

			//Get the vale from the cache
			element := getElement(req.URL.Path[1:])
		
			if element==nil {
				//Return a not-found				
				w.WriteHeader(404)
			} else {
				//Write the response to the client
				b, err := json.Marshal(element.Value.(node).V)
				if err!=nil {
					if enablePrint {fmt.Println("Error geting Key: ",req.URL.Path[1:],err)}
					w.WriteHeader(404)
				} else {
					w.Write([]byte(b))
				}			
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
 * Search the jsons that has the key with the specified value
 */
func search(key string, value string) []interface{} {

	arr := make([]interface{},0)	
		
	//Search the Map for the value
	for _, v := range mapa {
		//TODO: This is absolutely un-efficient, we are creating a new array for each iteration. Fix this.
		//Is this possible to have something like java ArrayLists  ?
		nod := v.Value.(node)
		if nod.V[key]==value {
			arr = append(arr,nod)
		}
	}

	return arr
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

	b := []byte(valor)
	var f interface{}
	err := json.Unmarshal(b, &f)

	if err != nil {
		fmt.Println("Error:", err)
		return
	} 
	
	m := f.(map[string]interface{})

	var elemento *list.Element
	elemento = mapa[clave]

	if elemento == nil {

		//Add the value to the list and get a pointer to the node	
		n := node{clave, m}
	
		lisChan <- 1 
		elemento = lista.PushFront(n)
		<- lisChan	

		//Save the node in the map
		mapChan <- 1 
		mapa[clave] = elemento
		<- mapChan

		//Increase the memory counter in a diffetet gorutine
		go func(){
			//Increments the memory counter (Key + Value in LRU, + Key in MAP)
			memBytes += int64((len(clave)*2)+len(m))

			if enablePrint {fmt.Println("Inc Bytes: ",memBytes)}

			//Purge de LRU
			purgeLRU()
		}()
		
	} else {		

		//Store the previous bytes
		var prevBytes int = len(elemento.Value.(node).V)

		//Update the element, creating a new node (I dont't know how to update only the node value)
		elemento.Value = node{clave, m}

		//Move the element to the front of the LRU
		go moveFront(elemento)

		//Remove the element from the list in a separated gorutine
		go func(){
			//Update the Bytes counter
			memBytes = memBytes - int64(prevBytes) + int64(len(m))

			if enablePrint {fmt.Println("Upd Bytes: ",memBytes)}

			//Purge LRU
			purgeLRU()
		}()		
	
	}

}

/*
 * Purge the LRU List deleting the last element
 */
func purgeLRU(){

	//Checks the memory limit and decrease it if it's necessary
	for memBytes>maxMemBytes {

		//Get the last element to remove it. Sync is not needed because nothing 
		//happens if the element is moved in the middle of this rutine, at last it will be removed
		lastElement := lista.Back()
		
		//Delete the element from the map
		var key string = lastElement.Value.(node).K		
		var removeBytes int = len(lastElement.Value.(node).V)+(len(key)*2)
		
		mapChan <- 1
		delete(mapa, key)		
		<- mapChan
	
		//Delete the element from the LRU
		lisChan <- 1 
		lista.Remove(lastElement)		
		<- lisChan

		memBytes -= int64(removeBytes)

		if enablePrint {fmt.Println("Purge Done: ",memBytes)}
	}

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
	go moveFront(elemento)

	//Return the element
	return elemento

}

/*
 * Move the element to the front of the LRU, because it was readed or updated
 */
func moveFront(elemento *list.Element){
	//Move the element
	lisChan <- 1 
	lista.MoveToFront(elemento)
	<- lisChan
	if enablePrint {fmt.Println("LRU Updated")}
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

	//Delete the key in the map
	mapChan <- 1
	delete(mapa, clave)
	<- mapChan	

	//Remove the element from the list in a separated gorutine
	go func(){
		
		//Delete the element in the LRU List 
		lisChan <- 1 
		lista.Remove(elemento)
		<- lisChan 

		//Decrement the byte counter, decrease the Key * 2 + Value
		var n node = elemento.Value.(node)
		memBytes -= int64((len(n.K)*2)+len(n.V))

		if enablePrint {fmt.Println("Dec Bytes: ",memBytes)}

		//Print message
		if enablePrint {fmt.Println("Delete successfull, ID: ",clave)}
	}()

	return true

}


