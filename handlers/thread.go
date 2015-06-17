package handlers

import (
	"io"
	"fmt"
	"bufio"
	"log"
	"os"
	//"regexp"
	"os/exec"
	"time"
)

func check(err error) {
	if err != nil {
		log.Fatal("foo %s", err)
	}
}

func countingMode(line string) {
	

}

func Start(prefix string, dir string) {
	player_count := 0
	os.Chdir(dir)
	cmd := exec.Command("./dontstarve_dedicated_server_nullrenderer", "-console")
	
	fmt.Println("got here?")
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
		for scanner.Scan() {
			if (scanner.Text() == "ConsoleInput: \"c_listallplayers()\"") {
				// Detect player number
			}
			fmt.Printf("%s: %s\n", prefix, scanner.Text())
		}
		//r.ReadString("\n")
		//l, err := stdout.Reader.ReadBytes("\n")
		//io.Copy(os.Stdout, stdout)
	}()
	go func() {
		io.Copy(os.Stderr, stderr)
	}()

	go func() {
		for {
			stdin.Write([]byte("c_listallplayers()\n"))
			time.Sleep(3 * time.Second)
		}
	}()
	fmt.Println("yo")
	err = cmd.Wait()
	if err != nil {
		fmt.Println("CLEAN UP")
	}
	fmt.Println("exited %s\n\n", prefix)
}
