package models

import (
	"encoding/json"
	"fmt"

	"github.com/mattbaird/elastigo/lib"
)

const UIDPrefix = "UID_"

type ELKManager struct {
	rm   *ResourceManager
	conn *elastigo.Conn
}

func NewELKManager(rm *ResourceManager) *ELKManager {
	return &ELKManager{rm: rm}
}

func (e *ELKManager) Dispose() {
	e.rm = nil
	e.conn = nil
}

func (e ELKManager) LiteralSearchELK(index, _type string) (int, error) {
	err := e.validateParams(index, _type)
	if err != nil {
		return -1, err
	}

	err = e.initialize()
	if err != nil {
		return -1, err
	}

	query, err := e.rm.LoadQuizQuery()
	if err != nil {
		return -1, err
	}

	result, err := e.conn.Search(index, _type, nil, query)
	if err != nil {
		return -1, err
	}

	return e.parseQueryResult(result)
}

func (e ELKManager) ExistsRecordELK(index, _type string, record ELKRecord) (bool, error) {
	err := e.validateParams(index, _type)
	if err != nil {
		return false, err
	}

	err = e.initialize()
	if err != nil {
		return false, err
	}

	id := fmt.Sprintf("%s%s", UIDPrefix, record.Nickname)

	return e.conn.ExistsBool(index, _type, id, nil)
}

func (e ELKManager) RecordSuccessELK(index, _type string, record ELKRecord) error {
	err := e.validateParams(index, _type)
	if err != nil {
		return err
	}

	err = e.initialize()
	if err != nil {
		return err
	}

	result, err := e.tryInsert(index, _type, record)
	if err != nil {
		return err
	}

	return e.verifySuccess(result)
}

func (e *ELKManager) initialize() error {
	if e.rm == nil {
		return fmt.Errorf("ResourceManager isn't initialized")
	}

	if e.conn == nil {
		c := elastigo.NewConn()
		c.Domain = e.rm.GetELKAddress()
		c.Port = e.rm.GetELKPort()

		e.conn = c
	}
	return nil
}

func (_ ELKManager) validateParams(params ...string) error {
	for _, param := range params {
		if param == "" {
			return fmt.Errorf("Given parameters cannot be empty: %v", params)
		}
	}
	return nil
}

func (_ ELKManager) parseQueryResult(result elastigo.SearchResult) (int, error) {
	if len(result.Hits.Hits) == 1 {
		var m map[string][]int
		err := json.Unmarshal(*result.Hits.Hits[0].Fields, &m)
		if err != nil {
			return -1, err
		}
		return m["author.id"][0], nil
	} else {
		return -1, fmt.Errorf("Number of results expected: %v, got %v",
			1, len(result.Hits.Hits))
	}
}

func (e ELKManager) tryInsert(index, _type string, record ELKRecord) ([]byte, error) {
	url := fmt.Sprintf("/%s/%s/%s%s", index, _type, UIDPrefix, record.Nickname)

	return e.conn.DoCommand("POST", url, nil, record)
}

func (e ELKManager) verifySuccess(result []byte) error {
	var m map[string]interface{}
	err := json.Unmarshal(result, &m)
	if err != nil {
		return err
	}

	success := m["created"].(bool)
	if !success {
		return fmt.Errorf("Unfortunately, this nickname is already taken.")
	} else {
		return nil
	}
}
