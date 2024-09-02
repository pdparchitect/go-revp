package main

import (
    "flag"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
)

func main() {
    var (
        target string
        port   string
    )

    flag.StringVar(&target, "url", "", "The URL to forward requests to")
    flag.StringVar(&port, "port", "8080", "The port to run the proxy server on")

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
    }

    http.Handle("/", proxy)

    log.Printf("Starting proxy server on port %s forwarding to %s", port, target)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}
