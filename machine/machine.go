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

// Server ...
type Server struct {
	router *mux.Router
}

// NewServer ...
func NewServer() *Server {
	server := new(Server)
	server.router = mux.NewRouter().StrictSlash(true)
	return server
}

// Run start a http server on X port
func (s *Server) Run(port int) {
	s.router.Methods("GET").Path("/").Name("hello").Handler(http.HandlerFunc(helloworld))
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), s.router)
	if err != nil {
		fmt.Print("init fail :", err)
	}
}

func helloworld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello world")
}

// AddHandlers ...
func (s *Server) AddHandlers(routes []Route) error {
	for _, route := range routes {
		s.router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return nil
}
