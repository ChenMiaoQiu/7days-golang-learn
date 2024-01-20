package consistenthash

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	hash := New(3, func(data []byte) uint32 {
		i, _ := strconv.Atoi(string(data))
		return uint32(i)
	})

	//2, 12, 22, 4, 14, 24, 6, 16, 26
	hash.Add("6", "4", "2")

	testCase := map[string]string{
		"2":  "2",
		"5":  "6",
		"27": "2",
		"23": "4",
		"1":  "2",
	}

	for k, v := range testCase {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should output %s, but output %s\n", k, v, hash.Get(k))
		}
	}

	//8, 18, 28
	hash.Add("8")
	testCase["27"] = "8"
	for k, v := range testCase {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should output %s, but output %s\n", k, v, hash.Get(k))
		}
	}
}
