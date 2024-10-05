package main

import (
	"fmt"
	"log"
	"os"
)









func main() {

	// set errors login to errors file
	logFile, err := os.OpenFile("errors.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        log.Fatalf("Could not open log file: %v", err)
    }
  defer logFile.Close()

	log.SetOutput(logFile)

	// recovering any panic in this goroutine 
	defer func() {
    if x := recover(); x != nil {
        // recovering from a panic; x contains whatever was passed to panic()
        log.Printf("run time panic: %v", x)

        // if you just want to log the panic, panic again
        panic(x)
    }
	}()


	s, e := NewScreen()

	if e != nil {
		fmt.Println("new screen error")
		os.Exit(1)
	}
	

	s.Start()	



}
