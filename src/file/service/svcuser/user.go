package svcuser

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/whoisnian/k8s-example/src/file/global"
	"github.com/whoisnian/k8s-example/src/file/model"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
)

var (
	httpClient *http.Client
	baseURL    *url.URL
)

func Setup(endpoint string) {
	httpClient = &http.Client{
		Timeout:   time.Second * 5,
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	var err error
	if baseURL, err = url.Parse(endpoint); err != nil {
		panic(err)
	}
}

func UserInfo(c *gin.Context) (*model.User, error) {
	ctx, span := global.Tracer.Start(c.Request.Context(), "svcuser.UserInfo")
	defer span.End()

	u := baseURL.ResolveReference(&url.URL{Path: "/internal/user/info"})
	req := (&http.Request{
		Method:     http.MethodGet,
		URL:        u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Body:       nil,
		Host:       u.Host,
	}).WithContext(ctx)
	if c.GetHeader("Cookie") != "" {
		req.Header.Set("Cookie", c.GetHeader("Cookie"))
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	user := &model.User{}
	if resp.StatusCode == http.StatusUnauthorized {
		span.SetAttributes(attribute.Int64("user_id", user.ID), attribute.Int("status", resp.StatusCode))
		return user, nil // user.ID == 0 if unauthorized
	} else if resp.StatusCode != http.StatusOK {
		span.SetAttributes(attribute.Int64("user_id", -1), attribute.Int("status", resp.StatusCode))
		return nil, errors.New(resp.Status)
	} else {
		err = json.NewDecoder(resp.Body).Decode(user)
		span.SetAttributes(attribute.Int64("user_id", user.ID), attribute.Int("status", resp.StatusCode))
		return user, err
	}
}
