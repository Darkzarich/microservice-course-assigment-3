# syntax=docker/dockerfile:1
FROM golang:1.22 AS build
WORKDIR /src

COPY go.mod ./

RUN go mod download

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 go build -o /bin/app main.go

FROM gcr.io/distroless/base-debian12
COPY --from=build /bin/app /app
EXPOSE 8080 19090
ENTRYPOINT ["/app"]
