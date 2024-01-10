package gormpager

import (
	"math"

	"gorm.io/gorm"
)

// SelectPages get data using requested query
func (c *Page[T]) SelectPages(pager *GormPager, query *gorm.DB) error {
	var model T
	if res := query.Model(&model).Count(&c.TotalEntries); res.Error != nil {
		return res.Error
	}
	if err := c.setPageSizeAndTotalPages(pager.options.PageSizeLowerLimit, pager.options.PageSizeUpperLimit); err != nil {
		return err
	}
	if res := query.Model(&model).Scopes(func(db *gorm.DB) *gorm.DB {
		offset := c.getOffset()
		return db.Offset(offset).Limit(int(c.PageSize))
	}).Find(&c.Data); res.Error != nil {
		return res.Error
	}

	c.setNextPageAndEntriesCount()

	return nil
}

// SelectPagesRow is used when you have a raw sql and you want to paginate that result
func (c *Page[T]) SelectPagesRow(pager *GormPager, countRawQuery, rawQuery *gorm.DB) error {
	if res := countRawQuery.Scan(&c.TotalEntries); res.Error != nil {
		return res.Error
	}

	if err := c.setPageSizeAndTotalPages(pager.options.PageSizeLowerLimit, pager.options.PageSizeUpperLimit); err != nil {
		return err
	}

	if res := rawQuery.Offset(c.getOffset()).Limit(int(c.PageSize)).Scan(&c.Data); res.Error != nil {
		return res.Error
	}

	c.setNextPageAndEntriesCount()

	return nil
}

// HasNextPage says if there is o
func (c *Page[T]) HasNextPage() bool {
	return c.CurrentPage < c.TotalPages
}

func (c *Page[T]) setPageSizeAndTotalPages(lower, upper uint) error {
	setTotalPages := func() {
		var (
			totalEntriesF = float64(c.TotalEntries)
			pageSizeF     = float64(c.PageSize)
			totalPagesF   = math.Ceil(totalEntriesF / pageSizeF)
		)
		c.TotalPages = int64(totalPagesF)
		if c.TotalPages == 0 {
			c.TotalPages = 1
		}
	}

	checkPageSizeLimits := func(lower, upper uint) {
		if c.PageSize > int64(upper) {
			c.PageSize = int64(upper)
			return
		}
		if c.PageSize < int64(lower) {
			c.PageSize = int64(lower)
			return
		}
	}

	checkPageSizeLimits(lower, upper)
	setTotalPages()
	return c.checkPageCanBeCreated()
}

func (c *Page[T]) checkPageCanBeCreated() error {
	if c.CurrentPage > c.TotalPages {
		return ErrPageNotExists
	}
	return nil
}

func (c *Page[T]) setNextPageAndEntriesCount() {
	calculateNextPage := func() {
		if len(c.Data) == 0 {
			c.NextPage = 1
			return
		}
		if c.CurrentPage == c.TotalPages {
			c.NextPage = 1
			return
		}
		c.NextPage = c.CurrentPage + 1
	}
	setEntriesCount := func() {
		c.EntriesCount = int64(len(c.Data))
	}

	calculateNextPage()
	setEntriesCount()
}

func (c *Page[T]) getOffset() int {
	return int((c.CurrentPage - 1) * c.PageSize)
}
