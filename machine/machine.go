package machine

import (
	"fmt"
	"net/http"
)

// Route is basic struct to add to router
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// RunServer start a http server on X port
func RunServer(port int) {
	http.HandleFunc("/", helloword)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Print("init fail :", err)
	}
}

func helloword(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello world")
}

// AddHandlers ...
func AddHandlers(routes []Route) error {
	return nil
}
