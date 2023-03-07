package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"

	fsn "github.com/fsnotify/fsnotify"
)

func run(c *config) error {
	w, err := fsn.NewWatcher()
	if err != nil {
		return fmt.Errorf("watcher: %w", err)
	}

	// todo: for a first start, "bind address already in use"
	// may indicate that there is other process running; react to that

	for verb := "started"; ; {
		cmd := exec.Command(c.path, c.args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		// ^^ todo: wrap Stderr so that the "child: " prefix is shown

		err = cmd.Start()
		if err != nil {
			return err
		}
		t0 := time.Now()

		w.Add(c.path)
		log.Printf("child %s, pid: %d", verb, cmd.Process.Pid)

		done := make(chan struct{})
		go func() {
			cmd.Wait()
			if time.Since(t0).Seconds() < 2 {
				log.Printf("child exited early")
			}
			close(done)
		}()
		exited := func() bool {
			select {
			case <-done:
				return true
			default:
				return false
			}
		}

		select {
		case ev, ok := <-w.Events:
			if !ok {
				return fmt.Errorf("watch: event chan closed")
			}
			switch {
			case ev.Has(fsn.Create) || ev.Has(fsn.Write):
				if exited() {
					continue
				}
				goto restart
			default:
				log.Printf("WARN: watch op %s", ev.Op)
			}

		case err, ok := <-w.Errors:
			if !ok {
				return fmt.Errorf("watch: error chan closed")
			}
			return fmt.Errorf("watch: received error: %w", err)
		}
	restart:
		err = cmd.Process.Signal(os.Signal(syscall.SIGTERM))
		if err != nil {
			// signal was not supported
			cmd.Process.Kill()
		}
		<-done
		verb = "restarted"
		time.Sleep(200 * time.Millisecond) // better to wait for freed port
	}
}
