package client

import (
	"context"
	"github.com/portnyagin/practicum_project/internal/app/infrastructure"
	"go.uber.org/zap"
	"net/http"
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
	c.log.Debug("GophermartClient: ", zap.String("serviceAddress", c.serviceAddress))
	address := c.serviceAddress + GophermartClientURL + orderNum
	if c.serviceAddress == ":8080" {
		address = "http://localhost" + address
	} else if c.serviceAddress == "localhost:8080" {
		address = "http://" + address
	}

	c.log.Debug("GophermartClient: ", zap.String("url", address))
	req, err := http.NewRequest("POST", address, nil)
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
