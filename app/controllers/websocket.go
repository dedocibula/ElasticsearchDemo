package controllers

import (
	"ElasticsearchDemo/app/models"

	"github.com/revel/revel"
	"golang.org/x/net/websocket"
)

type Websocket struct {
	*revel.Controller
}

func (c Websocket) Records(ws *websocket.Conn) revel.Result {
	subscription := models.QuizMonitorInstance().Subscribe()
	defer models.QuizMonitorInstance().Unsubscribe(subscription)

	for {
		record := <-subscription.New
		if websocket.JSON.Send(ws, &record) != nil {
			return nil
		}
	}
}
