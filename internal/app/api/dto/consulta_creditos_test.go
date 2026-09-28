package dto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"api-creditos-consulta/internal/app/api/dto"
	"api-creditos-consulta/internal/core/domain"
)

func TestMapDomainToResponse(t *testing.T) {
	d := domain.ConsultaCreditosResponse{
		CodRespuesta: 0,
		DescError:    "",
		Creditos: []domain.Credito{
			{
				NumeroOperacion: "123",
				Participantes: []domain.Participante{
					{Rut: "1-9"},
				},
				Propiedades: []domain.Propiedad{
					{Direccion: "Calle 1"},
				},
				Politicas: domain.Politicas{DiasMora: 10},
			},
		},
	}

	resp := dto.MapDomainToResponse(d)
	assert.Equal(t, 0, resp.CodRespuesta)
	assert.Len(t, resp.Creditos, 1)
	assert.Equal(t, "123", resp.Creditos[0].NumeroOperacion)
	assert.Len(t, resp.Creditos[0].Participantes, 1)
	assert.Equal(t, "1-9", resp.Creditos[0].Participantes[0].Rut)
	assert.Len(t, resp.Creditos[0].Propiedades, 1)
	assert.Equal(t, "Calle 1", resp.Creditos[0].Propiedades[0].Direccion)
	assert.Equal(t, 10, resp.Creditos[0].Politicas.DiasMora)
}

func TestMapDomainToResponse_Empty(t *testing.T) {
	d := domain.ConsultaCreditosResponse{
		CodRespuesta: 0,
		DescError:    "",
		Creditos:     nil,
	}

	resp := dto.MapDomainToResponse(d)
	assert.Equal(t, 0, resp.CodRespuesta)
	assert.NotNil(t, resp.Creditos)
	assert.Len(t, resp.Creditos, 0)
}
