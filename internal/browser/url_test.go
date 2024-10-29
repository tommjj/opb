package browser

import "testing"

func TestBuildQuery(t *testing.T) {
	url, num := BuildQuery("https://www.google.com/search?q=$&code=$", "$")

	t.Log(url, num)

}
