package extract

type ModInformation struct {
	GameVersions       []string `json:"game_versions"`
	SupportedLoaders   []string `json:"loaders"`
	ProjectId          string   `json:"id"`
	ProjectTitle       string   `json:"title"`
	ProjectDescription string   `json:"description"`
	ProjectUpdated     string   `json:"updated"`
	ProjectDownloads   uint     `json:"downloads"`
	ProjectIconUrl     string   `json:"icon_url"`
}

type File struct {
	Hashes struct {
		Sha1   string `json:"sha1"`
		Sha512 string `json:"sha512"`
	} `json:"hashes"`
	URL      string `json:"url"`
	Filename string `json:"filename"`
	Size     int    `json:"size"`
}

type ModVersionInformation struct {
	GameVersions     []string `json:"game_versions"`
	SupportedLoaders []string `json:"loaders"`
	VersionId        string   `json:"id"`
	ProjectId        string   `json:"project_id"`
	VersionName      string   `json:"name"`
	VersionPublished string   `json:"date_published"`
	ProjectDownloads uint     `json:"downloads"`
	Files            []File   `json:"files"`
}
