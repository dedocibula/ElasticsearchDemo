package controllers

import (
	"ElasticsearchDemo/app/models"
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

var rm = models.ResourceManager{}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Lesson(id string) revel.Result {
	address := rm.LoadESAddress() + rm.LoadESPort() + "/_plugin/marvel/sense/index.html"
	param := rm.LoadESAddress() + "/App/Material/" + id
	return c.Redirect("%s?load_from=%s", address, param)
}

func (c App) Material(id int) revel.Result {
	str, err := rm.LoadMaterial(id)
	if err != nil {
		return c.NotFound("Couldn't find specified file")
	}
	return c.RenderText(str)
}
