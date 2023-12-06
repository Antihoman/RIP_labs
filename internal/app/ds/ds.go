package ds

import (
	"time"
)

type User struct {
	UUID      string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"-"`
	Moderator bool   `json:"moderator"`
	Login     string `gorm:"size:30;not null" json:"login"`
	Password  string `gorm:"size:40;not null" json:"-"`
	Name      string `gorm:"size:50;not null" json:"name"`
}

type Card struct {
	UUID        string  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"uuid" binding:"-"`
	ImageURL    *string `gorm:"type:varchar(100)" json:"image_url" binding:"-"`
	IsDeleted   bool    `gorm:"not null;default:false" json:"-" binding:"-"`
	Type        string  `gorm:"size:50;not null" form:"type" json:"type" binding:"required,max=50"`
	Name        string  `gorm:"size:50;not null" form:"name" json:"name" binding:"required,max=50"`
	Description string  `gorm:"size:200;not null" form:"description" json:"description" binding:"required,max=200"`
}

const DRAFT string = "черновик"
const FORMED string = "Сформирован"
const COMPELTED string = "Завершён"
const REJECTED string = "Отклонён"
const DELETED string = "Удален"

type Turn struct {
	UUID           string     `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Status         string     `gorm:"size:20;not null"`
	CreationDate   time.Time  `gorm:"not null;type:timestamp"`
	FormationDate  *time.Time `gorm:"type:timestamp"`
	CompletionDate *time.Time `gorm:"type:timestamp"`
	ModeratorId    *string    `json:"-"`
	CustomerId     string     `gorm:"not null"`
	Phase          string     `gorm:"size:50;not null"`

	Moderator *User
	Customer  User
}

type PlayedCards struct {
	TurnId   string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"turn_id"`
	CardId   string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"card_id"`
	TakeFood uint   `gorm:"not null"`

	Card *Card `gorm:"foreignKey:CardId" json:"card"`
	Turn *Turn `gorm:"foreignKey:TurnId" json:"turn"`
}
