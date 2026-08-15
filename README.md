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
- [Projeto em funcionamento](#projeto-em-funcionamento)
- [Como rodar o projeto](#como-rodar-o-projeto)
- [Colaboradores](#colaboradores)

## Descrição do projeto

### Funcionalidades Principais

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
├── models
├── repositories
├── security
├── services
└── shared
```

| Camada | Responsabilidade |
| ------ | ---------------- |

### Arquitetura do Software (macro)

## Projeto em funcionamento

## Como rodar o projeto

```bash
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
