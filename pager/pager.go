package pager

// This should be an independent dependency
// (actually it is, but I can't reference it here! )

// Pager represents a pager used to return results by page
type Pager interface {
	Enabled() bool
	Page() int
	// Limit = 0 means no limit so a zeroed Pager means no pagination
	Limit() int
	First() int
	Last() int
	IsVisible(index int) bool
	PageOf(index int) int
}

type defaultPager struct {
	page  int
	limit int
}

// New returns a new pager
func New(page int, limit int) Pager {
	return defaultPager{page, limit}
}

// NoPager builds a valid pager corresponding to "no pager"
// ie: a pager that returns all the results in one page
func NoPager() Pager {
	return defaultPager{0, 0}
}

func (p defaultPager) Enabled() bool {
	return p.limit > 0 && p.page > 0
}

func (p defaultPager) Page() int {
	return p.page
}

func (p defaultPager) Limit() int {
	return p.limit
}

func (p defaultPager) First() int {
	if !p.Enabled() {
		return 0
	}

	return p.limit * (p.page - 1)
}

func (p defaultPager) Last() int {
	if !p.Enabled() {
		return 0
	}

	return p.First() + p.limit - 1
}

func (p defaultPager) IsVisible(index int) bool {
	if !p.Enabled() {
		return true
	}

	return index >= p.First() && index <= p.Last()
}

func (p defaultPager) PageOf(index int) int {
	if !p.Enabled() {
		return 1
	}

	return index/p.limit + 1
}
