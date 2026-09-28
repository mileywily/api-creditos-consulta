package ports

import (
	"context"
	"api-creditos-consulta/internal/core/domain"
)

type FinnFlowClient interface {
	ObtenerCreditosPorRut(ctx context.Context, rut string) ([]domain.Credito, error)
}
