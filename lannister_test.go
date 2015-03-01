package main

import "testing"

func TestTranlatePostPath(t *testing.T) {
	cases := []struct {
		in, date, expected string
	}{
		{"/posts/2015-02-01-something.md", "2015-02-01 15:04", "2015/02/01/something.md"},
		{"/posts/2014-12-25-somethingelse.md", "2014-12-25 12:30", "2014/12/25/somethingelse.md"},
	}
	for _, c := range cases {
		p := new(Post)
		meta := map[interface{}]interface{}{
			"date": c.date,
		}
		p.Filename = c.in
		p.Metadata = meta
		actual, err := p.DestPath()
		if err != nil {
			t.Errorf("translatePostPath threw error: %s", err)
		}
		if actual != c.expected {
			t.Errorf("translatePostPath(%q) == %q, expected %q", c.in, actual, c.expected)
		}
	}
}
