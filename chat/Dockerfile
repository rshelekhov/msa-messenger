#######################################
# Stage 1: Build the application
#######################################
FROM golang:1.23.0-alpine3.19 AS build

WORKDIR /app

# Copy the go.mod and go.sum files
COPY chat/go.mod chat/go.sum ./
RUN go mod download -x && go mod verify

# Copy the source code for the chat service
COPY chat/ ./

# Build the chat service
RUN CGO_ENABLED=0 GOOS=linux go build -o chat ./cmd/main.go

#######################################
# Stage 2: Create the final image
#######################################
FROM alpine:3.19 AS final

WORKDIR /

# Install ca-certificates and tzdata
RUN apk --no-cache add ca-certificates tzdata

# Copy the chat binary from the build stage
COPY --from=build /app/chat .

# Copy config files
COPY --from=build /app/config ./config

# Expose the port that the chat service will run on
EXPOSE 8083

HEALTHCHECK --interval=30s --timeout=3s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8083/health || exit 1

CMD ["./chat"] 