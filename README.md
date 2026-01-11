# Calanggo URL Shortener 🦎

**Calanggo** é um encurtador de URLs de alta performance desenvolvido em Go, seguindo os princípios da **Arquitetura Hexagonal** (Ports & Adapters).

## 📋 Visão Geral

O projeto visa demonstrar uma implementação limpa e desacoplada de um serviço de encurtamento de links, separando a lógica de negócio (Core) das dependências externas (Adapters), como interfaces HTTP e persistência de dados.

### Funcionalidades Principais

*   **Encurtamento de URLs:** Gera códigos curtos e únicos para URLs longas.
*   **Redirecionamento:** Redireciona usuários da URL curta para a original.
*   **Contagem de Visitas:** Monitora o número de acessos de cada link encurtado.
*   **Documentação API:** Interface Swagger integrada.

## 🏗️ Arquitetura

O sistema utiliza a arquitetura Hexagonal para garantir testabilidade e manutenibilidade.

*   **Core (Hexagon):** Contém as regras de negócio (`LinkService`) e entidades (`Link`). Não depende de frameworks externos.
*   **Ports:** Interfaces que definem como o mundo externo interage com o Core (Inbound) e como o Core interage com recursos externos (Outbound).
*   **Adapters:** Implementações concretas das portas.
    *   *Driving (Primary):* Handler HTTP.
    *   *Driven (Secondary):* Repositório SQLite / Memória.

## 🚀 Como Executar

### Pré-requisitos

*   [Go](https://go.dev/) 1.21+ instalado.

### Passos

1.  Clone o repositório:
    ```bash
    git clone https://github.com/seu-usuario/calanggo-server.git
    cd calanggo-server
    ```

2.  Instale as dependências:
    ```bash
    go mod download
    ```

3.  Execute a aplicação:
    ```bash
    go run main.go
    ```
    *O servidor iniciará na porta `8080` por padrão.*

## 📚 Documentação da API (Swagger)

Com o servidor rodando, acesse a documentação interativa em:

👉 **[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)**

### Endpoints Principais

*   **POST** `/api/v1/shorten`
    *   Cria uma nova URL encurtada.
    *   Body: `{ "long_url": "https://exemplo.com" }`
*   **GET** `/{code}`
    *   Redireciona para a URL original.

## 🛠️ Estrutura do Projeto

```text
/
├── docs/               # Documentação Swagger e Specs
├── internal/
│   ├── adapters/       # Implementações (HTTP, SQLite)
│   └── core/           # Lógica de Negócio (Services, Domain)
├── pkg/                # Pacotes utilitários (Base62)
├── tests/              # Testes de integração
└── main.go             # Ponto de entrada
```

## 🧪 Testes

Para rodar os testes do projeto:

```bash
go test ./...
```
