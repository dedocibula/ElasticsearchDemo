package models

import (
	"fmt"
	"io/ioutil"

	"github.com/revel/revel"
)

type ResourceManager struct{}

func (_ ResourceManager) LoadMaterial(materialId int) (string, error) {
	dat, err := ioutil.ReadFile(
		fmt.Sprintf("%s/public/materials/lesson-%v.txt",
			revel.BasePath,
			materialId))
	return string(dat), err
}

func (_ ResourceManager) LoadESAddress() string {
	return fmt.Sprintf("http://%s",
		revel.Config.StringDefault("es.address", "localhost"))
}

func (_ ResourceManager) LoadESPort() string {
	return fmt.Sprintf(":%s",
		revel.Config.StringDefault("es.port", "9200"))
}
