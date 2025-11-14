FROM golang:1.25-trixie AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build --ldflags '-linkmode external -extldflags=-static' -o spamscale .

FROM gcr.io/distroless/static-debian12

WORKDIR /mnt

COPY --from=builder /app/spamscale /spamscale

ENTRYPOINT ["/spamscale"]
CMD ["serve"]