package app

import (
	"log"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"lab1/internal/app/config"
	"lab1/internal/app/dsn"
	"lab1/internal/app/redis"
	"lab1/internal/app/repository"
	"lab1/internal/app/role"

	_ "lab1/docs"
	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // gin-swagger middleware
	_ "lab1/docs"
)

type Application struct {
	repo        *repository.Repository
	minioClient *minio.Client
	config      *config.Config
	redisClient *redis.Client
}

func (app *Application) Run() {
	log.Println("Server start up")

	r := gin.Default()

	r.Use(ErrorHandler())

	// Услуги - получатели
	api := r.Group("/api")
	{
		res := api.Group("/cards")
		{
			res.GET("/", app.WithAuthCheck(role.NotAuthorized, role.Customer, role.Moderator), app.GetAllCards)                     // Список с поиском
			res.GET("/:id", app.WithAuthCheck(role.NotAuthorized, role.Customer, role.Moderator), app.GetCard)            // Одна услуга
			res.DELETE("/:id", app.WithAuthCheck(role.Moderator), app.DeleteCard)                         				// Удаление
			res.PUT("/:id", app.WithAuthCheck(role.Moderator), app.ChangeCard)                            				// Изменение
			res.POST("", app.WithAuthCheck(role.Moderator), app.AddCard)                                           				// Добавление
			res.POST("/:id/add_to_turn", app.WithAuthCheck(role.Customer,role.Moderator), app.AddToTurn) 		// Добавление в заявку
		}

		// Заявки - уведомления
		n := api.Group("/turns")
		{
			n.GET("/", app.WithAuthCheck(role.Customer, role.Moderator), app.GetAllTurns)                                         				  // Список (отфильтровать по дате формирования и статусу)
			n.GET("/:id",app.WithAuthCheck(role.Customer, role.Moderator),  app.GetTurn)                             				  // Одна заявка
			n.PUT("", app.WithAuthCheck(role.Customer, role.Moderator), app.UpdateTurn)                                	  // Изменение (добавление транспорта)
			n.DELETE("", app.WithAuthCheck(role.Customer,role.Moderator), app.DeleteTurn)                                      				  // Удаление
			n.DELETE("/delete_recipient/:id", app.WithAuthCheck(role.Customer, role.Moderator), app.DeleteFromTurn) 	  // Изменеие (удаление услуг)
			n.PUT("/user_confirm", app.WithAuthCheck(role.Customer, role.Moderator), app.UserConfirm)                                    				  // Сформировать создателем
			n.PUT("/:id/moderator_confirm", app.WithAuthCheck(role.Moderator), app.ModeratorConfirm)                         				  // Завершить отклонить модератором
			n.PUT("/:id/sending", app.Sending)
		}

		// Пользователи (авторизация)
		u := api.Group("/user")
		{
			u.POST("/sign_up", app.Register)
			u.POST("/login", app.Login)
			u.POST("/logout", app.Logout)
		}

		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		r.Run(fmt.Sprintf("%s:%d", app.config.ServiceHost, app.config.ServicePort))

		log.Println("Server down")
	}
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

	app.minioClient, err = minio.New(app.config.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4("", "", ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	app.redisClient, err = redis.New(app.config.Redis)
	if err != nil {
		return nil, err
	}

	return &app, nil
}