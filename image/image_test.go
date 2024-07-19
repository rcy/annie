package image

import (
	"context"
	"fmt"
	"testing"
)

func TestGenerateRunpod(t *testing.T) {
	ctx := context.Background()
	prompt := "A beautiful sunset over the mountains."
	gi, err := GenerateRunpod(ctx, prompt)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(gi)
}
