package app

import (
	"fmt"
	"net/http"

	_ "lab1/docs"
	"lab1/internal/app/ds"
	"lab1/internal/app/schemes"

	"github.com/gin-gonic/gin"
)

// @Summary		Получить все карты
// @Tags		карты
// @Description	Возвращает всех доуступных карт с опциональной фильтрацией по названию
// @Produce		json
// @Param		name query string false "Название для фильтрации"
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
	response := schemes.GetAllCardsResponse{DraftTurn: nil, Cards: cards}
	if userId, exists := c.Get("userId"); exists {
		draftTurn, err := app.repo.GetDraftTurn(userId.(string))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if draftTurn != nil {
			response.DraftTurn = &draftTurn.UUID
		}
	}
	c.JSON(http.StatusOK, response)
}

// @Summary		Получить одну карту
// @Tags		карты
// @Description	Возвращает более подробную информацию об одной карте
// @Produce		json
// @Param		id path string true "id карты"
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
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("карта не найден"))
		return
	}
	c.JSON(http.StatusOK, card)
}

// @Summary		Удалить карту
// @Tags		карты
// @Description	Удаляет карту по id
// @Param		id path string true "id карты"
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
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("карта не найден"))
		return
	}	
	if card.ImageURL != nil {
		if err := app.deleteImage(c, card.UUID); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
	card.ImageURL = nil
	card.IsDeleted = true
	if err := app.repo.SaveCard(card); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

// @Summary		Добавить карту
// @Tags		карты
// @Description	Добавить новую карту
// @Accept		mpfd
// @Param     	image formData file false "Изображение карты"
// @Param     	name formData string true "Название" format:"string" maxLength:100
// @Param     	type formData string true "Тип" format:"string" maxLength:100
// @Param     	need_food formData int true "Нужно еды" format:"int"
// @Param     	description formData string true "Описание" format:"string" maxLength:100
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

// @Summary		Изменить карту
// @Tags		карты
// @Description	Изменить данные полей о карте
// @Accept		mpfd
// @Produce		json
// @Param		id path string true "Идентификатор карты" format:"uuid"
// @Param     	image formData file false "Изображение карты"
// @Param     	name formData string true "Название" format:"string" maxLength:100
// @Param     	type formData string true "Тип" format:"string" maxLength:100
// @Param     	need_food formData int true "Нужно еды" format:"int"
// @Param     	description formData string true "Описание" format:"string" maxLength:100
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
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("карта не найден"))
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

// @Summary		Добавить в ход
// @Tags		карты
// @Description	Добавить выбранную карту в черновик хода
// @Produce		json
// @Param		id path string true "id карты"
// @Success		200
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
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("карта не найден"))
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
