package slicex

// 两个切片的结构体互相转换
func SliceMap[Src any, Dst any](src []Src, fn func(idx int, src Src) Dst) []Dst {
	dst := make([]Dst, len(src))
	for i, val := range src {
		dst = append(dst, fn(i, val))
	}
	return dst
}
