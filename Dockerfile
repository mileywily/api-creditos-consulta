# Etapa 1: Build
FROM golang:1.27-alpine AS builder

# Habilitar Go modules y configurar proxy si es necesario
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Descargar dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el codigo fuente
COPY . .

# Compilar el binario
RUN go build -ldflags="-w -s" -o api-bin ./cmd/api/main.go

# Etapa 2: Imagen productiva (Scratch o Alpine ligero)
FROM alpine:3.19

# Agregar certificados raiz para llamadas HTTPS a APM / Finnflow
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copiar el binario desde la etapa de builder
COPY --from=builder /app/api-bin .

# Exponer el puerto por defecto
EXPOSE 8080

# Ejecutar la aplicacion
CMD ["./api-bin"]
