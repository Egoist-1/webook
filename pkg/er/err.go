package er

type Err struct {
	//错误码
	code ErrCode
	//日志定位
	msg string
}

func (e Err) Error() string {
	return string(e.code) + e.code.String() + e.msg
}

func (e Err) Code() ErrCode {
	return e.code
}

// msg 定位作用
func NewErr(code ErrCode, msg string, errToString string) error {
	return Err{
		code: code,
		msg:  msg + errToString,
	}
}
func NewServerErr(msg string, errToString string) error {
	return Err{
		code: ServerErr,
		msg:  msg + errToString,
	}
}

type ErrCode int

func (e ErrCode) ToInt() int {
	return int(e)
}

//go:generate stringer -type ErrCode -linecomment

const (
	ServerErr ErrCode = 50001 //系统错误
)
