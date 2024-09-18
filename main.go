package main

import (
    "errors"
    "flag"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
    "strings"
)

func main() {
    var (
        target  string
        port    string
        headers http.Header = http.Header{}
    )

    flag.StringVar(&target, "url", "", "The URL to forward requests to")
    flag.StringVar(&port, "port", "8080", "The port to run the proxy server on")
    
    flag.Func("header", "Add headers in key:value format (can be used multiple times)", func(s string) error {
        parts := strings.SplitN(s, ":", 2)

        if len(parts) != 2 {
            return errors.New("invalid header format")
        }

        key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])

        headers.Add(key, value)

        return nil
    })
    
	flag.Parse()

    if target == "" {
        log.Fatal("You must specify a URL with the -url flag.")
    }

    targetURL, err := url.Parse(target)
    if err != nil {
        log.Fatalf("Failed to parse target URL: %s", err)
    }

    proxy := httputil.NewSingleHostReverseProxy(targetURL)

    originalDirector := proxy.Director
    proxy.Director = func(req *http.Request) {
        originalDirector(req)

        req.Host = targetURL.Host
        req.URL.Host = targetURL.Host
        req.URL.Scheme = targetURL.Scheme
        
        for key, values := range headers {
            for _, value := range values {
                req.Header.Set(key, value)
            }
        }
    }

    http.Handle("/", proxy)

    log.Printf("Starting proxy server on port %s forwarding to %s", port, target)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}
