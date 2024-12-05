package utils

import (
	"github.com/go-xuan/quanx/types/stringx"
	"github.com/go-xuan/sqlx/consts"
)

// FirstIndexOfKeyword 获取sql中关键字首次出现的下标
func FirstIndexOfKeyword(sql, key string) int {
	kl, loop, index := len(key), true, 0
	for loop {
		if newIndex := stringx.Index(sql, key, 1); newIndex >= 0 {
			sl := len(sql)
			if newIndex == 0 && sql[kl:kl+1] == consts.Blank {
				index, loop = index+newIndex, false
			} else if newIndex == sl-kl && sql[newIndex-1:newIndex] == consts.Blank {
				index, loop = index+newIndex, false
			} else if sql[newIndex-1:newIndex] == consts.Blank && sql[newIndex+kl:newIndex+kl+1] == consts.Blank {
				index, loop = index+newIndex, false
			} else {
				// 当前index无效则缩减原sql继续loop
				index = newIndex + kl
				sql = sql[index:]
			}
		} else {
			index, loop = -1, false // 没找到直接跳出
		}
	}
	return index
}
