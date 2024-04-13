FROM golang:alpine AS builder

WORKDIR /build

ADD go.mod .

COPY . .

RUN go mod download && go build -o it-revolution-backend .

FROM alpine

WORKDIR /app

COPY --from=builder /build/it-revolution-backend /app/it-revolution-backend

CMD ["./it-revolution-backend"]
