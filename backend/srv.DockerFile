FROM golang:1.23 AS builder

ARG ZYG_SRV_PORT=5000
ARG DATABASE_URL
ARG SUPABASE_JWT_SECRET
ARG RESEND_API_KEY
ARG ZYG_DB_QUERY_DEBUG=0

WORKDIR /usr/src/app

# Copy only go.mod and go.sum first for better layer caching
COPY go.mod ./
COPY go.sum ./
RUN go mod download && go mod verify

# Copy the rest of the source code
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o server ./cmd/srv/main.go

# Build the runtime container image from scratch, copying what is needed from the previous stage.
FROM alpine

# Copy the binary to the production image from the builder stage.
COPY --from=builder /usr/src/app/server /usr/local/bin/server

ENV ZYG_SRV_PORT=${ZYG_SRV_PORT}
ENV DATABASE_URL=${DATABASE_URL}
ENV SUPABASE_JWT_SECRET=${SUPABASE_JWT_SECRET}
ENV RESEND_API_KEY=${RESEND_API_KEY}
ENV ZYG_DB_QUERY_DEBUG=${ZYG_DB_QUERY_DEBUG}

EXPOSE ${ZYG_SRV_PORT}

CMD server -host 0.0.0.0 -port $ZYG_SRV_PORT
