package main

import (
	"github.com/gin-gonic/gin"
)

func recurFib(n int) int {
	if n == 0 || n == 1 {
		return n
	}
	return recurFib(n-1) + recurFib(n-2)
}

func memoFib(n int, memo map[int]int) int {
	if n == 0 || n == 1 || memo[n] != 0 {
		return memo[n]
	}
	memo[n] = memoFib(n - 1, memo) + memoFib(n -2, memo)
	return memo[n]
}

func iter(idx int) int {
	if idx == 0 || idx == 1 { return idx }

	nums := make([]int, 0)
	nums = append([]int{0, 1}, nums...)

	for i := 2; i < idx; i++ {
		nums = append(nums, nums[i-1] + nums[i-2])
	}
	return nums[idx - 1]
}
