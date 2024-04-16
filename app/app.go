package app

func checkWinner(b Board) string {
	// 승리 조건을 나타내는 조합들
	win := [][]int{
		{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, // 수평
		{1, 4, 7}, {2, 5, 8}, {3, 6, 9}, // 수직
		{1, 5, 9}, {3, 5, 7}, // 대각선
	}

	// 각 위치에 대해 X와 O의 위치를 정렬
	xPositions := b.O
	oPositions := b.X

	// 각 승리 조건을 검사하여 승자를 확인
	for _, condition := range win {
		a, b, c := condition[0], condition[1], condition[2]
		if containsAll(xPositions, a, b, c) {
			return "O"
		} else if containsAll(oPositions, a, b, c) {
			return "X"
		}
	}

	// 승자가 없는 경우
	return ""
}

func containsAll(slice []int, values ...int) bool {
	for _, value := range values {
		if !contains(slice, value) {
			return false
		}
	}
	return true
}

// contains 함수는 슬라이스에 특정 요소가 포함되어 있는지 확인합니다.
func contains(slice []int, value int) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
