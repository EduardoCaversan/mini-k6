# 🔥 Load Test CLI em Go

Projeto profissional de **teste de carga e performance** escrito em **Go**, com foco em simplicidade, extensibilidade e resultados exportáveis para ferramentas como **Power BI**, **Excel** e **Grafana**.

## 📌 Visão Geral

Este projeto é uma ferramenta de linha de comando que executa múltiplas requisições HTTP concorrentes contra um endpoint configurado, mede a performance das respostas e gera um relatório estatístico com:

- Total de requisições
- Sucessos e falhas
- Tempo médio de resposta
- Tempo total de execução
- Status codes agrupados
- Requisições por segundo (`RPS`)
- Resultados exportáveis em JSON

## 🚀 Tecnologias Utilizadas

- [Golang](https://golang.org/)
- `net/http` para as requisições
- `sync.WaitGroup` e `goroutines` para concorrência
- JSON estruturado como saída (ideal para Power BI)
- `context.WithTimeout` para controle de tempo de resposta

## ⚙️ Como Usar

### Pré-requisitos

- Go 1.18+
- Internet para alcançar os endpoints

### Execução

Compile ou execute diretamente com:

```bash
go run main.go
````

### Exemplo de Entrada

```json
{
  "concurrent_users": 5,
  "duration_seconds": 10,
  "max_requests": 50,
  "requests": [
    {
      "method": "GET",
      "url": "https://jsonplaceholder.typicode.com/posts"
    },
    {
      "method": "GET",
      "url": "https://jsonplaceholder.typicode.com/posts/1"
    },
    {
      "method": "POST",
      "url": "https://jsonplaceholder.typicode.com/posts",
      "headers": {
        "Content-Type": "application/json"
      },
      "body": {
        "title": "foo",
        "body": "bar",
        "userId": 1
      }
    }
  ]
}
```

### Exemplo de Saída

```json
{
  "total": 100,
  "success": 98,
  "failures": 2,
  "average_seconds": 0.2123,
  "total_seconds": 21.23,
  "requests_per_second": 4.70,
  "by_status_code": {
    "200": 98,
    "500": 2
  }
}
```

## 📊 Análise dos Resultados

Esses dados podem ser facilmente exportados para Power BI, Excel ou qualquer outra ferramenta de BI para visualizações personalizadas.

* `average_seconds`: tempo médio por requisição
* `total_seconds`: soma de todas as durações
* `requests_per_second`: performance geral (quanto maior, melhor)
* `by_status_code`: breakdown por status HTTP

## 🧪 Implementações Técnicas

* Uso de `http.Client` com `Timeout`
* Tratamento de erros com mensagens detalhadas
* Cálculo de métricas com precisão de segundos (`float64`)
* Respostas armazenadas em slices com logs completos
* JSON de resumo no final com `Summary` agregado

## 📦 Estrutura do Projeto

```
.
├── main.go             # Código principal do executável
├── results.go          # Structs para armazenar os dados
├── README.md           # Este arquivo
```

## ✍️ Autor

**Eduardo Caversan**
[📧 educaversan.dev@gmail.com](mailto:educaversan.dev@gmail.com)
🌐 Desenvolvedor Fullstack | Software Engineer

---

Sinta-se à vontade para contribuir ou adaptar este projeto conforme suas necessidades!
