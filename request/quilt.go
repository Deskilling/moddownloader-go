package request

import "encoding/json"

func GetLatestQuiltVersion() string {
	response, err := ModrinthWebRequest("https://meta.quiltmc.org/v3/versions/loader")
	if err != nil {
		panic(err)
	}

	var quiltVersions []Version
	err = json.Unmarshal([]byte(response), &quiltVersions)
	if err != nil {
		panic(err)
	}
	return quiltVersions[0].Version
}
