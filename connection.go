package gormpager

import (
	"gorm.io/gorm"
)

type GormPager struct {
	*gorm.DB
	options Options
}

// WrapGormDB set default options to wrapper
// PageSizeUpperLimit = 100
// PageSizeLowerLimit = 10
func WrapGormDB(conn *gorm.DB) *GormPager {
	return &GormPager{
		DB: conn,
		options: Options{
			PageSizeUpperLimit: defaultUpperPageSize,
			PageSizeLowerLimit: defaultLowerPageSize,
		},
	}
}

func WrapGormDBWithOptions(conn *gorm.DB, options Options) *GormPager {
	switch {
	case options.PageSizeLowerLimit == 0:
		options.PageSizeLowerLimit = defaultLowerPageSize
		fallthrough
	case options.PageSizeUpperLimit == 0:
		options.PageSizeUpperLimit = defaultUpperPageSize
	}

	return &GormPager{
		DB:      conn,
		options: options,
	}
}
