package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("Hello, World!")

	store,err:=NewPostGresStorage();
	if err!=nil{
		log.Fatal(err)
	}

	fmt.Print(store)

	if err:=store.init(); err!=nil{
		log.Fatal(err)
}

	server := NewAPIserver(":8080", store)
	server.Start()
}

