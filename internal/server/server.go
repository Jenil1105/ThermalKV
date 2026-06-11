package server

import (
	"net"
)

type KVService interface {
	Set(key string, value string)
	Get(key string) (string, bool)
	Delete(key string)
	SetTTL(key string, seconds int)
	CoolKey(key string) error
	Count() int
	Exists(key string) bool
	Keys() []string
	GetInfo() []string
}

type Server struct {
	Listener net.Listener
	DB       KVService
}

func New(db KVService) (*Server, error) {
	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		return nil, err
	}
	return &Server{
		Listener: listener,
		DB:       db,
	}, nil
}
