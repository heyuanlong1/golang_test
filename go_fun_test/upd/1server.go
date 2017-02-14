package main

import(
"fmt"
"net"
"os"
)

//go run 1service.go
func main() {
	service := ":5000"
	udpAddr,err := net.ResolveUDPAddr("udp4",service)
	checkErr(err)
	conn,err := net.ListenUDP("udp",udpAddr)
	checkErr(err)
	defer conn.Close()
	handleClient(conn)
}


func handleClient(conn *net.UDPConn) {
	var buf [512]byte
	for {
		n,rAddr,err := conn.ReadFromUDP(buf[0:])
		if err !=nil {
			return
		}
		fmt.Println("Receive from client", rAddr.String(), string(buf[0:n]))
		_,err2 := conn.WriteToUDP([]byte("Welcome client!"),rAddr)
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