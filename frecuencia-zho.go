package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
)

func main() {
    file, err := os.Open("/Users/fscasserra/code/go/src/github.com/Fersca/Go-practice/caracteres.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    contador:=1
    for scanner.Scan() {
        linea := scanner.Text()
        pos := contador
        caracter:= linea[46]
        fmt.Println(caracter, " - ", pos)
        contador++
        if contador==10 {
            return
        }
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}