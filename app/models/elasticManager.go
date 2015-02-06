package models

import (
	"github.com/mattbaird/elastigo/lib"
)

type ElasticManager struct {
	ELKHost string
	ELKPort string

	conn *elastigo.Conn
}
