package main

import (
	"context"
	"net/http"
	"time"

	"github.com/hertz-contrib/sse"

	"github.com/cloudwego/hertz/pkg/network"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func main() {
	h := server.Default()

	h.GET("/sse", func(ctx context.Context, c *app.RequestContext) {
		// client can tell server last event it received with Last-Event-ID header
		lastEventID := sse.GetLastEventID(c)
		hlog.CtxInfof(ctx, "last event ID: %s", lastEventID)

		// you must set status code and response headers before first render call
		c.SetStatusCode(http.StatusOK)

		sse.Stream(ctx, c, func(ctx context.Context, w network.ExtWriter) {
			// send a timestamp event to client with current time every second
			for t := range time.NewTicker(1 * time.Second).C {
				event := &sse.Event{
					Event: "timestamp",
					Data:  t.Format(time.RFC3339),
				}
				err := event.Render(w)
				if err != nil {
					return
				}
			}
		})
	})

	h.Spin()
}