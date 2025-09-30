package handler

import (
	"net/http"

	"my-tdd-app/backend/openapi"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ApiHandler struct {
	db *gorm.DB
}

// コンストラクタ
func NewApiHandler(db *gorm.DB) *ApiHandler {
	return &ApiHandler{db: db}
}

// GetUsers は oapi-codegen で生成された ServerInterface の一部を実装
func (h *ApiHandler) GetUsers(c echo.Context) error {
	var users []openapi.User // 上記のテストコードで使ったUserモデルと同じもの

	if err := h.db.Find(&users).Error; err != nil {
		// エラーハンドリング
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "database error"})
	}

	// 成功したらユーザー一覧をJSONで返す
	return c.JSON(http.StatusOK, users)
}

// ... PostUsersなどの他のメソッドもここに実装していく ...
