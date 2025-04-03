package web

type UnpublishList struct {
}

type PubDetailVo struct {
	Id         int
	Title      string
	Content    string
	ReadCnt    int64
	LikeCnt    int64
	CollectCnt int64
	Liked      bool
	Collected  bool
	Ctime      int64
}
