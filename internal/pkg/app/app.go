package app

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"lab1/internal/app/config"
	"lab1/internal/app/dsn"
	"lab1/internal/app/repository"
)

type Application struct {
	repo        *repository.Repository
	minioClient *minio.Client
	config      *config.Config
}

func (app *Application) Run() {
	log.Println("Server start up")

	r := gin.Default()

	r.Use(ErrorHandler())

	r.GET("/api/cards", app.GetAllCards)
	r.GET("/api/cards/:card_id", app.GetCard)
	r.DELETE("/api/cards/:card_id", app.DeleteCard)
	r.PUT("/api/cards/:card_id", app.ChangeCard)
	r.POST("/api/cards", app.AddCard)
	r.POST("/api/cards/:card_id/add_to_turn", app.AddToTurn)

	r.GET("/api/turns", app.GetAllTurns)
	r.GET("/api/turns/:turn_id", app.GetTurn)
	r.PUT("/api/turns/:turn_id/update", app.UpdateTurn)
	r.DELETE("/api/turns/:turn_id", app.DeleteTurn)
	r.DELETE("/api/turns/:turn_id/delete_card/:card_id", app.DeleteFromTurn)
	r.PUT("/api/turns/:turn_id/user_confirm", app.UserConfirm)
	r.PUT("/api/turns/:turn_id/moderator_confirm", app.ModeratorConfirm)

	r.Static("/image", "./resources/images")
	r.Static("/css", "./resources/css")
	r.Run("localhost:8080")
	log.Println("Server down")
}

func New() (*Application, error) {
	var err error
	loc, _ := time.LoadLocation("UTC")
	time.Local = loc
	app := Application{}
	app.config, err = config.NewConfig()
	if err != nil {
		return nil, err
	}

	app.repo, err = repository.New(dsn.FromEnv())
	if err != nil {
		return nil, err
	}

	app.minioClient, err = minio.New(app.config.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4("", "", ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	return &app, nil
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, err := range c.Errors {
			log.Println(err.Err)
		}
		lastError := c.Errors.Last()
		if lastError != nil {
			switch c.Writer.Status() {
			case http.StatusBadRequest:
				c.JSON(-1, gin.H{"error": "wrong request"})
			case http.StatusNotFound:
				c.JSON(-1, gin.H{"error": lastError.Error()})
			case http.StatusMethodNotAllowed:
				c.JSON(-1, gin.H{"error": lastError.Error()})
			default:
				c.Status(-1)
			}
		}
	}
}