package models

import (
	"fmt"
	"github.com/revel/revel"
	"io/ioutil"
)

type ResourceManager struct{}

func (r ResourceManager) LoadMaterial(materialId int) (string, error) {
	dat, err := ioutil.ReadFile(
		fmt.Sprintf("%s/public/materials/lesson-%v.txt",
			revel.BasePath,
			materialId))
	return string(dat), err
}

func (r ResourceManager) LoadESAddress() string {
	return fmt.Sprintf("http://%s",
		revel.Config.StringDefault("es.address", "localhost"))
}

func (r ResourceManager) LoadESPort() string {
	return fmt.Sprintf(":%s",
		revel.Config.StringDefault("es.port", "9200"))
}
