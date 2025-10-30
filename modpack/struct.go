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
