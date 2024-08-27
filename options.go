package rigo

import "os"

type Options struct {
	Port string
}

func (o *Options) init() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "6379"
	}

	o.Port = ":" + port
}

func NewOptions() *Options {
	opt := &Options{}
	opt.init()
	return opt
}
