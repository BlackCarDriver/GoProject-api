package util

// 数学相关小工具

import (
	"crypto/md5"
	"fmt"
)

// BalancedSplit 均衡地分割一个数组到多个。 array=原数组, maxSize=分割后每个数组最大容量
// 例如 (array)=100, maxSize=40 时, 得到的数组长度分别为{34,33,33}
func BalancedSplit(array []int64, maxSize int) (result [][]int64, err error) {
	if len(array) <= 0 || maxSize <= 1 {
		return make([][]int64, 0), fmt.Errorf("unexpect params: len(array)=%v maxSize=%d", len(array), maxSize)
	}
	defer func() {
		msg, ok := recover().(interface{})
		if ok {
			var errptr = &err
			*errptr = fmt.Errorf("catch panic: err=%v", msg)
			result = make([][]int64, 0)
		}
	}()
	nums := int(len(array) / maxSize) // 分割切片数量
	if maxSize > len(array) || len(array)%maxSize > 0 {
		nums++
	}
	resSize := len(array) % nums       // 余数
	baseSize := int(len(array) / nums) // 最短切片长度
	result = make([][]int64, nums)

	fmt.Printf("sliceNums:%d  baseLen:%d  resSize:%d \n", nums, baseSize, resSize)

	index := 0
	for i := 0; i < nums; i++ {
		endIndex := index + baseSize
		if i < resSize {
			endIndex++
		}
		result[i] = array[index:endIndex]
		index = endIndex
	}
	return result, nil
}

// GetMD5KeyAndMod 计算某项数据的md5值并按照长度截取,n太小或太大会返回完整的结果
func GetMD5KeyAndMod(anyMsg interface{}, maxLen int) string {
	md5Encoder := md5.New()
	md5Encoder.Write([]byte(fmt.Sprint(anyMsg)))
	md5Value := fmt.Sprintf("%x", md5Encoder.Sum(nil))
	if maxLen <= 0 || maxLen >= len(md5Value) {
		return md5Value
	}
	return md5Value[0:maxLen]
}
