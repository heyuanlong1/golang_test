package main

import(
"fmt"
"net"
"os"
)

//go run 1service.go
func main() {
	listen_sock,err := net.Listen("tcp","localhost:5000")
	checkErr(err)
	defer listen_sock.Close()
	for{
		new_conn,err := listen_sock.Accept()
		if err != nil {
			continue
		}
		go handleClient(new_conn)
	}
}


func handleClient(conn net.Conn) {
	defer conn.Close()
	var buf [512]byte
	for {
		n,err := conn.Read(buf[0:])
		if err !=nil {
			return
		}
		rAddr := conn.RemoteAddr()
		fmt.Println("Receive from client", rAddr.String(), string(buf[0:n]))
		_,err2 := conn.Write([]byte("Welcome client!"))
		if err2 !=nil {
			return
		}
	}
	
}

func checkErr(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
        os.Exit(1)
    }
}