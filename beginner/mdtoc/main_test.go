package main

import (
	"fmt"
	"testing"
)

func Test_extractHeadline(t *testing.T) {
	tests := []struct {
		from  string
		level int
		title string
	}{
		{"foo", 0, ""},
		{"# foo", 1, "foo"},
		{"something # foo", 0, ""},
		{"## foo bar baz", 2, "foo bar baz"},
		{"##    foo bar baz   ", 2, "foo bar baz"},
		{"### foo bar baz", 3, "foo bar baz"},
		{"#### foo bar baz", 4, "foo bar baz"},
		{"##### foo bar baz", 5, "foo bar baz"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%02d", i+1), func(t *testing.T) {
			level, title := extractHeadline(tt.from)
			if tt.level != level || tt.title != title {
				t.Errorf("expected `%s` not to be level=%d, title=`%s`, but got level=%d, title=`%s`", tt.from, tt.level, tt.title, level, title)
			}
		})
	}
}

func Test_titleAsLink(t *testing.T) {
	tests := []struct {
		from string
		to   string
	}{
		{"foo", "foo"},
		{"FOO", "foo"},
		{"foo Bar", "foo-bar"},
		{"foo Bar BAZ", "foo-bar-baz"},
		{"foo    Bar    BAZ", "foo-bar-baz"},
		{"foo + Bar - BAZ = z??oing", "foo-bar-baz-z-oing"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%02d", i+1), func(t *testing.T) {
			if v := titleAsLink(tt.from); v != tt.to {
				t.Errorf("expect `%s` -> `%s` but got `%s`", tt.from, tt.to, v)
			}
		})
	}
}
