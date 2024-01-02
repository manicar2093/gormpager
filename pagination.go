package gormpager

import (
	"math"

	"gorm.io/gorm"
)

// SelectPages get data using requested query
func (c *Page[T]) SelectPages(pager *GormPager, query *gorm.DB) error {
	var model T
	query.Model(&model).Count(&c.TotalEntries)
	c.checkPageSizeLimits(pager.options.PageSizeLowerLimit, pager.options.PageSizeUpperLimit)
	c.setTotalPages()
	if err := c.checkPageCanBeCreated(); err != nil {
		return err
	}
	if res := query.Model(&model).Scopes(func(db *gorm.DB) *gorm.DB {
		offset := c.getOffset()
		return db.Offset(offset).Limit(int(c.PageSize))
	}).Find(&c.Data); res.Error != nil {
		return res.Error
	}
	c.calculateNextPage()
	c.setEntriesCount()
	return nil
}

// SelectPagesRow is used when you have a raw sql and you want to paginate that result
func (c *Page[T]) SelectPagesRow(pager *GormPager, countRawQuery, rawQuery *gorm.DB) error {
	countRawQuery.Scan(&c.TotalEntries)
	c.checkPageSizeLimits(pager.options.PageSizeLowerLimit, pager.options.PageSizeUpperLimit)
	c.setTotalPages()
	if err := c.checkPageCanBeCreated(); err != nil {
		return err
	}

	if res := rawQuery.Offset(c.getOffset()).Limit(int(c.PageSize)).Scan(&c.Data); res.Error != nil {
		return res.Error
	}
	c.calculateNextPage()
	c.setEntriesCount()
	return nil
}

// HasNextPage says if there is o
func (c *Page[T]) HasNextPage() bool {
	return c.CurrentPage < c.TotalPages
}

func (c *Page[T]) setTotalPages() {
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

func (c *Page[T]) checkPageSizeLimits(lower, upper uint) {
	if c.PageSize > int64(upper) {
		c.PageSize = int64(upper)
		return
	}
	if c.PageSize < int64(lower) {
		c.PageSize = int64(lower)
		return
	}
}

func (c *Page[T]) checkPageCanBeCreated() error {
	if c.CurrentPage > c.TotalPages {
		return ErrPageNotExists
	}
	return nil
}

func (c *Page[T]) calculateNextPage() {
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

func (c *Page[T]) setEntriesCount() {
	c.EntriesCount = int64(len(c.Data))
}

func (c *Page[T]) getOffset() int {
	return int((c.CurrentPage - 1) * c.PageSize)
}
