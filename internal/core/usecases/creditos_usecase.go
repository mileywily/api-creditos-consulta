package usecases

import (
	"context"

	"api-creditos-consulta/internal/core/domain"
	"api-creditos-consulta/internal/core/ports"
)

type creditosUseCase struct {
	finnflowClient ports.FinnFlowClient
}

func NewCreditosUseCase(client ports.FinnFlowClient) ports.CreditosUseCase {
	return &creditosUseCase{
		finnflowClient: client,
	}
}

func (uc *creditosUseCase) ConsultarCreditos(ctx context.Context, req domain.ConsultaCreditosRequest) (domain.ConsultaCreditosResponse, error) {
	// The business rule says if no credits are found, return an empty array with success, not 500.
	creditos, err := uc.finnflowClient.ObtenerCreditosPorRut(ctx, req.Rut)
	if err != nil {
		return domain.ConsultaCreditosResponse{
			CodRespuesta: 1, // error code
			DescError:    err.Error(),
			Creditos:     []domain.Credito{},
		}, nil // We still return 200 OK conceptually, handled at HTTP layer, or return error to trigger 500 depending on legacy behavior. Let's return error to handler.
	}

	if creditos == nil {
		creditos = []domain.Credito{}
	}

	return domain.ConsultaCreditosResponse{
		CodRespuesta: 0,
		DescError:    "",
		Creditos:     creditos,
	}, nil
}
