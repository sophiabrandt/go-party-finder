This is an alternative branch where I swapped the custom HandlerFunc in the web package for a custom Handler (compatible with http.Handler).

This doesn't solve the problem of using third party middleware though.

```diff
diff --git a/business/app/handlers/handlers.go b/business/app/handlers/handlers.go
index 4b6810f..0ef737d 100644
--- a/business/app/handlers/handlers.go
+++ b/business/app/handlers/handlers.go
@@ -12,7 +12,7 @@
 	"github.com/sophiabrandt/go-party-finder/foundation/web"
 )
 
-// Router  creates a new http.Handler with all routes.
+// Router creates a new http.Handler with all routes.
 func Router(build string, shutdown chan os.Signal, apc *apc.AppContext, staticFilesDir string, log *log.Logger, db *sqlx.DB) http.Handler {
 	// Creates a new web application with all routes and middleware.
 	app := web.NewApp(shutdown, apc, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics(log))
@@ -22,18 +22,18 @@ func Router(build string, shutdown chan os.Signal, apc *apc.AppContext, staticFi
 		build: build,
 		db:    db,
 	}
-	app.HandleDebug(http.MethodGet, "/readiness", cg.readiness)
-	app.HandleDebug(http.MethodGet, "/liveness", cg.liveness)
+	app.HandleDebug(http.MethodGet, "/readiness", web.Handler{cg.readiness})
+	app.HandleDebug(http.MethodGet, "/liveness", web.Handler{cg.liveness})
 
 	// index route and parties routes
 	pg := partyGroup{
 		party: party.New(log, db),
 	}
-	app.Handle(http.MethodGet, "/", pg.query)
-	app.Handle(http.MethodGet, "/parties/{page}/{rows}", pg.query)
-	app.Handle(http.MethodGet, "/parties/{id}", pg.queryByID)
-	app.Handle(http.MethodGet, "/parties/create", pg.createForm)
-	app.Handle(http.MethodPost, "/parties/create", pg.create)
+	app.Handle(http.MethodGet, "/", web.Handler{pg.query})
+	app.Handle(http.MethodGet, "/parties/{page}/{rows}", web.Handler{pg.query})
+	app.Handle(http.MethodGet, "/parties/{id}", web.Handler{pg.queryByID})
+	app.Handle(http.MethodGet, "/parties/create", web.Handler{pg.createForm})
+	app.Handle(http.MethodPost, "/parties/create", web.Handler{pg.create})
 
 	// static file server
 	filesDir := http.Dir(staticFilesDir)
diff --git a/business/mid/errors.go b/business/mid/errors.go
index b9e715b..c46505f 100644
--- a/business/mid/errors.go
+++ b/business/mid/errors.go
@@ -24,7 +24,7 @@ func Errors(log *log.Logger) web.Middleware {
 			}
 
 			// Run the next handler and catch any propagated error.
-			if err := handler(ctx, w, r); err != nil {
+			if err := handler.H(ctx, w, r); err != nil {
 
 				// Log the error.
 				log.Printf("%s: ERROR: %v", v.TraceID, err)
@@ -45,7 +45,7 @@ func Errors(log *log.Logger) web.Middleware {
 			return nil
 		}
 
-		return h
+		return web.Handler{h}
 	}
 
 	return m
diff --git a/business/mid/logger.go b/business/mid/logger.go
index 9fe860f..bffaa8f 100644
--- a/business/mid/logger.go
+++ b/business/mid/logger.go
@@ -28,7 +28,7 @@ func Logger(log *log.Logger) web.Middleware {
 			)
 
 			// Call the next handler.
-			err := handler(ctx, w, r)
+			err := handler.H(ctx, w, r)
 
 			log.Printf("%s: completed: %s %s -> %s (%d) (%s)",
 				v.TraceID,
@@ -39,7 +39,7 @@ func Logger(log *log.Logger) web.Middleware {
 			// Return the error so it can be handled further up the chain.
 			return err
 		}
-		return h
+		return web.Handler{h}
 	}
 	return m
 }
diff --git a/business/mid/metrics.go b/business/mid/metrics.go
index 353fe4f..936d9b8 100644
--- a/business/mid/metrics.go
+++ b/business/mid/metrics.go
@@ -30,11 +30,11 @@ func Metrics() web.Middleware {
 			// Don't count anything on /debug routes towards metrics.
 			// Call the next handler to continue processing.
 			if strings.HasPrefix(r.URL.Path, "/debug") {
-				return handler(ctx, w, r)
+				return handler.H(ctx, w, r)
 			}
 
 			// Call the next handler.
-			err := handler(ctx, w, r)
+			err := handler.H(ctx, w, r)
 
 			// Increment the request counter.
 			m.req.Add(1)
@@ -53,7 +53,7 @@ func Metrics() web.Middleware {
 			return err
 		}
 
-		return h
+		return web.Handler{h}
 	}
 
 	return m
diff --git a/business/mid/panics.go b/business/mid/panics.go
index 5fdf8fe..e9bd06e 100644
--- a/business/mid/panics.go
+++ b/business/mid/panics.go
@@ -36,10 +36,10 @@ func Panics(log *log.Logger) web.Middleware {
 			}()
 
 			// Call the next handler and set its return value in he err variable.
-			return handler(ctx, w, r)
+			return handler.H(ctx, w, r)
 		}
 
-		return h
+		return web.Handler{h}
 	}
 
 	return m
diff --git a/foundation/web/web.go b/foundation/web/web.go
index 7e16027..f52de2a 100644
--- a/foundation/web/web.go
+++ b/foundation/web/web.go
@@ -35,14 +35,22 @@ type Values struct {
 // could try to be registered more than once, causing a panic.
 var registered = make(map[string]bool)
 
-// A Handler is a type that handles an http request within our own little mini
+// A HandlerFunc is a function that handles an http request within our own little mini
 // framework.
-type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error
+type HandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request) error
+
+// Handler is the custom web handler for the application.
+type Handler struct {
+	H HandlerFunc
+}
 
 // ServeHTTP is a wrapper to make the Handler compliant with the http.Handler interface.
 func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
 	ctx := r.Context()
-	h(ctx, w, r)
+	err := h.H(ctx, w, r)
+	if err != nil {
+		panic(err)
+	}
 }
 
 // App is the entrypoint for the web application.
@@ -119,7 +127,7 @@ func (a *App) handle(debug bool, method string, path string, handler Handler, mw
 		ctx = context.WithValue(ctx, KeyValues, &v)
 
 		// Call the wrapped handler functions.
-		if err := handler(ctx, w, r); err != nil {
+		if err := handler.H(ctx, w, r); err != nil {
 			a.SignalShutdown()
 			return
 		}
@@ -139,7 +147,7 @@ func (a *App) handle(debug bool, method string, path string, handler Handler, mw
 		return
 	}
 
-	a.mux.MethodFunc(method, path, h)
+	a.mux.Method(method, path, http.HandlerFunc(h))
 }
 
 // neuteredFileSystem disallows directory listings for a static file server.
```
