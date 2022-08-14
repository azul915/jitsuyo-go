package sec11

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"time"
)

func Prac() {
	resp, err := http.Get("http://example.com")
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Println("something went wrong!!")
	}
	resp.Body.Close()

	type User struct {
		Name string
		Addr string
	}

	u := User{
		Name: "O'Reilly Japan",
		Addr: "東京都新宿区四谷坂町",
	}
	payload, err := json.Marshal(u)
	if err != nil {
		log.Fatal(err)
	}
	res, err := http.Post("http://example.com/", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Fatal(err)
	}
	res.Body.Close()

	cl := &http.Client{
		Timeout:   10 * time.Second,
		Transport: http.DefaultTransport,
	}

	nr, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
	if err != nil {
		log.Fatal(err)
	}
	r, err := cl.Do(nr)
	if err != nil {
		log.Fatal(err)
	}
	r.Body.Close()

	ctx := context.Background()
	ncr, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com", nil)
	if err != nil {
		log.Fatal(err)

	}
	ncr.Header.Add("Authorization", "foobar")
	rr, err := cl.Do(ncr)
	if err != nil {
		log.Fatal(err)
	}
	rr.Body.Close()

	client := &http.Client{
		Transport: &customRoundTripper{
			base: http.DefaultTransport,
		},
	}
	req, err := http.NewRequestWithContext(ctx, "GET", "http://example.com", nil)
	if err != nil {
		log.Fatal(err)
	}
	client.Do(req)

	client = &http.Client{
		Transport: &loggingRoundTripper{
			transport: http.DefaultTransport,
			logger:    nil,
		},
	}
	req, err = http.NewRequestWithContext(ctx, "GET", "http://example.com", nil)
	if err != nil {
		log.Fatal(err)
	}
	client.Do(req)

	client = &http.Client{
		Transport: &basicAuthRoundTripper{
			username: "username",
			password: "password",
			base:     http.DefaultTransport,
		},
	}
	req, _ = http.NewRequestWithContext(ctx, "GET", "http://example.com", nil)
	client.Do(req)
}

type customRoundTripper struct {
	base http.RoundTripper
}

func (c customRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := c.base.RoundTrip(req)
	return resp, err
}

type loggingRoundTripper struct {
	transport http.RoundTripper
	logger    func(string, ...interface{})
}

func (t *loggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.logger == nil {
		t.logger = log.Printf
	}

	start := time.Now()
	resp, err := t.transport.RoundTrip(req)
	if resp != nil {
		t.logger("%s %s %d %s, duration: %d", req.Method, req.URL.String(), resp.StatusCode, http.StatusText(resp.StatusCode), time.Since(start))
	}
	return resp, err
}

type basicAuthRoundTripper struct {
	base     http.RoundTripper
	username string
	password string
}

func (rt *basicAuthRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(rt.username, rt.password)
	return rt.base.RoundTrip(req)
}

type retryableRoundTripper struct {
	base     http.RoundTripper
	attempts int
	waitTime time.Duration
}

func (rt *retryableRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var (
		resp *http.Response
		err  error
	)
	for count := 0; count < rt.attempts; count++ {
		resp, err = rt.base.RoundTrip(req)
		if !rt.shouldRetry(resp, err) {
			return resp, err
		}
		select {
		case <-req.Context().Done():
			return nil, req.Context().Err()
		case <-time.After(rt.waitTime):
		}
	}
	return resp, err
}
func (rt *retryableRoundTripper) shouldRetry(resp *http.Response, err error) bool {
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Temporary() {
			return true
		}
	}
	if resp != nil {
		if resp.StatusCode == 429 || (500 <= resp.StatusCode && resp.StatusCode <= 504) {
			return true
		}
	}
	return false
}
