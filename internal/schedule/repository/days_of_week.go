package repository

func toSmallIntDays(days []int) []int16 {
	result := make([]int16, 0, len(days))
	for _, day := range days {
		result = append(result, int16(day))
	}

	return result
}

func fromSmallIntDays(days []int16) []int {
	result := make([]int, 0, len(days))
	for _, day := range days {
		result = append(result, int(day))
	}

	return result
}
