package e2e_test

import (
	"bytes"
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
	"api-creditos-consulta/internal/core/usecases"
	"api-creditos-consulta/internal/infra/client"
)

// setupE2ERouter inicializa la aplicación real tal como lo hace main.go
func setupE2ERouter(finnflowURL string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Registrar validaciones personalizadas
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("customRutPattern", func(fl validator.FieldLevel) bool {
			matched, _ := regexp.MatchString(`^[0-9]{1,2}\.[0-9]{3}\.[0-9]{3}-[0-9kK]$`, fl.Field().String())
			if !matched {
				matched, _ = regexp.MatchString(`^[0-9]+-[0-9kK]$`, fl.Field().String())
			}
			return matched
		})
	}

	// Inyectar dependencias reales
	finnClient := client.NewFinnFlowClient(finnflowURL, "e2e-key", "e2e-secret")
	uc := usecases.NewCreditosUseCase(finnClient)
	h := handler.NewCreditosHandler(uc)

	router.POST("/api/creditos/consulta/", h.ConsultarCreditos)
	return router
}

func TestE2E_ConsultaCreditos_HappyPath(t *testing.T) {
	// 1. Levantar servidor Mock de FinnFlow (Backend real simulado)
	mockFinnFlow := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "e2e-key", r.Header.Get("X-Client-Id"))
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"data": [
				{
					"numeroOperacionOriginal": "OP-12345",
					"participantes": [{"rut": "25671646-1"}],
					"propiedades": [{"direccion": "Calle E2E"}],
					"politicas": {"diasMora": 10}
				}
			]
		}`))
	}))
	defer mockFinnFlow.Close()

	// 2. Levantar el Router Gin real apuntando al Mock de FinnFlow
	router := setupE2ERouter(mockFinnFlow.URL)

	// 3. Disparar petición E2E (Simulando a Salesforce o Postman)
	reqBody := []byte(`{"Rut": "25671646-1"}`)
	req, _ := http.NewRequest("POST", "/api/creditos/consulta/", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 4. Validar resultado de Extremo a Extremo
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, float64(0), response["codRespuesta"])
	creditos := response["creditos"].([]interface{})
	assert.Len(t, creditos, 1)

	primerCredito := creditos[0].(map[string]interface{})
	assert.Equal(t, "OP-12345", primerCredito["numeroOperacionOriginal"])
	
	politicas := primerCredito["politicas"].(map[string]interface{})
	assert.Equal(t, float64(10), politicas["DiasMora"]) // Verifica que el mapeo DTO interno funcionó
}

func TestE2E_ConsultaCreditos_FaltaToken(t *testing.T) {
	router := setupE2ERouter("http://no-importa-porque-no-llegara.com")

	reqBody := []byte(`{"Rut": "25671646-1"}`)
	req, _ := http.NewRequest("POST", "/api/creditos/consulta/", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	// SIN HEADER DE AUTORIZACION

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code) // E2E de seguridad validado
}

func TestE2E_ConsultaCreditos_RutInvalidoLegacyBug(t *testing.T) {
	router := setupE2ERouter("http://no-importa-porque-falla-antes.com")

	reqBody := []byte(`{"Rut": "RUT-INVALIDO"}`)
	req, _ := http.NewRequest("POST", "/api/creditos/consulta/", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	
	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	
	// Validamos que todo el stack convirtió el error de validación 402 en un 500 legacy.
	assert.Contains(t, response["descError"], "402 PAYMENT_REQUIRED")
	assert.Contains(t, response["descError"], "invalid_format")
}
