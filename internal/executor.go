package internal

import (
	"fmt"
	"sync"
	"time"

	"mini-k6/models"
	"mini-k6/results"

	"github.com/go-resty/resty/v2"
)

func ExecuteScenario(scenario models.TestScenario) []results.TestResult {
	var wg sync.WaitGroup
	resultsChan := make(chan results.TestResult, len(scenario.Requests)*scenario.ConcurrentUsers*10)

	fmt.Printf("🔧 Iniciando cenário com %d usuários concorrentes por %ds...\n", scenario.ConcurrentUsers, scenario.DurationSeconds)

	for i := 0; i < scenario.ConcurrentUsers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			client := resty.New()
			client.SetTimeout(30 * time.Second)
			end := time.Now().Add(time.Duration(scenario.DurationSeconds) * time.Second)
			reqCount := 0

			fmt.Printf("👤 Usuário #%d começou\n", id)

			for time.Now().Before(end) {
				if scenario.MaxRequests > 0 && reqCount >= scenario.MaxRequests {
					break
				}
				for _, req := range scenario.Requests {
					if scenario.MaxRequests > 0 && reqCount >= scenario.MaxRequests {
						break
					}

					start := time.Now()
					status, errMsg := makeRequest(client, req)
					duration := time.Since(start)

					resultsChan <- results.TestResult{
						Method:     req.Method,
						URL:        req.URL,
						StatusCode: status,
						Error:      errMsg,
						Duration:   duration,
					}

					fmt.Printf("🛰️  [%s] %s -> %d (%v)\n", req.Method, req.URL, status, duration)

					reqCount++
				}
				time.Sleep(10 * time.Millisecond)
			}

			fmt.Printf("✅ Usuário #%d finalizou (total de requisições: %d)\n", id, reqCount)
		}(i + 1)
	}

	wg.Wait()
	close(resultsChan)

	var finalResults []results.TestResult
	for r := range resultsChan {
		finalResults = append(finalResults, r)
	}

	fmt.Println("📊 Teste finalizado.")
	return finalResults
}

func makeRequest(client *resty.Client, req models.APIRequest) (int, string) {
	r := client.R()

	for k, v := range req.Headers {
		r.SetHeader(k, v)
	}

	if req.Body != nil {
		r.SetBody(req.Body)
	}

	var resp *resty.Response
	var err error

	switch req.Method {
	case "GET":
		resp, err = r.Get(req.URL)
	case "POST":
		resp, err = r.Post(req.URL)
	case "PUT":
		resp, err = r.Put(req.URL)
	case "DELETE":
		resp, err = r.Delete(req.URL)
	default:
		return 0, "❌ Método inválido"
	}

	if err != nil {
		return 0, err.Error()
	}

	return resp.StatusCode(), ""
}
