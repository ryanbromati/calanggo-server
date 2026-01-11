# Guia da Standard Library do Go (Baterias Inclusas)
*As bibliotecas nativas utilizadas no projeto calanggo-server e seus equivalentes no mundo .NET.*

## 1. Web & Networking

### `net/http`
**O que é:** O coração da web em Go. Funciona tanto como **Servidor** (como o Kestrel/IIS) quanto como **Cliente** (HttpClient). É robusto o suficiente para produção (usado pelo Cloudflare, Google, etc).
**C# Eq:** `Microsoft.AspNetCore.*` + `System.Net.Http`.
**Uso no projeto:**
- Subir o servidor (`ListenAndServe`).
- Roteamento (`NewServeMux`).
- Manipular Requests e Responses (`ResponseWriter`, `Request`).

### `encoding/json`
**O que é:** Serialização e Deserialização de JSON via Reflection e Tags.
**C# Eq:** `System.Text.Json` ou `Newtonsoft.Json`.
**Uso no projeto:**
- Ler o corpo do request (`NewDecoder(r.Body).Decode(...)`).
- Escrever a resposta (`NewEncoder(w).Encode(...)`).
- Tags nas structs: `` `json:"nome_campo"` ``.

---

## 2. Dados & Persistência

### `database/sql`
**O que é:** Uma interface genérica para bancos de dados SQL (similar ao JDBC do Java ou ADO.NET base). Ele gerencia o **Connection Pool** automaticamente. Ele **não** contém a lógica do banco específico (Postgres, SQLite), apenas a abstração.
**C# Eq:** `System.Data.Common` / ADO.NET puro.
**Uso no projeto:**
- Executar queries (`ExecContext`, `QueryRowContext`).
- Gerenciar transações e conexões.
- **Nota:** Exige um driver importado com `_` (como fizemos com o `go-sqlite`) para funcionar.

### `context`
**O que é:** Essencial para controle de ciclo de vida de requisições. Carrega **Deadlines** (Timeouts), **Sinais de Cancelamento** e valores de escopo (como User ID em middlewares).
**C# Eq:** `CancellationToken` (principalmente) e `HttpContext.Items`.
**Uso no projeto:**
- Passado do Handler -> Service -> Repository.
- Permite que o banco de dados pare uma query na metade se o usuário cancelar o request HTTP.

---

## 3. Concorrência & Sincronização

### `sync`
**O que é:** Primitivas de sincronização de baixo nível para memória compartilhada.
**C# Eq:** `System.Threading` (`Monitor`, `lock`, `Interlocked`).
**Uso no projeto:**
- `sync.RWMutex`: Usado no `memory.go` para proteger o mapa de acessos simultâneos (Leitura vs Escrita).

### `time`
**O que é:** Manipulação de tempo, datas e duração.
**C# Eq:** `System.DateTime`, `System.TimeSpan`, `Task.Delay`.
**Uso no projeto:**
- Marcar `CreatedAt`.
- Gerar seed para o rand (`time.Now().UnixNano()`).
- Definir timeouts.

---

## 4. Testes

### `testing`
**O que é:** Framework de testes unitários integrado. Simples, sem assertions mágicas.
**C# Eq:** `xUnit`, `NUnit`.
**Uso no projeto:**
- Rodar testes com `go test`.
- Objeto `t *testing.T` para reportar erros.

### `net/http/httptest`
**O que é:** Utilitários para testar handlers HTTP sem precisar abrir portas de rede reais.
**C# Eq:** `Microsoft.AspNetCore.TestHost`.
**Uso no projeto:**
- Criar um `NewServer` fake que simula a API completa nos testes de integração.

---

## 5. Utilitários Gerais

### `fmt` (Format)
**O que é:** Entrada e Saída formatada (I/O).
**C# Eq:** `Console.WriteLine`, `string.Format`.

### `log`
**O que é:** Logger simples (data + mensagem) para stderr. Em produção, prefira libs estruturadas como `slog` ou `zap`.
**C# Eq:** `Console.WriteLine` (ou `ILogger` básico).

### `os`
**O que é:** Interface com o Sistema Operacional (Arquivos, Variáveis de Ambiente, Args).
**C# Eq:** `System.Environment`, `System.IO`.

### `errors`
**O que é:** Criação e manipulação de erros simples.
**C# Eq:** `new Exception(...)`.

### `strings`
**O que é:** Manipulação eficiente de strings.
**C# Eq:** `System.String`, `StringBuilder`.
**Uso no projeto:**
- `strings.Builder`: Usado no Base62 para concatenar caracteres sem alocar memória excessiva (Performance).
- `strings.Split`, `strings.Join`, etc.
