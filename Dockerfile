# Etapa 1: Compilación
FROM golang:1.22-alpine AS builder


WORKDIR /app

# Descargamos paquetes necesarios
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compilamos el binario en modo producción 
RUN go build -o server ./web


FROM alpine:latest

WORKDIR /app

# Copiamos solo el ejecutable desde la etapa de compilación
COPY --from=builder /app/server .

# Exponemos el puerto
EXPOSE 8080

# Ejecutamos el binario
CMD ["./server"]
