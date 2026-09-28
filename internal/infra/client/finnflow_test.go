package client_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"api-creditos-consulta/internal/infra/client"
)

func TestObtenerCreditosPorRut_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/core/v1/creditos", r.URL.Path)
		assert.Equal(t, "key", r.Header.Get("X-Client-Id"))
		assert.Equal(t, "secret", r.Header.Get("X-Client-Secret"))
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data": [{"NumeroOperacion": "123"}]}`))
	}))
	defer server.Close()

	c := client.NewFinnFlowClient(server.URL, "key", "secret")
	res, err := c.ObtenerCreditosPorRut(context.Background(), "1-9")

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "123", res[0].NumeroOperacion)
}

func TestObtenerCreditosPorRut_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	c := client.NewFinnFlowClient(server.URL, "key", "secret")
	res, err := c.ObtenerCreditosPorRut(context.Background(), "1-9")

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res, 0)
}

func TestObtenerCreditosPorRut_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`Internal Server Error`))
	}))
	defer server.Close()

	c := client.NewFinnFlowClient(server.URL, "key", "secret")
	res, err := c.ObtenerCreditosPorRut(context.Background(), "1-9")

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "FinnFlow consulta failed: 500")
}

func TestObtenerCreditosPorRut_DecodeError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`))
	}))
	defer server.Close()

	c := client.NewFinnFlowClient(server.URL, "key", "secret")
	res, err := c.ObtenerCreditosPorRut(context.Background(), "1-9")

	assert.Error(t, err)
	assert.Nil(t, res)
}
