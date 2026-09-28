package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"api-creditos-consulta/internal/core/domain"
	"api-creditos-consulta/internal/core/ports"
)

type finnFlowClient struct {
	baseURL    string
	key        string
	secret     string
	httpClient *http.Client
}

func NewFinnFlowClient(baseURL, key, secret string) ports.FinnFlowClient {
	return &finnFlowClient{
		baseURL: baseURL,
		key:     key,
		secret:  secret,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// ObtenerCreditosPorRut acts as HTTP client towards ASICOM
func (f *finnFlowClient) ObtenerCreditosPorRut(ctx context.Context, rut string) ([]domain.Credito, error) {
	reqBody, _ := json.Marshal(map[string]string{"rut": rut})
	
	url := fmt.Sprintf("%s/api/core/v1/creditos", f.baseURL)
	slog.Debug("Requesting creditos from FinnFlow at: "+url, "logger", "cl.bancofalabella.mortgage.consulta.service.ThirdPartyApiService", "thread", "http-nio-8080-exec-1")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request to FinnFlow: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Client-Id", f.key)
	req.Header.Set("X-Client-Secret", f.secret)

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error calling FinnFlow: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		slog.Debug("No creditos found for RUT in FinnFlow (404)", "logger", "cl.bancofalabella.mortgage.consulta.service.ThirdPartyApiService", "thread", "http-nio-8080-exec-1")
		return []domain.Credito{}, nil
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		// For legacy parity: any HTTP error is wrapped as a specific formatted string
		return nil, fmt.Errorf("FinnFlow consulta failed: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		Data []domain.Credito `json:"data"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding FinnFlow response: %w", err)
	}

	slog.Debug(fmt.Sprintf("Consulta completed successfully. Found %d creditos.", len(result.Data)), "logger", "cl.bancofalabella.mortgage.consulta.service.ThirdPartyApiService", "thread", "http-nio-8080-exec-1")

	return result.Data, nil
}
