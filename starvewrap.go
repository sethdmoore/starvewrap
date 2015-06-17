package main

import (
	"fmt"
	"github.com/sethdmoore/starvewrap/handlers"
)

/*
func check(err error) {
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}
*/

func main() {
	dir := "/home/steam/steamapps/DST/bin"
	//servers := []string{"1", "2", "3"}
	servers := []string{"1"}
	messages := make(chan int)

	fmt.Println("hello")
	for idx, num := range servers {
		go func(i string) {
			handlers.Start(i, dir)
			messages <- idx + 1
		}(num)
	}

	for i := 0; i < len(servers); i++ {
		fmt.Printf("%v\n", messages)
		<- messages
	}
	fmt.Printf("exited")
}
