package extract

import (
	"encoding/json"
	"fmt"
)

func Mod(modData string) (ModInformation, error) {
	var info ModInformation
	// Umarshal converts the json data into a Go Struct.
	// byte converts the modData json into a byte slice, this is required for Unmarshal
	// Then there is the pointer to the struct which then searches the byte slice and sets the values in the struct
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
