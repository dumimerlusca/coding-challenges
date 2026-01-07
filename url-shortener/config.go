package main

type HttpConfig struct {
	Addr string
}

type DbConfig struct {
	ConnString string
}

type Config struct {
	Http HttpConfig
	Db   DbConfig
	Url  string
}
