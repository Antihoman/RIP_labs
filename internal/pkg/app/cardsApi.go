package app

import (
	"fmt"
	"net/http"

	_ "lab1/docs"
	"lab1/internal/app/ds"
	"lab1/internal/app/schemes"


	"github.com/gin-gonic/gin"
)

// @Summary		Получить всех получателей
// @Tags		Получатели
// @Description	Возвращает всех доуступных получателей с опциональной фильтрацией по ФИО
// @Produce		json
// @Param		fio query string false "ФИО для фильтрации"
// @Success		200 {object} schemes.GetAllCardsResponse
// @Router		/api/cards [get]
func (app *Application) GetAllCards(c *gin.Context) {
	var request schemes.GetAllCardsRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	cards, err := app.repo.GetCardByName(request.Name)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var draftTurn *ds.Turn = nil
	if userId, exists := c.Get("userId"); exists {
		draftTurn, err = app.repo.GetDraftTurn(userId.(string))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
	response := schemes.GetAllCardsResponse{DraftTurn: nil, Cards: cards}
	if draftTurn != nil {
		response.DraftTurn = &schemes.TurnShort{UUID: draftTurn.UUID}
		cardsCount, err := app.repo.CountCards(draftTurn.UUID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		response.DraftTurn.CardCount = int(cardsCount)
	}
	c.JSON(http.StatusOK, response)
}

// @Summary		Получить одного получателя
// @Tags		Получатели
// @Description	Возвращает более подробную информацию об одном получателе
// @Produce		json
// @Param		id path string true "id получателя"
// @Success		200 {object} ds.Card
// @Router		/api/cards/{id} [get]
func (app *Application) GetCard(c *gin.Context) {
	var request schemes.CardRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	card, err := app.repo.GetCardByID(request.CardId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if card == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("получатель не найден"))
		return
	}
	c.JSON(http.StatusOK, card)
}

// @Summary		Удалить получателя
// @Tags		Получатели
// @Description	Удаляет получателя по id
// @Param		id path string true "id получателя"
// @Success		200
// @Router		/api/cards/{id} [delete]
func (app *Application) DeleteCard(c *gin.Context) {
	var request schemes.CardRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	card, err := app.repo.GetCardByID(request.CardId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if card == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("получатель не найден"))
		return
	}
	card.ImageURL = nil
	card.IsDeleted = true
	if card.ImageURL != nil {
		if err := app.deleteImage(c, card.UUID); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	c.Status(http.StatusOK)
}

// @Summary		Добавить получателя
// @Tags		Получатели
// @Description	Добавить нового получателя
// @Accept		mpfd
// @Param     	image formData file false "Изображение получателя"
// @Param     	fio formData string true "ФИО" format:"string" maxLength:100
// @Param     	email formData string true "Почта" format:"string" maxLength:100
// @Param     	age formData int true "Возраст" format:"int"
// @Param     	adress formData string true "Адрес" format:"string" maxLength:100
// @Success		200
// @Router		/api/cards [post]
func (app *Application) AddCard(c *gin.Context) {
	var request schemes.AddCardRequest
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	card := ds.Card(request.Card)
	if err := app.repo.AddCard(&card); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if request.Image != nil {
		imageURL, err := app.uploadImage(c, request.Image, card.UUID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		card.ImageURL = imageURL
	}
	if err := app.repo.SaveCard(&card); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusCreated)
}

// @Summary		Изменить получателя
// @Tags		Получатели
// @Description	Изменить данные полей о получателе
// @Accept		mpfd
// @Produce		json
// @Param		id path string true "Идентификатор получателя" format:"uuid"
// @Param		fio formData string false "ФИО" format:"string" maxLength:100
// @Param		email formData string false "Почта" format:"string" maxLength:100
// @Param		age formData int false "Возраст" format:"int"
// @Param		image formData file false "Изображение получателя"
// @Param		adress formData string false "Адрес" format:"string" maxLength:100
// @Router		/api/cards/{id} [put]
func (app *Application) ChangeCard(c *gin.Context) {
	var request schemes.ChangeCardRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	card, err := app.repo.GetCardByID(request.CardId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if card == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("получатель не найден"))
		return
	}

	if request.Name != nil {
		card.Name = *request.Name
	}
	if request.Image != nil {
		if card.ImageURL != nil {
			if err := app.deleteImage(c, card.UUID); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}
		imageURL, err := app.uploadImage(c, request.Image, card.UUID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		card.ImageURL = imageURL
	}
	if request.Type != nil {
		card.Type = *request.Type
	}
	if request.NeedFood != nil {
		card.NeedFood = *request.NeedFood
	}
	if request.Description != nil {
		card.Description = *request.Description
	}

	if err := app.repo.SaveCard(card); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, card)
}

// @Summary		Добавить в уведомление
// @Tags		Получатели
// @Description	Добавить выбранного получателя в черновик уведомления
// @Produce		json
// @Param		id path string true "id получателя"
// @Success		200 {object} schemes.AddToTurnResp
// @Router		/api/cards/{id}/add_to_turn [post]
func (app *Application) AddToTurn(c *gin.Context) {
	var request schemes.AddToTurnRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var err error

	card, err := app.repo.GetCardByID(request.CardId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if card == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("получатель не найден"))
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
		turn, err = app.repo.CreateDraftTurn(userId)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	if err = app.repo.AddToTurn(turn.UUID, request.CardId); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	cardsCount, err := app.repo.CountCards(turn.UUID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, schemes.AddToTurnResp{CardsCount: cardsCount})
}