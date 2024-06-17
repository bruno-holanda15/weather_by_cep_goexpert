FROM golang:1.22.1 as build
WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o weather-by-cep ./cmd/
ENTRYPOINT [ "./weather-by-cep" ]

# FROM scratch
# WORKDIR /app
# COPY --from=build /app/weather-by-cep .
# COPY --from=build /app/.env .
# ENTRYPOINT [ "./weather-by-cep" ]