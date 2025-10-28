package extract

import (
	"encoding/json"
	"fmt"
)

func Mod(modData string) (ModInformation, error) {
	var info ModInformation
	if err := json.Unmarshal([]byte(modData), &info); err != nil {
		return ModInformation{}, fmt.Errorf("failed to unmarshal mod information: %w", err)
	}

	return info, nil
}

func Version(modVersionData string) ([]ModVersionInformation, error) {
	var versionsInfo []ModVersionInformation
	if err := json.Unmarshal([]byte(modVersionData), &versionsInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal mod version information: %w", err)
	}

	return versionsInfo, nil
}

func VersionHash(modVersionData string) (ModVersionInformation, error) {
	var fileInfo ModVersionInformation
	if err := json.Unmarshal([]byte(modVersionData), &fileInfo); err != nil {
		return ModVersionInformation{}, fmt.Errorf("failed to unmarshal mod version information: %w", err)
	}

	return fileInfo, nil
}
