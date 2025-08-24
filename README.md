# extract_todos

## 🧩 Webserver Lifecycle (Detached Mode)

The webserver is managed as a detached subprocess for independent execution and lifecycle control.

### ✅ Start Webserver (Detached)

The `StartWebServerDetached()` function compiles and starts the `./cmd/webserver` binary in detached mode:

```go
cmd := exec.Command("./webserver")
cmd.Stdout = logging.WebServerLogWriter
cmd.Stderr = logging.WebServerLogWriter
cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

if err := cmd.Start(); err != nil {
    return err
}
go func() {
    _ = cmd.Wait() // prevent zombie process
}()

PID is saved to webserver.pid

Logs go to a centralized WebServerLogger

Zombie processes are prevented via cmd.Wait() in a goroutine

🛑 Stop Webserver
To stop the webserver:

Read the PID from webserver.pid

Send SIGTERM

Wait for process exit

Remove the binary and PID file




 ┌───────────────┐
 │     TUI       │
 │ (open DB conn)│
 │  - show menu  │
 │  - select proj│
 └───────┬───────┘
         │ start
         v
 ┌───────────────┐
 │  Webserver    │   <-- separate process
 │ (own DB conn) │
 │  - serve HTTP │
 │  - GitHub sync│
 └───────────────┘
         │
         │ stop (when extraction runs)
         v
 ┌───────────────┐
 │  Extraction   │
 │ (exclusive DB │
 │   read/write) │
 └───────────────┘
         │
         │ restart webserver
         v
 ┌───────────────┐
 │  Webserver    │  (new DB conn)
 └───────────────┘


 ┌───────────────┐
 │     TUI       │
 │ (open DB conn)│
 │  - show menu  │
 │  - select proj│
 └───────┬───────┘
         │ start
         v
 ┌───────────────┐
 │  Webserver    │   <-- separate process
 │ (calls       │
 │  ProviderFactory)
 │  to get DB    │
 │  connection) │
 │  - serve HTTP │
 │  - GitHub sync│
 └───────────────┘
         │
         │ stop (when extraction runs)
         v
 ┌───────────────┐
 │  Extraction   │
 │ (calls       │
 │  ProviderFactory)
 │  to get DB    │
 │  - read/write │
 └───────────────┘
         │
         │ restart webserver
         v
 ┌───────────────┐
 │  Webserver    │  (new DB conn via ProviderFactory)
 └───────────────┘
