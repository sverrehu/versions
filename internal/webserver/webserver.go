package webserver

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/sverrehu/gotest/versions/internal/config"
	"github.com/sverrehu/gotest/versions/internal/repos"
	"github.com/sverrehu/goutils/lrumap"
)

type handler struct {
	target  string
	handler http.Handler
}

var handlers []handler

// shared cache used by all handler instances
var cache *lrumap.LRUMap

type commonReleasesHandler struct {
	h repos.ReleasesFetcher
}

func (h *commonReleasesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("request: %s %s", r.Method, r.URL)
	w.Header().Set("Content-Type", "application/json")
	var jsonReleases []byte = nil
	cached := cache.Get(r.URL.Path)
	if cached == nil {
		pkg := r.PathValue("package")
		releases, err := h.h.GetReleases(pkg)
		if err != nil {
			var re *repos.ReleasesFetcherError
			ok := errors.As(err, &re)
			if ok && re.IsParameterError {
				sendBadRequest(w, err.Error(), r.URL)
			} else {
				sendInternalServerError(w, err, r.URL)
			}
			return
		}
		jsonReleases, err = json.Marshal(releases)
		if err != nil {
			sendInternalServerError(w, err, r.URL)
			return
		}
		cache.Put(r.URL.Path, jsonReleases)
	} else {
		log.Printf("cache hit for %s", r.URL.Path)
		jsonReleases = cache.Get(r.URL.Path).([]byte)
	}
	_, err := w.Write(jsonReleases)
	if err != nil {
		log.Printf("error writing response for url: %v: %v", r.URL, err.Error())
	}
}

func sendInternalServerError(w http.ResponseWriter, err error, url *url.URL) {
	log.Printf("internal server error for url: %v: %v", url, err.Error())
	w.WriteHeader(http.StatusInternalServerError)
}

func sendBadRequest(w http.ResponseWriter, message string, url *url.URL) {
	log.Printf("bad request: %s", message)
	w.WriteHeader(http.StatusBadRequest)
	err := json.NewEncoder(w).Encode(map[string]string{
		"message": message,
	})
	if err != nil {
		sendInternalServerError(w, err, url)
	}
}

func Run() error {
	mux := http.NewServeMux()
	for _, h := range handlers {
		log.Printf("Adding handler for %s\n", h.target)
		mux.Handle(h.target+"/{package...}", h.handler)
	}
	port := config.Cfg().WebServer.Port
	log.Printf("Starting server at port %d\n", port)
	err := http.ListenAndServe(":"+strconv.Itoa(port), mux)
	return err
}

func init() {
	cache = lrumap.New(3000, 2*time.Hour)
	handlers = []handler{
		{target: "/maven", handler: &commonReleasesHandler{&repos.MavenReleasesFetcher{}}},
		{target: "/dockerhub", handler: &commonReleasesHandler{&repos.OCIReleasesFetcher{}}},
		{target: "/github-releases", handler: &commonReleasesHandler{&repos.GitHubReleasesFetcher{}}},
		{target: "/gitlab-releases", handler: &commonReleasesHandler{&repos.GitLabReleasesFetcher{}}},
	}
}
