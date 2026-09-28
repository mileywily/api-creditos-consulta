package usecases_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"api-creditos-consulta/internal/core/domain"
	"api-creditos-consulta/internal/core/usecases"
)

type mockFinnFlowClient struct {
	creditos []domain.Credito
	err      error
}

func (m *mockFinnFlowClient) ObtenerCreditosPorRut(ctx context.Context, rut string) ([]domain.Credito, error) {
	return m.creditos, m.err
}

func TestConsultarCreditos_Success(t *testing.T) {
	mockClient := &mockFinnFlowClient{
		creditos: []domain.Credito{
			{NumeroOperacion: "123"},
		},
		err: nil,
	}

	uc := usecases.NewCreditosUseCase(mockClient)
	resp, err := uc.ConsultarCreditos(context.Background(), domain.ConsultaCreditosRequest{Rut: "1-9"})

	assert.NoError(t, err)
	assert.Equal(t, 0, resp.CodRespuesta)
	assert.Len(t, resp.Creditos, 1)
}

func TestConsultarCreditos_Empty(t *testing.T) {
	mockClient := &mockFinnFlowClient{
		creditos: nil, // Should be transformed to empty array
		err:      nil,
	}

	uc := usecases.NewCreditosUseCase(mockClient)
	resp, err := uc.ConsultarCreditos(context.Background(), domain.ConsultaCreditosRequest{Rut: "1-9"})

	assert.NoError(t, err)
	assert.Equal(t, 0, resp.CodRespuesta)
	assert.NotNil(t, resp.Creditos)
	assert.Len(t, resp.Creditos, 0)
}

func TestConsultarCreditos_Error(t *testing.T) {
	mockClient := &mockFinnFlowClient{
		creditos: nil,
		err:      assert.AnError,
	}

	uc := usecases.NewCreditosUseCase(mockClient)
	resp, err := uc.ConsultarCreditos(context.Background(), domain.ConsultaCreditosRequest{Rut: "1-9"})

	assert.NoError(t, err) // UseCase returns error in response object, not as raw error
	assert.Equal(t, 1, resp.CodRespuesta)
	assert.Equal(t, assert.AnError.Error(), resp.DescError)
	assert.NotNil(t, resp.Creditos)
	assert.Len(t, resp.Creditos, 0)
}
