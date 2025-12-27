package utils

import (
	"strconv"
	"strings"
)

func ParseFromStrToArray(s string) ([]int, error) {
	splitedArray := strings.Split(s, ",")
	resultArray := make([]int, len(splitedArray))
	
	for idx, splitedStr := range splitedArray {
		curNumber, err := strconv.Atoi(splitedStr)
		if err != nil {
			GetLogger().Error(err.Error())
			return nil, err
		}

		resultArray[idx] = curNumber
	}

	return resultArray, nil
}

func ParseFromArrayToStr(arr []int) string {
	strList := make([]string, len(arr))

	for idx, number := range arr {
		curStr := strconv.Itoa(number)
		strList[idx] = curStr
	}

	return strings.Join(strList, ",")
}





