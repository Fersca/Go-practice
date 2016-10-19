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
        caracter := ""
        for i, c := range linea {
            s := fmt.Sprintf("%c", c)
            caracter = s
            if i==2 {
                break
            }
        }        
        fmt.Println(caracter, " - ", pos)
        contador++
        if contador==100 {
            return
        }
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}