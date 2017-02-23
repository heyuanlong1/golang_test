package main

import ( 
    "fmt" 
    "io" 
    "net" 
"net/http"
"net/rpc"
)

type  Watcher int

func (w *Watcher) Getinfo(arg int,result *int) error {
	*result = 200
	return nil
}

func Webget(w http.ResponseWriter, r *http.Request) { 
    io.WriteString(w, "<html><body>Webget</body></html>") 
}

func main() {
	http.HandleFunc("/Webget",Webget)
	watcher := new(Watcher)
	rpc.Register(watcher)
	rpc.HandleHTTP()

	l , err := net.Listen("tcp",":6000")
	 if err != nil { 
        fmt.Println("监听失败，端口可能已经被占用") 
    } 
    fmt.Println("正在监听6000端口") 
    http.Serve(l,nil)
}