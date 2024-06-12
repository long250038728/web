package algorithm

import "fmt"

var sortData = []int32{
	300, 200, 100, 50, 20,
}

// bubbleSort 冒泡排序 子循环确定第x次循环后倒数第x位第x大
// 第1次循环确定倒数第1位第1大
// 第2次循环确定倒数第2位第2大
// 第3次循环确定倒数第3位第3大
func bubbleSort() {
	//第一个循环 遍历所有元素
	//第二个循环 遍历子集
	//  第一次（对比第一个与最后一个）保证第一大
	//    第一次第0个与第1个对比
	//    第二次第1个与第2个对比
	//             ....
	//    第N次最大的值会在最后
	//
	//  第二次（对比第一个与最后第二个）保证第二大
	//    第一次第0个与第1个对比
	//    第二次第1个与第2个对比
	//             ....
	//    第N次确保第二大的在倒数第二

	for i := 0; i < len(sortData)-1; i++ {
		swapped := false

		for j := 0; j < len(sortData)-1-i; j++ {
			if sortData[j] > sortData[j+1] {
				sortData[j], sortData[j+1] = sortData[j+1], sortData[j]
				swapped = true
			}
		}

		//从头遍历到尾巴，如果全程都无需要交换，则代表就是有序的
		// 1 跟 2 对比发现不需要替换
		// 2 跟 3 对比发现不需要替换
		//  ...
		// n-1 跟 n 对比发现不需要替换
		// 全程有序，无需要再对比了
		if !swapped {
			break
		}
	}
	fmt.Println(sortData)
}

func insertSort() {
	for i := 1; i < len(sortData)-1; i++ {
		//从第二位开始
		value := sortData[i]
		j := i - 1
		//遍历自己左边的位置到0计数排序
		for ; j >= 0; j-- {
			if sortData[j] > value {
				sortData[j+1] = sortData[j] // 数据移动
			} else {
				break
			}
		}
		sortData[j+1] = value // 插入数据
	}
	fmt.Println(sortData)
}

func selectSort() {
	arr := sortData
	n := len(arr)
	for i := 0; i < n-1; i++ {
		//假设自己是最小的
		minIndex := i

		//遍历自己右边的所有数据
		for j := i + 1; j < n; j++ {

			// 如果比我小，就把minIndex进行赋值
			if arr[j] < arr[minIndex] {
				minIndex = j
			}
		}

		// 如果不相等，则代表自己不是最小的
		if minIndex != i {
			arr[i], arr[minIndex] = arr[minIndex], arr[i]
		}
	}
	fmt.Println(arr)
}
