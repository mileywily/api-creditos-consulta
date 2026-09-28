# 1. Leer el puerto directamente de la variable de entorno (o usar 8080 si no existe para no fallar)
$puerto = if ($env:PORT) { $env:PORT } else { "8080" }

# 2. Buscar si hay algún proceso usando ese puerto y matarlo
Write-Host "Buscando procesos en el puerto $puerto..." -ForegroundColor Yellow
$conexion = Get-NetTCPConnection -LocalPort $puerto -ErrorAction SilentlyContinue

if ($conexion) {
    Write-Host "Proceso encontrado (PID: $($conexion.OwningProcess)). Matando proceso..." -ForegroundColor Red
    Stop-Process -Id $conexion.OwningProcess -Force
    Write-Host "Proceso terminado exitosamente." -ForegroundColor Green
} else {
    Write-Host "El puerto $puerto está libre." -ForegroundColor Green
}

# 3. Limpiar la caché de pruebas de Go
Write-Host "Limpiando caché de pruebas..." -ForegroundColor Yellow
go clean -testcache

# 4. Ejecutar toda la suite de pruebas unitarias
Write-Host "Ejecutando suite de pruebas unitarias..." -ForegroundColor Yellow
go test -v ./...

# 5. Generar reporte de cobertura (Coverage)
Write-Host "Calculando cobertura de código..." -ForegroundColor Yellow
go test ./internal/... -coverprofile coverfile_out > $null
go tool cover -func coverfile_out