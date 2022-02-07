package main

import (
	"log"
	"time"
)

// Main function
func main() {

	delay := time.NewTimer(time.Second * 3)
	log.Println("start")

	select {
	case <-delay.C:
		log.Println("end")

	}

}
