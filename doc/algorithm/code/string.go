package code

func BF(main, sub string) int {
	for i := 0; i < len(main)-1; i++ {
		num := 0
		for j := 0; j < len(sub); j++ { // 这里是遍历字符串，低级语言对比 O(N*M)
			if main[num+i] == sub[j] {
				num += 1
			} else {
				break
			}
		}
		if num == len(sub) {
			return i
		}
	}

	return -1
}

func RK(main, sub string) int {
	num := 0

	for len(sub)+num <= len(main) {
		pk := main[num : len(sub)+num] //这里是直接用高级语言的string对比，通过算出这个区间的hash值与sub的hash值对比。O(N)
		if pk == sub {
			return num
		}
		num += 1
	}

	return -1
}
