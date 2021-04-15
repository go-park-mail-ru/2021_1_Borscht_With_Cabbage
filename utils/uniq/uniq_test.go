package uniq

import (
	"context"
	"testing"
)

func TestGetUniqFilename(t *testing.T) {
	ctx := new(context.Context)
	filename := "file.jpg"

	_, err := GetUniqFilename(*ctx, filename)
	if err != nil {
		t.Errorf("GetUniqFilename error")
		return
	}
}
