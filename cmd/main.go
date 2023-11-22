package main

import (
	"lab1/iternal/api"
	
	"log"
)

func main() {
	log.Println("Application start")
	api.StartServer()
	log.Println("Application terminated")
}