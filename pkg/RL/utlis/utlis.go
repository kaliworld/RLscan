package utlis

import "strings"

// UniqueInts 数组去重
func UniqueInts(input []int) []int {
	unique := make([]int, 0, len(input)) // Allocate memory for the new slice
	set := make(map[int]bool)            // A set to store the encountered elements

	for _, elem := range input {
		if !set[elem] { // If the element is not yet in the set, then it is unique
			set[elem] = true              // Mark the element as encountered
			unique = append(unique, elem) // Append the unique element to the result slice
		}
	}
	return unique // Return the result slice containing unique elements
}

// RemoveZeros 函数会移除切片中所有的零
func RemoveZeros(slice []int) []int {
	// 创建一个新的切片来存放结果，这个切片最开始是空的
	var result []int

	// 遍历原始切片
	for _, value := range slice {
		// 只将非零元素追加到结果切片中
		if value != 0 {
			result = append(result, value)
		}
	}

	// 返回结果切片
	return result
}

// LinesFromString 将字符串拆分为行，并返回字符串切片
func LinesFromString(str string) []string {
	lines := strings.Split(str, "\n")
	for i, line := range lines {
		lines[i] = strings.Replace(line, " open", "", -1) // `-1` 代表替换所有出现
	}

	return lines
}

// MergeSlicesExcludeDuplicates 合并两个字符数组，并排除重复性和空值
func MergeSlicesExcludeDuplicates(slice1, slice2 []string) []string {
	// 使用map帮助排除重复项和空字符串
	uniqueItems := make(map[string]bool)

	// 定义一个把切片加入map的函数
	addSliceToMap := func(slice []string) {
		for _, item := range slice {
			if item != "" && !uniqueItems[item] {
				uniqueItems[item] = true
			}
		}
	}

	// 合并两个切片
	addSliceToMap(slice1)
	addSliceToMap(slice2)

	// 从map创建一个结果切片
	mergedSlice := make([]string, 0, len(uniqueItems))
	for item := range uniqueItems {
		mergedSlice = append(mergedSlice, item)
	}

	return mergedSlice
}
