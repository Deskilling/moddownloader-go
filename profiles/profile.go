package profiles

import (
	"fmt"
	"github.com/deskilling/moddownloader-go/request"
)

func profileLastVersionId(loader string, version string) string {
	if loader == "fabric" {
		return fmt.Sprintf("fabric-loader-%s-%s", request.GetLatestFabricVersion(), version)
	} else if loader == "quilt" {
		return fmt.Sprintf("quilt-loader-%s-%s", request.GetLatestQuiltVersion(), version)
	} else if loader == "forge" {
		return fmt.Sprintf("%s-forge-%s", version, request.GetLatestForgeVersion(version))
	}

	return "unknown"
}

func latestProfile(loader string, version string) string {
	if loader == "fabric" {
		return fmt.Sprintf("fabric-loader-%s", version)
	} else if loader == "quilt" {
		return fmt.Sprintf("Quilt %s", version)
	} else if loader == "forge" {
		// TODO - Check
		return "forge"
	}
	return "unknown"
}
