package schemes

import (
	"lab1/internal/app/ds"
)

type AllCardsResponse struct {
	Cards []ds.Card `json:"recipients"`
}

type TurnShort struct {
	UUID      string `json:"uuid"`
	CardCount int    `json:"recipient_count"`
}

type GetAllCardsResponse struct {
	DraftTurn *TurnShort `json:"draft_turn"`
	Cards     []ds.Card  `json:"recipients"`
}

type AllTurnsResponse struct {
	Turns []TurnOutput `json:"turns"`
}

type TurnResponse struct {
	Turn  TurnOutput `json:"turn"`
	Cards []ds.Card  `json:"recipients"`
}

type UpdateTurnResponse struct {
	Turn TurnOutput `json:"turns"`
}

type TurnOutput struct {
	UUID           string  `json:"uuid"`
	Status         string  `json:"status"`
	CreationDate   string  `json:"creation_date"`
	FormationDate  *string `json:"formation_date"`
	CompletionDate *string `json:"completion_date"`
	Moderator      *string `json:"moderator"`
	Customer       string  `json:"customer"`
	TurnPhase      string  `gorm:"size:50;not null"`
}

func ConvertTurn(turn *ds.Turn) TurnOutput {
	output := TurnOutput{
		UUID:         turn.UUID,
		Status:       turn.Status,
		CreationDate: turn.CreationDate.Format("0001-01-01 01:01:01"),
		TurnPhase:    turn.Phase,
		Customer:     turn.Customer.Name,
	}

	if turn.FormationDate != nil {
		formationDate := turn.FormationDate.Format("0001-01-01 01:01:01")
		output.FormationDate = &formationDate
	}

	if turn.CompletionDate != nil {
		completionDate := turn.CompletionDate.Format("0001-01-01 01:01:01")
		output.CompletionDate = &completionDate
	}

	if turn.Moderator != nil {
		output.Moderator = &turn.Moderator.Name
	}

	return output
}
