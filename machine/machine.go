package machine

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
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
	router := mux.NewRouter().StrictSlash(true)
	router.Methods("GET").Path("/").Name("hello").Handler(http.HandlerFunc(helloworld))
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		fmt.Print("init fail :", err)
	}
}

func helloworld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello world")
}

// AddHandlers ...
func AddHandlers(routes []Route) error {
	return nil
}
