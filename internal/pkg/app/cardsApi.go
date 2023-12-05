package app

import (
	"fmt"
	"net/http"

	"lab1/internal/app/ds"
	"lab1/internal/app/schemes"


	"github.com/gin-gonic/gin"
)

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
	draftTurn, err := app.repo.GetDraftTurn(app.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	response := schemes.GetAllCardsResponse{DraftTurn: nil, Card: cards}
	if draftTurn != nil {
		response.DraftTurn = &schemes.TurnShort{UUID: draftTurn.UUID}
		containers, err := app.repo.GetTurnContent(draftTurn.UUID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		response.DraftTurn.CardCount = len(containers)
	}
	c.JSON(http.StatusOK, response)
}

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
	card.IsDeleted = true
	if err := app.repo.SaveCard(card); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

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

	c.Status(http.StatusOK)
}

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
	if request.Description != nil {
		card.Description = *request.Description
	}

	if err := app.repo.SaveCard(card); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, card)
}

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
	turn, err = app.repo.GetDraftTurn(app.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if turn == nil {
		turn, err = app.repo.CreateDraftTurn(app.getCustomer())
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	if err = app.repo.AddToTurn(turn.UUID, request.CardId); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var cards []ds.Card
	cards, err = app.repo.GetTurnContent(turn.UUID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, schemes.AllCardsResponse{Card: cards})
}