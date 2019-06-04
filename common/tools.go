package common

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func ExecCommand(command string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}

func ExecLiveCommand(command string) {
	cmd := exec.Command(command)
	cmd.Stdin = os.Stdin
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	go func() {
		for {
			l, err := out.ReadString('\n')
			if err != nil && err.Error() != "EOF" {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			fmt.Print(l)
			time.Sleep(100 * time.Millisecond)
		}
	}()
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
