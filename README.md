# GormPager

![gormpager](gopher.png "GormPager")

A page handler for [GORM](https://github.com/go-gorm/gorm)

## Requirements

Use Go 1.20+ is a must. GormPager uses Generics.

## Description

I know, GORM is great but, how do a pager implementation? Well, here it is

## Example

Create your gorm connection:

```golang
db, err := gorm.Open(postgres.Open(os.Getenv("DB_DNS")))
if err != nil {
    log.Panicln(err)
}
```

Wrap it with GormPager

```golang
pager := gormpager.WrapGormDB(db)
```

or

```golang
pager := gormpager.WrapGormDBWithOptions(db, gormpager.Options{
    PageSizeUpperLimit: 1000,
    PageSizeLowerLimit: 5,
})
```

**Note:** Gormpager has default values for PageSize which are 10 and 100. If you need change this default values you can do it througth this gormpager.WrapGormDBWithOptions so you can request a small page.

Start using it!

```golang
page = gormpager.Page[User]{
    CurrentPage: expectedCurrentPage,
}
if err := page.SelectPages(pager, db.Where("user_id = 3")); err != nil {
    log.Panic(err)
}
```

Happy Coding!
