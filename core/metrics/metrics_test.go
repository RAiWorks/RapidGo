package metrics

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() { gin.SetMode(gin.TestMode) }

// setup creates an isolated Metrics instance backed by a fresh registry,
// and returns the Metrics, a Gin engine wired with the middleware, and a
// handler function that scrapes the isolated registry.
func setup() (*Metrics, *gin.Engine, gin.HandlerFunc) {
	reg := prometheus.NewRegistry()
	m := newMetrics(reg)
	e := gin.New()
	e.Use(m.Middleware())

	h := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	scrape := func(c *gin.Context) { h.ServeHTTP(c.Writer, c.Request) }
	e.GET("/metrics", scrape)

	return m, e, scrape
}

// scrape performs a GET /metrics and returns the body text.
func scrape(e *gin.Engine) string {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	e.ServeHTTP(w, req)
	return w.Body.String()
}

// T01: New returns non-nil Metrics with collectors registered
func TestNew(t *testing.T) {
	reg := prometheus.NewRegistry()
	m := newMetrics(reg)
	if m == nil {
		t.Fatal("newMetrics returned nil")
	}
	if m.requests == nil || m.duration == nil || m.size == nil {
		t.Fatal("one or more collectors are nil")
	}
	families, err := reg.Gather()
	if err != nil {
		t.Fatalf("gather failed: %v", err)
	}
	// No observations yet, but metric families should be registered.
	// Counter/Histogram families only appear after first observation,
	// so just verify the registration didn't panic.
	_ = families
}

// T02: Middleware increments http_requests_total
func TestMiddleware_IncrementsCounter(t *testing.T) {
	_, e, _ := setup()
	e.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })

	w := httptest.NewRecorder()
	e.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/ping", nil))

	body := scrape(e)
	if !strings.Contains(body, "http_requests_total") {
		t.Fatal("http_requests_total not found in metrics output")
	}
}

// T03: Middleware records http_request_duration_seconds
func TestMiddleware_RecordsDuration(t *testing.T) {
	_, e, _ := setup()
	e.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })

	w := httptest.NewRecorder()
	e.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/ping", nil))

	body := scrape(e)
	if !strings.Contains(body, "http_request_duration_seconds") {
		t.Fatal("http_request_duration_seconds not found in metrics output")
	}
}

// T04: Middleware records http_response_size_bytes
func TestMiddleware_RecordsSize(t *testing.T) {
	_, e, _ := setup()
	e.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })

	w := httptest.NewRecorder()
	e.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/ping", nil))

	body := scrape(e)
	if !strings.Contains(body, "http_response_size_bytes") {
		t.Fatal("http_response_size_bytes not found in metrics output")
	}
}

// T05: Counter has correct method label
func TestMiddleware_LabelsMethod(t *testing.T) {
	_, e, _ := setup()
	e.POST("/submit", func(c *gin.Context) { c.String(200, "ok") })

	w := httptest.NewRecorder()
	e.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/submit", nil))

	body := scrape(e)
	if !strings.Contains(body, `method="POST"`) {
		t.Fatal("method label POST not found")
	}
}

// T06: Counter has correct path label using route template
func TestMiddleware_LabelsPath(t *testing.T) {
	_, e, _ := setup()
	e.GET("/users/:id", func(c *gin.Context) { c.String(200, "user") })

	w := httptest.NewRecorder()
	e.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/users/42", nil))

	body := scrape(e)
	if !strings.Contains(body, `path="/users/:id"`) {
		t.Fatal("path label /users/:id not found")
	}
	if strings.Contains(body, `path="/users/42"`) {
		t.Fatal("raw path /users/42 should not appear — must use route template")
	}
}

// T07: Counter has correct status label
func TestMiddleware_LabelsStatus(t *testing.T) {
	_, e, _ := setup()
	e.GET("/notfound", func(c *gin.Context) { c.String(404, "nope") })

	w := httptest.NewRecorder()
	e.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/notfound", nil))

	body := scrape(e)
	if !strings.Contains(body, `status="404"`) {
		t.Fatal("status label 404 not found")
	}
}

// T08: Unmatched routes use "unmatched" as path label
func TestMiddleware_UnmatchedPath(t *testing.T) {
	_, e, _ := setup()
	// Don't register /nowhere — it will be unmatched

	w := httptest.NewRecorder()
	e.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/nowhere", nil))

	body := scrape(e)
	if !strings.Contains(body, `path="unmatched"`) {
		t.Fatal("unmatched path label not found")
	}
}

// T09: Handler returns Prometheus text format with registered metrics
func TestHandler(t *testing.T) {
	_, e, _ := setup()
	e.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })

	// Generate some data
	w := httptest.NewRecorder()
	e.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/ping", nil))

	body := scrape(e)
	if !strings.Contains(body, "http_requests_total") {
		t.Fatal("handler output missing http_requests_total")
	}
	if !strings.Contains(body, "http_request_duration_seconds") {
		t.Fatal("handler output missing http_request_duration_seconds")
	}
	if !strings.Contains(body, "http_response_size_bytes") {
		t.Fatal("handler output missing http_response_size_bytes")
	}
}

// T10: Multiple requests correctly accumulate counters
func TestMiddleware_MultipleRequests(t *testing.T) {
	_, e, _ := setup()
	e.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })

	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		e.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/ping", nil))
	}

	body := scrape(e)
	// Counter should show 5.0 for /ping GET 200
	// The /metrics scrape itself also generates a request, but it's a separate label set.
	if !strings.Contains(body, `http_requests_total{method="GET",path="/ping",status="200"} 5`) {
		t.Fatalf("expected counter value 5 for /ping, got:\n%s", body)
	}
}

// T11: Different status codes produce distinct label sets
func TestMiddleware_DifferentStatusCodes(t *testing.T) {
	_, e, _ := setup()
	e.GET("/ok", func(c *gin.Context) { c.String(200, "ok") })
	e.GET("/fail", func(c *gin.Context) { c.String(500, "error") })

	w := httptest.NewRecorder()
	e.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/ok", nil))
	w = httptest.NewRecorder()
	e.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/fail", nil))

	body := scrape(e)
	if !strings.Contains(body, `status="200"`) {
		t.Fatal("status 200 label not found")
	}
	if !strings.Contains(body, `status="500"`) {
		t.Fatal("status 500 label not found")
	}
}

// T12: Handler output includes default Go runtime metrics
func TestHandler_IncludesGoMetrics(t *testing.T) {
	// Use the default registry which includes Go collector
	e := gin.New()
	e.GET("/metrics", Handler())

	w := httptest.NewRecorder()
	e.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/metrics", nil))

	body := w.Body.String()
	if !strings.Contains(body, "go_goroutines") {
		t.Fatal("go_goroutines not found in default handler output")
	}
}
