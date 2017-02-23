package main
 
import (
 "io"
 "log"
 "net/http"
)
 
func helloHandler(w http.ResponseWriter, r *http.Request) {
 io.WriteString(w, "hello world!")
}
 
func main() {
 http.HandleFunc("/hello", helloHandler)
 err := http.ListenAndServeTLS(":6000", "server.crt", "server.key", nil)
 if err != nil {
 log.Fatal("ListenAndServeTLS:", err.Error())
 }
}

/*

http://studygolang.com/articles/2946

生成秘钥和证书
openssl genrsa -out server.key 2048
openssl req -new -x509 -key server.key -out server.crt -days 3650

https://112.74.27.195:6000/hello
*/