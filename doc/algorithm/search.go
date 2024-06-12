package algorithm

import (
	"errors"
)

func binaryArraySearch(data []int32, value int32) (int32, error) {
	dataLen := int32(len(data))
	if dataLen == 0 {
		return 0, errors.New("数据为空")
	}
	var left, right int32 = 0, dataLen - 1

	for left <= right {
		mid := left + (right-left)/2

		midValue := data[mid]
		if midValue == value {
			return mid, nil
		} else if midValue > value {
			right = mid - 1
		} else if midValue < value {
			left = mid + 1
		}
	}

	return 0, errors.New("数据不存在")
}
