# Short URL Service

A simple URL shortener built with Go, PostgreSQL and the chi router. Provides a small REST API to create short links, redirect to the original URL, and fetch basic statistics.

## Tech
- Go (net/http + chi)
- PostgreSQL
- SQL migrations 

## Features
- Create short URLs for given targets
- HTTP redirect from short code to original URL
- Log access and usage activity

## Project layout 
- cmd/server — application entrypoint
- internal/handler — HTTP handlers
- internal/repository — DB access
- internal/service — Services
- migrations — SQL migrations
## 🛠️ API Endpoints

### 🔐 Auth
Public endpoints for access management and session control.
* **POST** `/signup` — Register a new user account.
* **POST** `/signin` — Authenticate user and receive an Access Token.
* **POST** `/refresh-token` — Refresh an expired Access Token using a valid Refresh Token.

### 👥 User
Management of user profiles and data.
* **GET** `/users` — Retrieve a list of all registered users.
* **POST** `/users` — Manually create a new user.
* **GET** `/users/{id}` — Find specific user details by their unique ID.
* **DELETE** `/users/{id}` — Permanently remove a user from the system.

### 🔗 Link
Core functionality for managing shortened URLs.
* **GET** `/links` — Fetch all links belonging to the authenticated user.
* **POST** `/links` — Generate a new shortened link.
* **GET** `/links/{id}` — Retrieve detailed information for a specific link.
* **PATCH** `/links/{id}` — Update existing link data (e.g., target URL).
* **DELETE** `/links/{id}` — Delete a shortened link.

### ⚡ Redirect
The primary public entry point for link redirection.
* **GET** `/{code}` — Redirect users to the original URL based on the unique short code.

### 📊 Analytics
Detailed performance metrics and visitor insights.
* **GET** `/links/{id}/analytics` — General summary (Total clicks, Unique clicks, etc.).
* **GET** `/links/{id}/analytics/clicks` — Detailed click analytics based on time series (Hourly/Daily).
* **GET** `/links/{id}/analytics/country` — Click distribution based on the visitor's country of origin.

---

## 🔑 Authentication
Protected endpoints require a valid JSON Web Token (JWT) passed in the header:
> **Authorization**: Bearer `<your_jwt_token>`

## 📝 Integration Notes
* **Content-Type**: All requests containing a body must use `application/json`.
* **Response Format**: All API responses are returned in standard JSON format.
