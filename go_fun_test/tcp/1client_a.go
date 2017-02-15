package main

import(
"fmt"
"net"
"os"
)

/*

http://www.cnblogs.com/gaopeng527/p/6128290.html
http://blog.csdn.net/qq_15437667/article/details/51042366

*/


//go run 1client_a.go 127.0.0.1:5000
func main() {
	var buf[512]byte
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr,"Usage: %s host:port ",os.Args[0])
		os.Exit(1)
	}
	service := os.Args[1]
	tcpAddr,err := net.ResolveTCPAddr("tcp4",service)
	checkErr(err)
	conn,err := net.DialTCP("tcp",nil ,tcpAddr)
	defer conn.Close()
	checkErr(err)
	rAddr := conn.RemoteAddr()
	n ,err := conn.Write([]byte("Hello servie"))
	checkErr(err)
	n, err =conn.Read(buf[0:])
	checkErr(err)
	fmt.Println("Reply from server ", rAddr.String(), string(buf[0:n]))
    os.Exit(0)
}

func checkErr(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
        os.Exit(1)
    }
}
