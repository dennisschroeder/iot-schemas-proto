package main
import (
    "fmt"
    "net"
    "time"
)
func main() {
    conn, err := net.DialTimeout("tcp", "192.168.178.88:50000", 2*time.Second)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Success:", conn.LocalAddr(), "->", conn.RemoteAddr())
    conn.Close()
}