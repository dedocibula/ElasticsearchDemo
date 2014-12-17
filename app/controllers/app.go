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

func (c App) Sense() revel.Result {
	address := rm.LoadESAddress() + rm.LoadESPort() + "/_plugin/marvel/sense/index.html"
	return c.Redirect(address)
}

func (c App) Lesson(id int) revel.Result {
	str, err := rm.LoadMaterial(id)
	if err != nil {
		return c.NotFound("Couldn't find specified file")
	}
	return c.RenderText(str)
}
