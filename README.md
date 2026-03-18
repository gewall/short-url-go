# Short URL Service

A simple URL shortener built with Go, PostgreSQL and the chi router. Provides a small REST API to create short links, redirect to the original URL, and fetch basic statistics.

## Tech
- Go (net/http + chi)
- PostgreSQL
- SQL migrations 

## Features
- Create short URLs for given targets
- HTTP redirect from short code to original URL

## Project layout (suggested)
- cmd/server — application entrypoint
- internal/handler — HTTP handlers
- internal/repository — DB access
- internal/service — Services
- migrations — SQL migrations
