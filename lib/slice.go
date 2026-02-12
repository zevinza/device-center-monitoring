package lib

import (
	"sort"

	"github.com/spf13/cast"
)

// Index return value of slice at index
//
// this is safe to use to avoid index out of range
func Index[T any](slice []T, index int) T {
	var zero T
	if index < 0 || index >= len(slice) {
		return zero
	}

	return slice[index]
}

// FillSlice fill in slice with data in given index
func FillSlice[T any](slice []T, value T, index int) []T {
	if index == 0 {
		var first []T
		first = append(first, value)
		first = append(first, slice...)
		return first
	}

	if index == len(slice) {
		slice = append(slice, value)
		return slice
	}

	// todo : insert in middle of array

	return slice
}

type IntTypes interface {
	string | float64 | float32 | int64 | int32 | int16 | int8 | int | uint64 | uint32 | uint16 | uint8 | uint
}

// SliceToInt convert slice if supported data type to []int
func SliceToInt[T IntTypes](slice []T) []int {
	var res []int
	for _, s := range slice {
		res = append(res, cast.ToInt(s))
	}
	return res
}

// LastLoop return true if this is the last iteration
func LastLoop[T any](loop []T, i int) bool {
	return i == len(loop)-1
}

// FirstLoop return true if this is the first iteration
func FirstLoop(i int) bool {
	return i == 0
}

// RemoveDuplicate to remove duplicate data in slice
func RemoveDuplicate[T comparable](slice []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range slice {
		if _, ok := allKeys[item]; !ok {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

// FindInSlice return true when given key found in slice
func FindInSlice[T comparable](key T, slice ...T) bool {
	for _, val := range slice {
		if key == val {
			return true
		}
	}

	return false
}

func NotIn[T comparable](key T, from ...T) bool {
	for _, s := range from {
		if key == s {
			return false
		}
	}
	return true
}

func RemoveNil[T any](v []*T) []*T {
	var new []*T
	for i := range v {
		if v[i] != nil {
			new = append(new, v[i])
		}
	}
	return new
}

func RemoveSlice[T comparable](slc []T, keys ...T) []T {
	var res []T
	for _, s := range slc {
		if !FindInSlice(s, keys...) {
			res = append(res, s)
		}
	}
	return res
}

// Pages to trim slice with given page and limit like in query db
func Pages[T any](v []T, page, limit int) []T {
	count := len(v)
	if limit > 0 {
		if page > 0 {
			page--
		}

		offset := page * limit
		if offset > count {
			return nil
		}
		if offset+limit > count {
			return v[offset:]
		}
		return v[offset : offset+limit]
	}
	return v
}

func IsValidSlicePtr[T any](v []*T) bool {
	for i := range v {
		if v[i] == nil {
			return false
		}
	}
	return true
}

type sortable interface {
	string | int | int64 | float32 | float64
}

// CompareSlice to compare 2 slice, index matter by default ([a, b] != [b, a])
// ignoreIndex to ignoring index and compare value only ([a, b] == [b, a])
func CompareSlice[T sortable](a, b []T, ignoreIndex ...bool) bool {
	if len(a) != len(b) {
		return false
	}

	if len(ignoreIndex) > 0 && ignoreIndex[0] {
		sort.Slice(a, func(i, j int) bool {
			return a[i] < a[j]
		})

		sort.Slice(b, func(i, j int) bool {
			return b[i] < b[j]
		})
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func SliceIntersection[T comparable](a, b []T) []T {
	m := make(map[T]bool)
	for _, v := range a {
		m[v] = true
	}

	var res []T
	for _, v := range b {
		if _, ok := m[v]; ok {
			res = append(res, v)
		}
	}

	return res
}

func SliceRev[T any](slice []*T) []T {
	var res []T
	for _, s := range slice {
		res = append(res, Rev(s))
	}
	return res
}

func SliceRemove[T any](slice []T, i int) []T {
	return append(slice[:1], slice[i+1:]...)
}

func ReverseArray[T any](slice []T) []T {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

// // MakeSlice make single data to slice
// func MakeSlice[T any](v T) []T {
// 	var slice []T
// 	slice = append(slice, v)

// 	return slice
// }
