$ErrorActionPreference = "Continue"

# 1. Leer el puerto directamente de la variable de entorno (o usar 8080 por defecto)
$puerto = if ($env:PORT) { $env:PORT } else { "8080" }

# 2. Buscar si hay algún proceso usando ese puerto y matarlo
Write-Host "Buscando procesos en el puerto $puerto..." -ForegroundColor Yellow
$conexiones = Get-NetTCPConnection -LocalPort $puerto -ErrorAction SilentlyContinue

if ($conexiones) {
    # Extraer PIDs únicos (por si hay más de una conexión usando el mismo puerto)
    $pids = $conexiones | Select-Object -ExpandProperty OwningProcess -Unique
    
    foreach ($id in $pids) {
        Write-Host "Proceso encontrado (PID: $id). Terminando..." -ForegroundColor Red
        Stop-Process -Id $id -Force -ErrorAction SilentlyContinue
    }
    Write-Host "Procesos terminados exitosamente." -ForegroundColor Green
} else {
    Write-Host "El puerto $puerto está libre." -ForegroundColor Green
}

# 3. Limpiar la caché de pruebas de Go
Write-Host "`nLimpiando caché de pruebas..." -ForegroundColor Yellow
go clean -testcache

# 4. Ejecutar toda la suite de pruebas unitarias
Write-Host "Ejecutando suite de pruebas unitarias..." -ForegroundColor Yellow
go test -v ./...

# 5. Generar reporte de cobertura (Coverage)
Write-Host "`nCalculando cobertura de código..." -ForegroundColor Yellow
go test ./internal/... -coverprofile coverfile_out | Out-Null
go tool cover -func coverfile_out
