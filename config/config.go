package config

import bolt "go.etcd.io/bbolt"

type Config struct {
	DB *bolt.DB
}

type Todo struct {
	File string
	Line int
	Text string
}
