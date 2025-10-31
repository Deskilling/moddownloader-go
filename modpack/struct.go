package modpack

type Modpack struct {
	Version string
	Loader  string
	Input   string
	Output  string
}

type Config struct {
	Modpacks map[string]Modpack
}

type Mrpack struct {
	Dependencies  map[string]string `json:"dependencies"`
	Files         []file            `json:"files"`
	FormatVersion int               `json:"formatVersion"`
	Game          string            `json:"game"`
	Name          string            `json:"name"`
	VersionId     string            `json:"versionId"`
}

type file struct {
	DownloadUrl []string          `json:"downloads"`
	Env         map[string]string `json:"env"`
	FileSize    int               `json:"fileSize"`
	Hashes      hashes            `json:"hashes"`
	Path        string            `json:"path"`
}

type hashes struct {
	Sha1   string `json:"sha1"`
	Sha512 string `json:"sha512"`
}
