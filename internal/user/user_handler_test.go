package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alanloffler/go-calth-api/internal/common/response"
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestUserHandler_GetAll_NotFound(t *testing.T) {
	router := setupTestRouter()

	mockRepo := &MockUserRepository{}
	mockRepo.On("GetAll", mock.Anything).Return([]sqlc.User{}, assert.AnError)

	handler := &UserHandler{repo: mockRepo}
	router.GET("/users", handler.GetAll)

	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var resp response.ApiResponse[any]
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "Usuarios no encontrados", resp.Message)
}

func TestUserHandler_GetAll_Success(t *testing.T) {
	router := setupTestRouter()

	users := []sqlc.User{}
	mockRepo := &MockUserRepository{}
	mockRepo.On("GetAll", mock.Anything).Return(users, nil)

	handler := &UserHandler{repo: mockRepo}
	router.GET("/users", handler.GetAll)

	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var resp response.ApiResponse[any]
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Usuarios encontrados", resp.Message)
	assert.NotNil(t, resp.Data)
}

func TestUserHandler_Update_InvalidID(t *testing.T) {
	router := setupTestRouter()

	handler := &UserHandler{}
	router.PATCH("/users/:id", handler.Update)

	req, _ := http.NewRequest("PATCH", "/users/invalid-uuid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var resp response.ApiResponse[any]
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_Update_ValidationError(t *testing.T) {
	router := setupTestRouter()

	mockRepo := &MockUserRepository{}
	mockRepo.On("Update", mock.Anything, mock.Anything).Return(sqlc.User{}, assert.AnError)

	handler := &UserHandler{repo: mockRepo}
	router.PATCH("/users/:id", handler.Update)

	body := `{"user": {}}`
	req, _ := http.NewRequest("PATCH", "/users/550e8400-e29b-41d4-a716-446655440000", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var resp response.ApiResponse[any]
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUserHandler_Delete_InvalidID(t *testing.T) {
	router := setupTestRouter()

	handler := &UserHandler{}
	router.DELETE("/users/:id", handler.Delete)

	req, _ := http.NewRequest("DELETE", "/users/invalid-uuid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var resp response.ApiResponse[any]
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
