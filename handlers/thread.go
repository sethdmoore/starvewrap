package handlers

import (
	"bufio"
	"fmt"
	"github.com/sethdmoore/starvewrap/commands"
	"github.com/sethdmoore/starvewrap/globals"
	"github.com/sethdmoore/starvewrap/signals"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

func check(err error) {
	if err != nil {
		log.Fatal("err: %v", err)
	}
}

func Start(prefix string, dir string, bin string, mainsig chan int) {
	safeToUpgrade := false
	token := "\x00"
	player_count := 0
	os.Chdir(dir + "/bin")

	fmt.Println("%s/%s", dir, bin)

	cmd := exec.Command("./"+bin, "-console")

	// ensure the child process is in its own process group
	// TODO: investigate if this works on Windows
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	stdout, err := cmd.StdoutPipe()
	check(err)

	stderr, err := cmd.StderrPipe()
	check(err)

	stdin_lock := false
	stdin_chan := make(chan bool)
	stdin, err := cmd.StdinPipe()
	check(err)

	err = cmd.Start()
	check(err)

	fmt.Printf("%s: %s\n", prefix, player_count)

	// stdout
	go func() {
		scanner := bufio.NewScanner(stdout)

		for scanner.Scan() {
			/*
				if (scanner.Text() == "ConsoleInput: \"c_listallplayers()\"") {
					// Detect player number
				}
			*/

			/*
				select {
				case outsig := <-stdout_sig:
					if outsig == signals.SIGINT {
						fmt.Printf("STDOUT GOT THE SIGNAL")
						break loop
					}
				default:
			*/
			// filter app spam from the logs
			if strings.Contains(scanner.Text(), commands.INPUT_TAG) {
				//continue
				// don't write to stdin before the server is up
				/*
					close will emit "zero values" from this chan
					'false' in this case, therefore disabling the lock
				*/
			} else if scanner.Text() == globals.INIT_SUCCESS {
				if stdin_lock == false {
					fmt.Println("YO")
					stdin_chan <- true
				}
			}
			fmt.Printf("%s: %s\n", prefix, scanner.Text())
			//}

		}
		stdout.Close()
		fmt.Printf("%s: STDOUT closed\n", prefix)
		return
	}()

	// stderr
	go func() {
		io.Copy(os.Stderr, stderr)
		fmt.Println(prefix + ": STDERR closed")
		stderr.Close()
		return
	}()

	/*
		Playerlist poll
	*/
	go func() {

	loop:
		for {
			// block input until server is up

		inner:
			for stdin_lock == false {
				// probably don't need a select if we handle
				// the timeout somewhere else
				select {
				case stdin_lock = <-stdin_chan:
					if stdin_lock {
						break inner
					} else {
						// signal failure, kill child
					}
				default:
					fmt.Println("blocked")
					time.Sleep(1 * time.Second)
					// implement timeout
				}
			}
			fmt.Println("I AM HERE")

			// handle shutdowns
			select {
			case insig := <-mainsig:
				switch {
				case insig == signals.SIGINT:
					commands.Exec(stdin, commands.SAVE_SHUTDOWN)
					break loop
				default:
					fmt.Println("Some other signal")
				}
			default:
				WritePlayerList(stdin, token)
				if GetNumPlayers(dir+"/data") == 0 {
					safeToUpgrade = true
					//fmt.Println("Safe to Upgrade")
				} else {
					safeToUpgrade = false
				}
			}

			time.Sleep(3 * time.Second)
		}

		// probably not necessary to close this...
		// I'm pretty sure it receives EOF from the cmd
		stdin.Close()
		return
	}()

	//<-main2

	// wait for server to exit
	err = cmd.Wait()
	if err != nil {
		fmt.Println("nonzero")
		fmt.Println("CLEAN UP")
	}

	fmt.Printf("%s: exited\n", prefix)
}
