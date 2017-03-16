package main

import (
	"net/http"
	"strconv"
	"log"
	"runtime/debug"

	CONF "conf"
	open_gc1 "open/gc1"
	user_gc1 "user/gc1"
)

func safeHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e, ok := recover().(error); ok {
				http.Error(w, e.Error(), http.StatusInternalServerError)
				// 或者输出自定义的50x错误页面
				// w.WriteHeader(http.StatusInternalServerError)
				// renderHtml(w, "error", e)
				// logging
				log.Println("WARN: panic in %v. - %v", fn, e)
				log.Println(string(debug.Stack()))
			}
		}()
		fn(w, r)
	}
}

func main() {
	http.HandleFunc("/open/gc1/uploadAvatar", safeHandler(open_gc1.UploadAvatar))   //设置访问的路由
	http.HandleFunc("/open/gc1/getAvatarInfo", safeHandler(open_gc1.GetAvatarInfo)) //设置访问的路由

	http.HandleFunc("/user/gc1/Sayhello", safeHandler(user_gc1.Sayhello)) //设置访问的路由

	port := strconv.Itoa(CONF.Server_post)
	err := http.ListenAndServe(":"+port, nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
