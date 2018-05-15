package gacha

import (
	"fmt"
	"net/http"
)

// RunServer ...
func RunServer(port int) {
	http.HandleFunc("/", helloword)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Print("init fail :", err)
	}
}

func helloword(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello world")
}
