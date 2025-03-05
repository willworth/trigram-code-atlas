package indexer

type FileEntry struct {
	Path     string   `json:"path"`
	Trigrams []string `json:"trigrams"`
	Mtime    string   `json:"mtime"`
}

type Atlas struct {
	Files    []FileEntry `json:"files"`
	Metadata struct {
		Version   string `json:"version"`
		Created   string `json:"created"`
		FileCount int    `json:"file_count"`
	} `json:"metadata"`
}
