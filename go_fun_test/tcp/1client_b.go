package main

import(
"fmt"
"net"
"os"
)

//go run 1client_b.go 
func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:5000")
    defer conn.Close()
    if err != nil {
        os.Exit(1)  
    }

    conn.Write([]byte("Hello world!"))  

    var msg [200]byte
    n, _ := conn.Read(msg[0:])

    fmt.Println("msg is:", string(msg[0:n]))
}
