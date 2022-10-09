package server

import (
	"fmt"
	"net/http"
	"os"
)

type Server struct{ *http.ServeMux }

func Create() *Server {
	s := http.NewServeMux()
	return &Server{s}
}

func (s *Server) Listen(port string) {
	fmt.Printf("Server starting on port %v", port)
	e := http.ListenAndServe(":"+port, s)

	if e != nil {
		fmt.Print("An error occured during server startup", e)
		os.Exit(2)
	}

}
