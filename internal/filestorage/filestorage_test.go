package filestorage

import (
	"testing"

	"github.com/tommjj/go-opb/internal/config"
)

func TestRW(t *testing.T) {
	file := New(".config.json")

	data := config.Config{}

	err := file.Load(&data)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", data)
}
