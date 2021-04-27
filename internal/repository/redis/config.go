package redis

type Config struct {
	Key         string
	Expire      int64
	RedisServer Server
}

type Server struct {
	Address string
	Port    string
	Db      int
}
