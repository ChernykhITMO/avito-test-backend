FROM golang:1.25-bookworm AS build

WORKDIR /src

COPY go.mod ./
COPY go.sum ./
COPY cmd ./cmd
COPY docs ./docs
COPY internal ./internal
COPY migrations ./migrations

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/app ./cmd/api

FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /app

COPY --from=build /out/app /app/app
COPY --from=build /src/migrations /app/migrations

EXPOSE 8080

ENTRYPOINT ["/app/app"]
