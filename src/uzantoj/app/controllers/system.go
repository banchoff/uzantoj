package controllers

import (
	"github.com/revel/revel"
)

type System struct {
	*revel.Controller
}

func (c System) Index() revel.Result {
	return c.Render()
}
