package main

import (
	"time"
)

type tomlConfig struct {
	Title   string
	Owner   ownerInfo
	DB      database `toml:"database"`
	Servers map[string]server
	Clients clients
	Logs    logs
	Api     api
	Smtp    smtp
}

type ownerInfo struct {
	Name string
	Org  string `toml:"organization"`
	Bio  string
	DOB  time.Time
	URL  string
}

type database struct {
	Server   string
	Port     string
	User     string
	Password string
	Dbname   string
	ConnMax  int `toml:"connection_max"`
	Enabled  bool
}

type server struct {
	IP string
	DC string
}

type clients struct {
	Data  [][]interface{}
	Hosts []string
}

type logs struct {
	Location   string
	Flag       int
	Permission int
}

type api struct {
	Listener string
}

type smtp struct {
	From     string
	CC       string
	Host     string
	Port     int
	Username string
	Password string
}
