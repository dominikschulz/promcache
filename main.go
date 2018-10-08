package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"

	"github.com/lox/httpcache"
	"github.com/lox/httpcache/httplog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	defaultListen = ":9091"
	defaultDir    = "/tmp/cachedata"
)

var (
	listen   string
	useDisk  bool
	dir      string
	dumpHttp bool
	verbose  bool
	upstream string
	duration int
)

func main() {
	flag.StringVar(&listen, "listen", defaultListen, "the host and port to bind to")
	flag.StringVar(&dir, "dir", defaultDir, "the dir to store cache data in, implies -disk")
	flag.BoolVar(&useDisk, "disk", false, "whether to store cache data to disk")
	flag.BoolVar(&verbose, "v", false, "show verbose output and debugging")
	flag.BoolVar(&dumpHttp, "dumphttp", false, "dumps http requests and responses to stdout")
	flag.StringVar(&upstream, "upstream", "127.0.0.1:9090", "upstream host to connect to")
	flag.IntVar(&duration, "duration", 60, "forced cache duration")
	flag.Parse()

	if verbose {
		httpcache.DebugLogging = true
	}

	req := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "promcache_http_requests_total",
			Help: "Number of HTTP requests to the caching proxy",
		},
		[]string{"cache"},
	)
	prometheus.MustRegister(req)

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":9092", mux)
	}()

	proxy := &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = "http"
			r.URL.Host = upstream
			r.Host = r.URL.Host
			if verbose {
				fmt.Printf("UPSTREAM-REQUEST TO: %s\n", r.URL.String())
			}
			req.WithLabelValues("miss").Inc()
		},
		ModifyResponse: func(r *http.Response) error {
			r.Header.Add("Cache-Control", fmt.Sprintf("max-age=%d, public", duration))
			return nil
		},
	}
	handler := httpcache.NewHandler(httpcache.NewMemoryCache(), proxy)
	handler.Shared = true

	respLogger := httplog.NewResponseLogger(handler)
	respLogger.DumpRequests = dumpHttp
	respLogger.DumpResponses = dumpHttp
	respLogger.DumpErrors = dumpHttp

	log.Printf("proxy listening on http://%s", listen)
	log.Fatal(http.ListenAndServe(listen, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/v1/query") {
			respLogger.ServeHTTP(w, r)
			return
		}
		req.WithLabelValues("request").Inc()
		r.Header.Del("Pragma")
		r.Header.Del("Cache-Control")
		val := r.URL.Query()
		start, err := strconv.Atoi(val.Get("start"))
		if err != nil {
			fmt.Println(err)
		} else {
			nStart := strconv.Itoa(start - (start % duration))
			val.Set("start", nStart)
		}
		end, err := strconv.Atoi(val.Get("end"))
		if err != nil {
			fmt.Println(err)
		} else {
			nEnd := strconv.Itoa(end - (end % duration) + duration)
			val.Set("end", nEnd)
		}
		r.URL.RawQuery = val.Encode()

		respLogger.ServeHTTP(w, r)
	})))
}
