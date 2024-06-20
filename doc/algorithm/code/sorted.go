package code

import (
	"fmt"
)

var sortData = []int32{
	100, 300, 200, 400, 150,
}

// 冒泡排序，相邻的两个对比，对比大小
// 选择排序，我与我右边的对比，找出最小
// 插入排序，我与我左边的对比，如果比我大就右挪动

//-----------------------------------------------------------------

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

//-----------------------------------------------------------------

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

//-----------------------------------------------------------------

func selectSort() {
	arr := sortData
	n := len(arr)
	for i := 0; i < n-1; i++ {
		//假设自己是最小的
		minIndex := i

		//遍历自己右边的所有数据与我对比（左边已经是最小了）
		for j := i + 1; j < n; j++ {

			// 如果比我小，就把minIndex进行赋值
			if arr[j] < arr[minIndex] {
				minIndex = j
			}
		}

		// i是自己，如果最小的不是我，那我就跟最小的替换
		if minIndex != i {
			arr[i], arr[minIndex] = arr[minIndex], arr[i]
		}
	}
	fmt.Println(arr)
}

//-----------------------------------------------------------------

func MergeSort() {
	// 100   300    200   400   150
	// 第一次 left: 100  300       right: 200   400   150
	// 第二次 left: 100	 		  right: 300
	// 第三次 merge 100 300
	//
	// 第四次: left: 200	 	      right: 400   150
	// 第五次: left: 400	 	      right: 150
	// 第六次: merge 150 400
	//
	// 第七次: left: 200	 	      right: 150 400
	// 第八次: merge 150 200 400
	//
	// 第九次: left: 100 300	 	      right: 150 200 400
	// 第十次: merge 100  150 200 300 400
	//
	//	  [100   300    200   400   150]
	//	  left:[100  300]						right:[200  400  150]
	//    left:[100]   right:[300]
	//	  merge [100  300]
	//   										left:[200]   right: [400 150]
	//	  													 left:[40]   right:[150]
	//														 merge [40 150]
	//											merge [40 150 200]
	//	  merge [40 100 150 200 300]
	fmt.Println(subMergeSort(sortData))
}

func subMergeSort(arr []int32) []int32 {
	if len(arr) <= 1 {
		return arr
	}
	// 将数组分成两半
	mid := len(arr) / 2
	left := subMergeSort(arr[:mid])
	right := subMergeSort(arr[mid:])

	// 合并两个有序子数组
	return merge(left, right)
}

var merge = func(left, right []int32) []int32 {
	result := make([]int32, 0, len(left)+len(right))
	var i, j int32 = 0, 0

	for i < int32(len(left)) && j < int32(len(right)) {
		if left[i] <= right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}

	// 将剩余的元素添加到结果中
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)
	fmt.Println(result)
	return result
}

//-----------------------------------------------------------------

func quickSort() {
	subQuickSort(sortData, 0, int32(len(sortData)-1))
	fmt.Println(sortData)
}

// quickSort 函数用于递归地对数组进行排序
func subQuickSort(arr []int32, low, high int32) {
	if low < high {
		pi := partition(arr, low, high) // 获取分区索引
		// 递归排序分区的两部分
		subQuickSort(arr, low, pi-1)
		subQuickSort(arr, pi+1, high)
	}
}

// partition 函数用于将数组分成两部分
func partition(arr []int32, low, high int32) int32 {
	pivot := arr[high] // 选择最后一个元素作为基准
	i := low - 1       // i 用于追踪小于 pivot 的元素的索引

	for j := low; j < high; j++ {
		if arr[j] < pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i] // 交换元素
		}
	}

	// 将 pivot 放到正确的位置
	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}

//-----------------------------------------------------------------

func bucketSort() {
	//每10为一个桶，每个桶都是一个数组
	buckets := make([][]int32, 10, 10)
	for i := 0; i < len(sortData); i++ {
		val := sortData[i]
		bucketNum := val/10 - 1

		// 动态扩展 buckets 切片
		if bucketNum >= int32(len(buckets)) {
			newBuckets := make([][]int32, bucketNum+1, bucketNum+1)
			copy(newBuckets, buckets)
			buckets = newBuckets
		}

		if buckets[bucketNum] == nil {
			buckets[bucketNum] = make([]int32, 0, 100)
		}
		buckets[bucketNum] = append(buckets[bucketNum], 0)

		//桶内有序（冒泡排序）
		for i := 0; i < len(buckets[bucketNum])-1; i++ {
			isSwp := false //是否交换

			for j := 0; j < len(buckets[bucketNum])-1-i; j++ {
				if buckets[bucketNum][j] > buckets[bucketNum][j+1] {
					isSwp = true //已经交换
					buckets[bucketNum][j], buckets[bucketNum][j+1] = buckets[bucketNum][j+1], buckets[bucketNum][j]
				}
			}

			if !isSwp { //全程无需交互，代表全部有序
				break
			}
		}

	}
	fmt.Println(buckets)
}

func countingSort() {
	//每一个值都是一个桶，每个桶都是一个数组
	counting := make([][]int32, 100, 100)
	for i := 0; i < len(sortData); i++ {
		val := sortData[i]

		// 动态扩展 buckets 切片
		if val >= int32(len(counting)) {
			newCounting := make([][]int32, val+1, val+1)
			copy(newCounting, counting)
			counting = newCounting
		}

		if counting[val] == nil {
			counting[val] = make([]int32, 0, 100)
		}
		counting[val] = append(counting[val], val)
	}

	fmt.Println(counting)
}
