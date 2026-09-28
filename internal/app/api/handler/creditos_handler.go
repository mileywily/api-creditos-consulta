package handler

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"api-creditos-consulta/internal/app/api/dto"
	"api-creditos-consulta/internal/core/domain"
	"api-creditos-consulta/internal/core/ports"
)

type CreditosHandler struct {
	useCase ports.CreditosUseCase
}

func NewCreditosHandler(uc ports.CreditosUseCase) *CreditosHandler {
	return &CreditosHandler{
		useCase: uc,
	}
}

func (h *CreditosHandler) ConsultarCreditos(c *gin.Context) {
	// Validación de existencia de Token (Autorización)
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		h.handleDomainError(c, fmt.Errorf("token_not_valid"))
		return
	}

	// Replicando log exacto del controller legado
	bodyBytes, _ := io.ReadAll(c.Request.Body)
	rawRequest := string(bodyBytes)
	slog.Info("Received consulta request: "+rawRequest, "logger", "cl.bancofalabella.mortgage.consulta.controller.ConsultaCreditosController", "thread", "http-nio-8080-exec-1")

	// Restaurar el body para Gin
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var req dto.ConsultaCreditosRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleValidationError(c, err)
		return
	}

	domainReq := domain.ConsultaCreditosRequest{
		Rut: req.Rut,
	}

	resp, err := h.useCase.ConsultarCreditos(c.Request.Context(), domainReq)
	if err != nil {
		h.handleDomainError(c, err)
		return
	}

	dtoResp := dto.MapDomainToResponse(resp)

	slog.Info(fmt.Sprintf("Consulta processed for RUT: %s with status: %d and error: %s", req.Rut, dtoResp.CodRespuesta, dtoResp.DescError), "logger", "cl.bancofalabella.mortgage.consulta.controller.ConsultaCreditosController", "thread", "http-nio-8080-exec-1")

	c.JSON(http.StatusOK, dtoResp)
}

func (h *CreditosHandler) handleValidationError(c *gin.Context, err error) {
	// HttpMessageNotReadableException equivalent
	if err.Error() == "EOF" || regexp.MustCompile(`cannot unmarshal`).MatchString(err.Error()) {
		c.JSON(http.StatusBadRequest, gin.H{
			"codRespuesta": 1,
			"descError":    "Request body is required",
			"creditos":     []dto.Credito{},
		})
		return
	}

	var fieldName string
	var code, message string
	var status int

	if errs, ok := err.(validator.ValidationErrors); ok && len(errs) > 0 {
		firstErr := errs[0]
		fieldName = firstErr.Field()
		tag := firstErr.Tag()

		if tag == "customRutPattern" || tag == "customFechaPattern" {
			status = http.StatusPaymentRequired // 402
			code = "invalid_format"
			message = fmt.Sprintf("%s with invalid format", fieldName)
		} else {
			status = http.StatusBadRequest // 400
			code = "missing_field"
			message = fmt.Sprintf("Missing %s field", fieldName)
		}
	} else {
		status = http.StatusBadRequest
		code = "bad_request"
		message = "Invalid request format"
	}

	// REPLICATE LEGACY BUG: Squash to 500 with the exact String that Java produced
	statusText := ""
	switch status {
	case 400:
		statusText = "400 BAD_REQUEST"
	case 402:
		statusText = "402 PAYMENT_REQUIRED"
	default:
		statusText = fmt.Sprintf("%d ERROR", status)
	}

	legacyMensaje := fmt.Sprintf(`FinnFlow consulta failed: %s - {"code":"%s","message":"%s"}`, statusText, code, message)

	c.JSON(http.StatusInternalServerError, dto.ConsultaCreditosResponse{
		CodRespuesta: 500, // Usando 500 en codRespuesta o manteniendo 1 según DAD, pero emulamos el error de java
		DescError:    legacyMensaje,
		Creditos:     []dto.Credito{},
	})
}

func (h *CreditosHandler) handleDomainError(c *gin.Context, err error) {
	msg := err.Error()

	if msg == "token_not_valid" {
		c.JSON(http.StatusUnauthorized, dto.ConsultaCreditosResponse{
			CodRespuesta: 401,
			DescError:    "Token is not valid",
			Creditos:     []dto.Credito{},
		})
		return
	}

	c.JSON(http.StatusInternalServerError, dto.ConsultaCreditosResponse{
		CodRespuesta: 500,
		DescError:    "Internal server error: " + msg,
		Creditos:     []dto.Credito{},
	})
}

