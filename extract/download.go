package extract

import (
	"slices"
)

type Download struct {
	Url      string
	Filename string
	Sha1     string
	Sha512   string
}

func GetDownload(mvi []ModVersionInformation, version, loader string) (*Download, error) {
	for _, v := range mvi {
		if slices.Contains(v.GameVersions, version) && slices.Contains(v.SupportedLoaders, loader) {
			if len(v.Files) > 0 {
				var dl = Download{
					Url:      v.Files[0].URL,
					Filename: v.Files[0].Filename,
					Sha1:     v.Files[0].Hashes.Sha1,
					Sha512:   v.Files[0].Hashes.Sha512,
				}

				return &dl, nil
			}
		}
	}

	return nil, nil
}
