package handlers

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	//"regexp"
	"os/exec"
	"time"
)

func check(err error) {
	if err != nil {
		log.Fatal("foo %s", err)
	}
}

func Start(prefix string, dir string, bin string) {
	//safeToUpgrade := false
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
		for scanner.Scan() {
			/*
				if (scanner.Text() == "ConsoleInput: \"c_listallplayers()\"") {
					// Detect player number
				}
			*/
			if strings.Contains(scanner.Text(), "ConsoleInput: ") {
				continue
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
			WritePlayerList(stdin, token)
			GetNumPlayers(dir + "/data")
			time.Sleep(3 * time.Second)
		}
	}()

	err = cmd.Wait()
	if err != nil {
		fmt.Println("CLEAN UP")
	}
	fmt.Println("exited %s\n\n", prefix)
}
