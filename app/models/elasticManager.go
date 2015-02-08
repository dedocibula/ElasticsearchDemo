package models

import (
	"encoding/json"
	"fmt"

	"github.com/mattbaird/elastigo/lib"
)

type ElasticManager struct {
	RM *ResourceManager

	conn *elastigo.Conn
}

func (e *ElasticManager) Initialize() bool {
	if e.RM == nil {
		return false
	}

	if e.conn == nil {
		c := elastigo.NewConn()
		c.Domain = e.RM.GetELKAddress()
		c.Port = e.RM.GetELKPort()

		e.conn = c
	}
	return true
}

func (e *ElasticManager) LiteralQueryELK() (int, error) {
	e.Initialize()

	query, err := e.getQuizQuery()
	if err != nil {
		return -1, err
	}

	result, err := e.conn.Search("dba", "question", nil, query)
	if err != nil {
		return -1, err
	}

	if len(result.Hits.Hits) == 1 {
		var m map[string][]int
		err = json.Unmarshal(*result.Hits.Hits[0].Fields, &m)
		if err != nil {
			return -1, err
		}
		return m["author.id"][0], nil
	} else {
		return -1, fmt.Errorf("Number of results expected: %v, got %v",
			1, len(result.Hits.Hits))
	}
}

func (e *ElasticManager) Dispose() {
	e.RM = nil
	e.conn = nil
}

func (e ElasticManager) getQuizQuery() (string, error) {
	if e.RM == nil {
		return "", fmt.Errorf("ResourceManager isn't initialized")
	}
	return e.RM.LoadQuizQuery()
}
