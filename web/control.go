package web

import (
	"fmt"
	"net"
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
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	getLog().Printf("Building webserver...")
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("failed to build webserver binary: %w", err)
	}
	getLog().Printf("Build succeeded, starting detached process...")
	// Open a log file for the detached process
	logFile, err := os.OpenFile("webserver.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open webserver.log: %w", err)
	}
	//start compiled binary in detached mode
	//script := "./webserver > webserver.log 2>&1 &" //echo $!"
	//script := "./webserver &" //> webserver.log 2>&1 &" //echo $!"
	//cmd := exec.Command("bash", "-c", script)
	cmd := exec.Command("./webserver")
	cmd.Stdout = logFile //logging.WebServerLogWriter
	cmd.Stderr = logFile //logging.WebServerLogWriter
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	//output, err := cmd.Output()
	if err := cmd.Start(); err != nil {
		getLog().Printf("Failed to start webserver detached process: %v", err)
		return fmt.Errorf("failed to start webserver: %w", err)
	}
	getLog().Printf("Detached webserver PID: %d", cmd.Process.Pid)

	//write PID to file
	pidStr := strconv.Itoa(cmd.Process.Pid)
	if err := os.WriteFile(pidFile, []byte(pidStr), 0644); err != nil {
		return fmt.Errorf("failed to write PID file: %w", err)
	}

	// if err := WaitForWebServer(8080, 5*time.Second); err != nil {
	// 	getLog().Printf("Webserver did not start in time: %v", err)
	// 	_ = cmd.Process.Kill() // Kill the process if it didn't start correctly
	// 	return fmt.Errorf("webserver did not start in time: %w", err)
	// }

	// //Zombie processes can occur if the parent process exits before the child
	// go func() {
	// 	_ = cmd.Wait() // Wait for the command to finish
	// 	getLog().Println("Webserver detached process exited - Say NO! to Zombies!")
	// }()

	// Wait in a goroutine to prevent zombies
	go func() {
		err := cmd.Wait()
		if err != nil {
			fmt.Fprintf(logFile, "Webserver process exited with error: %v\n", err)
		} else {
			fmt.Fprintln(logFile, "Webserver process exited cleanly")
		}
		logFile.Close()
	}()

	return nil

}

func WaitForWebServer(port int, timeout time.Duration) error {
	//logging.WebServerLoggerPrintf("Webserver started with PID %s", pidStr)
	//getLog().Printf("Started webserver detached process with PID %d", cmd.Process.Pid)
	//return nil
	// Wait until port 8080 is accepting connections
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	for {
		select {
		case <-ticker.C:

			conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
			if err == nil {
				_ = conn.Close()
				getLog().Printf("Webserver started successfully on :8080")
				return nil
			}
		case <-timer.C:
			return fmt.Errorf("timeout waiting for webserver to listen on :8080")
		}
	}
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
	//if err := process.Signal(syscall.SIGTERM); err != nil {
	if err := syscall.Kill(-pid, syscall.SIGTERM); err != nil {
		return fmt.Errorf("Failed to send interrupt signal to webserver with PID %d: %v", pid, err)
	}
	getLog().Printf("Sent SIGTERM to process %d\n", pid)

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
				getLog().Printf("Process %d has exited", pid)
				cleanup()
				return nil
			}
			if isZombie(pid) {
				getLog().Printf("Process %d is a zombie, cleaning up", pid)
				cleanup()
				return nil
			}
		}
		getLog().Printf("tick")
	}
}

func cleanup() {
	_ = os.Remove(pidFile) // Clean up PID file
	_ = os.Remove("webserver")
}

func isZombie(pid int) bool {
	statusFile := fmt.Sprintf("/proc/%d/status", pid)
	data, err := os.ReadFile(statusFile)
	if err != nil {
		return false
	}
	return strings.Contains(string(data), "State:\tZ")
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

	if process.Signal(syscall.Signal(0)) != nil {
		getLog().Printf("Webserver process with PID %d is not running", pid)
		return false
	}

	if isZombie(pid) {
		getLog().Printf("Webserver process with PID %d is a zombie", pid)
		return false
	}

	// conn, err := net.DialTimeout("tcp", "127.0.0.1:8080", 1*time.Second)
	// if err != nil {
	// 	getLog().Printf("Webserver is not running (PID %d): %v", pid, err)
	// 	return false
	// }
	// _ = conn.Close()
	return true
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
