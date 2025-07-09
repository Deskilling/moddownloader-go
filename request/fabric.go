package request

import "encoding/json"

func GetLatestFabricVersion() string {
	response, err := ModrinthWebRequest("https://meta.fabricmc.net/v2/versions/loader")
	if err != nil {
		panic(err)
	}

	var fabricVersion []Version
	err = json.Unmarshal([]byte(response), &fabricVersion)
	if err != nil {
		panic(err)
	}
	return fabricVersion[0].Version
}
