package domain

type Interactive struct {
	Biz        string
	BizId      int64
	ReadCnt    int64
	LikeCnt    int64
	CollectCnt int64
	Liked      bool
	Collected  bool
	Ctime      int64
	Utime      int64
}

type Collection struct {
	Biz   string
	BizId int64
	//文件夹id
	CollectionId int64
	//文件夹名称
	CollectionName string
	Ctime          int64
	Utime          int64
}
