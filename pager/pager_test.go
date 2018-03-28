package pager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPagerProperties(t *testing.T) {
	assert := assert.New(t)

	fixtures := []struct {
		page    int
		limit   int
		first   int
		last    int
		enabled bool
	}{
		{
			page:    1,
			limit:   1,
			first:   0,
			last:    0,
			enabled: true,
		},
		{
			page:    3,
			limit:   20,
			first:   40,
			last:    59,
			enabled: true,
		},
		{
			page:    5,
			limit:   3,
			first:   12,
			last:    14,
			enabled: true,
		},
		{
			page:    0,
			limit:   0,
			first:   0,
			last:    0,
			enabled: false,
		},
	}

	for _, f := range fixtures {
		p := New(f.page, f.limit)

		assert.Equal(f.page, p.Page())
		assert.Equal(f.limit, p.Limit())
		assert.Equal(f.first, p.First())
		assert.Equal(f.last, p.Last())
		assert.Equal(f.enabled, p.Enabled())
	}
}

func TestPagerIsVisible(t *testing.T) {
	assert := assert.New(t)

	fixtures := []struct {
		page    int
		limit   int
		index   int
		visible bool
	}{
		{
			page:    1,
			limit:   1,
			index:   0,
			visible: true,
		},
		{
			page:    1,
			limit:   1,
			index:   1,
			visible: false,
		},
		{
			page:    2,
			limit:   10,
			index:   9,
			visible: false,
		},
		{
			page:    2,
			limit:   10,
			index:   10,
			visible: true,
		},
		{
			page:    2,
			limit:   10,
			index:   11,
			visible: true,
		},
		{
			page:    2,
			limit:   10,
			index:   19,
			visible: true,
		},
		{
			page:    2,
			limit:   10,
			index:   20,
			visible: false,
		},
		// limit: 0 disables the pager
		{
			page:    10,
			limit:   0,
			index:   9999,
			visible: true,
		},
	}

	for _, f := range fixtures {
		p := New(f.page, f.limit)

		assert.Equal(f.visible, p.IsVisible(f.index))
	}
}

func TestPagerPageOf(t *testing.T) {
	assert := assert.New(t)

	fixtures := []struct {
		limit  int
		index  int
		pageOf int
	}{
		{
			limit:  1,
			index:  0,
			pageOf: 1,
		},
		{
			limit:  1,
			index:  1,
			pageOf: 2,
		},
		{
			limit:  10,
			index:  9,
			pageOf: 1,
		},
		{
			limit:  10,
			index:  10,
			pageOf: 2,
		},
		{
			limit:  10,
			index:  11,
			pageOf: 2,
		},
		{
			limit:  10,
			index:  19,
			pageOf: 2,
		},
		{
			limit:  10,
			index:  20,
			pageOf: 3,
		},
		// limit: 0 disables the pager
		{
			limit:  0,
			index:  9999,
			pageOf: 1,
		},
	}

	for _, f := range fixtures {
		p := New(2, f.limit)

		assert.Equal(f.pageOf, p.PageOf(f.index))
	}
}

func TestNoPager(t *testing.T) {
	assert := assert.New(t)

	p := NoPager()
	assert.Equal(false, p.Enabled())
}
