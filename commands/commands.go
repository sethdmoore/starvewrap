package commands

import (
	"io"
)

const (
	SAVE          string = "c_save()"
	SAVE_SHUTDOWN string = "c_shutdown(true)"
	SHUTDOWN      string = "c_shutdown(false)"
	INPUT_TAG     string = "--STARVEWRAP"
)

func command(cmd string) []byte {
	exec := cmd + " " + INPUT_TAG + "\n"
	return []byte(exec)
}

func Exec(writer io.WriteCloser, cmd string) {
	writer.Write(command(cmd))
}
