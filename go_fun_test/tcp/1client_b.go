package main

import(
"fmt"
"net"
"os"
)

//go run 1client_b.go 127.0.0.1:5000
func main() {
	conn, err := net.Dial("udp", "127.0.0.1:5000")
    defer conn.Close()
    if err != nil {
        os.Exit(1)  
    }

    conn.Write([]byte("Hello world!"))  

    var msg [20]byte
    n, _ := conn.Read(msg[0:])

    fmt.Println("msg is:", string(msg[0:n]))
}
