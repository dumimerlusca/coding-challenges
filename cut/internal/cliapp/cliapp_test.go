package app

import (
	"testing"
)

func TestParseSelectValue(t *testing.T) {
	got, err := parseSelectValue("1-3,8-10,20")
	if err != nil {
		t.Error(err)
	}

	want := []int{1, 2, 3, 8, 9, 10, 20}

	for i := range want {
		if want[i] != got[i] {
			t.Errorf("want %v, got %v", want[i], got[i])
		}
	}

}
