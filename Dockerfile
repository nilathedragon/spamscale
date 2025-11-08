FROM golang:1.25-trixie AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o spamscale .

FROM gcr.io/distroless/static-debian12

WORKDIR /

COPY --from=builder /app/spamscale /spamscale

ENTRYPOINT ["/spamscale"]
