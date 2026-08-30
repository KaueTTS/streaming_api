<h1 align="center">Padrões Docker</h1>

### Tópicos

- [Principais comandos do Docker](#principais-comandos-do-docker)
- [Padrões de configuração do dockerfila](#padrões-de-configuração-do-dockerfile)
    - [Golang](#-golang)

## Principais comandos do Docker

- `docker compose up --build` - Construir imagem do container.

- `docker compose down` - Parar e remover containers, redes e volumes criados.

- `docker compose stop` - Parar container.

- `docker compose start` - Iniciar container.

- `docker compose restart` - Reiniciar container.

- `docker compose ps` - Listar containers que estão rodando no projeto.

- `docker compose exec <servico_container> <comando>` - Roda um comando dentro do container.

- `docker compose logs` - Mostra os logs dos containers (se colocar `-f` no final, ele mostra os logs em tempo real).

- `docker compose config` - Configuração do docker.

- `docker system prune -a --volumes` - Realiza uma limpeza profunda e completa (Removendo containers, redes e volumes parados, e imagens que não estão sendo usadas).

## Padrões de configuração do dockerfile

### # Golang

```dockerfile
# Estágio de compilação
FROM golang:<version>-alpine AS builder
WORKDIR /app

# Cache das dependências
COPY go.mod go.sum* ./
RUN go mod download

# Cópia do código
COPY . .

# Compilação estática e otimizada
RUN CGO_ENABLED=0 GOOS=linux go build -o <nome-do-projeto> ./main.go

# Estágio final (imagem leve)
FROM alpine:3.21
WORKDIR /app

# Adiciona certificados e cria usuário seguro
RUN apk add --no-cache ca-certificates && \
    adduser -D -g '' appuser

# Copia o binário do builder
COPY --from=builder /app/to-do-list-api ./to-do-list-api

# Troca para o usuário seguro antes de rodar a aplicação
USER appuser

# Expõe a porta e define o comando padrão
EXPOSE 8080
CMD ["./<nome-do-projeto>"]
```
