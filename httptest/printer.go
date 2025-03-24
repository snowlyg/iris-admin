package httptest

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/gorilla/websocket"
)

// DebugPrinter implements Printer and WebsocketPrinter.
// Uses net/http/httputil to dump both requests and responses.
// Also prints all websocket messages.
type DebugPrinter struct {
	logger httpexpect.Logger
	body   bool
}

// NewDebugPrinter returns a new DebugPrinter given a logger and body
// flag. If body is true, request and response body is also printed.
func NewDebugPrinter(logger httpexpect.Logger, body bool) DebugPrinter {
	return DebugPrinter{logger, body}
}

// Request implements Printer.Request.
func (p DebugPrinter) Request(req *http.Request) {
	if req == nil {
		return
	}

	if req.Header.Get("Content-Type") == "text/html" || req.URL.Path == "/api/v1/file/upload" {
		p.body = false
	}

	dump, err := httputil.DumpRequest(req, p.body)
	if err != nil {
		panic(err)
	}
	p.logger.Logf("%s", dump)
}

// Response implements Printer.Response.
func (p DebugPrinter) Response(resp *http.Response, duration time.Duration) {
	if resp == nil {
		return
	}
	if resp.Header.Get("Content-Type") == "image/jpeg" || resp.Header.Get("Content-Type") == "image/png" || resp.Header.Get("Content-Type") == "text/html" {
		p.body = false
	}

	dump, err := httputil.DumpResponse(resp, p.body)
	if err != nil {
		panic(err)
	}

	text := strings.Replace(string(dump), "\r\n", "\n", -1)
	lines := strings.SplitN(text, "\n", 2)

	p.logger.Logf("%s %s\n%s", lines[0], duration, lines[1])
}

// WebsocketWrite implements WebsocketPrinter.WebsocketWrite.
func (p DebugPrinter) WebsocketWrite(typ int, content []byte, closeCode int) {
	b := &bytes.Buffer{}
	fmt.Fprintf(b, "-> Sent: %s", wsMessageType(typ))
	if typ == websocket.CloseMessage {
		fmt.Fprintf(b, " %s", wsCloseCode(closeCode))
	}
	fmt.Fprint(b, "\n")
	if len(content) > 0 {
		if typ == websocket.BinaryMessage {
			fmt.Fprintf(b, "%v\n", content)
		} else {
			fmt.Fprintf(b, "%s\n", content)
		}
	}
	fmt.Fprintf(b, "\n")
	p.logger.Logf(b.String())
}

// WebsocketRead implements WebsocketPrinter.WebsocketRead.
func (p DebugPrinter) WebsocketRead(typ int, content []byte, closeCode int) {
	b := &bytes.Buffer{}
	fmt.Fprintf(b, "<- Received: %s", wsMessageType(typ))
	if typ == websocket.CloseMessage {
		fmt.Fprintf(b, " %s", wsCloseCode(closeCode))
	}
	fmt.Fprint(b, "\n")
	if len(content) > 0 {
		if typ == websocket.BinaryMessage {
			fmt.Fprintf(b, "%v\n", content)
		} else {
			fmt.Fprintf(b, "%s\n", content)
		}
	}
	fmt.Fprintf(b, "\n")
	p.logger.Logf(b.String())
}

type wsMessageType int

func (wmt wsMessageType) String() string {
	s := "unknown"

	switch wmt {
	case websocket.TextMessage:
		s = "text"
	case websocket.BinaryMessage:
		s = "binary"
	case websocket.CloseMessage:
		s = "close"
	case websocket.PingMessage:
		s = "ping"
	case websocket.PongMessage:
		s = "pong"
	}

	return fmt.Sprintf("%s(%d)", s, wmt)
}

type wsCloseCode int

// https://developer.mozilla.org/en-US/docs/Web/API/CloseEvent/code
func (wcc wsCloseCode) String() string {
	s := "Unknown"

	switch wcc {
	case 1000:
		s = "NormalClosure"
	case 1001:
		s = "GoingAway"
	case 1002:
		s = "ProtocolError"
	case 1003:
		s = "UnsupportedData"
	case 1004:
		s = "Reserved"
	case 1005:
		s = "NoStatusReceived"
	case 1006:
		s = "AbnormalClosure"
	case 1007:
		s = "InvalidFramePayloadData"
	case 1008:
		s = "PolicyViolation"
	case 1009:
		s = "MessageTooBig"
	case 1010:
		s = "MandatoryExtension"
	case 1011:
		s = "InternalServerError"
	case 1012:
		s = "ServiceRestart"
	case 1013:
		s = "TryAgainLater"
	case 1014:
		s = "BadGateway"
	case 1015:
		s = "TLSHandshake"
	}

	return fmt.Sprintf("%s(%d)", s, wcc)
}
