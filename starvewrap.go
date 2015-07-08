package main

import (
	"fmt"
	"github.com/sethdmoore/starvewrap/handlers"
	"github.com/sethdmoore/starvewrap/signals"
	"os"
	"os/signal"
)

/*
func check(err error) {
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}
*/

func main() {
	dir := "/home/steam/steamapps/DST"
	bin := "dontstarve_dedicated_server_nullrenderer"
	steamcmd := "/home/steam/steamcmd/steamcmd.sh"

	//server_num := []string{"1", "2", "3"}
	main := make(chan os.Signal)

	// number of server threads to spin up
	server_num := []string{"1"}

	// block until all threads complete
	semaphore := make(chan int)

	// allow us to signal any number of servers
	server_sig := make([]chan int, len(server_num))

	// dummy for now, update the base game
	handlers.Update(steamcmd)

	// spin up the server_num
	for idx, num := range server_num {
		if server_sig[idx] == nil {
			server_sig[idx] = make(chan int)
		}
		go func(id string) {
			handlers.Start(id, dir, bin, server_sig[idx])
			semaphore <- idx + 1
		}(num)
	}

	signal.Notify(main, os.Interrupt)
	sig := <-main

	if sig == os.Interrupt {
		for idx, _ := range server_num {
			go func(x int) {
				server_sig[x] <- signals.SIGINT
			}(idx)
		}
		//handlers.Cleanup()
	}

	// semaphore, wait to read from n number of chans
	for i := 0; i < len(server_num); i++ {
		fmt.Printf("%v\n", semaphore)
		<-semaphore
	}

	fmt.Printf("exited")
}
