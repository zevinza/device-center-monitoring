package lib

import (
	"encoding/json"
	"log"
	"strings"
	"testing"
	"time"
)

func TestPassword(t *testing.T) {
	plain := "password"

	for range 10 {
		log.Println(PasswordEncrypt(plain))
	}
}

func TestPasswordCompare(t *testing.T) {
	plain := "ieWmDVZk"
	enc := "$2a$10$Nu3WmB5t19nteY9ul99hz.X1Di1VfxUrTH6ouVlhEpwhgjuJZO5Mi"

	log.Println(PasswordCompare(enc, plain))
}

func TestSwap(t *testing.T) {
	log.Println(1)
	a := "foo"
	b := "      fo       o  "

	log.Println(b)
	b = TrimSpace(b)
	log.Println(b)
	log.Println(strings.EqualFold(a, b))

	log.Println(a, b)

	Swap(&a, &b)

	log.Println(a, b)
}

func TestPointer(t *testing.T) {
	var (
		a   *string
		b   *int
		ti  *time.Time
		c   = Pointer("s")
		d   = Pointer(9)
		now = Pointer(time.Now())
	)

	var g []string = nil

	log.Println(json.Marshal(g))

	ti = nil

	log.Println(Rev(a))
	log.Println(Rev(b))
	log.Println(Rev(ti))
	log.Println(Rev(c))
	log.Println(Rev(d))
	log.Println(Rev(now))
}

func TestOperator(t *testing.T) {
	a := 7
	log.Println(If(a%2 == 0, "even", "odd"))
}

func TestRemoveSlice(t *testing.T) {
	slc := []int{1, 2, 3, 4, 5, 6, 6, 7}
	log.Println(RemoveSlice(slc, 1, 6))
}

func TestIsStarted(t *testing.T) {
	n := time.Now().AddDate(0, 0, -1)
	log.Println(IsStarted(n))
}

func TestRemoveSpecialChars(t *testing.T) {
	log.Println(RemoveSpecialChars("JKFK.J(DU234!@#$#@$#)"))
}

func TestIsVald(t *testing.T) {
	a := ""
	log.Println(IsValidPtr(&a))
}
