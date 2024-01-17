package http

import (
	"context"
	"fmt"

	"github.com/valyala/fasthttp"

	loggerApi "github.com/somatech1/mikros/apis/logger"
	"github.com/somatech1/mikros/components/definition"
	"github.com/somatech1/mikros/components/logger"
	"github.com/somatech1/mikros/components/plugin"
)

type Client struct {
	isAuthenticated bool
	plugin.Entry
}

func New() *Client {
	return &Client{}
}

func (c *Client) CanBeInitialized(options *plugin.CanBeInitializedOptions) bool {
	_, ok := options.Definitions.ServiceTypes()[definition.ServiceType_HTTP]
	return ok
}

func (c *Client) Initialize(_ context.Context, opt *plugin.InitializeOptions) error {
	c.isAuthenticated = !opt.Definitions.HTTP.DisableAuth
	return nil
}

func (c *Client) AddResponseHeader(ctx context.Context, key, value string) {
	if c.IsEnabled() {
		return
	}

	if c, ok := ctx.(*fasthttp.RequestCtx); ok {
		// We only accept a string 'value' here to avoid doing conversion
		// inside the handler.
		c.SetUserValue(fmt.Sprintf("handler-attribute-%s", key), value)
	}
}

func (c *Client) SetResponseCode(ctx context.Context, code int) {
	if c.IsEnabled() {
		return
	}

	if c, ok := ctx.(*fasthttp.RequestCtx); ok {
		c.SetUserValue("handler-response-code", code)
	}
}

func (c *Client) Fields() []loggerApi.Attribute {
	return []loggerApi.Attribute{
		logger.String("svc.http.auth", fmt.Sprintf("%t", c.isAuthenticated)),
	}
}
