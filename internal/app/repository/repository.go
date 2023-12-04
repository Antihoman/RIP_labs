package repository

//вопрос владу
import (
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"lab1/internal/app/ds"
)

type Repository struct {
	db *gorm.DB
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) GetCardByID(id string) (*ds.Card, error) { // ?
	card := &ds.Card{}
	err := r.db.Where("card_id = ?", id).First(card).Error
	if err != nil {
		return nil, err
	}

	return card, nil
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

func (r *Repository) DeleteCard(id string) error {
	err := r.db.Exec("UPDATE cards SET is_deleted = ? WHERE card_id = ?", true, id).Error
	if err != nil {
		return err
	}

	return nil
}