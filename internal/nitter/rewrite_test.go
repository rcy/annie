package nitter

import "testing"

func TestRewrite(t *testing.T) {
	for _, tc := range []struct {
		url  string
		want string
	}{
		{
			url:  "https://twitter.com/status/1234",
			want: "https://nitter.net/status/1234",
		},
		{
			url:  "https://x.com/status/1234",
			want: "https://nitter.net/status/1234",
		},
		{
			url:  "https://fox.com/story",
			want: "https://fox.com/story",
		},
		{
			url:  "twitter.com/foo",
			want: "nitter.net/foo",
		},
	} {
		got := Rewrite(tc.url)
		if got != tc.want {
			t.Errorf("expected %s, got %s", tc.want, got)
		}
	}
}
