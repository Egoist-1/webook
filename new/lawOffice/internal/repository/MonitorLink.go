package repository

type Record struct {
	Id       int64
	Device   string
	AID      string //广告计划id
	AID_NAME string //广告计划名称
	SL       string //语言
	IDFA     string //IOS 6+的设备id字段，32位
	IDFA_MD5 string //IOS 6+的设备id的md5摘要，32位
	OS       string //操作系统平台
	MAC      string //移动设备mac地址,转换成大写字母,去掉“:”，并且取md5摘要后的结果
	MAC1     string // 移动设备 mac 地址,转换成大写字母,并且取md5摘要后的结果，32位
	MODEL    string //手机型号
	IP       string
	UA       string
	Ctime    int64 //
	TS       int64 //发送点击事件的时间戳
	//
	CALLBACK_PARAM string
	CALLBACK_URL   string
	Utime          int64
}
