package extract

import (
	"encoding/json"

	"github.com/charmbracelet/log"
)

func Mod(modData string) (*ModInformation, error) {
	var mInfo ModInformation
	if err := json.Unmarshal([]byte(modData), &mInfo); err != nil {
		log.Error("failed to unmarshal", "err", err)
		return nil, err
	}

	return &mInfo, nil
}

func AllVersions(modVersionData string) (*[]ModVersionInformation, error) {
	var vInfo []ModVersionInformation
	if err := json.Unmarshal([]byte(modVersionData), &vInfo); err != nil {
		log.Error("failed to unmarshal", "err", err)
		return nil, err
	}

	return &vInfo, nil
}

func Version(modVersionData string) (*ModVersionInformation, error) {
	var fInfo ModVersionInformation
	if err := json.Unmarshal([]byte(modVersionData), &fInfo); err != nil {
		log.Error("failed to unmarshal", "err", err)
		return nil, err
	}

	return &fInfo, nil
}
