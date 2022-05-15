package update

type ManifestData struct {
	DL       string                       `json:"dl"`
	Only     []string                     `json:"only"`
	Hash     string                       `json:"hash"`
	CompHash string                       `json:"compHash"`
	Patches  map[string]map[string]string `json:"patches"`
}

type UpdateFile struct {
	Name string
	Path string
	URL  string
}
