package book

type Book struct {
	Title  string
	Author string
	Pages  int
}

const (
	SHORT_STORY = "SHORT STORY"
	NOVEL       = "NOVEL"
)

func (b Book) CategoryByLength() string {
	switch {
	case b.Pages > 300:
		return NOVEL
	default:
		return SHORT_STORY
	}
}
