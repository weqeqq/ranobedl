package format

import (
	"fmt"
	"ranobedl/cachemgr"
	"ranobedl/format/internal/builder"
	"ranobedl/format/internal/fb2"
	"ranobedl/format/internal/nodehandler"
	"ranobedl/schema"
)

type Format int

const (
	FB2 Format = iota
	Epub
)

func newBuilder(format Format) builder.Builder {
	switch format {
	case FB2:
		return fb2.NewBuilder()
	case Epub:
		panic("Unimplemented")
	default:
		panic("Unreachable")
	}
}
func getRenderInlineFn(format Format) nodehandler.RenderInline {
	switch format {
	case FB2:
		return fb2.RenderInline
	case Epub:
		panic("Unimplemented")
	default:
		panic("Unreachable")
	}
}

type exporter struct {
	Builder        builder.Builder
	RenderInlineFn nodehandler.RenderInline
	RanobeProvider cachemgr.RanobeProvider
	UniqueName     string
}

func newExporter(ranobeProvider cachemgr.RanobeProvider, uniqueName string) *exporter {
	return &exporter{
		RanobeProvider: ranobeProvider,
		UniqueName:     uniqueName,
		Builder:        newBuilder(FB2),
		RenderInlineFn: getRenderInlineFn(FB2),
	}
}

func (e *exporter) prepare() error {
	if ranobeInfo, err := cachemgr.LoadRanobeInfo(
		e.RanobeProvider,
		e.UniqueName); err != nil {
		return err
	} else {
		e.Builder.SetTitle(ranobeInfo.Name)
		return nil
	}
}
func (e *exporter) pushChapter(chapterPath string, number string, volume string) error {
	e.Builder.PushChapter(fmt.Sprintf("n%sv%s", number, volume))

	if node, err := schema.FromFile(chapterPath); err != nil {
		return err
	} else {
		return nodehandler.PushBlock(e.Builder, e.RenderInlineFn, node)
	}
}
func (e *exporter) Export(outputPath string) error {
	if err := e.prepare(); err != nil {
		return err
	}
	pathInfo, err := cachemgr.LoadPathInfo(
		e.RanobeProvider,
		e.UniqueName,
	)
	if err != nil {
		return err
	}
	for _, chapter := range pathInfo.Data {
		if err := e.pushChapter(chapter.Path, chapter.Number, chapter.Volume); err != nil {
			return err
		}
	}
	return e.Builder.Build(outputPath)
}

func Export(ranobeProvider cachemgr.RanobeProvider, uniqueName string, outputPath string) error {
	return newExporter(ranobeProvider, uniqueName).Export(outputPath)
}
