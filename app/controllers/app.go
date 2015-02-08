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

var (
	rm = models.NewResourceManager()
	lm = models.NewLessonManager()
)

func (c App) Index() revel.Result {
	lessons := lm.GenerateLessons()
	return c.Render(lessons)
}

func (c App) Sense() revel.Result {
	senseUrl := fmt.Sprintf("http://%s:%s/%s",
		rm.GetELKAddress(),
		rm.GetELKPort(),
		SenseUri)
	return c.Redirect(senseUrl)
}

func (c App) Lesson(id int) revel.Result {
	str, err := rm.LoadMaterial(id)
	if err != nil {
		return c.NotFound("Couldn't find specified file")
	}
	return c.RenderText(str)
}
