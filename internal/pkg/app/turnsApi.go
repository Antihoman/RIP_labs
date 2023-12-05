package app

import (
	"fmt"
	"net/http"
	"time"

	"lab1/internal/app/ds"
	"lab1/internal/app/schemes"


	"github.com/gin-gonic/gin"
)

func (app *Application) GetAllTurns(c *gin.Context) {
	var request schemes.GetAllTurnsRequst
	if err := c.ShouldBindQuery(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	turns, err := app.repo.GetAllTurns(request.FormationDateStart, request.FormationDateEnd, request.Status)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	outputTurns := make([]schemes.TurnOutput, len(turns))
	for i, turn := range turns {
		outputTurns[i] = schemes.ConvertTurn(&turn)
	}
	c.JSON(http.StatusOK, schemes.AllTurnsResponse{Turns: outputTurns})
}

func (app *Application) GetTurn(c *gin.Context) {
	var request schemes.TurnRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	turn, err := app.repo.GetTurnById(request.TurnId, app.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if turn == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("уведомление не найдено"))
		return
	}

	cards, err := app.repo.GetTurnContent(request.TurnId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, schemes.TurnResponse{Turn: schemes.ConvertTurn(turn), Card: cards})
}

func (app *Application) UpdateTurn(c *gin.Context) {
	var request schemes.UpdateTurnRequest
	if err := c.ShouldBindUri(&request.URI); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	turn, err := app.repo.GetTurnById(request.URI.TurnId, app.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if turn == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("уведомление не найдено"))
		return
	}
	turn.Phase = request.TurnPhase
	if app.repo.SaveTurn(turn); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, schemes.UpdateTurnResponse{Turn:schemes.ConvertTurn(turn)})
}

func (app *Application) DeleteTurn(c *gin.Context) {
	var request schemes.TurnRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	turn, err := app.repo.GetTurnById(request.TurnId, app.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if turn == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("увдомление не найдено"))
		return
	}
	turn.Status = ds.DELETED

	if err := app.repo.SaveTurn(turn); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

func (app *Application) DeleteFromTurn(c *gin.Context) {
	var request schemes.DeleteFromTurnRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	turn, err := app.repo.GetTurnById(request.TurnId, app.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if turn == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("уведомление не найдено"))
		return
	}
	if turn.Status != ds.DRAFT {
		c.AbortWithError(http.StatusMethodNotAllowed, fmt.Errorf("нельзя редактировать уведомление со статусом: %s", turn.Status))
		return
	}

	if err := app.repo.DeleteFromTurn(request.TurnId, request.CardId); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	cards, err := app.repo.GetTurnContent(request.TurnId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, schemes.AllCardsResponse{Card: cards})
}

func (app *Application) UserConfirm(c *gin.Context) {
	var request schemes.UserConfirmRequest
	if err := c.ShouldBindUri(&request.URI); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	turn, err := app.repo.GetTurnById(request.URI.TurnId, app.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if turn == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("уведомление не найдено"))
		return
	}
	if turn.Status != ds.DRAFT {
		c.AbortWithError(http.StatusMethodNotAllowed, fmt.Errorf("нельзя сформировать уведомление со статусом %s", turn.Status))
		return
	}
	if request.Confirm {
		turn.Status = ds.FORMED
		now := time.Now()
		turn.FormationDate = &now
	} else {
		turn.Status = ds.DELETED
	}

	if err := app.repo.SaveTurn(turn); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

func (app *Application) ModeratorConfirm(c *gin.Context) {
	var request schemes.ModeratorConfirmRequest
	if err := c.ShouldBindUri(&request.URI); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	turn, err := app.repo.GetTurnById(request.URI.TurnId, app.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if turn == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("уведомление не найдено"))
		return
	}
	if turn.Status != ds.FORMED {
		c.AbortWithError(http.StatusMethodNotAllowed, fmt.Errorf("нельзя изменить статус с \"%s\" на \"%s\"", turn.Status,  ds.FORMED))
		return
	}
	if request.Confirm {
		turn.Status = ds.COMPELTED
		now := time.Now()
		turn.CompletionDate = &now
	
	} else {
		turn.Status = ds.REJECTED
	}
	turn.ModeratorId = app.getModerator()
	
	if err := app.repo.SaveTurn(turn); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}