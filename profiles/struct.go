package profiles

type Settings struct {
	EnableAnalytics  bool   `json:"enableAnalytics"`
	EnableAdvanced   bool   `json:"enableAdvanced"`
	KeepLauncherOpen bool   `json:"keepLauncherOpen"`
	SoundOn          bool   `json:"soundOn"`
	ShowMenu         bool   `json:"showMenu"`
	EnableSnapshots  bool   `json:"enableSnapshots"`
	EnableHistorical bool   `json:"enableHistorical"`
	EnableReleases   bool   `json:"enableReleases"`
	ProfileSorting   string `json:"profileSorting"`
	ShowGameLog      bool   `json:"showGameLog"`
	CrashAssistance  bool   `json:"crashAssistance"`
}

type Profile struct {
	LastUsed      string `json:"lastUsed"`
	LastVersionId string `json:"lastVersionId"`
	Created       string `json:"created"`
	Icon          string `json:"icon"`
	Name          string `json:"name"`
	Type          string `json:"type"`
}

type Config struct {
	//Settings Settings           `json:"settings"`
	Profiles map[string]Profile `json:"profiles"`
	Version  int                `json:"version"`
}
