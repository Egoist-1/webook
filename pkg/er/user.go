package er

const (
	UserExist                ErrCode = 1001 //用户已存在
	UserInvalidInput         ErrCode = 1003 //输入错误
	UserAuthFailed           ErrCode = 1004 //账号密码错误
	UserOperationTooFrequent ErrCode = 1005 //操作太频繁
)
