# Especificación y Caso de Uso: SRV02 - Consulta de Créditos Hipotecarios Asociados

## 1. Información General del Servicio (Contexto DAD)
* **ID Servicio:** SRV02
* **Nombre:** Consulta de créditos hipotecarios asociados
* **Acción:** Crear (Nuevo Servicio)
* **Objetivo Principal:** Permitir que Salesforce consulte todos los créditos hipotecarios vigentes asociados a un cliente específico. Este servicio es el **punto de entrada de la Fase de Elegibilidad** para una renegociación.
* **Paradigma de Responsabilidad:** FinnFlow actúa únicamente como proveedor de datos "crudos" (core). Salesforce será quien ejecute las reglas de negocio para determinar si un crédito califica o no para renegociación (Regla RN03).

---

## 2. Caso de Uso (CU)

### **CU: Buscar y Listar Operaciones del Cliente**
* **Actor Principal:** Sistema Salesforce (Ejecutivo del Banco).
* **Actor Secundario:** Sistema ASICOM (FinnFlow).
* **Precondiciones:**
  1. El cliente se encuentra autenticado o identificado en la sucursal/plataforma con su RUT.
  2. Salesforce tiene un token válido de autorización (`Bearer`).
* **Flujo Principal:**
  1. El Ejecutivo (a través de Salesforce) inicia el flujo de Renegociación e ingresa el RUT del cliente.
  2. Salesforce envía un request `POST` al endpoint `/api/creditos/consulta/` incluyendo el RUT.
  3. FinnFlow recibe el RUT y consulta en la Base de Datos Core todas las operaciones asociadas a ese cliente.
  4. FinnFlow arma el payload cruzando las tablas de: Créditos, Participantes (codeudores/avales), Propiedades en garantía y Políticas (días de mora).
  5. FinnFlow retorna la lista completa (`creditos[]`) y un `codRespuesta = 0` (Éxito).
  6. Salesforce recibe los datos, ejecuta su motor interno de elegibilidad y le muestra al Ejecutivo solo los créditos aptos para renegociar.
* **Excepciones:**
  * **El cliente no tiene créditos:** FinnFlow retorna un array vacío de créditos, sin lanzar un error 500 (Http 200 OK con `creditos: []`).

---

## 3. Regla de Negocio Asociada
* **RN03 (Responsable de elegibilidad):** La elegibilidad **no** se determina en ASICOM. Salesforce aplica las reglas de negocio utilizando la información entregada por este servicio. FinnFlow debe retornar la data sin aplicar filtros de estado que impidan visualizar la realidad del crédito.

---

## 4. Especificación Técnica y Contrato API

### **Endpoint:**
`POST /api/creditos/consulta/`

### **Headers:**
`Authorization: Bearer {{token}}`
`Content-Type: application/json`

### **Request Body (Entrada):**
```json
{
  "Rut": "25671646-1"
}
```

### **Response Body (Salida / DAD):**
El servicio debe devolver la siguiente jerarquía de datos obligatoria:

| Bloque / Objeto | Campo | Tipo | Requerido | Descripción |
| :--- | :--- | :--- | :--- | :--- |
| **Raíz** | `codRespuesta` | INT | Sí | Código interno de respuesta (Ej. 0 = Éxito). |
| **Raíz** | `descError` | VARCHAR | Sí | Descripción de error en caso de fallo, vacío si hay éxito. |
| **creditos[]** | `NumeroOperacion` | VARCHAR | Sí | Identificador de la operación actual. |
| **creditos[]** | `NumeroOperacionOriginal` | VARCHAR | Sí | Identificador de la operación original asociada. |
| **creditos[]** | `Estado` | VARCHAR | Sí | Estado del crédito (Ej. "Vigente"). |
| **creditos[]** | `FechaActivacion` | DATE | Sí | Fecha de otorgamiento del crédito. |
| **creditos[]** | `Producto` | INT | Sí | Código de catálogo Producto. |
| **creditos[]** | `DescripcionProducto` | VARCHAR | Sí | Glosa del producto. |
| **creditos[]** | `Objetivo` | INT | Sí | Código de catálogo Objetivo. |
| **creditos[]** | `DescripcionObjetivo` | VARCHAR | Sí | Glosa del objetivo. |
| **creditos[]** | `Destino` | INT | Sí | Código de catálogo Destino. |
| **creditos[]** | `DescripcionDestino` | VARCHAR | Sí | Glosa del destino. |
| **creditos[]** | `TipoGarantia` | INT | Sí | Tipo de garantía asociado (Catálogo compartido). |
| **creditos[]** | `ValorGarantia` | DECIMAL | Sí | Valor de tasación/garantía. |
| **creditos[]** | `SubsidioOriginal` | INT | Sí | Subsidio original. |
| **participantes[]** | `TipoParticipacion` | VARCHAR | Sí | Ej. Deudor, Codeudor, Aval. |
| **participantes[]** | `Rut`, `Nombre`, `ApellidoPaterno`, `ApellidoMaterno`, `FechaNacimiento`, `Mail` | Varios | Sí | Datos personales del participante. |
| **propiedades[]** | `TipoInmueble`, `Antiguedad`, `Direccion`, `Numero`, `Depto`, `Comuna`, `ValorPropiedad` | Varios | Sí | Datos del activo en garantía (Depto es opcional). |
| **politicas** | `DiasMora`, `DividendosEnMora`, `DividendosPagados`, `CampanaHipotecaria` | INT | Sí | Bloque clave para evaluar elegibilidad (mora). |

### **Mock (Ejemplo JSON Completo):**
```json
{
  "codRespuesta": 0,
  "creditos": [
    {
      "NumeroOperacion": "123456789",
      "numeroOperacionOriginal": "123456789",
      "Estado": "Vigente",
      "FechaActivacion": "2019-03-10",
      "Producto": 1,
      "DescripcionProducto": "Mutuo Recursos Propios",
      "Objetivo": 1,
      "DescripcionObjetivo": "Vivienda",
      "Destino": 1,
      "DescripcionDestino": "Compraventa",
      "TipoGarantia": 1,
      "ValorGarantia": 2300.0000,
      "SubsidioOriginal": 1,
      "participantes": [
        {
          "TipoParticipacion": "1",
          "Rut": "18015051-K",
          "Nombre": "Pedro",
          "ApellidoPaterno": "Perez",
          "ApellidoMaterno": "Perez",
          "FechaNacimiento": "1992-07-04",
          "Mail": "micorreo@gmail.com"
        }
      ],
      "propiedades": [
        {
          "TipoInmueble": "0",
          "Antiguedad": 2,
          "Direccion": "Av. Principal",
          "Numero": "1234",
          "Depto": null,
          "Comuna": "6400",
          "ValorPropiedad": 2300.0000
        }
      ],
      "politicas": {
        "DiasMora": 45,
        "DividendosEnMora": 6,
        "DividendosPagados": 5,
        "CampanaHipotecaria": 0
      }
    }
  ],
  "descError": ""
}
```
