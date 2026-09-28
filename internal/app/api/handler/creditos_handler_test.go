package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"

	"api-creditos-consulta/internal/app/api/handler"
	"api-creditos-consulta/internal/core/domain"
)

type mockCreditosUseCase struct {
	resp domain.ConsultaCreditosResponse
	err  error
}

func (m *mockCreditosUseCase) ConsultarCreditos(ctx context.Context, req domain.ConsultaCreditosRequest) (domain.ConsultaCreditosResponse, error) {
	return m.resp, m.err
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("customRutPattern", func(fl validator.FieldLevel) bool {
			matched, _ := regexp.MatchString(`^[0-9]{1,2}\.[0-9]{3}\.[0-9]{3}-[0-9kK]$`, fl.Field().String())
			if !matched {
				matched, _ = regexp.MatchString(`^[0-9]+-[0-9kK]$`, fl.Field().String())
			}
			return matched
		})
	}
	return router
}

func TestConsultarCreditosHandler_Success(t *testing.T) {
	uc := &mockCreditosUseCase{
		resp: domain.ConsultaCreditosResponse{
			CodRespuesta: 0,
			DescError:    "",
			Creditos:     []domain.Credito{},
		},
		err: nil,
	}
	
	h := handler.NewCreditosHandler(uc)
	router := setupTestRouter()
	router.POST("/api/creditos/consulta/", h.ConsultarCreditos)

	body := []byte(`{"Rut": "1-9"}`)
	req, _ := http.NewRequest("POST", "/api/creditos/consulta/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer mock-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["codRespuesta"])
	assert.NotNil(t, response["creditos"])
}

func TestConsultarCreditosHandler_NoToken(t *testing.T) {
	h := handler.NewCreditosHandler(&mockCreditosUseCase{})
	router := setupTestRouter()
	router.POST("/api/creditos/consulta/", h.ConsultarCreditos)

	body := []byte(`{"Rut": "1-9"}`)
	req, _ := http.NewRequest("POST", "/api/creditos/consulta/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	// No agregamos el header Authorization
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestConsultarCreditosHandler_BindError(t *testing.T) {
	h := handler.NewCreditosHandler(&mockCreditosUseCase{})
	router := setupTestRouter()
	router.POST("/api/creditos/consulta/", h.ConsultarCreditos)

	// Missing Rut field in request
	body := []byte(`{"invalid": "data"}`)
	req, _ := http.NewRequest("POST", "/api/creditos/consulta/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer mock-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code) // Replicated legacy bug changes bad request to 500
}

func TestConsultarCreditosHandler_UseCaseError(t *testing.T) {
	uc := &mockCreditosUseCase{
		resp: domain.ConsultaCreditosResponse{},
		err:  assert.AnError,
	}
	h := handler.NewCreditosHandler(uc)
	router := setupTestRouter()
	router.POST("/api/creditos/consulta/", h.ConsultarCreditos)

	body := []byte(`{"Rut": "1-9"}`)
	req, _ := http.NewRequest("POST", "/api/creditos/consulta/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer mock-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
