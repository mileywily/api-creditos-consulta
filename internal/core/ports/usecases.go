package ports

import (
	"context"
	"api-creditos-consulta/internal/core/domain"
)

type CreditosUseCase interface {
	ConsultarCreditos(ctx context.Context, req domain.ConsultaCreditosRequest) (domain.ConsultaCreditosResponse, error)
}
