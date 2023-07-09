package utils

import "goList/types"

func IsEmptyList(tasks []types.Task) bool {
	return len(tasks) == 0
}