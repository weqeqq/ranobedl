package builder

type Builder interface {
	SetTitle(name string)
	SetAuthor(author string)

	PushChapter(chapterTitle string) error
	PushParagraph(text string) error
	PushImage(imagePath string) error

	Build(filename string) error
}
