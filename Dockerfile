FROM golang:1.23-alpine AS builder
WORKDIR /my-domain
COPY . .
ENV GOCACHE=/root/.cache/go-build
RUN go install github.com/google/wire/cmd/wire@latest && \
    go mod download && \
    wire ./cmd
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o ./cmd/app ./cmd

FROM alpine:3.20
COPY --from=builder /my-domain/cmd/app /
EXPOSE 8000
ENTRYPOINT ["./app"]
