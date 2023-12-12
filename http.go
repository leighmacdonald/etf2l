package etf2l

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/time/rate"
	"io"
	"net/http"
	"time"
)

const (
	maxBucket     = 60
	limitInterval = 10 * time.Second
)

var (
	ErrNotFound = errors.New("Not found (404)")
)

type LimiterClient struct {
	*http.Client
	*rate.Limiter
}

func (c *LimiterClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	if errWait := c.Wait(ctx); errWait != nil {
		return nil, errors.Wrap(errWait, "Failed to wait for request")
	}

	resp, errDo := c.Client.Do(req)
	if errDo != nil {
		return nil, errors.Wrap(errDo, "Failed to make request")
	}

	return resp, nil
}

func newHTTPClient() *LimiterClient {
	return &LimiterClient{
		Client:  http.DefaultClient,
		Limiter: rate.NewLimiter(rate.Every(limitInterval), maxBucket),
	}
}

func fullURL(path string) string {
	return fmt.Sprintf("https://api-v2.etf2l.org%s", path)
}

type Client struct {
	*LimiterClient
}

func New() *Client {
	return &Client{LimiterClient: newHTTPClient()}
}

func (client *Client) call(ctx context.Context, method string, path string, body any, receiver any) error {
	var reqBody io.Reader

	if body != nil {
		rb, errMarshal := json.Marshal(body)
		if errMarshal != nil {
			return errors.Wrap(errMarshal, "Failed to marshal payload")
		}

		reqBody = bytes.NewReader(rb)
	}

	req, errReq := http.NewRequestWithContext(ctx, method, fullURL(path), reqBody)
	if errReq != nil {
		return errors.Wrap(errReq, "Failed to create request")
	}

	req.Header.Add("Content-Type", `application/json`)
	req.Header.Add("Accept", "application/json")

	resp, errResp := client.Do(ctx, req)
	if errResp != nil {
		return errors.Wrap(errResp, "Failed to call endpoint")
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusTooManyRequests {
		return errors.New("Rate limited")
	}

	if !(resp.StatusCode >= http.StatusOK && resp.StatusCode <= http.StatusIMUsed) {
		if resp.StatusCode == http.StatusNotFound {
			return ErrNotFound
		}

		return errors.Errorf("Invalid status code: %s", resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	if errJSON := decoder.Decode(&receiver); errJSON != nil {
		return errors.Wrap(errJSON, "Failed to unmarshal json payload")
	}

	return nil
}
