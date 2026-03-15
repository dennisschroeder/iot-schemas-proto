package main
import (
    "fmt"
    "net/http"
    "io"
    "strings"
)
func main() {
    body := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}`
    req, _ := http.NewRequest("POST", "http://homeassistant.local:8123/api/mcp", strings.NewReader(body))
    req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiIxZTk4NmE2MDQwYTg0NWVlYWY0ZjhlMDg1YzM0NGE3MiIsImlhdCI6MTc3MjIxMDYwMywiZXhwIjoyMDg3NTcwNjAzfQ.ikYcKsPcXCTKA_D6w8LHb5_JEo2Wncragv4OZevFyFk")
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")
    
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer resp.Body.Close()
    respBytes, _ := io.ReadAll(resp.Body)
    fmt.Println("Status:", resp.StatusCode)
    fmt.Println("Body:", string(respBytes))
}