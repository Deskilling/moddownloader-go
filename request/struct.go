package request

type Version struct {
	Version     string `json:"version"`
	VersionType string `json:"version_type"`
	//Major       bool   `json:"major"`

	// Fabric Specific
	Build  int  `json:"build"`
	Stable bool `json:"stable"`
}

type EndpointMap map[string]string

var ModrinthEndpoint = EndpointMap{
	"default":               "https://api.modrinth.com",
	"modInformation":        "https://api.modrinth.com/v2/project/%s",
	"modVersionInformation": "https://api.modrinth.com/v2/project/%s/version",
	"versionFileHash":       "https://api.modrinth.com/v2/version_file/%s",
	"versionUpdate":         "https://api.modrinth.com/v2/version_file/{hash}/update",
	"availableVersions":     "https://api.modrinth.com/v2/tag/game_version",
	"availableLoaders":      "https://api.modrinth.com/v2/tag/loader",

	// "search": "https://api.modrinth.com/v2/search",
}
