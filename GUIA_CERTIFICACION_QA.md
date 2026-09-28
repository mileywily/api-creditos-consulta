# Guía de Certificación QA: SRV02 - Consulta de Créditos Hipotecarios

Este documento consolida los procedimientos técnicos necesarios para que el equipo de Calidad (QA) y Certificación pueda validar, ejecutar y aprobar el microservicio **SRV02** de integración entre Salesforce y ASICOM (FinnFlow).

---

## 1. Preparación del Entorno (Limpieza de Puertos)
Es imperativo asegurar que los puertos locales (`8080`, `8081` u `8082`) estén libres antes de levantar el aplicativo o correr pruebas para evitar colisiones de tipo `listen tcp: bind` o `ECONNREFUSED`.

Se ha proveído un script automático en la raíz del proyecto (`run_tests.ps1`). Para ejecutar una limpieza manual en **PowerShell**, corra el siguiente bloque:

```powershell
# Definir el puerto que usa la API
$puerto = if ($env:PORT) { $env:PORT } else { "8082" }

# Matar cualquier proceso huérfano en ese puerto
$conexiones = Get-NetTCPConnection -LocalPort $puerto -ErrorAction SilentlyContinue
if ($conexiones) {
    $pids = $conexiones | Select-Object -ExpandProperty OwningProcess -Unique
    foreach ($id in $pids) { Stop-Process -Id $id -Force -ErrorAction SilentlyContinue }
    Write-Host "Puerto $puerto liberado." -ForegroundColor Green
}
```

---

## 2. Ejecución de Pruebas Automatizadas (Código)
El proyecto cuenta con un 100% de cobertura funcional validada bajo **Go 1.27+**.

### 2.1 Pruebas Unitarias y de Integración
Limpian la caché, compilan y corren la suite de pruebas que incluyen Mocks aislados y Servidores HTTP Simulados (`httptest`) para el cliente de FinnFlow.
```powershell
go clean -testcache
go test -v ./internal/...
```

### 2.2 Pruebas End-to-End (E2E) Automatizadas
Ejecutan el ecosistema completo. Levantan el router de Gin, inyectan el cliente hacia un Mock Server local en código y disparan peticiones HTTP completas emulando a Salesforce.
```powershell
go test -v ./test/e2e/...
```

### 2.3 Reporte de Cobertura (Coverage)
Para garantizar >95% de cobertura estipulada para certificación:
```powershell
go test ./internal/... ./test/e2e/... -coverprofile coverfile_out
go tool cover -func coverfile_out
```

---

## 3. Pruebas Manuales y de Caja Negra (Postman)
Para certificar el microservicio manualmente visualizando los Requests y Responses exactos exigidos por el **DAD (Documento de Análisis de Requerimientos)**.

### Paso A: Levantar el Backend Simulado (Mock FinnFlow)
Abra una consola PowerShell y levante el servidor que simula ser ASICOM (correrá en el puerto 9090):
```powershell
cd C:\HEXAGONAL-SRV02
go run test/mock_finnflow/main.go
```

### Paso B: Levantar el Microservicio SRV02
Abra una **SEGUNDA** consola PowerShell, inyecte las variables de entorno sin hardcoding y corra el API:
```powershell
cd C:\HEXAGONAL-SRV02
$env:PORT="8082"
$env:FINNFLOW_URL="http://localhost:9090"
$env:FINNFLOW_KEY="qa-key"
$env:FINNFLOW_SECRET="qa-secret"

go run cmd/api/main.go
```

### Paso C: Ejecutar la Colección de Postman
Importe en Postman el archivo **`SRV02_Postman_Collection.json`** ubicado en la raíz del proyecto. Encontrará 4 casos listos para testear (haga clic en *Send* en cada uno):

| Caso de Uso | Payload (Body) | Headers Necesarios | Resultado Esperado (HTTP) |
| :--- | :--- | :--- | :--- |
| **1. Happy Path** | `{"Rut": "25671646-1"}` | `Authorization: Bearer <token>` | **200 OK** (Devuelve Array de Créditos) |
| **2. Falta de Token** | `{"Rut": "25671646-1"}` | *(Ninguno)* | **401 Unauthorized** (Token is not valid) |
| **3. Campo Faltante** | `{}` | `Authorization: Bearer <token>` | **500 Internal Error** (Legacy 400 - missing_field) |
| **4. Rut Inválido** | `{"Rut": "RUT-INVALIDO"}` | `Authorization: Bearer <token>` | **500 Internal Error** (Legacy 402 - invalid_format) |

---

## 4. Trazabilidad de Reglas de Negocio (Compliance)
Para certificar la **RN03**, el equipo de QA validará que el servicio **no descarta ni filtra créditos** por elegibilidad, sino que mapea y entrega fielmente todo el catálogo asociado al RUT, delegando correctamente la lógica de decisión hacia Salesforce.
