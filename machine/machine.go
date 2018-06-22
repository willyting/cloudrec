package machine

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
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
	server *http.Server
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
	s.server = &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: s.router}
	err := s.server.ListenAndServe()
	if err != nil {
		fmt.Print("init fail :", err)
	}
}

func (s *Server) stop() {
	s.server.Close()
}

func helloworld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello world")
}

// AddHandlers ...
func (s *Server) AddHandlers(routes []Route) error {
	for _, route := range routes {
		handler := cors.AllowAll().Handler(route.HandlerFunc)
		s.router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
		headHhandler := cors.AllowAll().Handler(http.HandlerFunc(helloworld))
		s.router.
			Methods(http.MethodOptions).
			Path(route.Pattern).
			Name(route.Name + "Options").
			Handler(headHhandler)
	}
	return nil
}
