package main

import (
	"my-tdd-app/backend/handler"
	"my-tdd-app/backend/openapi"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=db user=user password=password dbname=mydatabase port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}

	// 3. APIサーバーのハンドラーをセットアップ
	//    (これはAPI本体の実装に依存する)
	//    handler := NewApiHandler(db) // 仮のAPIハンドラー
	//    api.server = httptest.NewServer(handler) // テストサーバーを起動
	e := echo.New()
	apiHandler := handler.NewApiHandler(db)
	openapi.RegisterHandlers(e, apiHandler)

	// 8080でEchoサーバーを起動
	e.Logger.Fatal(e.Start(":8080"))
}
