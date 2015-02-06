package models

import (
	"fmt"
	"io/ioutil"

	"github.com/revel/revel"
)

const (
	elkAddressKey = "elk.address"
	elkPortKey    = "elk.port"
)

type ResourceManager struct {
	elkAddress, elkPort string
}

func (_ ResourceManager) LoadMaterial(materialId int) (string, error) {
	dat, err := ioutil.ReadFile(
		fmt.Sprintf("%s/public/materials/lesson-%v.txt",
			revel.BasePath,
			materialId))
	return string(dat), err
}

func (r *ResourceManager) GetELKAddress() string {
	if r.elkAddress == "" {
		r.elkAddress = r.loadFromConfig(elkAddressKey, "localhost")
	}
	return r.elkAddress
}

func (r *ResourceManager) GetELKPort() string {
	if r.elkPort == "" {
		r.elkPort = r.loadFromConfig(elkPortKey, "9200")
	}
	return r.elkPort
}

func (_ ResourceManager) loadFromConfig(key, defaultValue string) string {
	return revel.Config.StringDefault(key, defaultValue)
}
