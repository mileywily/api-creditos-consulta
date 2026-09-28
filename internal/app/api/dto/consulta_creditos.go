package dto

import "api-creditos-consulta/internal/core/domain"

type ConsultaCreditosRequest struct {
	Rut string `json:"Rut" binding:"required,customRutPattern"`
}

type Participante struct {
	TipoParticipacion string `json:"TipoParticipacion"`
	Rut               string `json:"Rut"`
	Nombre            string `json:"Nombre"`
	ApellidoPaterno   string `json:"ApellidoPaterno"`
	ApellidoMaterno   string `json:"ApellidoMaterno"`
	FechaNacimiento   string `json:"FechaNacimiento"`
	Mail              string `json:"Mail"`
}

type Propiedad struct {
	TipoInmueble   string  `json:"TipoInmueble"`
	Antiguedad     int     `json:"Antiguedad"`
	Direccion      string  `json:"Direccion"`
	Numero         string  `json:"Numero"`
	Depto          *string `json:"Depto,omitempty"`
	Comuna         string  `json:"Comuna"`
	ValorPropiedad float64 `json:"ValorPropiedad"`
}

type Politicas struct {
	DiasMora           int `json:"DiasMora"`
	DividendosEnMora   int `json:"DividendosEnMora"`
	DividendosPagados  int `json:"DividendosPagados"`
	CampanaHipotecaria int `json:"CampanaHipotecaria"`
}

type Credito struct {
	NumeroOperacion         string         `json:"NumeroOperacion"`
	NumeroOperacionOriginal string         `json:"numeroOperacionOriginal"`
	Estado                  string         `json:"Estado"`
	FechaActivacion         string         `json:"FechaActivacion"`
	Producto                int            `json:"Producto"`
	DescripcionProducto     string         `json:"DescripcionProducto"`
	Objetivo                int            `json:"Objetivo"`
	DescripcionObjetivo     string         `json:"DescripcionObjetivo"`
	Destino                 int            `json:"Destino"`
	DescripcionDestino      string         `json:"DescripcionDestino"`
	TipoGarantia            int            `json:"TipoGarantia"`
	ValorGarantia           float64        `json:"ValorGarantia"`
	SubsidioOriginal        int            `json:"SubsidioOriginal"`
	Participantes           []Participante `json:"participantes"`
	Propiedades             []Propiedad    `json:"propiedades"`
	Politicas               Politicas      `json:"politicas"`
}

type ConsultaCreditosResponse struct {
	CodRespuesta int       `json:"codRespuesta"`
	DescError    string    `json:"descError"`
	Creditos     []Credito `json:"creditos"`
}

// MapDomainToResponse maps the domain response to the DTO response
func MapDomainToResponse(d domain.ConsultaCreditosResponse) ConsultaCreditosResponse {
	// For simplicity in this migration, we map identically since they match the DAD contract
	// In a real scenario, we might map field by field. Here we use a JSON marshal/unmarshal or manual mapping.
	var creditos []Credito
	for _, c := range d.Creditos {
		var participantes []Participante
		for _, p := range c.Participantes {
			participantes = append(participantes, Participante{
				TipoParticipacion: p.TipoParticipacion,
				Rut:               p.Rut,
				Nombre:            p.Nombre,
				ApellidoPaterno:   p.ApellidoPaterno,
				ApellidoMaterno:   p.ApellidoMaterno,
				FechaNacimiento:   p.FechaNacimiento,
				Mail:              p.Mail,
			})
		}
		var propiedades []Propiedad
		for _, p := range c.Propiedades {
			propiedades = append(propiedades, Propiedad{
				TipoInmueble:   p.TipoInmueble,
				Antiguedad:     p.Antiguedad,
				Direccion:      p.Direccion,
				Numero:         p.Numero,
				Depto:          p.Depto,
				Comuna:         p.Comuna,
				ValorPropiedad: p.ValorPropiedad,
			})
		}
		creditos = append(creditos, Credito{
			NumeroOperacion:         c.NumeroOperacion,
			NumeroOperacionOriginal: c.NumeroOperacionOriginal,
			Estado:                  c.Estado,
			FechaActivacion:         c.FechaActivacion,
			Producto:                c.Producto,
			DescripcionProducto:     c.DescripcionProducto,
			Objetivo:                c.Objetivo,
			DescripcionObjetivo:     c.DescripcionObjetivo,
			Destino:                 c.Destino,
			DescripcionDestino:      c.DescripcionDestino,
			TipoGarantia:            c.TipoGarantia,
			ValorGarantia:           c.ValorGarantia,
			SubsidioOriginal:        c.SubsidioOriginal,
			Participantes:           participantes,
			Propiedades:             propiedades,
			Politicas: Politicas{
				DiasMora:           c.Politicas.DiasMora,
				DividendosEnMora:   c.Politicas.DividendosEnMora,
				DividendosPagados:  c.Politicas.DividendosPagados,
				CampanaHipotecaria: c.Politicas.CampanaHipotecaria,
			},
		})
	}
	
	if creditos == nil {
		creditos = []Credito{} // Ensure it's empty array, not null
	}

	return ConsultaCreditosResponse{
		CodRespuesta: d.CodRespuesta,
		DescError:    d.DescError,
		Creditos:     creditos,
	}
}
