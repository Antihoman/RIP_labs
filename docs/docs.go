// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/cards": {
            "get": {
                "description": "Возвращает всех доуступных получателей с опциональной фильтрацией по ФИО",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Получатели"
                ],
                "summary": "Получить всех получателей",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ФИО для фильтрации",
                        "name": "fio",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schemes.GetAllCardsResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Добавить нового получателя",
                "consumes": [
                    "multipart/form-data"
                ],
                "tags": [
                    "Получатели"
                ],
                "summary": "Добавить получателя",
                "parameters": [
                    {
                        "type": "file",
                        "description": "Изображение получателя",
                        "name": "image",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "ФИО",
                        "name": "fio",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Почта",
                        "name": "email",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Возраст",
                        "name": "age",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Адрес",
                        "name": "adress",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/api/cards/{id}": {
            "get": {
                "description": "Возвращает более подробную информацию об одном получателе",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Получатели"
                ],
                "summary": "Получить одного получателя",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id получателя",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/ds.Card"
                        }
                    }
                }
            },
            "put": {
                "description": "Изменить данные полей о получателе",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Получатели"
                ],
                "summary": "Изменить получателя",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Идентификатор получателя",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ФИО",
                        "name": "fio",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "Почта",
                        "name": "email",
                        "in": "formData"
                    },
                    {
                        "type": "integer",
                        "description": "Возраст",
                        "name": "age",
                        "in": "formData"
                    },
                    {
                        "type": "file",
                        "description": "Изображение получателя",
                        "name": "image",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "Адрес",
                        "name": "adress",
                        "in": "formData"
                    }
                ],
                "responses": {}
            },
            "delete": {
                "description": "Удаляет получателя по id",
                "tags": [
                    "Получатели"
                ],
                "summary": "Удалить получателя",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id получателя",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/api/cards/{id}/add_to_turn": {
            "post": {
                "description": "Добавить выбранного получателя в черновик уведомления",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Получатели"
                ],
                "summary": "Добавить в уведомление",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id получателя",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schemes.AddToTurnResp"
                        }
                    }
                }
            }
        },
        "/api/turns": {
            "get": {
                "description": "Возвращает все ходы с фильтрацией по статусу и дате формирования",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Ходы"
                ],
                "summary": "Получить все ходы",
                "parameters": [
                    {
                        "type": "string",
                        "description": "статус ходы",
                        "name": "status",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "начальная дата формирования",
                        "name": "formation_date_start",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "конечная дата формирвания",
                        "name": "formation_date_end",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schemes.AllTurnsResponse"
                        }
                    }
                }
            },
            "put": {
                "description": "Позволяет изменить тип чернового ходы и возвращает обновлённые данные",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Ходы"
                ],
                "summary": "Указать тип ходы",
                "parameters": [
                    {
                        "description": "Тип ходы",
                        "name": "turn_type",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/app.SwaggerUpdateTurnRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schemes.TurnOutput"
                        }
                    }
                }
            },
            "delete": {
                "description": "Удаляет черновое карту",
                "tags": [
                    "Ходы"
                ],
                "summary": "Удалить черновое карту",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/api/turns/delete_recipient/{id}": {
            "delete": {
                "description": "Удалить получателя из черновово ходы",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Ходы"
                ],
                "summary": "Удалить получателя из черновово ходы",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id получателя",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schemes.AllCardsResponse"
                        }
                    }
                }
            }
        },
        "/api/turns/user_confirm": {
            "put": {
                "description": "Сформировать карту пользователем",
                "tags": [
                    "Ходы"
                ],
                "summary": "Сформировать ход",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schemes.TurnOutput"
                        }
                    }
                }
            }
        },
        "/api/turns/{id}": {
            "get": {
                "description": "Возвращает подробную информацию об ходы и его типе",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Ходы"
                ],
                "summary": "Получить одно карту",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id ходы",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schemes.TurnResponse"
                        }
                    }
                }
            }
        },
        "/api/turns/{id}/moderator_confirm": {
            "put": {
                "description": "Подтвердить или отменить карту модератором",
                "tags": [
                    "Ходы"
                ],
                "summary": "Подтвердить карту",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id ходы",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "подтвердить",
                        "name": "confirm",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "boolean"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schemes.TurnOutput"
                        }
                    }
                }
            }
        },
        "/api/user/login": {
            "post": {
                "description": "Авторизует пользователя по логиню, паролю и отдаёт jwt токен для дальнейших запросов",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Авторизация"
                ],
                "summary": "Авторизация",
                "parameters": [
                    {
                        "description": "login and password",
                        "name": "user_credentials",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/schemes.LoginReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schemes.SwaggerLoginResp"
                        }
                    }
                }
            }
        },
        "/api/user/logout": {
            "post": {
                "description": "Выход из аккаунта",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Авторизация"
                ],
                "summary": "Выйти из аккаунта",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/api/user/sign_up": {
            "post": {
                "description": "Регистрация нового пользователя",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Авторизация"
                ],
                "summary": "Регистрация",
                "parameters": [
                    {
                        "description": "login and password",
                        "name": "user_credentials",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/schemes.RegisterReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schemes.RegisterResp"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "app.SwaggerUpdateTurnRequest": {
            "type": "object",
            "properties": {
                "turn_type": {
                    "type": "string"
                }
            }
        },
        "ds.Card": {
            "type": "object",
            "required": [
                "description",
                "name",
                "needfood",
                "type"
            ],
            "properties": {
                "description": {
                    "type": "string",
                    "maxLength": 200
                },
                "image_url": {
                    "type": "string"
                },
                "name": {
                    "type": "string",
                    "maxLength": 50
                },
                "needfood": {
                    "type": "integer"
                },
                "type": {
                    "type": "string",
                    "maxLength": 50
                },
                "uuid": {
                    "type": "string"
                }
            }
        },
        "schemes.AddToTurnResp": {
            "type": "object",
            "properties": {
                "card_count": {
                    "type": "integer"
                }
            }
        },
        "schemes.AllCardsResponse": {
            "type": "object",
            "properties": {
                "cards": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/ds.Card"
                    }
                }
            }
        },
        "schemes.AllTurnsResponse": {
            "type": "object",
            "properties": {
                "turns": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/schemes.TurnOutput"
                    }
                }
            }
        },
        "schemes.GetAllCardsResponse": {
            "type": "object",
            "properties": {
                "cards": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/ds.Card"
                    }
                },
                "draft_turn": {
                    "$ref": "#/definitions/schemes.TurnShort"
                }
            }
        },
        "schemes.LoginReq": {
            "type": "object",
            "required": [
                "login",
                "password"
            ],
            "properties": {
                "login": {
                    "type": "string",
                    "maxLength": 30
                },
                "password": {
                    "type": "string",
                    "maxLength": 30
                }
            }
        },
        "schemes.RegisterReq": {
            "type": "object",
            "required": [
                "login",
                "password"
            ],
            "properties": {
                "login": {
                    "type": "string",
                    "maxLength": 30
                },
                "password": {
                    "type": "string",
                    "maxLength": 30
                }
            }
        },
        "schemes.RegisterResp": {
            "type": "object",
            "properties": {
                "ok": {
                    "type": "boolean"
                }
            }
        },
        "schemes.SwaggerLoginResp": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "expires_in": {
                    "type": "integer"
                },
                "token_type": {
                    "type": "string"
                }
            }
        },
        "schemes.TurnOutput": {
            "type": "object",
            "properties": {
                "completion_date": {
                    "type": "string"
                },
                "creation_date": {
                    "type": "string"
                },
                "customer": {
                    "type": "string"
                },
                "formation_date": {
                    "type": "string"
                },
                "moderator": {
                    "type": "string"
                },
                "sending_status": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "takefood": {
                    "type": "integer"
                },
                "turn_type": {
                    "type": "string"
                },
                "uuid": {
                    "type": "string"
                }
            }
        },
        "schemes.TurnResponse": {
            "type": "object",
            "properties": {
                "cards": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/ds.Card"
                    }
                },
                "turn": {
                    "$ref": "#/definitions/schemes.TurnOutput"
                }
            }
        },
        "schemes.TurnShort": {
            "type": "object",
            "properties": {
                "card_count": {
                    "type": "integer"
                },
                "uuid": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "127.0.0.1:8080",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "Electronic notifications",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
