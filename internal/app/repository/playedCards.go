package repository

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	"lab1/internal/app/ds"
)

func (r *Repository) GetAllTurns(customerId *string, formationDateStart, formationDateEnd *time.Time, status string) ([]ds.Turn, error) {
	var turns []ds.Turn
	query := r.db.Preload("Customer").Preload("Moderator").
		Where("LOWER(status) LIKE ?", "%"+strings.ToLower(status)+"%").
		Where("status != ?", ds.StatusDeleted).
		Where("status != ?", ds.StatusDraft)

	if formationDateStart != nil && formationDateEnd != nil {
		query = query.Where("formation_date BETWEEN ? AND ?", *formationDateStart, *formationDateEnd)
	} else if formationDateStart != nil {
		query = query.Where("formation_date >= ?", *formationDateStart)
	} else if formationDateEnd != nil {
		query = query.Where("formation_date <= ?", *formationDateEnd)
	}
	if err := query.Find(&turns).Error; err != nil {
		return nil, err
	}
	return turns, nil
}

func (r *Repository) GetDraftTurn(customerId string) (*ds.Turn, error) {
	turn := &ds.Turn{}
	err := r.db.First(turn, ds.Turn{Status: ds.StatusDraft, CustomerId: customerId}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return turn, nil
}

func (r *Repository) CreateDraftTurn(customerId string) (*ds.Turn, error) {
	turn := &ds.Turn{CreationDate: time.Now(), CustomerId: customerId, Status: ds.StatusDraft}
	err := r.db.Create(turn).Error
	if err != nil {
		return nil, err
	}
	return turn, nil
}

func (r *Repository) GetTurnById(turnId string, userId *string) (*ds.Turn, error) {
	turn := &ds.Turn{}
	query := r.db.Preload("Moderator").Preload("Customer").
		Where("status != ?", ds.StatusDeleted)
	if userId != nil {
		query = query.Where("customer_id = ?", userId)
	}
	err := query.First(turn, ds.Turn{UUID: turnId}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return turn, nil
}

func (r *Repository) GetTurnContent(turnId string) ([]ds.Card, error) {
	var cards []ds.Card

	err := r.db.Table("played_cards").
		Select("cards.*").
		Joins("JOIN cards ON played_cards.card_id = cards.uuid").
		Where(ds.PlayedCards{TurnId: turnId}).
		Scan(&cards).Error

	if err != nil {
		return nil, err
	}
	return cards, nil
}

func (r *Repository) SaveTurn(turn *ds.Turn) error {
	err := r.db.Save(turn).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteFromTurn(turnId, cardId string) error {
	err := r.db.Delete(&ds.PlayedCards{TurnId: turnId, CardId: cardId}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) CountCards(turnId string) (int64, error) {
	var count int64
	err := r.db.Model(&ds.PlayedCards{}).
		Where("turn_id = ?", turnId).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}