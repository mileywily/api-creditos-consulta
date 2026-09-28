package domain

type ConsultaCreditosRequest struct {
	Rut string
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
