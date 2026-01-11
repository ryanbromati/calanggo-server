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

### "While"
```go
for contador < 10 { ... }
```

### LINQ? Não existe.
Não procure `.Where().Select()`. A filosofia do Go é "Loops explícitos são melhores que mágica oculta".
*   **Filter:** Faça um `for` com `if` e `append`.
*   **Map:** Faça um `for` com `append`.

> **Dica:** O pacote `slices` (Go 1.21+) tem `slices.Contains`, `slices.Sort`, etc.

---

## 4. Estruturas e Dados

### `make` vs `new`
*   **`make([]int, 0, 10)`**: Inicializa estruturas complexas (Slices, Maps, Channels). É o `new List<int>(10)`.
*   **Struct Tags (`json:"..."`)**: Equivalente aos Attributes `[JsonPropertyName]`. Define metadados para bibliotecas (JSON, ORMs).

### Slice (Lista Dinâmica)
Um Slice é uma "janela" para um array. É leve e rápido.
```go
lista := make([]string, 0, 5) // Len 0, Cap 5
lista = append(lista, "Item") // Adiciona e retorna a nova referência slice
```

---

## 5. Web & Concorrência

### Mux (`http.NewServeMux`)
É o Roteador (`app.MapControllers`).
*   **Wildcard `{id}`**: `mux.HandleFunc("GET /{id}", ...)` -> `r.PathValue("id")`.

### Fire-and-Forget
Executar tarefa sem travar o request.
*   **C#:** `Task.Run()` (Cuidado com Context).
*   **Go:** `go func() { ... }()` ou Channels.
    *   **Importante:** Use `context.Background()` na goroutine, senão ela morre quando o request HTTP termina!

---

## 6. O Que Você Vai Encontrar em Breve (Antecipando Dúvidas)

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
Se sua struct tiver os métodos da interface, ela **automaticamente** implementa a interface. Isso permite criar interfaces para classes que você nem é dono!

### C. `defer`
Executa algo no final da função (útil para `Dispose`/`finally`).
```go
f, _ := os.Open("arq")
defer f.Close() // Será executado quando a função sair, não importa como
```

### D. Pacotes Públicos vs Privados
*   **Letra Maiúscula:** Público (Exportado). `func Save()`
*   **Letra Minúscula:** Privado (Internal). `func save()`
Isso vale para funcões, structs, campos e constantes.

---

## 7. Testes: Adeus Assertions Mágicas (`*testing.T`)

Em C#, usamos `Assert.Equal(a, b)`. Em Go, nós escrevemos a lógica de comparação (`if a != b`).
O objeto `t` é o seu **Controlador de Teste** injetado.

| C# (xUnit/NUnit) | Go (`testing` stdlib) | Comportamento |
| :--- | :--- | :--- |
| `Assert.Equal(exp, act)` | `if exp != act { t.Errorf(...) }` | **Soft Fail:** Marca erro, mas continua executando o teste. |
| `Assert.NotNull(obj)` | `if obj == nil { t.Fatalf(...) }` | **Hard Fail:** Para o teste na hora (útil se o próximo passo causaria Panic). |
| `Console.WriteLine()` | `t.Logf(...)` | Só aparece no console se o teste falhar ou usar `go test -v`. |
| `[Theory] / [InlineData]` | `t.Run("nome", func(t...) {})` | **Table-Driven Tests:** Cria sub-testes nomeados dentro de um loop. |

---

## 8. Zero Values (O "Null" Diferente)

Em C#, se você não inicializa um objeto, ele é `null`. Em Go, variáveis sempre começam com um valor "utilizável", nunca lixo de memória.

| Tipo | Zero Value em Go | Em C# seria... |
| :--- | :--- | :--- |
| `int` / `float` | `0` | `0` |
| `bool` | `false` | `false` |
| `string` | `""` (vazia, não nil) | `null` |
| `pointer` | `nil` | `null` |
| `struct` | Todos os campos zerados | `null` |

**Cuidado:** Como saber se `Preco == 0` significa "Grátis" ou "Não informado"?
*   **Solução:** Use ponteiro (`*int`). Se for `nil`, não foi informado. Se for `0` (apontando para valor), é grátis.

---

## 9. Contexto: Timeouts e Cancelamento

O `context.Context` é vital para backends robustos. É o equivalente ao `CancellationToken` do C#, mas turbinado.

```go
// Define um limite de 2 segundos para tudo que usar este ctx
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel() // Sempre chame cancel para liberar recursos!

// Passa o ctx para o banco. Se demorar > 2s, o Go aborta a conexão.
err := repo.Save(ctx, dados) 
```
Se você não usar Context com Timeouts, seu servidor vai acumular conexões presas ("hanging") até cair.

---

## 10. Generics (Go 1.18+)

C# tem Generics poderosos. Go agora tem o básico. A sintaxe usa colchetes `[]`.

**Exemplo: Função que imprime qualquer slice**
```go
// T any = T pode ser qualquer coisa
func PrintSlice[T any](s []T) {
    for _, v := range s {
        fmt.Println(v)
    }
}

// Uso
PrintSlice[int]([]int{1, 2})
PrintSlice[string]([]string{"a", "b"})
```
*Nota:* Em Go, tente resolver primeiro com **Interfaces**. Use Generics apenas quando estiver escrevendo estruturas de dados ou algoritmos genéricos.
