package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func main() {

	for {
		fmt.Printf("\033[" + strconv.Itoa(rand.Intn(24)) + ";" + strconv.Itoa(rand.Intn(80)) + "H\x1b[3" + strconv.Itoa(rand.Intn(8)) + ";1m" + "*")
		time.Sleep(20 * time.Millisecond)
	}

}
