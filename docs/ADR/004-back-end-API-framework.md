# ADR-00X-XXX

* Creation Date: 07/01/2022
* Status: Accepted

## Context

We chose to use Go as our back-end language and we will need to build an HTTP API. We need library that support path with variables (`/with/{parameter}/in/path`)

## Decision

Gorilla/Mux was selected since it match all our needs (Native net/http server was not) and we are familiar with it.

## Options

### 1. Gorilla/Mux

```go
    type SomeHandler struct {
        Cfg          *core.Config
        SessionStore *core.SessionStore
    }
    
    func (h SomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    ...
    w.WriteHeader(200)
    }
    
    ...

    r := mux.NewRouter()
    r.PathPrefix("/api/:some/URL").Handler(SomeHandler{Cfg: cfg, SessionStore: sessionStore}).Methods("POST")
    srv := &http.Server{
        Handler: r,
        Addr:    fmt.Sprintf("0.0.0.0:%v", cfg.Port_HTTP),
    }
    srv.ListenAndServe()
```

#### Benefits
* Express-like API
* Several extensions (Sessions, middlewares)
* Easy to use for everyone
* We already have experiences with this library
#### Drawbacks
* Looking for a new maintainer

### 2. Gin

```go

    router := gin.Default()

    router.GET("/user/:name", func(c *gin.Context) {
    	name := c.Param("name")
    	c.String(http.StatusOK, "Hello %s", name)
    })

    router.Run()

```

#### Benefits
* Popular
* Many features and examples
* Claims to be fast
#### Drawbacks
* Custom HTTP router

### 3. Native HTTP server

```go

    type myHandler struct{}
    
    func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
        io.WriteString(w, "URL: "+r.URL.String())
    }

    mux := http.NewServeMux()

    mux.Handle("/",&myHandler{})
	
    _ := http.ListenAndServe(":8080", mux)
```

#### Benefits
* Minimal dependencies
* Simple API
#### Drawbacks
* Doesn't support path variables (ex: `/some/:arg/in/path` where :arg is a variable part)

## Technical resources
- [Architecture Decision Record](https://github.com/joelparkerhenderson/architecture-decision-record/blob/main/examples/programming-languages/index.md)
- [Go net/http vs Gin](https://www.stephengream.com/go-nethttp-vs-gin)
- [Different approaches to HTTP routing in Go](https://benhoyt.com/writings/go-routing/)
