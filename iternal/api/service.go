package api

import (
	"strings"
)

type Card struct {
	ID 			int
	Name        string
	Description string
	Phase       string
	ImageURL    string
}

type FilteredCards struct {
	Cards []Card
	Filter   string
}

var cards = []Card{
	{0, "Пугливое", "Описание карточки 1", "Развитие", "/image/img1.png"},
	{1, "r-стратегия", "Описание карточки 2", "Установление кормовой базы", "/image/img2.png"},
	{2, "Теплокровность", "Описание карточки 3", "Питание", "/image/img3.png"},
	{3, "Пугливое", "Описание карточки 1", "Развитие", "/image/img4.png"},
	{4, "r-стратегия", "Описание карточки 2", "Установление кормовой базы", "/image/img5.png"},
	{5, "Теплокровность", "Описание карточки 3", "Питание", "/image/img6.png"},
	{6, "Пугливое", "Описание карточки 1", "Развитие", "/image/img7.png"},
	{7, "r-стратегия", "Описание карточки 2", "Установление кормовой базы", "/image/img8.png"},
	{8, "Теплокровность", "Описание карточки 3", "Питание", "/image/img9.png"},
}

func GetAllServices(filter string) FilteredCards {
	var filteredCards []Card

	if filter == "" {
		filteredCards = cards
	} else {
		for _, s := range cards {
			if containsIgnoreCase(s.Name, filter) {
				filteredCards = append(filteredCards, s)
			}
		}
	}

	return FilteredCards{
		Cards: filteredCards,
		Filter:   filter,
	}
}

func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func GetServiceByIndex(index int) Card {
	return cards[index]
}
