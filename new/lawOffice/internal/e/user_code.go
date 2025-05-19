package e

const (
	UserExist                errCode = 1001 //用户已存在
	UserInvalidInput         errCode = 1003 //输入错误
	UserAuthFailed           errCode = 1004 //账号密码错误
	UserOperationTooFrequent errCode = 1005 //操作太频繁
)

const (
	Code_NotFind                     errCode = 2001 //验证码不存在,请重新发送
	Code_VerifyFail                  errCode = 2002 //验证失败
	Code_TooManyVerificationAttempts errCode = 2003 //验证次数过多,请重新发送
)
