/* 
 * Mango Cache.
 *
 * 2013 - Fernando Scasserra - twitter: @fersca.
 *
 * Mango is a persistance cache system written in golang that performs in constant time.
 * It keeps a MAP to store the object internally, and a Double Linked list to purge the LRU elements.
 * 
 * LRU updates are done in backgrounds gorutines.
 * LRU and MAP modifications are performed through channels in order to keep them synchronized.
 * Bytes stored are counted in order to limit the amount of memory used by the application.
 *
 */

package main

import (
	"fmt"
	"container/list"
	"net"
	"net/http"
	"runtime"
	"encoding/json"
	"strings"
	"strconv"
	"io/ioutil"
	"os"
	//"time"
)

//Create the list to support the LRU List
var lista list.List

//Max byte in memory (Key + Data), today set to 100KB
const maxMemBytes int64 = 50//1048576
var memBytes int64 = 0
var sequence int = 0
const pointerLen int = 4+8 //Bytes of pointer in 32bits machines plus int64 for the key of element in hashmemBytes

//Channes to sync the List, map
var lisChan chan int

//chennel to acces to the collection map
var collectionChan chan int

//Print information
const enablePrint bool = false

//Struct to hold the value and the key in the LRU
type node struct {
	V map[string]interface{}
	Swap bool
	col string
	key string
}
//Struct to hold the value and the key in the LRU
type searchNode struct {
	Id string
	Document map[string]interface{}
}

//Holds the relation between the diferent collections of element with the corresponding channel to write it
type collectionChannel struct {
	Mapa map[string]*list.Element
	Canal chan int
}

//Create the map that stores the list of collections
var collections map[string]collectionChannel

/*
 * Init the system variables
 */
func init(){

	//Welcome Message
	fmt.Println("Starting Mango...")

	//Set the thread quantity based on the number of CPU's
	coreNum := runtime.NumCPU()
	fmt.Println("Core numbers: ",coreNum)
	fmt.Println("Max memory defined as: ",maxMemBytes," bytes")
	runtime.GOMAXPROCS(coreNum)

	//Create a new doble-linked list to act as LRU
	lista  = *list.New()

	//Create the channels
	lisChan = make(chan int,1)
	collectionChan = make(chan int,1)

	collections = make(map[string]collectionChannel)

	fmt.Println("Ready.")
}

/*
 * Create the server
 */
func main() {

	//Start the console
	go console()

	//Create the webserver
	http.Handle("/", http.HandlerFunc(processRequest))
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		fmt.Printf("Mango ListenAndServe Error",err)
	}

}

/*
 * Start the command console
 */
func console(){

	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			continue
		}
		go handleTCPConnection(conn)
	}

}

/*
 * Process each HTTP connection
 */
func handleTCPConnection(conn net.Conn){

	fmt.Println("Connection stablished")

	//Create the array to hold the command
	var command []byte = make([]byte,100)

	for {
		//Read from connection waiting for a command
		cant, err := conn.Read(command)
		if err == nil {

			//read the command and create the string
			var commandStr string = string(command)

			//Exit the connection
			if commandStr[0:4] == "exit" {
				fmt.Println("Cerrando Conexion")
				conn.Close()
				return
			}

			//Get the element
			if commandStr[0:3] == "get" {

				comandos := strings.Split(commandStr[:cant-2]," ")

				fmt.Println("Collection: ",comandos[1], " - ",len(comandos[1]))
				fmt.Println("Id: ",comandos[2]," - ",len(comandos[2]))

				//b,err := getElement(comandos[1],atoi(comandos[2]))
				b,err := getElement(comandos[1],comandos[2])

				if b!=nil {
					conn.Write(b)
					conn.Write([]byte("\n"))
				} else {
					if err==nil{
						conn.Write([]byte("Key not found\n"))
					} else {
						fmt.Println("Error: ", err)
					}
				}
				continue
			}

			//Get the total quantity of elements
			if commandStr[0:8] == "elements" {

				comandos := strings.Split(commandStr[:cant-2]," ")

				fmt.Println("Collection: ",comandos[1], " - ",len(comandos[1]))

				b, err := getElements(comandos[1])
				if err==nil {
					conn.Write(b)
					conn.Write([]byte("\n"))
				} else {
					fmt.Println("Error: ", err)
				}
				continue
			}

			//return the bytes used
			if commandStr[0:6] == "memory" {

				//result := "Uses: "+strconv.FormatInt(memBytes,10)+"bytes, "+ strconv.FormatInt((memBytes/(maxMemBytes/100)),10)+"%\n"
				result := "Uses: "+strconv.FormatInt(memBytes,10)+"bytes\n"
				conn.Write([]byte(result))

				continue
			}


			//POST elements
			if commandStr[0:4] == "post" {

				comandos := strings.Split(commandStr[:cant-2]," ")

				fmt.Println("Collection: ",comandos[1], " - ",len(comandos[1]))	
				fmt.Println("Key: ",comandos[2], " - ",len(comandos[2]))
				fmt.Println("JSON: ",comandos[3]," - ",len(comandos[3]))

				err := createElement(comandos[1],comandos[2],comandos[3])

				var result string
				if err!=nil{
					fmt.Println(err)
				} else {
					//result = "Element Created: "+strconv.Itoa(id)+"\n"
					result = "Element Created: "+comandos[2]+"\n"
					conn.Write([]byte(result))
				}

				continue
			}

			if commandStr[0:6] == "delete" {

				comandos := strings.Split(commandStr[:cant-2]," ")

				//Get the vale from the cache
				//result := deleteElement(comandos[1],atoi(comandos[2]))
				result := deleteElement(comandos[1],comandos[2])

				if result==false {
					//Return a not-found				
					conn.Write([]byte("Key not found"))
				} else {
					//Return a Ok
					response := "Key: "+comandos[2]+" from: "+comandos[1]+" deleted\n"
					conn.Write([]byte(response))
				}

				continue

			}

			if commandStr[0:6] == "search" {

				comandos := strings.Split(commandStr[:cant-2]," ")

				result, err := search(comandos[1],comandos[2],comandos[3])

				if err!=nil {
					fmt.Println(result)
					conn.Write([]byte("Error searching\n"))
				} else {
					conn.Write([]byte(result))
				}
				continue
			}

			//Exit the connection
			if commandStr[0:4] == "help" {
				result := showHelp()
				conn.Write([]byte(result))
				continue
			}

			//Default Message
			fmt.Println("Comando no definido: ", commandStr)
			conn.Write([]byte("Unknown Command\n"))

		} else {
			fmt.Println("Error: ", err)
		}

	}

}

/*
 * Help
 */
func showHelp() string {

	var help string = "FerCacher Help\n\n"

	help += "Available commands:\n\n"

	help  += "exit 					- Close the connection.\n"
	help  += "get {collection} {id}			- Get the JSON document from the specified collection.\n"
	help  += "elements {collection}			- Get the total elemets from the specified collection.\n"
	help  += "memory 				- Get the total ammount of memory used.\n"
	help  += "post {collection} {json}		- Save a new JSON document in the specified collection.\n"
	help  += "delete {collection} {id}		- Delete the JSON document from the specified collection.\n"
	help  += "search {collection} {key} {value}	- Search in the specified collection for Json documents with keys with the indicated value.\n"
	return help

}

/*
 * Process the commands recived from internet
 */
func processRequest(w http.ResponseWriter, req *http.Request){
	//Get the headers map	
	headerMap := w.Header()
	//Add the new headers
	headerMap.Add("System","Mango 0.1")
	//PrintInformation
	printRequest(req)

	comandos := strings.Split(req.URL.Path[1:],"/")

	//Performs action based on the request Method
	switch req.Method {

		case "GET":

			if comandos[1]=="elements" {
				b, err := getElements(comandos[0])
				if err!=nil {
					fmt.Println(b)
					w.WriteHeader(500)
					return
				}
				w.Write([]byte(b))
				return
			}

			if req.URL.Path[1:]=="search" {
				col := req.FormValue("col")
				key := req.FormValue("key")
				value := req.FormValue("value")
				result, err := search(col,key, value)
				if err!=nil {
					fmt.Println(result)
					w.WriteHeader(500)
					return
				}
				w.Write(result)
				return
			}

			//Get the vale from the cache
			//element, err := getElement(comandos[0],atoi(comandos[1]))
			element, err := getElement(comandos[0],comandos[1])

			if element!=nil {
				//Write the response to the client
				w.Write([]byte(element))
			} else {
				if err==nil {
					//Return a not-found				
					w.WriteHeader(404)
				} else {
					w.WriteHeader(500)
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
			err := createElement(comandos[0],comandos[1],string(p))

			if err!=nil{
				fmt.Println(err)
				w.WriteHeader(500)
			} else {
				//headerMap.Add("element_id",strconv.Itoa(id))
				headerMap.Add("element_id",comandos[1])
				//Response the 201 - created to the client
				w.WriteHeader(201)
			}

		case "DELETE":
			//Get the vale from the cache
			//result := deleteElement(comandos[0],atoi(comandos[1]))
			result := deleteElement(comandos[0],comandos[1])

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
 * Create the element in the collection
 */
func createElement(col string, id string, valor string) (error) {

	//Create the Json element
	b := []byte(valor)
	var f interface{}
	err := json.Unmarshal(b, &f)

	if err != nil {
		return err
	}

	//transform it to a map
	m := f.(map[string]interface{})

	//Add the value to the list and get a pointer to the node	
	n := node{m,false,col,id}

	//create the list element
	var elemento *list.Element
	lisChan <- 1
	elemento = lista.PushFront(n)
	<- lisChan

	//get the collection-channel relation
	cc := collections[col]
	var createDir bool = false

	if cc.Mapa==nil {

		fmt.Println("Creating new collection: ",col)
		//Create the new map and the new channel
		var newMapa map[string]*list.Element
		var newMapChann chan int
		newMapa = make(map[string]*list.Element)
		newMapChann = make(chan int,1)

		newCC := collectionChannel{newMapa, newMapChann}
		newCC.Mapa[id] = elemento

		//The collection doesn't exist, create one
		collectionChan <- 1
		collections[col] = newCC
		<- collectionChan
		createDir = true

	} else {
		fmt.Println("Using collection: ",col)
		//Save the node in the map
		cc.Canal <- 1
		cc.Mapa[id] = elemento
		<- cc.Canal
	}

	//Increase the memory counter in a diffetet gorutine
	go func(){
		//Increments the memory counter (Key + Value in LRU + len of col name, + Key in MAP)
		memBytes += int64(len(b))

		if enablePrint {fmt.Println("Inc Bytes: ",memBytes)}

		//Save the Json to disk
		saveJsonToDisk(createDir, col, id, valor)

		//Purge de LRU
		purgeLRU()
	}()

	return nil
}

func saveJsonToDisk(createDir bool, col string, id string, valor string) {

	if createDir {
		os.Mkdir("data/"+col,0777)
	}

	err := ioutil.WriteFile("data/"+col+"/"+id+".json", []byte(valor), 0644)
	if err!=nil {
		fmt.Println(err)
	}
}

func deleteJsonFromDisk(col string, clave string){
	os.Remove("data/"+col+"/"+clave+".json")
}

func readJsonFromDisK(col string, clave string) []byte {
	fmt.Println("Read from disk: ", col," - ",clave)
	content, err := ioutil.ReadFile("data/"+col+"/"+clave+".json")
	if err!=nil {
		fmt.Println(err)
	}
	fmt.Println("en funcion: ",content)
	return content
}

/*
 * Get the element from the Map and push the element to the first position of the LRU-List 
*/
func getElement(col string, id string) ([]byte, error) {

	cc := collections[col]

	//Get the element from the map
	elemento := cc.Mapa[id]

	//checks if the element exists in the cache
	if elemento==nil {
		return nil, nil
	}

	//Move the element to the front of the LRU-List using a goru
	go moveFront(elemento)

	//Verifica si esta swapeado
	if elemento.Value.(node).Swap==true {

		//Read the swapped json from disk
		b:=readJsonFromDisK(col, id)

		var f interface{}
		err := json.Unmarshal(b, &f)

		if err != nil {
			return nil,err
		}

		m := f.(map[string]interface{})

		//save the map in the node, mark it as un-swapped
		var unswappedNode node
		unswappedNode.V = m
		unswappedNode.Swap = false
		elemento.Value=unswappedNode

		//increase de memory counter
		memBytes += int64(len(b))

		//as we have load content from disk, we have to purge LRU
		go purgeLRU()
	}

	//Return the element
	b, err := json.Marshal(elemento.Value.(node).V)
	return b, err

}

func atoi(value string) int {
	number, _ := strconv.Atoi(value)
	return number
}

/*
 * Get the number of elements
 */
func getElements(col string) ([]byte, error){
	cc := collections[col]
	b, err := json.Marshal(len(cc.Mapa))

	return b, err
}

/*
 * Search the jsons that has the key with the specified value
 */
func search(col string, key string, value string) ([]byte, error) {

	arr := make([]interface{},0)
	cc := collections[col]

	//Search the Map for the value
	for id, v := range cc.Mapa {
		//TODO: This is absolutely inefficient, I'm creating a new array for each iteration. Fix this.
		//Is this possible to have something like java ArrayLists  ?
		nod := v.Value.(node)
		sNode := searchNode{id,nod.V}
		if nod.V[key]==value {
			arr = append(arr,sNode)
		}
	}

	//Create the Json object
	b, err := json.Marshal(arr)

	return b, err

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
 * Purge the LRU List deleting the last element
 */
func purgeLRU(){

	//Checks the memory limit and decrease it if it's necessary
	for memBytes>maxMemBytes {

		fmt.Println("Max memory reached! swapping", memBytes)

		fmt.Println("LRU Elements: ", lista.Len())

		//Get the last element and remove it. Sync is not needed because nothing 
		//happens if the element is moved in the middle of this rutine, at last it will be removed
		lastElement := lista.Back()
		if lastElement==nil {
			fmt.Println("Empty LRU")
			return
		}

		//Remove the element from the LRU
		deleteElementFromLRU(lastElement)

		//Set element as "S"wapped node
		var swappedNode node
		swappedNode.V = nil
		swappedNode.Swap = true
		lastElement.Value=swappedNode
		//it would be better to replace the content of the node instead of create a new one
		//but I cant get it done

		fmt.Println("quedo: ",lastElement.Value.(node).Swap)

		//Print a purge
		if enablePrint {fmt.Println("Purge Done: ",memBytes)}
	}

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
func deleteElement(col string, clave string) bool {

	cc := collections[col]

	//Get the element from the map
	elemento := cc.Mapa[clave]

	//checks if the element exists in the cache
	if elemento!=nil {

		//Delete the key in the map
		cc.Canal <- 1
		delete(cc.Mapa, clave)
		<- cc.Canal

		//Remove the element from the list in a separated gorutine
		go func(){

			deleteElementFromLRU(elemento)

			deleteJsonFromDisk(col, clave)

			//Print message
			if enablePrint {fmt.Println("Delete successfull, ID: ",clave)}
		}()

	} else {
		return false
	}

	return true

}

/*
 * Delete the element from de LRU and decrement the counters
 */
func deleteElementFromLRU(elemento *list.Element){

		//Delete the element in the LRU List 
		lisChan <- 1
		lista.Remove(elemento)
		<- lisChan

		//Decrement the byte counter, decrease the Key * 2 + Value
		var n node = elemento.Value.(node)
		b, _ := json.Marshal(n.V)
		memBytes -= int64(len(b))

		fmt.Println("Borrado: ",b)
		fmt.Println("Dec Bytes: ",len(b))

}

