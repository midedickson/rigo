package rigo

import (
	"net"
)

type Server struct {
	net.Listener
}

func NewServer(opt *Options) (*Server, error) {
	opt.init()
	listener, err := net.Listen("tcp", opt.Port)
	if err != nil {
		return nil, err
	}
	return &Server{
		Listener: listener,
	}, nil
}
