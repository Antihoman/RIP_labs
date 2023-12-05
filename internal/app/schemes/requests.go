package schemes

import (
	"lab1/internal/app/ds"

	"mime/multipart"
	"time"
)

type CardRequest struct {
	CardId string `uri:"card_id" binding:"required,uuid"`
}

type GetAllCardsRequest struct {
	Name string `form:"name"`
}

type AddCardRequest struct {
	ds.Card
	Image *multipart.FileHeader `form:"image" json:"image"`
}

type ChangeCardRequest struct {
	CardId      string                `uri:"card_id" binding:"required,uuid"`
	Image       *multipart.FileHeader `form:"image" json:"image"`
	Type        *string                `form:"type" json:"type" binding:"omitempty,max=30"`
	Name        *string                `form:"name" json:"name" binding:"omitempty,max=50"`
	Description *string                `form:"description" json:"description" binding:"omitempty,max=200"`
}

type AddToTurnRequest struct {
	CardId string `uri:"card_id" binding:"required,uuid"`
}

type GetAllTurnsRequst struct {
	FormationDateStart *time.Time `form:"formation_date_start" json:"formation_date_start" time_format:"2006-01-02 15:04:05"`
	FormationDateEnd   *time.Time `form:"formation_date_end" json:"formation_date_end" time_format:"2006-01-02 15:04:05"`
	Status             string     `form:"status"`
}

type TurnRequest struct {
	TurnId string `uri:"turn_id" binding:"required,uuid"`
}

type UpdateTurnRequest struct {
	URI struct {
		TurnId string `uri:"turn_id" binding:"required,uuid"`
	}
	TurnPhase string `form:"turn_phase" json:"turn_phase" binding:"required,max=50"`
}

type DeleteFromTurnRequest struct {
	TurnId string `uri:"turn_id" binding:"required,uuid"`
	CardId string `uri:"card_id" binding:"required,uuid"`
}

type UserConfirmRequest struct {
	URI struct {
		TurnId string `uri:"turn_id" binding:"required,uuid"`
	}
	Confirm bool `form:"confirm" binding:"required"`
}

type ModeratorConfirmRequest struct {
	URI struct {
		TurnId string `uri:"turn_id" binding:"required,uuid"`
	}
	Confirm bool `form:"confirm" binding:"required"`
}
