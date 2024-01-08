package schemes

import (
	"fmt"
	"lab1/internal/app/ds"
)

type AllCardsResponse struct {
	Cards []ds.Card `json:"cards"`
}

type GetAllCardsResponse struct {
	DraftTurn *string   `json:"draft_turn"`
	Cards     []ds.Card `json:"cards"`
}

type AllTurnsResponse struct {
	Turns []TurnOutput `json:"turns"`
}

type TurnResponse struct {
	Turn  TurnOutput `json:"turn"`
	Cards []ds.Card  `json:"cards"`
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
	SendingStatus  *string `json:"sending_status"`
	TakeFood       *string `json:"takefood"`
}

func ConvertTurn(turn *ds.Turn) TurnOutput {
	output := TurnOutput{
		UUID:          turn.UUID,
		Status:        turn.Status,
		CreationDate:  turn.CreationDate.Format("2006-01-02 15:04:05"),
		TakeFood:      turn.TakeFood,
		SendingStatus: turn.SendingStatus,
		Customer:      turn.Customer.Login,
	}

	if turn.FormationDate != nil {
		formationDate := turn.FormationDate.Format("2006-01-02 15:04:05")
		output.FormationDate = &formationDate
	}

	if turn.CompletionDate != nil {
		completionDate := turn.CompletionDate.Format("2006-01-02 15:04:05")
		output.CompletionDate = &completionDate
	}

	if turn.Moderator != nil {
		fmt.Println(turn.Moderator.Login)
		output.Moderator = &turn.Moderator.Login
		fmt.Println(*output.Moderator)
	}

	return output
}

type AddToTurnResp struct {
	CardsCount int64 `json:"card_count"`
}

type AuthResp struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}
