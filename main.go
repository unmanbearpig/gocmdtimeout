package main

import (
	"context"
	"os/exec"
	"log"
	"fmt"
	"syscall"
	"time"
)

// https://stackoverflow.com/questions/67750520/golang-context-withtimeout-doesnt-work-with-exec-commandcontext-su-c-command



// The timeout only applies to the process started by exec, it won't kill any child processes. In your case it will kill the su but not the next python3 process.
//
// To kill all children started by a given process you can start it in a new process group and kill the entire group by sending SIGKILL to -pid (negative pid), like so:
//
// ctx, cancel := context.WithTimeout(context.Background(), 1000 * time.Millisecond)
// process := exec.CommandContext(ctx, "su", "-", "myuser", "-c", "python3 main.py")
//
// process.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
// go func() {
//     <-ctx.Done()
//     if ctx.Err() == context.DeadlineExceeded {
//         syscall.Kill(-process.Process.Pid, syscall.SIGKILL)
//     }
// }()
//
// processOutBytes, err := process.Output()
// cancel()
//
// if ctx.Err() == context.DeadlineExceeded {
//     fmt.Println("Timeout")
// }
//
// Note also that the code relying on syscall isn't portable; it won't even compile on Windows, for example.
//

func testCase() {
	ctx, cancel := context.WithTimeout(context.Background(), 1000 * time.Millisecond)
	defer cancel()

	process := exec.CommandContext(ctx, "./infloop_runner")

	type comp struct {
		output []byte
		err    error
	}
	done := make(chan comp)

	go func() {
		out, err := process.Output()
		done <- comp{out, err}
	}()

	select {
	case comp := <- done:
		log.Printf("got comp %v", comp)
	case _ = <- ctx.Done():
		log.Printf("got done")
	}

	// log.Printf("done. Err = %v", err)
	if ctx.Err() == context.DeadlineExceeded {
		log.Println("Timeout")
	}
}


func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1000 * time.Millisecond)
	defer cancel()

	process := exec.CommandContext(ctx, "./infloop_runner")
	process.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	go func() {
		<-ctx.Done()
		if ctx.Err() == context.DeadlineExceeded {
			syscall.Kill(-process.Process.Pid, syscall.SIGKILL)
		}
	}()

	_, err := process.Output()
	log.Println("err: ", err)
	cancel()

	if ctx.Err() == context.DeadlineExceeded {
		fmt.Println("Timeout")
	}
}
