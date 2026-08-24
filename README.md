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
  - [Autenticação](#autenticação)
  - [Perfis](#perfis)
  - [Conteúdos](#conteúdos)
  - [Favoritos](#favoritos)
- [Projeto em funcionamento](#projeto-em-funcionamento)
- [Como rodar o projeto](#como-rodar-o-projeto)
- [Colaboradores](#colaboradores)

## Descrição do projeto

O Streaming API é uma solução backend em Go para gestão de uma plataforma de streaming. A aplicação foi desenvolvida para oferecer uma API REST confiável, fácil de evoluir e preparada para uso em ambientes locais.

O projeto incorpora boas práticas de arquitetura, separando responsabilidades em camadas bem definidas, como rotas, controladores, serviços, repositórios e modelos. Além disso, a API conta com autenticação via JWT, proteção contra excesso de requisições, integração com SQLite e Redis, documentação com Swagger e monitoramento de operações por tracing.

### Funcionalidades Principais

> **_Autenticação:_** <br>
> **_Gerenciamento de Perfis do usuário:_** <br>
> **_Listagem de Filmes e Séries:_** <br>
> **_Gerenciamento de Favoritos:_** <br>
> **_Observabilidade com Jaeger:_** <br>
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

graph TD
    Client[Cliente/Frontend] -->|HTTP Request| Fiber[Rotas - Fiber]
    Fiber --> Middlewares[Middlewares JWT/Rate Limit]
    Middlewares --> Controller[Controller]
    Controller --> Service[Service - Regras de Negócio]
    Service --> Repository[Repository - Persistência]
    Repository --> SQLite[(SQLite - Dados Opcionais)]
    Repository --> Redis[(Redis - Cache)]

## Endpoints

A API está organizada sob o prefixo `/v1` e, em sua maior parte, utiliza autenticação por token JWT no header `Authorization`.

### Autenticação

| Método | Rota | Descrição |
| ------ | ---- | --------- |
| `POST` | `/v1/auth/register` | Cria um novo usuário na aplicação |
| `POST` | `/v1/auth/login` | Realiza login e retorna o token de acesso |
| `GET`  | `/v1/auth/me` | Retorna os dados do usuário autenticado |

### Perfis

> Todos os endpoints abaixo exigem autenticação.

| Método | Rota | Descrição |
| ------ | ---- | --------- |
| `GET` | `/v1/profiles` | Lista os perfis do usuário autenticado |
| `POST` | `/v1/profiles` | Cria um novo perfil associado ao usuário |
| `PUT` | `/v1/profiles/:id` | Atualiza um perfil específico |
| `DELETE` | `/v1/profiles/:id` | Remove um perfil do usuário |

### Conteúdos

> Todos os endpoints abaixo exigem autenticação.

| Método | Rota | Descrição |
| ------ | ---- | --------- |
| `GET` | `/v1/contents` | Lista os conteúdos disponíveis para o usuário autenticado |
| `GET` | `/v1/contents/search` | Busca conteúdos por termo ou filtro específico |

### Favoritos

> Todos os endpoints abaixo exigem autenticação.

| Método | Rota | Descrição |
| ------ | ---- | --------- |
| `GET` | `/v1/favorites` | Lista os itens favoritos do usuário |
| `POST` | `/v1/favorites` | Adiciona um item à lista de favoritos |
| `DELETE` | `/v1/favorites/:id` | Remove um item favorito pelo identificador |

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
