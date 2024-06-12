FROM public.ecr.aws/docker/library/golang:1.22-bullseye as BUILD

ENV PATH="/go/bin:${PATH}"
ENV WORKDIR="/app"

WORKDIR $WORKDIR

COPY . $WORKDIR

RUN apt-get update && \
  go install github.com/githubnemo/CompileDaemon@latest && \
  go mod download

ENTRYPOINT ["CompileDaemon", "-directory=./", "-build=go build -o weather-by-cep /app/cmd/main.go", "-command=./weather-by-cep"]