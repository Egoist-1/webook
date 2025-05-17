package dao

type Article struct {
	Id       int64
	Title    string
	Content  string
	Status   uint
	AuthorId int64
	Ctime    int64
	Utime    int64
}

type ArticlePublish Article
