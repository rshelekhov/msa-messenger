#######################################
# Stage 1: Build the application
#######################################
FROM golang:1.24.0-alpine AS build

WORKDIR /app

# Copy the go.mod and go.sum files
COPY gateway/go.mod gateway/go.sum ./
RUN go mod download -x && go mod verify

# Copy the source code for the gateway service
COPY gateway/ ./

# Build the gateway service
RUN CGO_ENABLED=0 GOOS=linux go build -o gateway ./cmd/main.go

#######################################
# Stage 2: Create the final image
#######################################
FROM alpine:3.19 AS final

WORKDIR /

# Install ca-certificates and tzdata
RUN apk --no-cache add ca-certificates tzdata

# Copy the gateway binary from the build stage
COPY --from=build /app/gateway .

# Expose the port that the gateway service will run on
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./gateway"] 