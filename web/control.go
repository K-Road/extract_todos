package web

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const pidFile = "webserver.pid"

func StartWebServerDetached() error {
	//Compile the binary if it doesn't exist
	buildCmd := exec.Command("go", "build", "-o", "webserver", "./cmd/webserver")
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("failed to build webserver binary: %w", err)
	}

	//start compiled binary in detached mode
	script := "./webserver > webserver.log 2>&1 & echo $!"
	cmd := exec.Command("bash", "-c", script)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to start webserver: %w", err)
	}

	//write PID to file
	pidStr := strings.TrimSpace(string(output))
	if err = os.WriteFile(pidFile, []byte(pidStr), 0644); err != nil {
		return fmt.Errorf("failed to write PID file: %w", err)
	}

	log.Printf("Webserver started with PID %s", pidStr)
	return nil
}

func StopWebServer() error {
	pid, err := readPID()
	if err != nil {
		return err
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("Failed to find webserver process with PID %d: %v", pid, err)
	}

	//send SIGINT
	if err := process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("Failed to send interrupt signal to webserver with PID %d: %v", pid, err)
	}
	log.Printf("Sent SIGTERM to process %d\n", pid)

	//Wait for process to exit
	const maxWait = 5 * time.Second
	timeout := time.After(maxWait)
	tick := time.Tick(200 * time.Millisecond)

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout waiting for process %d to exit", pid)
		case <-tick:
			if err := process.Signal(syscall.Signal(0)); err != nil {
				log.Printf("Process %d has exited", pid)
				_ = os.Remove(pidFile) // Clean up PID file
				_ = os.Remove("webserver")
				return nil
			}
		}
	}
}

func IsWebServerRunning() bool {
	pid, err := readPID()
	if err != nil {
		return false
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	return process.Signal(syscall.Signal(0)) == nil
}

// helper to read PID
func readPID() (int, error) {
	data, err := os.ReadFile(pidFile)
	if err != nil {
		return 0, fmt.Errorf("failed to read PID file: %v", err)
	}

	pidStr := strings.TrimSpace(string(data))
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return 0, fmt.Errorf("invalid PID in file: %v", err)
	}

	return pid, nil
}
