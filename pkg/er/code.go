package er

const (
	Code_NotFind                     ErrCode = 2001 //验证码不存在,请重新发送
	Code_VerifyFail                  ErrCode = 2002 //验证失败
	Code_TooManyVerificationAttempts ErrCode = 2003 //验证次数过多,请重新发送
)
