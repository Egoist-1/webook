package domain

type Article struct {
	Id      int64
	Title   string
	Content string
	Status  ArticleStatus
	Author  Author
	Ctime   int64
	Utime   int64
}
type Author struct {
	Id   int64
	Name string
}

type ArticleStatus uint

//go:generate stringer -type ArticleStatus -linecomment
const (
	// ArticleStatusUnknown 为了避免零值之类的问题
	ArticleStatusUnknown     ArticleStatus = iota //未知错误
	ArticleStatusUnpublished                      //文章未发布
	ArticleStatusPublished                        //文章已发布
	ArticleStatusPrivate                          //文章不可见
)
