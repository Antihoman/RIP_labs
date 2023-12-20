package app

import (
	"fmt"
	"net/http"
	"time"

	"lab1/internal/app/ds"
	"lab1/internal/app/role"
	"lab1/internal/app/schemes"

	"github.com/gin-gonic/gin"
)

// @Summary		Получить все уведомления
// @Tags		Уведомления
// @Description	Возвращает все уведомления с фильтрацией по статусу и дате формирования
// @Produce		json
// @Param		status query string false "статус уведомления"
// @Param		formation_date_start query string false "начальная дата формирования"
// @Param		formation_date_end query string false "конечная дата формирвания"
// @Success		200 {object} schemes.AllTurnsResponse
// @Router		/api/notifications [get]
func (app *Application) GetAllTurns(c *gin.Context) {
	var request schemes.GetAllTurnsRequst
	var err error
	if err = c.ShouldBindQuery(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userId := getUserId(c)
	userRole := getUserRole(c)
	fmt.Println(userId, userRole)
	var notifications []ds.Turn
	if userRole == role.Customer {
		notifications, err = app.repo.GetAllTurns(&userId, request.FormationDateStart, request.FormationDateEnd, request.Status)
	} else {
		notifications, err = app.repo.GetAllTurns(nil, request.FormationDateStart, request.FormationDateEnd, request.Status)
	}
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	outputTurns := make([]schemes.TurnOutput, len(notifications))
	for i, notification := range notifications {
		outputTurns[i] = schemes.ConvertTurn(&notification)
	}
	c.JSON(http.StatusOK, schemes.AllTurnsResponse{Turns: outputTurns})
}

// @Summary		Получить одно уведомление
// @Tags		Уведомления
// @Description	Возвращает подробную информацию об уведомлении и его типе
// @Produce		json
// @Param		id path string true "id уведомления"
// @Success		200 {object} schemes.TurnResponse
// @Router		/api/notifications/{id} [get]
func (app *Application) GetTurn(c *gin.Context) {
	var request schemes.TurnRequest
	var err error
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userId := getUserId(c)
	userRole := getUserRole(c)
	var notification *ds.Turn
	if userRole == role.Moderator {
		notification, err = app.repo.GetTurnById(request.TurnId, nil)
	} else {
		notification, err = app.repo.GetTurnById(request.TurnId, &userId)
	}
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if notification == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("уведомление не найдено"))
		return
	}

	cards, err := app.repo.GetTurnContent(request.TurnId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, schemes.TurnResponse{Turn: schemes.ConvertTurn(notification), Cards: cards})
}

type SwaggerUpdateTurnRequest struct {
	TurnPhase string `json:"notification_type"`
}

// @Summary		Указать тип уведомления
// @Tags		Уведомления
// @Description	Позволяет изменить тип чернового уведомления и возвращает обновлённые данные
// @Access		json
// @Produce		json
// @Param		notification_type body SwaggerUpdateTurnRequest true "Тип уведомления"
// @Success		200 {object} schemes.TurnOutput
// @Router		/api/notifications [put]
func (app *Application) UpdateTurn(c *gin.Context) {
	var request schemes.UpdateTurnRequest
	var err error
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// Получить черновую заявку
	var turn *ds.Turn

	userId := getUserId(c)
	turn, err = app.repo.GetDraftTurn(userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if turn == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("уведомление не найдено"))
		return
	}

	// Добавить тип
	turn.Phase = &request.TurnPhase
	if app.repo.SaveTurn(turn); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, schemes.ConvertTurn(turn))
}

// @Summary		Удалить черновое уведомление
// @Tags		Уведомления
// @Description	Удаляет черновое уведомление
// @Success		200
// @Router		/api/notifications [delete]
func (app *Application) DeleteTurn(c *gin.Context) {
	var err error
	// Получить черновую заявку
	var turn *ds.Turn
	userId := getUserId(c)
	turn, err = app.repo.GetDraftTurn(userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if turn == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("увдомление не найдено"))
		return
	}

	turn.Status = ds.StatusDeleted

	if err := app.repo.SaveTurn(turn); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}
// @Summary		Удалить получателя из черновово уведомления
// @Tags		Уведомления
// @Description	Удалить получателя из черновово уведомления
// @Produce		json
// @Param		id path string true "id получателя"
// @Success		200 {object} schemes.AllRecipientsResponse
// @Router		/api/notifications/delete_recipient/{id} [delete]
func (app *Application) DeleteFromTurn(c *gin.Context) {
	var request schemes.DeleteFromTurnRequest
	var err error
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// Получить черновую заявку
	var turn *ds.Turn
	userId := getUserId(c)
	turn, err = app.repo.GetDraftTurn(userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if turn == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("уведомление не найдено"))
		return
	}

	if err := app.repo.DeleteFromTurn(turn.UUID, request.CardId); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	cards, err := app.repo.GetTurnContent(turn.UUID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, schemes.AllCardsResponse{Cards: cards})
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
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("ход не найдено"))
		return
	}
	if turn.Status != ds.DRAFT {
		c.AbortWithError(http.StatusMethodNotAllowed, fmt.Errorf("нельзя сформировать ход со статусом %s", turn.Status))
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
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("ход не найдено"))
		return
	}
	if turn.Status != ds.FORMED {
		c.AbortWithError(http.StatusMethodNotAllowed, fmt.Errorf("нельзя изменить статус с \"%s\" на \"%s\"", turn.Status, ds.FORMED))
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