#######################################
# Stage 1: Build the application
#######################################
FROM golang:1.23.0-alpine3.19 AS build

WORKDIR /app

# Copy the go.mod and go.sum files
COPY notification/go.mod notification/go.sum ./
RUN go mod download -x && go mod verify

# Copy the source code for the notification service
COPY notification/ ./

# Build the notification service
RUN CGO_ENABLED=0 GOOS=linux go build -o notification ./cmd/main.go

#######################################
# Stage 2: Create the final image
#######################################
FROM alpine:3.19 AS final

WORKDIR /

# Install ca-certificates and tzdata
RUN apk --no-cache add ca-certificates tzdata

# Copy the notification binary from the build stage
COPY --from=build /app/notification .

# Expose the port that the notification service will run on
EXPOSE 8085

HEALTHCHECK --interval=30s --timeout=3s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8085/health || exit 1

CMD ["./notification"] 