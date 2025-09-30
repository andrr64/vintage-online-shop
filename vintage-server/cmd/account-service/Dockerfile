# Tahap 1: Builder
FROM golang:1.25-alpine AS builder

# Terima argumen nama service dari docker-compose.yml
# Ini adalah kuncinya agar Dockerfile ini bisa dipakai ulang
ARG SERVICE_NAME

# Tentukan direktori kerja
WORKDIR /app

# Salin file dependensi dari root proyek
COPY go.mod go.sum ./
# Download dependensi
RUN go mod download

# Salin semua kode dari root proyek
# Ini diperlukan agar import ke 'internal' atau 'pkg' bisa ditemukan
COPY . .

# Build HANYA service yang spesifik
# Perintah ini akan menjalankan sesuatu seperti:
# go build -o /app/main ./cmd/account-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/${SERVICE_NAME}

# -------------------------------------------------------------------

# Tahap 2: Final
FROM alpine:latest

WORKDIR /root/

# Salin HANYA binary yang sudah di-compile dari tahap 'builder'
COPY --from=builder /app/main .

# Expose port (opsional, tapi praktik yang baik)
# Port ini akan di-set melalui environment variable di docker-compose
EXPOSE 8080

# Jalankan binary
CMD ["./main"]