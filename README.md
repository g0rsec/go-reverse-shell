# goReverseShell

First ever GO project aimed to learn GO basics. PoC of a GO reverse shell, not expected to be used in real conditions.  
- Windows / Unix support
- Unsecure communications
- No AV bypass || anti reverse solution
- Only one shell process opened. Commands are processed by remote STDIN. STDOUT, STDERR are sent back.

```go 
go run server.go  
go build -ldflags -H=windowsgui client.go
```