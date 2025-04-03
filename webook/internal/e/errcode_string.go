// Code generated by "stringer -type errCode -linecomment"; DO NOT EDIT.

package e

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ServerErr-50001]
	_ = x[UserExist-1001]
	_ = x[UserInvalidInput-1003]
	_ = x[UserAuthFailed-1004]
	_ = x[UserOperationTooFrequent-1005]
	_ = x[Code_NotFind-2001]
	_ = x[Code_VerifyFail-2002]
	_ = x[Code_TooManyVerificationAttempts-2003]
}

const (
	_errCode_name_0 = "用户已存在"
	_errCode_name_1 = "输入错误账号密码错误操作太频繁"
	_errCode_name_2 = "验证码不存在,请重新发送验证失败验证次数过多,请重新发送"
	_errCode_name_3 = "系统错误"
)

var (
	_errCode_index_1 = [...]uint8{0, 12, 30, 45}
	_errCode_index_2 = [...]uint8{0, 34, 46, 80}
)

func (i errCode) String() string {
	switch {
	case i == 1001:
		return _errCode_name_0
	case 1003 <= i && i <= 1005:
		i -= 1003
		return _errCode_name_1[_errCode_index_1[i]:_errCode_index_1[i+1]]
	case 2001 <= i && i <= 2003:
		i -= 2001
		return _errCode_name_2[_errCode_index_2[i]:_errCode_index_2[i+1]]
	case i == 50001:
		return _errCode_name_3
	default:
		return "errCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
