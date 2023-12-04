package app

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"lab1/internal/app/config"
	"lab1/internal/app/ds"
	"lab1/internal/app/dsn"
	"lab1/internal/app/repository"
)

type Application struct {
	repo   *repository.Repository
	config *config.Config
	// dsn string
}

type GetCardBack struct {
	Cards []ds.Card
	Name  string
}

func (a *Application) Run() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/cards", func(c *gin.Context) {
		Name := c.Query("Name")
		cards, err := a.repo.GetCardByName(Name)

		if err != nil {
			log.Println("cant get recipients", err)
			c.Error(err)
			return
		}
		c.HTML(http.StatusOK, "index.tmpl", GetCardBack{
			Name:  Name,
			Cards: cards,
		})

	})

	r.GET("/cards/:id", func(c *gin.Context) {
		id := c.Param("id")
		card, err := a.repo.GetCardByID(id)
		if err != nil { // если не получилось
			log.Printf("cant get product by id %v", err)
			c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "detail.tmpl", *card)
	})

	r.POST("/cards", func(c *gin.Context) {
		id := c.PostForm("delete")

		a.repo.DeleteCard(id)

		cards, err := a.repo.GetCardByName("")
		if err != nil {
			log.Println("cant get recipients", err)
			c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "index.tmpl", GetCardBack{
			Name:  "",
			Cards: cards,
		})
	})

	r.Static("/image", "./resources/images")
	r.Static("/css", "./resources/css")
	r.Run("localhost:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	log.Println("Server down")
}

func New() (*Application, error) {
	var err error
	app := Application{}
	app.config, err = config.NewConfig()
	if err != nil {
		return nil, err
	}

	app.repo, err = repository.New(dsn.FromEnv())
	if err != nil {
		return nil, err
	}

	return &app, nil
}
