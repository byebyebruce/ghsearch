package ghsearch

import (
	"fmt"
	"testing"
)

func TestSearchRepo(t *testing.T) {
	r, err := SearchRepo("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", 1, "go", "grpc")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(r)
}
