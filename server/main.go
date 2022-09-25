package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type Server struct{ *http.ServeMux }

func Create() *Server {
	s := http.NewServeMux()
	return &Server{s}
}

func Listen(p int, s *Server) {
	ps := strconv.Itoa(p)
	fmt.Printf("Server starting on port %v", ps)
	e := http.ListenAndServe(":"+ps, s)

	if e != nil {
		fmt.Print("An error occured during server startup", e)
		os.Exit(2)
	}

}
