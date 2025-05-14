package dao

type Article struct {
	Id       int
	Title    string
	Content  string
	Status   uint
	AuthorId int
	Ctime    int64
	Utime    int64
}

type ArticlePublish Article

type Interactive struct {
	Id int
}
