<h1 align="center">Streaming API</h1>

<p align="center">
<img loading="lazy" src="https://img.shields.io/static/v1?label=STATUS&message=EM%20ANDAMENTO&color=blue&style=for-the-badge"/>
</p>

> [!IMPORTANT]
> _Esse projeto está em desenvolvimento._

### Tópicos

- [Descrição do projeto](#descrição-do-projeto)
  - [Funcionalidades Principais](#funcionalidades-principais)
- [Tecnologias](#tecnologias)
- [Arquitetura](#arquitetura)
  - [Arquitetura do Código (micro)](#arquitetura-do-código-micro)
  - [Arquitetura do Software (macro)](#arquitetura-do-software-macro)
- [Endpoints](#endpoints)
- [Projeto em funcionamento](#projeto-em-funcionamento)
- [Como rodar o projeto](#como-rodar-o-projeto)
- [Colaboradores](#colaboradores)

## Descrição do projeto

### Funcionalidades Principais

> **_Autenticação:_**
> **_Gerenciamento de Perfis do usuário:_**
> **_Listagem de Filmes e Séries:_**
> **_Gerenciamento de Favoritos:_**
> **_Observabilidade com Jaeger:_**
> **_Documentação com Swagger:_**

## Tecnologias

<details closed>
<summary>Linguagens e Frameworks</summary>
  <div width="140px">
      <img src="https://i.icoziv.workers.dev/icons?i=go,fiber,opentelemetry" />
  </div>
</details>

<details closed>
<summary>Bancos de Dados e Armazenamento</summary>
  <div width="140px">
      <img src="https://i.icoziv.workers.dev/icons?i=sqlite,redis" />
  </div>
</details>

<details closed>
<summary>DevOps, Cloud e Infraestrutura</summary>
  <div width="140px">
      <img src="https://i.icoziv.workers.dev/icons?i=docker" />
  </div>
</details>

<details closed>
<summary>Ferramentas</summary>
  <div width="140px">
      <img src="https://i.icoziv.workers.dev/icons?i=vscode,swagger,jwt,gorm" />
  </div>
</details>

## Arquitetura

### Arquitetura do Código (micro)

A estrutura de pastas segue uma abordagem baseada em **Domain-Driven Design (DDD)** e **Arquitetura Limpa**, garantindo um baixo acoplamento e facilitando os testes unitários.

```
src/
├── api
│   ├── routes
│   └── v1
│       ├── controllers
│       ├── dto
│       ├── validators
│       └── responses
├── configs
│   ├── db
│   ├── env
│   └── tracing
├── container
├── middlewares
├── mocks
│   ├── repositories
│   └── services
├── models
├── repositories
│   ├── http
│   ├── interfaces
│   └── sqlite
├── security
├── services
│   └── interfaces
└── shared
```

| Camada | Responsabilidade |
| ------ | ---------------- |
| API    |                  |
| Routes |                  |

### Arquitetura do Software (macro)

A aplicação funciona no modelo Client-Server. O cliente envia requisições HTTP que são roteadas pelo Fiber.

## Endpoints

## Projeto em funcionamento

Clique na imagem abaixo para assistir ao tutorial em vídeo!

[![Assista ao tutorial](./images/streaming_api.png "Projeto em funcionamento local")](semvideo.com)

**Descrição**: Este vídeo cobre todo o processo para visualizar o projeto em funcionamento, do início ao fim.

## Como rodar o projeto

```
< INSTALADORES >

go mod tidy


< INICIADORES >

docker compose up --build

OBS: Se preferir, pode rodar via debug pela sua IDE


< TESTES DE COVERAGE >

go test ./... -v -coverprofile=coverage.out
go tool cover -func=coverage.out
```

## Colaboradores

| [<img src="https://avatars.githubusercontent.com/u/69527468?v=4" width=115><br><sub>Kauê Bertaze de Oliveira</sub>](https://github.com/KaueTTS)<br><sub>Developer Full Stack</sub> |
| :---:
