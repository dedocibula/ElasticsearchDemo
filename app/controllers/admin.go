package controllers

import (
	"ElasticsearchDemo/app/models"

	"github.com/revel/revel"
	"golang.org/x/net/websocket"
)

type Admin struct {
	*revel.Controller
}

func (c Admin) Index() revel.Result {
	return c.Render()
}

func (c Admin) ResultEndpoint() revel.Result {
	qm := models.NewQuizManager()
	results := qm.GetResults()
	return c.RenderJson(results)
}

func (c Admin) RecordEndpoint(ws *websocket.Conn) revel.Result {
	subscription := models.QuizMonitorInstance().Subscribe()
	defer models.QuizMonitorInstance().Unsubscribe(subscription)

	for {
		record := <-subscription.New
		if websocket.JSON.Send(ws, &record) != nil {
			return nil
		}
	}
}