package client

import (
	"context"
	"encoding/json"
	"github.com/portnyagin/practicum_project/internal/app/dto"
	"github.com/portnyagin/practicum_project/internal/app/infrastructure"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	AccrualClientRequestTimeout = 25 * time.Second
	AccrualClientURL            = "/api/orders/"
)

type AccrualClient struct {
	serviceAddress string
	log            *infrastructure.Logger
	client         *http.Client
}

func NewAccrualClient(serviceAddress string, log *infrastructure.Logger) *AccrualClient {
	var target AccrualClient
	target.log = log
	target.serviceAddress = serviceAddress
	target.client = &http.Client{Timeout: AccrualClientRequestTimeout}

	return &target
}

func (c *AccrualClient) GetAccrual(ctx context.Context, orderNum string) (*dto.Accrual, error) {
	address := c.serviceAddress + AccrualClientURL + orderNum

	u, err := url.Parse(address)
	if err != nil {
		c.log.Error("AccrualClient: GetAccrual. Can't build url", zap.Error(err))
		return nil, err
	}
	if u.Scheme != "http" {
		u.Scheme = "http"
	}
	if u.Host == "" {
		u.Host = "localhost"
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		c.log.Error("AccrualClient: GetAccrual. Can't build request", zap.Error(err))
		return nil, err
	}
	req.Header.Add("Accept", `application/json`)
	resp, err := c.client.Do(req)
	if err != nil {
		c.log.Error("AccrualClient: GetAccrual. Can't execute request", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		var accrual dto.Accrual
		err = json.Unmarshal(body, &accrual)
		if err != nil {
			c.log.Error("AccrualClient: GetAccrual. Can't unmarshal request  body", zap.Error(err))
			return nil, err
		}
		return &accrual, nil
	} else if resp.StatusCode == http.StatusTooManyRequests {
		c.log.Error("AccrualClient: GetAccrual.Too many requests:", zap.Int("statusCode", resp.StatusCode))
		return nil, dto.ErrTooManyRequest
	} else {
		c.log.Error("AccrualClient: GetAccrual.Unexpected response from remote service:", zap.Int("statusCode", resp.StatusCode))
		return nil, dto.ErrRemoteServiceError
	}
}
