FROM docker.io/golang:1.26.5 AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -v -x -ldflags="-s -w" -trimpath -o app .

FROM scratch AS runner
WORKDIR /app

COPY --from=build /src/app /app/

ENTRYPOINT [ "/app/app" ]
EXPOSE 8080
