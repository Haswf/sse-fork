/*
 * Copyright 2023 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package sse

import (
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/network"
)

// Server-Sent Events
// Last Updated 31 March 2023
// https://html.spec.whatwg.org/multipage/server-sent-events.html#server-sent-events

const (
	ContentType  = "text/event-stream"
	noCache      = "no-cache"
	cacheControl = "Cache-Control"
	LastEventID  = "Last-Event-ID"
)

var fieldReplacer = strings.NewReplacer(
	"\n", "\\n",
	"\r", "\\r")

var dataReplacer = strings.NewReplacer(
	"\n", "\ndata:",
	"\r", "\\r")

type Event struct {
	Event string
	ID    string
	Retry uint64
	Data  []byte
}

// GetLastEventID retrieve Last-Event-ID header if present.
func GetLastEventID(c *app.RequestContext) string {
	return c.Request.Header.Get(LastEventID)
}

type Stream struct {
	w network.ExtWriter
}

// NewStream creates a new stream for publishing Event.
func NewStream(c *app.RequestContext) *Stream {
	c.Response.Header.SetContentType(ContentType)
	if c.Response.Header.Get(cacheControl) == "" {
		c.Response.Header.Set(cacheControl, noCache)
	}

	writer := NewStreamBodyWriter(&c.Response, c.GetWriter())
	c.Response.HijackWriter(writer)
	return &Stream{
		writer,
	}
}

// Publish push an event to client.
func (c *Stream) Publish(event *Event) error {
	err := Encode(c.w, event)
	if err != nil {
		return err
	}
	return c.w.Flush()
}
