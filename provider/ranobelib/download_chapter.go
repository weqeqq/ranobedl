package ranobelib

import (
	"fmt"
	"path/filepath"
	api "ranobedl/api/ranobelib"
	"ranobedl/cachemgr"
)

type chapterDownloader struct {
	*cachemgr.PathInfo

	UniqueName string
}

func (cd *chapterDownloader) chapterPath(number string, volume string) (string, error) {
	if ranobeDir, err := cachemgr.ConstructPath(provider, cd.UniqueName); err != nil {
		return "", err

	} else {
		return filepath.Join(
			ranobeDir,
			fmt.Sprintf("%s%s.json", volume, number)), nil
	}
}
func (cd *chapterDownloader) Download(number string, volume string) error {
	chapterContent, err := api.GetChapterContent(cd.UniqueName, number, volume)
	if err != nil {
		return err
	}
	schema, err := convertContent(cd.UniqueName, chapterContent)
	if err != nil {
		return err
	}
	chapterPath, err := cd.chapterPath(number, volume)
	if err != nil {
		return err
	}
	schema.ToFile(chapterPath)
	cd.PathInfo.Data = append(cd.PathInfo.Data, cachemgr.Chapter{
		Path:   chapterPath,
		Number: number,
		Volume: volume,
	})
	return nil
}

func downloadChapter(pathInfo *cachemgr.PathInfo, uniqueName string, number string, volume string) error {
	return (&chapterDownloader{
		PathInfo:   pathInfo,
		UniqueName: uniqueName,
	}).Download(number, volume)
}
