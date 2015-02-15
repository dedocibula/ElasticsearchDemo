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

func (c Quiz) Submit(attempt models.Attempt) revel.Result {
	c.Validation.Required(attempt.Name).Message("Your nick cannot be empty!")
	c.Validation.MinSize(attempt.Name, 5).Message("Your nick must be at least 5 characters long!")

	c.Validation.Required(attempt.Answer).Message("Your answer cannot be empty or 0!")
	c.Validation.Min(attempt.Answer, 0).Message("Your answer cannot be negative!")

	if c.Validation.HasErrors() {
		c.Validation.Keep()
	} else {
		c.validateAttempt(attempt)
	}
	c.FlashParams()

	return c.Redirect(Quiz.Index)
}

func (c Quiz) Results() revel.Result {
	results := qm.GetResults()
	return c.Render(results)
}

func (c Quiz) validateAttempt(attempt models.Attempt) {
	result := qm.Validate(attempt)
	switch result.Ok {
	case true:
		c.Flash.Success(result.Messages[0])
	case false:
		for _, message := range result.Messages {
			c.Validation.Error(message).Key(message)
		}
		c.Validation.Keep()
	}
}
