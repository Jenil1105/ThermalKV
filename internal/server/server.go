package server

import (
	"net"
	"thermalkv/internal/store"
)

type Server struct {
	Listener net.Listener
	DB       *store.Store
}

func New(db *store.Store) (*Server, error) {
	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		return nil, err
	}
	return &Server{
		Listener: listener,
		DB:       db,
	}, nil
}
