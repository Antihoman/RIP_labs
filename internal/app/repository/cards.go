package repository

import (
	"errors"
	"strings"

	"lab1/internal/app/ds"

	"gorm.io/gorm"
)

func (r *Repository) GetCardByID(id string) (*ds.Card, error) { // ?
	card := &ds.Card{UUID: id}
	err := r.db.First(card, "is_deleted = ?", false).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return card, nil
}

func (r *Repository) AddCard(card *ds.Card) error {
	err := r.db.Create(&card).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetAllCards() ([]ds.Card, error) {
	var card []ds.Card

	err := r.db.Find(&card).Error
	if err != nil {
		return nil, err
	}

	return card, nil
}

func (r *Repository) GetCardByName(Name string) ([]ds.Card, error) {
	var cards []ds.Card

	err := r.db.
		Where("LOWER(cards.name) LIKE ?", "%"+strings.ToLower(Name)+"%").Where("is_deleted = ?", false).
		Find(&cards).Error

	if err != nil {
		return nil, err
	}

	return cards, nil
}

func (r *Repository) SaveCard(card *ds.Card) error {
	err := r.db.Save(card).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) AddToTurn(turnId, cardId string) error {
	PCard := ds.PlayedCards{TurnId: turnId, CardId: cardId}
	err := r.db.Create(&PCard).Error
	if err != nil {
		return err
	}
	return nil
}