package main
import (
    "fmt"
    "net"
    "time"
)
func main() {
    conn, err := net.DialTimeout("tcp4", "192.168.178.88:50000", 2*time.Second)
    if err != nil {
        fmt.Println("tcp4 Error:", err)
    } else {
        fmt.Println("tcp4 Success:", conn.LocalAddr(), "->", conn.RemoteAddr())
        conn.Close()
    }
    
    conn2, err := net.DialTimeout("tcp6", "[::ffff:192.168.178.88]:50000", 2*time.Second)
    if err != nil {
        fmt.Println("tcp6 Error:", err)
    } else {
        fmt.Println("tcp6 Success:", conn2.LocalAddr(), "->", conn2.RemoteAddr())
        conn2.Close()
    }
}