# Design Specification: Calanggo URL Shortener
*Based on ByteByteGo System Design*

## 1. Requisitos
### Funcionais
1.  **Shortening:** Dado uma URL longa, retornar uma URL curta única.
2.  **Redirection:** Dado uma URL curta, redirecionar (HTTP 301 ou 302) para a original.
3.  **High Availability:** O sistema deve ser extremamente rápido na leitura (redirecionamento).

### Não-Funcionais
-   Baixa latência.
-   URLs encurtadas não devem ser previsíveis (desejável).

## 2. API Design (REST)

### A. Criar Short URL
**POST** `/api/v1/shorten`
-   **Request:** `{ "long_url": "https://www.google.com" }`
-   **Response:** `{ "short_url": "http://calang.go/7rB8u" }`

### B. Redirecionar
**GET** `/{short_code}`
-   **Response:** HTTP 301 (Permanent) ou 302 (Found) -> Location: long_url.
-   *Nota:* 301 é melhor para o servidor (browser cacheia), 302 é melhor para analytics (bate sempre no server). Usaremos **302** inicialmente para ver os logs de acesso.

## 3. Core Logic: Base62 Conversion
Em vez de usar Hash (MD5/SHA) que gera colisões, usaremos um **Gerador de ID Distribuído** (ou auto-incremento de banco para simplificar o MVP) convertido para **Base62**.

-   Caracteres: `[a-z, A-Z, 0-9]` = 62 caracteres.
-   Tamanho: 7 caracteres. $62^7 \approx 3.5$ trilhões de URLs.
-   Fluxo: `ID (int64) -> Base62 Encode -> ShortURL`.

## 4. Data Model
Embora o ByteByteGo sugira Relacional, Go trabalha muito bem com qualquer um.
Estrutura da Tabela `links`:
-   `id`: BIGINT (Primary Key, Auto Increment ou Snowflake ID)
-   `short_code`: VARCHAR(7) (Indexado, Unique)
-   `original_url`: VARCHAR(2048)
-   `created_at`: DATETIME

## 5. Estrutura do Projeto Go (Clean Arch Simplificada)
```text
/cmd
  /server       # Entry point (main.go)
/internal
  /core
    /domain     # Entidades (Link)
    /ports      # Interfaces (Repository, Service)
    /services   # Casos de uso (ShortenLink, GetOriginal)
  /adapters
    /db         # Implementação do Repository (Postgres/Memory)
    /http       # Handlers HTTP
/pkg
  /base62       # Algoritmo utilitário (reusável)
```
