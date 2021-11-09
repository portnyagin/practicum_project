package client

import (
	"context"
	"github.com/portnyagin/practicum_project/internal/app/infrastructure"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"time"
)

const (
	GophermartRequestTimeout = 25 * time.Second
	GophermartClientURL      = "/api/accrual/process/"
)

type GophermartClient struct {
	serviceAddress string
	log            *infrastructure.Logger
	client         *http.Client
}

func NewGophermartClient(serviceAddress string, log *infrastructure.Logger) *GophermartClient {
	var target GophermartClient
	target.log = log
	target.serviceAddress = serviceAddress
	target.client = &http.Client{Timeout: GophermartRequestTimeout}
	return &target
}

func (c *GophermartClient) ProcessRequest(ctx context.Context, orderNum string) bool {
	address := c.serviceAddress + GophermartClientURL + orderNum
	if c.serviceAddress == ":8080" {
		address = "localhost" + address
	}
	u, err := url.Parse(address)
	if err != nil {
		c.log.Error("GophermartClient: ProcessRequest. Can't build url", zap.Error(err))
		return false
	}
	if u.Scheme != "http" {
		u.Scheme = "http"
	}
	if u.Host == "" {
		u.Host = "localhost"
	}

	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		c.log.Error("GophermartClient: ProcessRequest. Can't build request", zap.Error(err))
		return false
	}
	req.Header.Add("Accept", `application/json`)
	resp, err := c.client.Do(req)
	if err != nil {
		c.log.Error("GophermartClient: ProcessRequest. Can't execute request", zap.Error(err))
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		return true
	}
	c.log.Error("GophermartClient: ProcessRequest. Can't process request", zap.Int("StatusCode", resp.StatusCode))
	return false
}
