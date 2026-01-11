# C# Developer's Guide to Go (Golang)
*Um guia prático baseado nas dúvidas reais surgidas durante a migração.*

## 1. Variáveis e Atribuição

### Declaração (`:=` vs `=`)
Em C#, `var` resolve tudo. Em Go, temos dois modos:

| Operador | Nome | Uso | Exemplo | Equivalente C# |
| :--- | :--- | :--- | :--- | :--- |
| **`:=`** | Short Declaration | **Criação** + Atribuição. Só dentro de funções. | `nome := "Ryan"` | `var nome = "Ryan";` |
| **`=`** | Assignment | Atribuição em variável **já existente**. | `nome = "João"` | `nome = "João";` |

> **Erro Comum:** Usar `:=` em variável que já existe gera erro "no new variables".

---

## 2. Ponteiros: O Conceito Chave (`*` e `&`)

Diferente de C# (onde `class` é ref e `struct` é value), em Go **tudo é Value (cópia)** por padrão.

### O Glossário
*   **`*Tipo`** (ex: `*Link`): "Eu guardo um **endereço** de memória, não o valor em si." (Reference Type).
*   **`&Variavel`** (ex: `&meuLink`): "Me dê o **endereço** onde essa variável está guardada."

### A Diferença: Contrato vs Entrega
*   **Assinatura (`func X() *Link`)**: O Contrato. "Prometo retornar um endereço".
*   **Retorno (`return &Link{}`)**: A Ação. "Aqui está o endereço do objeto que criei".

### Quando usar?
1.  **Performance:** Evitar copiar structs gigantes. Passe `*Struct`.
2.  **Mutabilidade:** Se você quer alterar o objeto original dentro de uma função, **tem** que passar o ponteiro.

```go
func (l *Link) Alterar() { l.Visits++ } // Altera o original
func (l Link) NaoAltera() { l.Visits++ } // Altera a cópia local
```

---

## 3. Controle de Fluxo: O Poder do `for`

Esqueça `while`, `do-while`, `foreach`. Em Go, tudo é `for`.

### "Foreach" (Range)
```go
nomes := []string{"A", "B"}
// i = índice, n = valor (cópia)
for i, n := range nomes {
    fmt.Println(i, n)
}
```

### O "Underscore" (`_`)
Go não deixa declarar variável e não usar. Se não precisa do índice:
```go
for _, n := range nomes { ... } // Ignora o índice
```

### LINQ? Não existe.
Não procure `.Where().Select()`. A filosofia do Go é "Loops explícitos são melhores que mágica oculta".
*   **Filter:** Faça um `for` com `if` e `append`.
*   **Map:** Faça um `for` com `append`.

> **Dica:** O pacote `slices` (Go 1.21+) tem `slices.Contains`, `slices.Sort`, etc.

---

## 4. Estruturas e Dados

### `make` vs `new`
*   **`make`**: Inicializa estruturas que precisam de configuração interna (Slices, Maps, Channels).
*   **`new`**: Aloca memória zerada e retorna ponteiro (raramente usado para maps/slices).

### Map (Dictionary)
`map[KeyType]ValueType`.
*   **C#:** `Dictionary<TKey, TValue>`.
*   **Inicialização:** `m := make(map[string]int)`.
    *   ⚠️ **Perigo:** Se usar apenas `var m map[string]int`, ele é `nil`. Escrever nele (`m["key"]=1`) causa **Panic**.
*   **Thread-Safety:** **NÃO** é seguro. Se duas goroutines escreverem ao mesmo tempo, app crasha. Use `sync.RWMutex` para proteger.

### Slice (Lista Dinâmica)
Uma "janela" para um array.
```go
// make([]Type, len, cap)
lista := make([]string, 0, 5) // Len 0, Cap 5 (Alocação prévia para performance)
lista = append(lista, "Item") 
```

---

## 5. Web & Concorrência

### Channels (`chan`)
O "Tubo" de comunicação entre Goroutines.
*   **C#:** `System.Threading.Channels` ou `BlockingCollection`.
*   **Sintaxe:** `make(chan Tipo, Buffer)`.
*   **Unbuffered:** `make(chan int)` -> Bloqueia o sender até alguém ler. (Sincronismo puro).
*   **Buffered:** `make(chan int, 100)` -> Aceita 100 itens sem bloquear o sender. Age como uma Fila.
*   **Uso:** 
    *   `ch <- valor` (Envia/Produz)
    *   `valor := <-ch` (Recebe/Consome)

### Mux (`http.NewServeMux`)
É o Roteador (`app.MapControllers`).
*   **Wildcard `{id}`**: `mux.HandleFunc("GET /{id}", ...)` -> `r.PathValue("id")`.

### Fire-and-Forget
Executar tarefa sem travar o request.
*   **C#:** `Task.Run()` (Cuidado com Context).
*   **Go:** `go func() { ... }()` ou Channels.
    *   **Importante:** Use `context.Background()` na goroutine, senão ela morre quando o request HTTP termina!

---

### A. Tratamento de Erro (`if err != nil`)
Go não tem `try/catch` para lógica de negócio. Funções retornam erros como valores.
```go
file, err := os.Open("teste.txt")
if err != nil {
    return err // Propague o erro explicitamente
}
```

### B. Interfaces Implícitas (Duck Typing)
Você não declara `class MeuServico implements IService`.
Se sua struct tiver os métodos da interface, ela **automaticamente** implementa a interface.

### C. `defer`
Executa algo no final da função (útil para `Dispose`/`finally`).
```go
f, _ := os.Open("arq")
defer f.Close() // Será executado quando a função sair
```

---

## 7. Testes: Adeus Assertions Mágicas (`*testing.T`)

Em C#, usamos `Assert.Equal(a, b)`. Em Go, nós escrevemos a lógica de comparação (`if a != b`).
O objeto `t` é o seu **Controlador de Teste** injetado.

| C# (xUnit/NUnit) | Go (`testing` stdlib) | Comportamento |
| :--- | :--- | :--- |
| `Assert.Equal(exp, act)` | `if exp != act { t.Errorf(...) }` | **Soft Fail:** Marca erro, mas continua executando o teste. |
| `Assert.NotNull(obj)` | `if obj == nil { t.Fatalf(...) }` | **Hard Fail:** Para o teste na hora (útil se o próximo passo causaria Panic). |

---

## 8. Zero Values (O "Null" Diferente)

Em C#, se você não inicializa um objeto, ele é `null`. Em Go, variáveis sempre começam com um valor "utilizável".

| Tipo | Zero Value em Go | Em C# seria... |
| :--- | :--- | :--- |
| `int` / `float` | `0` | `0` |
| `string` | `""` (vazia, não nil) | `null` |
| `pointer` / `interface` | `nil` | `null` |

---

## 9. Contexto: Timeouts e Cancelamento

O `context.Context` é vital. É o equivalente ao `CancellationToken` do C#.

```go
// Define um limite de 2 segundos
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel() 

// Se demorar > 2s, o Go aborta a conexão
err := repo.Save(ctx, dados) 
```

## 10. Generics (Go 1.18+)

Sintaxe com colchetes `[]`.

```go
// T any = T pode ser qualquer coisa
func PrintSlice[T any](s []T) { ... }
```
