package utils

// LCMOfNums is a function to find Least Common Multiple (LCM) via GCD (Greatest Common Divisor) of an array of unsigned integers.
func LCMOfNums(nums []uint64) uint64 {
	result := nums[0]
	for i := 1; i < len(nums); i++ {
		result = lcm(result, nums[i])
	}
	return result
}

// lcm is a function to find Least Common Multiple (LCM) via GCD (Greatest Common Divisor) of two unsigned integers.
func lcm(a, b uint64) uint64 {
	g := gcd(a, b)
	return a * b / g
}

// gcd is a function to find Greatest Common Divisor (GCD) of two unsigned integers.
func gcd(a, b uint64) uint64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
