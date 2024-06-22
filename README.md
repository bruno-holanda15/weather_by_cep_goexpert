## Overview

Projeto proposto na Pós Goexpert da Fullcycle para adição de telemetria utilizando o [OpenTelemetry](https://opentelemetry.io/docs/languages/go/) em conjunto com o [Zipkin](https://zipkin.io/) dentro de dois serviços.
Podemos avaliar a separação dos Serviços A e B pelos services do [docker-compose.yaml](docker-compose.yaml) com os nomes de go_wbc1 e go_wbc2 respectivamente.

Container go_wbc1 valida cep recebido na rota http://localhost:8081/weather, cria um span dentro do handler ValidateCep no handlers.go e aciona o usecase passando o contexto criado juntamente com o span, propagando sua transmissão ao fazer a request para o container/service go_wbc2.

Container go_wbc2 recebe a request do go_wbc2 e faz a chamada para pesquisar a temperatura do cep passado, pegando o contexto a partir da propagação passado no header da request.

Podemos visualizar o trace pelo serviço Zipkin na url http://localhost:9411

Adicione o arquivo .env na raiz do seu projeto clonado com a variável WEATHER_TOKEN="TOKEN_DA_WEATHERAPI" na raiz do projeto com o seu token para a api da WeatherAPI para busca de informações, para informações de como criar o token do serviço, acesso o site https://www.weatherapi.com/

[Documentação API](./docs/api.http)

## Objetivo

Aplicação de telemetria em sistemas independetes para avaliar tempo de execução dos mesmos.

## Requisitos
O sistema deve go_wbc1 receber um CEP válido de 8 digitos e acionar o go_wbc2.
O sistema go_wbc2 deve realizar a pesquisa do CEP e encontrar o nome da localização, a partir disso, deverá retornar as temperaturas e formata-lás em: Celsius, Fahrenheit, Kelvin.
O sistema deve responder adequadamente nos seguintes cenários:
Em caso de sucesso:
- Código HTTP: 200
- Response Body: { "temp_C": 28.5, "temp_F": 28.5, "temp_K": 28.5 }
Em caso de falha, caso o CEP não seja válido (com formato correto):
- Código HTTP: 422
- Mensagem: invalid zipcode
​​​Em caso de falha, caso o CEP não seja encontrado:
- Código HTTP: 404
- Mensagem: can not find zipcode