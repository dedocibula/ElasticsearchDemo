package controllers

import (
	"ElasticsearchDemo/app/models"

	"github.com/revel/revel"
)

var (
	em = models.NewELKManager(models.NewResourceManager())
)

type Quiz struct {
	*revel.Controller
}

func (c Quiz) Index() revel.Result {
	return c.Render()
}

func (c Quiz) Submit() revel.Result {
	answer, err := em.LiteralQueryELK()
	if err != nil {
		c.Flash.Error(err.Error())
	} else {
		c.Flash.Success("%v", answer)
	}
	c.FlashParams()
	return c.Redirect(Quiz.Index)
}
