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

	// Услуги - получатели
	r.GET("/api/cards", app.GetAllCards)                                     // Список с поиском
	r.GET("/api/cards/:card_id", app.GetCard)                           // Одна услуга
	r.DELETE("/api/cards/:card_id", app.DeleteCard)              // Удаление
	r.PUT("/api/cards/:card_id", app.ChangeCard)                 // Изменение
	r.POST("/api/cards", app.AddCard)                                    // Добавление
	r.POST("/api/cards/:card_id/add_to_turn", app.AddToTurn) // Добавление в заявку

	// Заявки - уведомления
	r.GET("/api/turns", app.GetAllTurns)                                                       // Список (отфильтровать по дате формирования и статусу)
	r.GET("/api/turns/:turn_id", app.GetTurn)                                          // Одна заявка
	r.PUT("/api/turns/:turn_id/update", app.UpdateTurn)                                // Изменение (добавление транспорта)
	r.DELETE("/api/turns/:turn_id", app.DeleteTurn)                             //Удаление
	r.DELETE("/api/turns/:turn_id/delete_card/:card_id", app.DeleteFromTurn) // Изменеие (удаление услуг)
	r.PUT("/api/turns/:turn_id/user_confirm", app.UserConfirm)                                 // Сформировать создателем
	r.PUT("/api/turns/:turn_id/moderator_confirm", app.ModeratorConfirm)                        // Завершить отклонить модератором

	r.Static("/image", "./resources/images")
	r.Static("/css", "./resources/css")
	r.Run("localhost:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
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