package util

// 字符串相关工具

import "unsafe"

// MustToString 字节转字符串,字节数组转字符串是非常消耗性能的，可考虑使用指针转换的方式
func MustToString(target []byte) string {
	return *(*string)(unsafe.Pointer(&target))
}

// GetRandomString 获取一个指定长度的随机字符串, 支持多线程
func GetRandomString(l int) string {
	randLock.Lock()
	defer randLock.Unlock()
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQESTUVWSYZ"
	bytes := []byte(str)
	var result []byte
	for i := 0; i < l; i++ {
		result = append(result, bytes[randSeed.Intn(len(bytes))])
	}
	return string(result)
}

// GetRandNumberByMod 生成一个随机数, 支持多线程
func GetRandNumberByMod(mod int64) (res int64) {
	randLock.Lock()
	defer randLock.Unlock()
	res = randSeed.Int63() % mod
	return
}
