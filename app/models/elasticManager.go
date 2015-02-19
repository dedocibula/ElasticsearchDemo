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

func (e ELKManager) ExistsRecordELK(index, _type, name string) (bool, error) {
	err := e.validateParams(index, _type)
	if err != nil {
		return false, err
	}

	err = e.initialize()
	if err != nil {
		return false, err
	}

	id := fmt.Sprintf("%s%s", UIDPrefix, name)

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

	return e.verifySuccess(result, "created", "Could not create a record")
}

func (e ELKManager) SelectRecordsELK(index, _type string) ([]ELKRecord, error) {
	err := e.validateParams(index, _type)
	if err != nil {
		return nil, err
	}

	err = e.initialize()
	if err != nil {
		return nil, err
	}

	result, err := elastigo.Search(index).Type(_type).Pretty().Query(
		elastigo.Query().All(),
	).Sort(
		elastigo.Sort("Timestamp").Asc(),
	).Result(e.conn)
	if err != nil {
		return nil, err
	}

	return e.buildELKRecords(result)
}

func (e ELKManager) ClearTypeELK(index, _type string) error {
	err := e.validateParams(index)
	if err != nil {
		return err
	}

	err = e.initialize()
	if err != nil {
		return err
	}

	result, err := e.tryClear(index, _type)
	if err != nil {
		return err
	}

	return e.verifySuccess(result, "acknowledged", "Could not clear records")
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

func (e ELKManager) tryClear(index, _type string) ([]byte, error) {
	url := fmt.Sprintf("/%s/%s", index, _type)

	return e.conn.DoCommand("DELETE", url, nil, nil)
}

func (e ELKManager) verifySuccess(result []byte, field, errMessage string) error {
	var m map[string]interface{}
	err := json.Unmarshal(result, &m)
	if err != nil {
		return err
	}

	success := m[field].(bool)
	if !success {
		return fmt.Errorf(errMessage)
	} else {
		return nil
	}
}

func (e ELKManager) buildELKRecords(result *elastigo.SearchResult) ([]ELKRecord, error) {
	var ers []ELKRecord
	if len(result.Hits.Hits) > 0 {
		for _, hit := range result.Hits.Hits {
			var er ELKRecord
			err := json.Unmarshal(*hit.Source, &er)
			if err != nil {
				return nil, err
			}
			ers = append(ers, er)
		}
	}
	return ers, nil
}
