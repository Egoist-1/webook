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
type ReadHistory struct {
	Id    int64
	Uid   int64
	Biz   string
	BizId int64
	Ctime int64
	Utime int64
}

type Collection struct {
	Id             int64
	Biz            string
	Uid            int64
	CollectionName string
	Ctime          int64
	Utime          int64
}

type UserCollectionBiz struct {
	Id    int64  `gorm:"primaryKey,autoIncrement"`
	Cid   int64  `gorm:"index"`
	BizId int64  `gorm:"uniqueIndex:biz_type_id_uid"`
	Biz   string `gorm:"type:varchar(128);uniqueIndex:biz_type_id_uid"`
	Ctime int64
	Utime int64
}
