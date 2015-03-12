package controllers

import (
	"ElasticsearchDemo/app/models"

	"code.google.com/p/go.net/websocket"
	"github.com/revel/revel"
)

type Admin struct {
	*revel.Controller
}

func (c Admin) Index() revel.Result {
	attemptFields := models.UIHelperInstance().GenerateAttemptFields()
	return c.Render(attemptFields)
}

func (c Admin) ClearResults() revel.Result {
	qm := models.NewQuizManager()
	if !qm.ClearResults() {
		c.Flash.Error("Couldn't clear the results. Please, perform manually.")
		c.FlashParams()
	}
	return c.Redirect(Admin.Index)
}

func (c Admin) ResultEndpoint() revel.Result {
	qm := models.NewQuizManager()
	results := qm.GetResults()
	return c.RenderJson(results)
}

func (c Admin) AttemptEndpoint(ws *websocket.Conn) revel.Result {
	subscription := models.QuizMonitorInstance().Subscribe()
	defer models.QuizMonitorInstance().Unsubscribe(subscription)

	for {
		attempt := <-subscription.New
		if websocket.JSON.Send(ws, &attempt) != nil {
			return nil
		}
	}
	return nil
}
