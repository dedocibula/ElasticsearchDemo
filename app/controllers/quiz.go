package controllers

import (
	"ElasticsearchDemo/app/models"

	"github.com/revel/revel"
)

var (
	qm = models.NewQuizManager()
)

type Quiz struct {
	*revel.Controller
}

func (c Quiz) Index() revel.Result {
	return c.Render()
}

func (c Quiz) Submit(answer int) revel.Result {
	c.Validation.Required(answer).Message("Your answer cannot be empty or 0!")
	c.Validation.Min(answer, 0).Message("Your answer cannot be negative!")

	if c.Validation.HasErrors() {
		c.Validation.Keep()
	} else {
		c.validateAnswer(answer)
	}
	c.FlashParams()

	return c.Redirect(Quiz.Index)
}

func (c Quiz) validateAnswer(answer int) {
	result := qm.Validate(answer)
	switch result.Ok {
	case true:
		c.Flash.Success(result.Message)
	case false:
		c.Flash.Error(result.Message)
	}
}
