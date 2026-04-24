package webserver

import (
	_ "embed"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/sverrehu/versions/internal/config"
	"github.com/sverrehu/versions/internal/repos"
	"github.com/sverrehu/versions/internal/state"
)

//go:embed index.html
var indexPage []byte

type handler struct {
	target  string
	handler http.Handler
}

var handlers []handler

type commonReleasesHandler struct {
	h repos.Fetcher
}

type indexHandler struct{}

func (h *commonReleasesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("request: %s %s", r.Method, r.URL)
	w.Header().Set("Content-Type", "application/json")
	var jsonReleases []byte = nil
	cached := state.GetCachedResponse(r.URL.Path)
	if cached == nil {
		pkg := r.PathValue("package")
		releases, err := h.h.GetReleases(pkg)
		if err != nil {
			var re *repos.FetcherError
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
		state.PutCachedResponse(r.URL.Path, jsonReleases)
	} else {
		log.Printf("cache hit for %s", r.URL.Path)
		jsonReleases = cached
	}
	_, err := w.Write(jsonReleases)
	if err != nil {
		log.Printf("error writing response for url: %v: %v", r.URL, err.Error())
	}
}

func (h *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	_, err := w.Write(indexPage)
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

func Run(port int) error {
	setupHandlers()
	mux := http.NewServeMux()
	for _, h := range handlers {
		log.Printf("Adding handler for %s\n", h.target)
		mux.Handle(h.target+"/{package...}", h.handler)
	}
	mux.Handle("/", &indexHandler{})
	log.Printf("Starting server at port %d", port)
	err := http.ListenAndServe(":"+strconv.Itoa(port), mux)
	return err
}

func setupHandlers() {
	datasourcesCfg := *config.Cfg().Datasources
	handlers = []handler{
		{target: "/github-releases", handler: &commonReleasesHandler{h: repos.NewGitHubReleasesFetcher(datasourcesCfg.GitHubReleasesDatasource)}},
		{target: "/github-tags", handler: &commonReleasesHandler{h: repos.NewGitHubTagsFetcher(datasourcesCfg.GitHubTagsDatasource)}},
		{target: "/gitlab-releases", handler: &commonReleasesHandler{h: repos.NewGitLabReleasesFetcher(datasourcesCfg.GitLabReleasesDatasource)}},
		{target: "/maven", handler: &commonReleasesHandler{h: repos.NewMavenReleasesFetcher(datasourcesCfg.MavenDatasource)}},
		{target: "/dockerhub", handler: &commonReleasesHandler{h: repos.NewOCIReleasesFetcher(datasourcesCfg.DockerhubDatasource)}},
	}
}
