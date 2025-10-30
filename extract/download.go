package extract

import (
	"slices"
)

type Download struct {
	url      string
	filename string
	sha1     string
	sha512   string
}

func GetDownload(mvi []ModVersionInformation, version, loader string) (*Download, error) {
	for _, v := range mvi {
		if slices.Contains(v.GameVersions, version) && slices.Contains(v.SupportedLoaders, loader) {
			if len(v.Files) > 0 {
				var dl = Download{
					url:      v.Files[0].URL,
					filename: v.Files[0].Filename,
					sha1:     v.Files[0].Hashes.Sha1,
					sha512:   v.Files[0].Hashes.Sha512,
				}

				return &dl, nil
			}
		}
	}

	return nil, nil
}
