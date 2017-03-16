package gc1

import (
"fmt"
"net/http"
)


func Sayhello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello go web user! \n") //这个写入到w的是输出到客户端的
}

