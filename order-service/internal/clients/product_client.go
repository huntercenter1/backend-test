package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ProductClient interface {
	Get(ctx context.Context, id string) (*Product, error)
	ApplyStockDelta(ctx context.Context, id string, delta int) (*Product, error)
}

type productClient struct {
	base string
	hc   *http.Client
}

type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

func NewProductClient(base string) ProductClient {
	return &productClient{
		base: base,
		hc:   &http.Client{ Timeout: 5 * time.Second },
	}
}

func (c *productClient) Get(ctx context.Context, id string) (*Product, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/products/%s", c.base, id), nil)
	res, err := c.hc.Do(req)
	if err != nil { return nil, err }
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("product get status %d", res.StatusCode)
	}
	var p Product
	return &p, json.NewDecoder(res.Body).Decode(&p)
}

func (c *productClient) ApplyStockDelta(ctx context.Context, id string, delta int) (*Product, error) {
	body, _ := json.Marshal(map[string]int{"delta": delta})
	req, _ := http.NewRequestWithContext(ctx, http.MethodPut, fmt.Sprintf("%s/products/%s/stock", c.base, id), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := c.hc.Do(req)
	if err != nil { return nil, err }
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("stock update status %d", res.StatusCode)
	}
	var p Product
	return &p, json.NewDecoder(res.Body).Decode(&p)
}
