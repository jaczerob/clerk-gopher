package update

import (
	"net/http"
	"net/url"
	"os"

	"github.com/jaczerob/clerk-gopher/internal/util"
)

type ManifestData struct {
	DL       string                       `json:"dl"`
	Only     []string                     `json:"only"`
	Hash     string                       `json:"hash"`
	CompHash string                       `json:"compHash"`
	Patches  map[string]map[string]string `json:"patches"`
}

type UpdateFile struct {
	Name   string
	Path   string
	URL    string
	Hash   string
	Exists bool

	out *os.File
}

type UpdateClient struct {
	executable *util.Executable
	http       *http.Client
	headers    map[string]string
	baseURL    *url.URL
}
