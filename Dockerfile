FROM golang:1.26.5 AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -trimpath -o app .

FROM scratch AS runner
WORKDIR /app

COPY --from=build /src/app /app/app

ENTRYPOINT [ "/app/app" ]
EXPOSE 8080
