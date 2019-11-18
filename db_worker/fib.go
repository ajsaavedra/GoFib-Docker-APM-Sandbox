package main

import (
	"strconv"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func recurFib(n int) int {
	if n == 0 || n == 1 {
		return n
	}
	return recurFib(n-1) + recurFib(n-2)
}

func memoFib(n int, memo map[int]int, parentSpan tracer.Span) int {
	name := "memoFib(" + strconv.Itoa(n) + ")"
	span := tracer.StartSpan(name, tracer.ChildOf(parentSpan.Context()))
	defer span.Finish()

	if n == 0 || n == 1 || memo[n] != 0 {
		span.SetTag("memoized_value", memo[n])
		return memo[n]
	}
	memo[n] = memoFib(n - 1, memo, span) + memoFib(n -2, memo, span)
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
