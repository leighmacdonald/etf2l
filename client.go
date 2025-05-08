package etf2l

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/pkg/errors"
)

var (
	ErrNotFound  = errors.New("Not found (404)")
	ErrEOF       = errors.New("End of results")
	ErrNoResults = errors.New("no rows in result set")
)

type Recursive interface {
	IsRecursive() bool
}

type PagedResult interface {
	NextURL(r Recursive) (string, error)
}

type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func fullURL(path string) string {
	return fmt.Sprintf("https://api-v2.etf2l.org%s", path)
}

type Client struct {
	sync.RWMutex
}

func New() *Client {
	return &Client{}
}

type HTTPExecutor interface {
	Do(req *http.Request) (*http.Response, error)
}

func (client *Client) call(ctx context.Context, httpClient HTTPExecutor, path string, body any, receiver any) error {
	client.Lock()
	defer client.Unlock()

	var reqBody io.Reader

	if body != nil {
		rb, errMarshal := json.Marshal(body)
		if errMarshal != nil {
			return errors.Wrap(errMarshal, "Failed to marshal payload")
		}

		reqBody = bytes.NewReader(rb)
	}

	req, errReq := http.NewRequestWithContext(ctx, http.MethodGet, fullURL(path), reqBody)
	if errReq != nil {
		return errors.Wrap(errReq, "Failed to create request")
	}

	req.Header.Add("Content-Type", `application/json`)
	req.Header.Add("Accept", "application/json")

	resp, errResp := httpClient.Do(req)
	if errResp != nil {
		return errors.Wrap(errResp, "Failed to call endpoint")
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusTooManyRequests {
		return errors.New("Rate limited")
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode > http.StatusIMUsed {
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
