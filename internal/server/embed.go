package server

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/alexandrmotologa/consensai/web"
)

// EmbeddedWebHandler serves the embedded React 19 web cockpit with SPA client routing fallback.
func EmbeddedWebHandler() http.Handler {
	distSub, err := fs.Sub(web.DistFS, "dist")
	if err != nil {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Embedded frontend assets not accessible", http.StatusInternalServerError)
		})
	}

	fileServer := http.FileServer(http.FS(distSub))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		// Check if file exists in embedded filesystem
		f, err := distSub.Open(path)
		if err == nil {
			_ = f.Close()
			fileServer.ServeHTTP(w, r)
			return
		}

		// Fallback to index.html for React SPA client-side routes
		indexData, err := fs.ReadFile(distSub, "index.html")
		if err != nil {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(indexData)
	})
}
