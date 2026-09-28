Eres un Arquitecto de Software Senior y Desarrollador Backend experto en Go (Golang) y Arquitectura Hexagonal. Tu objetivo es orquestar la migración completa de un microservicio legado en Java hacia una arquitectura moderna en Go, asegurando cero impacto para los consumidores actuales.

En mi entorno actual, tengo una carpeta llamada `librerias-requeridas` que contiene las dependencias, metadatos y scripts necesarios para cualquier nuevo proyecto.

Para lograr esto con éxito, utilizarás una **Estrategia Multi-Agente**. Asumirás o invocarás los siguientes roles:
1. **Agente Analista de Legado:** Audita el código Java y genera `docs/CASOS_DE_USO_LEGADO.md` con la matriz de mapeo de errores y endpoints.
2. **Agente Ingeniero Hexagonal (Tú):** Desarrolla el código limpio en Go.
3. **Agente QA de Paridad:** Dispara pruebas a ambos sistemas para certificar en `docs/PAYLOADS_CASOS_DE_USO.md` la paridad byte a byte (Cuerpos JSON, Headers, Status Codes).

### 🚨 REGLAS DE ORO Y LECCIONES APRENDIDAS
*   **Librerías Privadas:** Si la plantilla corporativa contiene librerías privadas inalcanzables localmente, **inicia el proyecto usando librerías estándar** (`gin`, `net/http`, `slog`). Al finalizar, crea una guía para migrar a las corporativas.
*   **Paridad Extrema:** Debes forzar la notación científica en floats si Java lo hacía (ej. `2.255E-4`), usar punteros para enviar nulos (`null`) exactos en errores, y replicar el mapeo de status codes de Java (incluso si devolvía 402 para "Not Found" o 200 OK para Timeouts).
*   **Configuración Estricta (12-Factor):** No "hardcodees" URLs (ej. `localhost:8080`) ni puertos. Falla rápido (`os.Exit(1)`) si falta una variable crítica.
*   **Estabilidad de Pipelines corporativos:** Corrige el `Makefile` corporativo cambiando `$(go list ...)` por `$(shell go list ...)`. Agrega siempre tests "dummies" en todas las subcarpetas de `/internal` para evitar que el pipeline falle por el error `go: no such tool "covdata"`. Si Windows Defender bloquea los tests locales, usa `go test -c`.
*   **Paridad de Observabilidad (Logs y Trazabilidad):** Replica exactamente la estructura JSON de los logs del legado (ej. Logback). Utiliza `slog.HandlerOptions{ReplaceAttr}` para mapear llaves (ej. `time` -> `timestamp`, `msg` -> `message`) e inyecta atributos nativos del legado (`thread`, `logger`) para que los indexadores corporativos no noten el cambio. Instancia siempre middlewares de APM (`dd-trace-go`).
*   **Manejo de Bugs Legados (Casteo de Excepciones):** Los frameworks Java (Spring/Jackson) suelen fallar al parsear respuestas inesperadas de terceros, envolviendo validaciones 400 en errores genéricos 500. En Go, atrapa los errores de dominio/validación y fuérzalos a replicar estos status codes y strings literales para garantizar que los clientes no se rompan.
*   **Orquestación de Testing Híbrida Agonóstica:** No asumas que Docker funciona sin fricciones o que los puertos estarán libres. Crea siempre un script nativo (ej. PowerShell) que destruya agresivamente procesos huérfanos (`Stop-Process`), levante el legado en Docker con imágenes ultra-ligeras (`alpine`), ejecute Go localmente para facilitar el debug, e inyecte todas las variables de entorno dinámicamente antes de lanzar la certificación de QA.

---

Por favor, ejecuta paso a paso el siguiente plan para crear el nuevo servicio:

### FASE 1: Discovery Inicial
0. Analiza el código fuente Java para generar el documento de casos de uso y matriz de errores.

### FASE 2: Inicialización y Estructura
1. Copia todos los archivos de la carpeta `librerias-requeridas` (`go.mod`, `go.sum`, `Makefile`, `.gitlab-ci.yml`, `sonar-project.properties`, `catalog-info.yaml`, `.gitignore`) a la raíz de mi nuevo proyecto.
2. Ejecuta el script de andamiaje correspondiente a mi sistema operativo (`scaffold.ps1` o `scaffold.sh`) para generar las carpetas de Arquitectura Hexagonal (`cmd`, `internal/app`, `internal/core`, `internal/infra`).
3. Modifica la primera línea del archivo `go.mod` para que el nombre del módulo coincida con el nombre de mi nuevo servicio.
4. Ejecuta el comando `go mod tidy` para enlazar las dependencias. *(Nota: Si falla por librerías privadas, aplica la regla de usar estándar temporalmente).*

### FASE 3: Adaptación de Metadatos
5. Abre el archivo `catalog-info.yaml` y actualiza los campos `name`, `description`, `owner`, y `system` con los datos del nuevo servicio proporcionados abajo.
6. Abre el archivo `sonar-project.properties` y actualiza `sonar.projectKey` y `sonar.projectName` con el nombre del nuevo servicio.
7. Abre el archivo `Makefile` y modifica la variable `API_NAME` para que coincida con el nuevo servicio, además de corregir el bug de `$(shell go list...)`.

### FASE 4: Desarrollo Inicial (Wiring y Mapeo Estricto)
8. Crea el archivo `cmd/api/main.go` instanciando Gin, configurando los middlewares de Datadog (`gintrace`), el Logger JSON, y preparando la inyección de dependencias (solo mediante Variables de Entorno).
9. Al crear los `RequestMapper` y `ResponseMapper` (DTOs), asegúrate de usar tags de JSON (`json:"nombre_raro"`) para calzar exactamente con el contrato legacy.

### FASE 5: Verificación y Testing Avanzado
10. Genera las pruebas Unitarias, de Integración (simulando Upstream) y E2E.
11. Dispara las comprobaciones del Agente QA para garantizar paridad.

---

### DATOS DEL NUEVO SERVICIO:
*   **Nombre del módulo / API:** [ej: api-bfcl-gestion-clientes]
*   **Descripción breve:** [ej: Microservicio para la administración y consulta de perfil de clientes]
*   **Equipo Dueño (Owner):** [ej: group:fiftech_clientes]
*   **Sistema al que pertenece:** [ej: customer-management-chile]
*   **Base de datos a utilizar:** [ej: Oracle, Postgres, MongoDB, N/A]
*   **Sistema Operativo que estoy usando:** [Windows o Linux/Mac]

Procede con la ejecución y avísame cuando estés listo para que te pase las definiciones del Dominio (Entidades), JSON requests/responses legados, y empecemos a crear los Endpoints manteniendo la paridad estricta.
### Arquitectura de Integración (Salesforce -> Go -> ASICOM/FinnFlow)
Para futuras referencias o creación de nuevos microservicios en esta célula, el flujo de datos y conexión es el siguiente:

1. **Salesforce (Frontend/CRM):** 
   - **Rol:** Actúa como el origen de los datos. Los ejecutivos hipotecarios llenan los formularios en Salesforce.
   - **Conexión:** Consume el microservicio de Go a través del puerto expuesto localmente por la variable de entorno \PORT\.

2. **Microservicio Go (Capa Adaptadora/Integration):** 
   - **Rol:** Recibe la solicitud de Salesforce, aplica validaciones estrictas (ej. formato de RUT, límites de caracteres como \NumeroOperacionOriginal\), descarta campos no mapeados gracias a \omitempty\ y formatea el JSON al estándar rígido que espera el Core.
   - **Conexión:** Sirve a Salesforce, pero a su vez actúa como cliente HTTP hacia ASICOM.

3. **ASICOM / FinnFlow (Core Bancario / Backend Legacy):**
   - **Rol:** Es el sistema central que consolida e inyecta la solicitud de crédito o renegociación definitivamente.
   - **Variables de Conexión:** El microservicio Go se autentica y transmite los datos hacia ASICOM utilizando obligatoriamente:
     - \FINNFLOW_URL\: Endpoint resoluble del Core.
     - \FINNFLOW_KEY\: Credencial de acceso (API Key o Client ID).
     - \FINNFLOW_SECRET\: Secreto para autorizar la transacción.

> ⚠️ Cualquier nuevo microservicio que deba inyectar datos hacia ASICOM desde Salesforce deberá replicar este mismo patrón de inyección de variables de entorno para evitar *hardcodear* credenciales en el código fuente.
