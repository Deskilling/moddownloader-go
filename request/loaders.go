package request

import (
	"encoding/json"
	"fmt"


	"github.com/charmbracelet/log"
)

func GetLatestFabricVersion() (string, error) {
	response, err := GetBody("https://meta.fabricmc.net/v2/versions/loader")
	if err != nil {
		log.Error("failed request", "err", err)
		return "", err
	}

	var fabricVersion []Version
	err = json.Unmarshal([]byte(response), &fabricVersion)
	if err != nil {
		log.Error("failed json umarshal", "err", err)
		return "", err
	}
	return fabricVersion[0].Version, nil
}

func GetLatestQuiltVersion() (string, error) {
	response, err := GetBody("https://meta.quiltmc.org/v3/versions/loader")
	if err != nil {
		log.Error("failed request", "err", err)
		return "", err
	}

	var quiltVersions []Version
	err = json.Unmarshal([]byte(response), &quiltVersions)
	if err != nil {
		log.Error("failed json umarshal", "err", err)
		return "", err
	}
	return quiltVersions[0].Version, nil
}

func GetLatestForgeVersion(version string) (string, error) {
	response, err := GetBody("https://files.minecraftforge.net/net/minecraftforge/forge/promotions_slim.json")
	if err != nil {
		log.Error("failed request", "err", err)
		return "", err
	}

	var forgePromotions map[string]map[string]string
	err = json.Unmarshal([]byte(response), &forgePromotions)
	if err != nil {
		log.Error("failed json umarshal", "err", err)
		return "", err
	}

	if promos, ok := forgePromotions["promos"]; ok {
		if latest, ok := promos[version+"-latest"]; ok {
			return latest, nil
		}
	}

	return "", fmt.Errorf("could not find latest forge version for %s", version)
}
