package controllers

import (
	"github.com/revel/revel"
)

type Quiz struct {
	*revel.Controller
}

func (c Quiz) Index() revel.Result {
	return c.Render()
}
