package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	for i := 0; i < 50; i++ {
		go star()
	}
	star()
}

func star() {

	row := 0
	col := 0
	color := 0

	for {
		fmt.Printf("\033[" + strconv.Itoa(row) + ";" + strconv.Itoa(col) + "H ")
		row = rand.Intn(24)
		col = rand.Intn(80)
		color = rand.Intn(8)
		fmt.Printf("\033[" + strconv.Itoa(row) + ";" + strconv.Itoa(col) + "H\x1b[3" + strconv.Itoa(color) + ";1m" + "*")
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}

}
