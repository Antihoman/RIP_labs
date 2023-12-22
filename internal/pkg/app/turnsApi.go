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

// @Summary		Получить все ходы
// @Tags		Ходы
// @Description	Возвращает все ходы с фильтрацией по статусу и дате формирования
// @Produce		json
// @Param		status query string false "статус ходы"
// @Param		formation_date_start query string false "начальная дата формирования"
// @Param		formation_date_end query string false "конечная дата формирвания"
// @Success		200 {object} schemes.AllTurnsResponse
// @Router		/api/turns [get]
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
	var turns []ds.Turn
	if userRole == role.Customer {
		turns, err = app.repo.GetAllTurns(&userId, request.FormationDateStart, request.FormationDateEnd, request.Status)
	} else {
		turns, err = app.repo.GetAllTurns(nil, request.FormationDateStart, request.FormationDateEnd, request.Status)
	}
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

// @Summary		Получить одно карту
// @Tags		Ходы
// @Description	Возвращает подробную информацию об ходы и его типе
// @Produce		json
// @Param		id path string true "id ходы"
// @Success		200 {object} schemes.TurnResponse
// @Router		/api/turns/{id} [get]
func (app *Application) GetTurn(c *gin.Context) {
	var request schemes.TurnRequest
	var err error
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userId := getUserId(c)
	userRole := getUserRole(c)
	var turn *ds.Turn
	if userRole == role.Moderator {
		turn, err = app.repo.GetTurnById(request.TurnId, nil)
	} else {
		turn, err = app.repo.GetTurnById(request.TurnId, &userId)
	}
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if turn == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("карту не найдено"))
		return
	}

	cards, err := app.repo.GetTurnContent(request.TurnId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, schemes.TurnResponse{Turn: schemes.ConvertTurn(turn), Cards: cards})
}

type SwaggerUpdateTurnRequest struct {
	TurnPhase string `json:"turn_type"`
}

// @Summary		Указать тип ходы
// @Tags		Ходы
// @Description	Позволяет изменить тип чернового ходы и возвращает обновлённые данные
// @Access		json
// @Produce		json
// @Param		turn_type body SwaggerUpdateTurnRequest true "Тип ходы"
// @Success		200 {object} schemes.TurnOutput
// @Router		/api/turns [put]
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
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("карту не найдено"))
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

// @Summary		Удалить черновое карту
// @Tags		Ходы
// @Description	Удаляет черновое карту
// @Success		200
// @Router		/api/turns [delete]
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
// @Summary		Удалить получателя из черновово ходы
// @Tags		Ходы
// @Description	Удалить получателя из черновово ходы
// @Produce		json
// @Param		id path string true "id получателя"
// @Success		200 {object} schemes.AllCardsResponse
// @Router		/api/turns/delete_recipient/{id} [delete]
func (app *Application) DeleteFromTurn(c *gin.Context) {
	var request schemes.DeleteFromTurnRequest
	var err error
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var turn *ds.Turn
	userId := getUserId(c)
	turn, err = app.repo.GetDraftTurn(userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if turn == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("карту не найдено"))
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

// @Summary		Сформировать ход
// @Tags		Ходы
// @Description	Сформировать карту пользователем
// @Success		200 {object} schemes.TurnOutput
// @Router		/api/turns/user_confirm [put]
func (app *Application) UserConfirm(c *gin.Context) {
	userId := getUserId(c)
	turn, err := app.repo.GetDraftTurn(userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if turn == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("карту не найдено"))
		return
	}
	
	// if err := sendingRequest(turn.UUID); err != nil {
	// 	c.AbortWithError(http.StatusInternalServerError, fmt.Errorf(`sending service is unavailable: {%s}`, err))
	// 	return
	// }

	sendingStatus := ds.SendingStarted
	turn.SendingStatus = &sendingStatus
	turn.Status = ds.StatusFormed
	now := time.Now()
	turn.FormationDate = &now

	if err := app.repo.SaveTurn(turn); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, schemes.ConvertTurn(turn))
}
// @Summary		Подтвердить карту
// @Tags		Ходы
// @Description	Подтвердить или отменить карту модератором
// @Param		id path string true "id ходы"
// @Param		confirm body boolean true "подтвердить"
// @Success		200 {object} schemes.TurnOutput
// @Router		/api/turns/{id}/moderator_confirm [put]
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

	userId := getUserId(c)
	turn, err := app.repo.GetTurnById(request.URI.TurnId,nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if turn == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("карту не найдено"))
		return
	}
	if turn.Status != ds.StatusFormed  {
		c.AbortWithError(http.StatusMethodNotAllowed, fmt.Errorf("нельзя изменить статус с \"%s\" на \"%s\"", turn.Status,  ds.StatusFormed ))
		return
	}


	if *request.Confirm {
		turn.Status = ds.StatusCompleted
		now := time.Now()
		turn.CompletionDate = &now
	
	} else {
		turn.Status = ds.StatusRejected
	}
	moderator, err := app.repo.GetUserById(userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	turn.Moderator = moderator
	
	if err := app.repo.SaveTurn(turn); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, schemes.ConvertTurn(turn))
}

func (app *Application) Sending(c *gin.Context) {
	var request schemes.SendingReq
	if err := c.ShouldBindUri(&request.URI); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	

	if request.Token != app.config.Token {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	turn, err := app.repo.GetTurnById(request.URI.TurnId, nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if turn == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("карту не найдено"))
		return
	}
	// if turn.Status != ds.StatusFormed || *turn.SendingStatus != ds.SendingStarted {
	// 	c.AbortWithStatus(http.StatusMethodNotAllowed)
	// 	return
	// }

	var sendingStatus string
	if *request.SendingStatus {
		sendingStatus = ds.SendingCompleted
	} else {
		sendingStatus = ds.SendingFailed
	}
	turn.SendingStatus = &sendingStatus

	if err := app.repo.SaveTurn(turn); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}