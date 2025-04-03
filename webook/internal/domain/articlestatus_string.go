// Code generated by "stringer -type ArticleStatus -linecomment"; DO NOT EDIT.

package domain

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ArticleStatusUnknown-0]
	_ = x[ArticleStatusUnpublished-1]
	_ = x[ArticleStatusPublished-2]
	_ = x[ArticleStatusPrivate-3]
}

const _ArticleStatus_name = "未知错误文章未发布文章已发布文章不可见"

var _ArticleStatus_index = [...]uint8{0, 12, 27, 42, 57}

func (i ArticleStatus) String() string {
	if i >= ArticleStatus(len(_ArticleStatus_index)-1) {
		return "ArticleStatus(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ArticleStatus_name[_ArticleStatus_index[i]:_ArticleStatus_index[i+1]]
}
