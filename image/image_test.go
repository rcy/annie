package image

import (
	"context"
	"fmt"
	"testing"
)

func TestGenerate(t *testing.T) {
	ctx := context.Background()
	prompt := "A beautiful sunset over the mountains."
	gi, err := Generate(ctx, prompt)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(gi)
}
