package ranobelib

import (
	api "ranobedl/api/ranobelib"
	"ranobedl/cachemgr"
)

type ranobeDownloader struct {
	UniqueName string
}

func (rd *ranobeDownloader) exportInfo() error {
	if ranobeInfo, err := api.GetRanobeInfo(rd.UniqueName); err != nil {
		return err

	} else {
		converted := cachemgr.RanobeInfo{
			Name:   ranobeInfo.Name,
			Author: ranobeInfo.Authors[0].Name,
		}
		return converted.Save(provider, rd.UniqueName)
	}
}
func (rd *ranobeDownloader) Download(callback func(current, total int)) error {
	if err := cachemgr.CreateRanobeDir(provider, rd.UniqueName); err != nil {
		return err
	}
	chapterInfo, err := api.GetChapterInfo(rd.UniqueName)
	if err != nil {
		return err
	}
	pathInfo := cachemgr.PathInfo{Data: []cachemgr.Chapter{}}

	for index, chapter := range chapterInfo {
		if err := downloadChapter(&pathInfo, rd.UniqueName, chapter.Number, chapter.Volume); err != nil {
			return err
		}
		callback(index, len(chapterInfo))
	}
	if err := rd.exportInfo(); err != nil {
		return err
	}
	return pathInfo.Save(provider, rd.UniqueName)
}

func DownloadRanobe(uniqueName string, callback func(current, total int)) error {
	return (&ranobeDownloader{UniqueName: uniqueName}).Download(callback)
}
