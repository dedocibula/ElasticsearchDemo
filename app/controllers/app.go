package controllers

import (
	"ElasticsearchDemo/app/models"
	"fmt"

	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

var (
	rm = models.ResourceManager{}
	lm = models.LessonManager{}
)

func (c App) Index() revel.Result {
	lessons := lm.GenerateLessons()
	return c.Render(lessons)
}

func (c App) Sense() revel.Result {
	senseUrl := fmt.Sprintf("http://%s:%s/%s",
		rm.GetELKAddress(),
		rm.GetELKPort(),
		rm.GetSenseUri())
	return c.Redirect(senseUrl)
}

func (c App) Lesson(id int) revel.Result {
	str, err := rm.LoadMaterial(id)
	if err != nil {
		return c.NotFound("Couldn't find specified file")
	}
	return c.RenderText(str)
}
