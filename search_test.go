package ghsearch

import (
	"fmt"
	"testing"
)

func TestSearchRepo(t *testing.T) {
	r, err := SearchRepo("xxx", 1, "go", "grpc")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(r)
}
