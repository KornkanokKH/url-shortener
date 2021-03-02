package main

import (
	"time"
)

// Config data model
type Config struct {
	Server
}

// Server data model
type Server struct {
	Name            string
	Env             string
	Port            int
	Timeout         time.Duration `conf:"default:5s"`
	ShutdownTimeout time.Duration `conf:"default:5s"`
}
