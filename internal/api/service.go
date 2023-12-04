package api

import (
	"strings"
)

type Service struct {
	ID 			int
	Name        string
	Description string
	Phase       string
	ImageURL    string
}

type FilteredServices struct {
	Services []Service
	Filter   string
}

var services = []Service{
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

func GetAllServices(filter string) FilteredServices {
	var filteredServices []Service

	if filter == "" {
		filteredServices = services
	} else {
		for _, s := range services {
			if containsIgnoreCase(s.Name, filter) {
				filteredServices = append(filteredServices, s)
			}
		}
	}

	return FilteredServices{
		Services: filteredServices,
		Filter:   filter,
	}
}

func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func GetServiceByIndex(index int) Service {
	return services[index]
}
