package lib

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// StrToInt make string to int
func StrToInt(s string) int {
	if i, ok := strconv.Atoi(s); ok == nil {
		return i
	}
	return 0
}

// StrToUUID make string to uuid, fallback is nil
func StrToUUID(s string) *uuid.UUID {
	if u, err := uuid.Parse(s); err == nil {
		return &u
	}

	return nil
}

func JsonString(v any) string {
	by, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(by)
}

func StrToBool(s string) bool {
	return strings.EqualFold(s, "true")
}

// IntToRoman make roman style from given int
func IntToRoman(num int) string {
	var roman string = ""
	var numbers = []int{1, 4, 5, 9, 10, 40, 50, 90, 100, 400, 500, 900, 1000}
	var romans = []string{"I", "IV", "V", "IX", "X", "XL", "L", "XC", "C", "CD", "D", "CM", "M"}
	var index = len(romans) - 1

	for num > 0 {
		for numbers[index] <= num {
			roman += romans[index]
			num -= numbers[index]
		}
		index -= 1
	}

	return roman
}

// Swap value from a to b and b to a
func Swap[T any](a, b *T) {
	var x T = *a
	*a = *b
	*b = x
}

// If will return ok value if condition is true, and not ok if condition false
//
// this func is used to fill data within one line
func If[T any](cond bool, ok, notOk T) T {
	if cond {
		return ok
	}
	return notOk
}

// LowerValue return lowest value of the list
func LowestValue[T sortable](val ...T) T {
	var lowest T
	for _, v := range val {
		if v < lowest {
			lowest = v
		}
	}
	return lowest
}

// HighestValue return highest value of the list
func HighestValue[T sortable](val ...T) T {
	var highest T
	for _, v := range val {
		if v > highest {
			highest = v
		}
	}
	return highest
}

func StructToMap[T any](v T) (map[string]interface{}, error) {
	var mapData map[string]interface{}
	by, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(by, &mapData); err != nil {
		return nil, err
	}
	return mapData, nil
}
