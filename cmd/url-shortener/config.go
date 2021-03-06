package main

import (
	"time"
	"url-shortener/internal/repository/redis"
)

// Config data model
type Config struct {
	Server
	Redis redis.Config
}

// Server data model
type Server struct {
	Name            string
	Env             string
	Port            int
	Timeout         time.Duration `conf:"default:5s"`
	ShutdownTimeout time.Duration `conf:"default:5s"`
}
