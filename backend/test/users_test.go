package test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"my-tdd-app/backend/handler"
	"my-tdd-app/backend/openapi"
	"net/http"
	"net/http/httptest"

	"github.com/cucumber/godog"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// User モデル（API本体のモデルと共用するのが理想）
type User struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique"`
}

// テストシナリオの状態を管理する構造体
type apiFeature struct {
	db       *gorm.DB
	server   *httptest.Server
	resp     *http.Response
	respBody []byte
}

// --- ステップ定義ここから ---

func (a *apiFeature) データベースに以下のユーザーが存在する(table *godog.Table) error {
	// テーブルのヘッダーを除いた各行をループ
	for i := 1; i < len(table.Rows); i++ {
		name := table.Rows[i].Cells[0].Value
		user := User{Name: name}
		// DBにユーザーを作成
		if err := a.db.Create(&user).Error; err != nil {
			return err
		}
	}
	return nil
}

func (a *apiFeature) usersにGETリクエストを送信する() error {
	// デバッグ用：サーバーURLを出力
	fmt.Printf("Requesting URL: %s/users\n", a.server.URL)
	// テストサーバーのURLに対してリクエスト
	resp, err := http.Get(a.server.URL + "/users")
	if err != nil {
		return err
	}
	a.resp = resp
	// レスポンスボディを読み取る（後で使うため）
	// ... (ボディ読み取り処理) ...
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	a.respBody = body
	return nil
}

func (a *apiFeature) レスポンスのステータスコードはXであるべき(expectedStatus int) error {
	if a.resp.StatusCode != expectedStatus {
		return fmt.Errorf("expected status %d, but got %d", expectedStatus, a.resp.StatusCode)
	}
	return nil
}

// ... 他のThenステップも同様に実装 ...
func (a *apiFeature) レスポンスのボディにはX人のユーザーが含まれているべき(expectedCount int) error {
	// JSONレスポンスを解析してユーザー数を確認
	var users []User
	if err := json.Unmarshal(a.respBody, &users); err != nil {
		return fmt.Errorf("failed to parse JSON response: %v", err)
	}

	if len(users) != expectedCount {
		return fmt.Errorf("expected %d users, but got %d", expectedCount, len(users))
	}
	return nil
}

func (a *apiFeature) ユーザー名Xがレスポンスに含まれているべき(arg1 string) error {
	// レスポンスのボディに特定のユーザー名が含まれているか確認
	var users []User
	if err := json.Unmarshal(a.respBody, &users); err != nil {
		return fmt.Errorf("failed to parse JSON response: %v", err)
	}

	for _, user := range users {
		if user.Name == arg1 {
			return nil
		}
	}
	return fmt.Errorf("user %q not found in response", arg1)
}

// --- ステップ定義ここまで ---

// InitializeScenarioでステップを登録する
func InitializeScenario(ctx *godog.ScenarioContext) {
	// apiFeature構造体を初期化
	api := &apiFeature{}

	// Scenario開始前に実行されるフック
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		// --- ここが重要！テスト用のDBセットアップ ---
		// 1. テスト用DBに接続
		//    (実際のDBとは別にするか、トランザクションを使うのが理想)
		dsn := "host=db user=user password=password dbname=mydatabase port=5432 sslmode=disable"
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return ctx, err
		}
		api.db = db

		// 2. テーブルをクリーンな状態にする
		db.Migrator().DropTable(&User{})
		db.AutoMigrate(&User{})

		// 3. APIサーバーのハンドラーをセットアップ
		//    (これはAPI本体の実装に依存する)
		//    handler := NewApiHandler(db) // 仮のAPIハンドラー
		//    api.server = httptest.NewServer(handler) // テストサーバーを起動
		e := echo.New()
		apiHandler := handler.NewApiHandler(db)
		openapi.RegisterHandlers(e, apiHandler)
		api.server = httptest.NewServer(e)
		// mux := http.NewServeMux()
		// mux.HandleFunc("/users", handler.GetUsers())
		// defer api.server.Close()

		return ctx, nil
	})

	// Scenario終了後に実行されるフック
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		// テストサーバーを閉じる
		if api.server != nil {
			api.server.Close()
		}
		return ctx, nil
	})

	// ステップと関数を結びつける
	ctx.Step(`^データベースに以下のユーザーが存在する:$`, api.データベースに以下のユーザーが存在する)
	ctx.Step(`^"/users" にGETリクエストを送信する$`, api.usersにGETリクエストを送信する)
	ctx.Step(`^レスポンスのステータスコードは(\d+)であるべき$`, api.レスポンスのステータスコードはXであるべき)
	ctx.Step(`^レスポンスのボディには(\d+)人のユーザーが含まれているべき$`, api.レスポンスのボディにはX人のユーザーが含まれているべき)
	ctx.Step(`^ユーザー名 "([^"]*)" がレスポンスに含まれているべき$`, api.ユーザー名Xがレスポンスに含まれているべき)
	// ... 他のステップも登録 ...
}
