# healthcheck-go
Pingdom-like monitoring solution

# UTILIZANDO

O healthcheck-go utiliza variáveis de ambiente como parâmetros 

URL_MONITORAMENTO = URL a qual você deseja monitorar (Default: http://www.ecloudc.com.br)

TIMEOUT_MONITORAMENTO = Tempo de timeout para a chamada GET (em segundos - Default: 5 segundos)

INTERVALO_MONITORAMENTO = Tempo de delay entre cada chamada (em segundos - Default: 60 segundos)

# TODO

Lista de e-mails configurável via parâmetro na aplicação

Implementar método POST

Implementar métricas

Implementar integração Telegram

# DESENVOLVENDO

1. Adicione as dependências do projeto

go get -v github.com/sendgrid/sendgrid-go 

go get -v github.com/gorilla/mux

2. Configure a sua API-KEY do Sendgrid

Edite o arquivo sendgrid.env

Adicione a API-KEY

Execute: source sendgrid.env

3. Rode o projeto

go run monitoramento.go

4. Build

go build monitoramento.go


# BUILD CONTAINER

docker build -t matheusbona/healthcheck-go .


# EXECUTAR CONTAINER

docker run -d --name healthcheck-go -e URL_MONITORAMENTO=http://ecloudc.com.br -e TIMEOUT_MONITORAMENTO=2 -e INTERVALO_MONITORAMENTO=5 matheusbona/healthcheck-go
