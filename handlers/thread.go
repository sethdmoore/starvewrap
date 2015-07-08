package handlers

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	//"regexp"
	"github.com/sethdmoore/starvewrap/signals"
	"os/exec"
	"time"
)

func check(err error) {
	if err != nil {
		log.Fatal("err: %v", err)
	}
}

func Start(prefix string, dir string, bin string, sig chan int) {
	//main := make(chan int)

	all := make(chan int)

	safeToUpgrade := false
	token := " <:_:> "
	player_count := 0
	os.Chdir(dir + "/bin")

	fmt.Println("%s/%s", dir, bin)
	cmd := exec.Command("./"+bin, "-console")

	stdout, err := cmd.StdoutPipe()
	check(err)
	stderr, err := cmd.StderrPipe()
	check(err)
	stdin, err := cmd.StdinPipe()
	check(err)

	err = cmd.Start()
	check(err)

	fmt.Printf("%s: %s\n", prefix, player_count)
	go func() {
		scanner := bufio.NewScanner(stdout)

		// label for exiting
	loop:
		for scanner.Scan() {
			/*
				if (scanner.Text() == "ConsoleInput: \"c_listallplayers()\"") {
					// Detect player number
				}
			*/

			select {

			case sig := <-all:
				if sig == signals.SIGINT {
					break loop
				}
			default:
			}

			if strings.Contains(scanner.Text(), "ConsoleInput: ") {
				continue
			}
			fmt.Printf("%s: %s\n", prefix, scanner.Text())
		}
		//r.ReadString("\n")
		//l, err := stdout.Reader.ReadBytes("\n")
		//io.Copy(os.Stdout, stdout)
		fmt.Println("Never get here?")
		return
	}()

	go func() {
		io.Copy(os.Stderr, stderr)
		fmt.Println("TODO: figure out how to break this")
		return
	}()

	/*
		Playerlist poll
	*/
	go func() {

	loop:
		for {

			WritePlayerList(stdin, token)
			if GetNumPlayers(dir+"/data") == 0 {
				safeToUpgrade = true
			}

			// handle shutdowns
			select {
			case sig := <-all:
				fmt.Println("%v", sig)
				stdin.Write([]byte("c_shutdown(true)\n"))
				break loop
				/*
					if sig == 0 {
						break loop
					}
				*/
			default:
				// no-op
			}

			time.Sleep(3 * time.Second)
		}
		return
	}()

	//<-main
	err = cmd.Wait()

	if err != nil {
		fmt.Println("CLEAN UP")
	}

	fmt.Println("exited %s\n\n", prefix)
}
