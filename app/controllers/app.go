package controllers

import (
	"ElasticsearchDemo/app/models"
	"fmt"

	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

const SenseUri = "_plugin/marvel/sense/index.html"

func (c App) Index() revel.Result {
	lm := models.NewLessonManager()
	lessons := lm.GenerateLessons()
	return c.Render(lessons)
}

func (c App) Sense() revel.Result {
	rm := models.NewResourceManager()
	senseUrl := fmt.Sprintf("http://%s:%s/%s",
		rm.GetELKAddress(),
		rm.GetELKPort(),
		SenseUri)
	return c.Redirect(senseUrl)
}

func (c App) Lesson(id int) revel.Result {
	rm := models.NewResourceManager()
	str, err := rm.LoadMaterial(id)
	if err != nil {
		return c.NotFound("Couldn't find specified file")
	}
	return c.RenderText(str)
}
