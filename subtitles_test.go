package subtitles

import (
	"io"
	"reflect"
	"strings"
	"testing"
	"time"
)

type collector struct {
	collection []Subtitle
}

func newCollector() *collector {
	return &collector{make([]Subtitle, 0)}
}

func (c *collector) cb() SubtitleCallback {
	return func(s Subtitle) (bool, error) {
		c.collection = append(c.collection, s)
		return false, nil
	}
}

func mustParseTime(t string) time.Time {
	if pt, err := parseSubRipTime(t); err != nil {
		panic(err)
	} else {
		return pt
	}
}

func TestParseSubRip(t *testing.T) {

	tests := []struct {
		name    string
		r       io.Reader
		want    []Subtitle
		wantErr bool
	}{
		{
			"works with empty stream",
			strings.NewReader(""),
			[]Subtitle{},
			false,
		},
		{
			"works with single subtitle",
			strings.NewReader(`1
00:00:35,684 --> 00:00:37,054
Lorem ipsum dolor sit amet,
consectetur adipiscing elit.`),
			[]Subtitle{
				{
					id:    1,
					start: mustParseTime("00:00:35,684"),
					end:   mustParseTime("00:00:37,054"),
					text: `Lorem ipsum dolor sit amet,
consectetur adipiscing elit.`,
				},
			},
			false,
		},
		{
			"works with multiple subtitles",
			strings.NewReader(`1
00:00:35,684 --> 00:00:37,054
Lorem ipsum dolor sit amet,
consectetur adipiscing elit.

2
00:00:37,184 --> 00:00:40,454
Donec aliquet arcu enim, quis bibendum felis cursus a.

`),
			[]Subtitle{
				{
					id:    1,
					start: mustParseTime("00:00:35,684"),
					end:   mustParseTime("00:00:37,054"),
					text: `Lorem ipsum dolor sit amet,
consectetur adipiscing elit.`,
				},
				{
					id:    2,
					start: mustParseTime("00:00:37,184"),
					end:   mustParseTime("00:00:40,454"),
					text:  `Donec aliquet arcu enim, quis bibendum felis cursus a.`,
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newCollector()
			if err := ParseSubRip(tt.r, c.cb()); (err != nil) != tt.wantErr {
				t.Errorf("ParseSubRip() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(c.collection, tt.want) {
				t.Errorf("parseSubRipTime() = %#v, want %#v", c.collection, tt.want)
			}
		})
	}
}

func Test_parseSubRipTime(t *testing.T) {
	type args struct {
		t string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSubRipTime(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSubRipTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseSubRipTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
