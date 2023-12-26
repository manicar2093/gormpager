package gormpager_test

import (
	"log"
	"os"
	"testing"

	"github.com/manicar2093/gormpager"

	"github.com/bxcodec/faker/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TestingModel struct {
	gorm.Model
	Name   string `json:"name,omitempty"`
	Age    uint   `json:"age,omitempty"`
	Hobbie string `json:"hobbies,omitempty"`
	UserID uint
}

func SliceGenerator[T any](count int, generator func() T) []T {
	res := []T{}
	for i := 0; i < count; i++ {
		res = append(res, generator())
	}
	return res
}

func TestPaginator(t *testing.T) {

	db, err := gorm.Open(postgres.Open(os.Getenv("TEST_DB_DNS")))
	if err != nil {
		log.Panicln(err)
	}
	if err := db.AutoMigrate(&TestingModel{}); err != nil {
		t.Error(err)
		t.FailNow()
	}

	pager := gormpager.WrapGormDB(db)

	t.Cleanup(func() {
		db.Exec("TRUNCATE TABLE testing_models")
	})

	t.Run("creates a new page", func(t *testing.T) {
		var (
			expectedTotalEntries int64 = 44
			savedData2                 = SliceGenerator(int(expectedTotalEntries), func() *TestingModel {
				return &TestingModel{
					Name:   faker.Name(),
					Age:    uint(faker.RandomUnixTime()),
					Hobbie: faker.Name(),
					UserID: 2,
				}
			})
			expectedPageSize    int64 = 15
			expectedLenData     int64 = expectedPageSize
			expectedTotalPages  int64 = 3
			expectedCurrentPage int64 = 1
			expectedNextPage    int64 = 2
			expectedHasNextPage bool  = true
			expectedPage              = gormpager.Page[TestingModel]{
				PageSize:    expectedPageSize,
				CurrentPage: expectedCurrentPage,
			}
		)
		db.Create(&savedData2)

		if err := expectedPage.SelectPages(pager, db.Where("user_id = 2")); err != nil {
			t.Error(err)
			t.FailNow()
		}

		if expectedPage.CurrentPage != expectedCurrentPage {
			t.Log("unexpected current page", expectedPage.CurrentPage)
			t.FailNow()
		}
		if expectedPage.TotalEntries != expectedTotalEntries {
			t.Log("unexpected total entries", expectedPage.TotalEntries)
			t.FailNow()
		}
		if expectedPage.PageSize != expectedPageSize {
			t.Log("unexpected page size", expectedPage.PageSize)
			t.FailNow()
		}
		if len(expectedPage.Data) != int(expectedLenData) {
			t.Log("unexpected data len", len(expectedPage.Data))
			t.FailNow()
		}
		if expectedPage.NextPage != expectedNextPage {
			t.Log("unexpected next page", expectedPage.NextPage)
			t.FailNow()
		}
		if expectedPage.HasNextPage() != expectedHasNextPage {
			t.Log("unexpected has next page", expectedPage.HasNextPage())
			t.FailNow()
		}
		if expectedPage.TotalPages != expectedTotalPages {
			t.Log("unexpected total pages", expectedPage.TotalPages)
			t.FailNow()
		}
		if expectedPage.EntriesCount != expectedLenData {
			t.Log("unexpected entries count", expectedPage.EntriesCount)
			t.FailNow()
		}
	})

	t.Run("when there is a single page with few entries", func(t *testing.T) {
		var (
			expectedTotalEntries int64 = 2
			savedData2                 = SliceGenerator(int(expectedTotalEntries), func() *TestingModel {
				return &TestingModel{
					Name:   faker.Name(),
					Age:    uint(faker.RandomUnixTime()),
					Hobbie: faker.Name(),
					UserID: 3,
				}
			})
			expectedPageSize    int64 = 15
			expectedLenData     int64 = 2
			expectedTotalPages  int64 = 1
			expectedCurrentPage int64 = 1
			expectedNextPage    int64 = 1
			expectedHasNextPage bool  = false
			expectedPage              = gormpager.Page[TestingModel]{
				PageSize:    expectedPageSize,
				CurrentPage: expectedCurrentPage,
			}
		)
		db.Create(&savedData2)

		if err := expectedPage.SelectPages(pager, db.Where("user_id = 3")); err != nil {
			t.Error(err)
			t.FailNow()
		}

		if expectedPage.CurrentPage != expectedCurrentPage {
			t.Log("unexpected current page", expectedPage.CurrentPage)
			t.FailNow()
		}
		if expectedPage.TotalEntries != expectedTotalEntries {
			t.Log("unexpected total entries", expectedPage.TotalEntries)
			t.FailNow()
		}
		if expectedPage.PageSize != expectedPageSize {
			t.Log("unexpected page size", expectedPage.PageSize)
			t.FailNow()
		}
		if len(expectedPage.Data) != int(expectedLenData) {
			t.Log("unexpected data len", len(expectedPage.Data))
			t.FailNow()
		}
		if expectedPage.NextPage != expectedNextPage {
			t.Log("unexpected next page", expectedPage.NextPage)
			t.FailNow()
		}
		if expectedPage.HasNextPage() != expectedHasNextPage {
			t.Log("unexpected has next page", expectedPage.HasNextPage())
			t.FailNow()
		}
		if expectedPage.TotalPages != expectedTotalPages {
			t.Log("unexpected total pages", expectedPage.TotalPages)
			t.FailNow()
		}
		if expectedPage.EntriesCount != expectedLenData {
			t.Log("unexpected entries count", expectedPage.EntriesCount)
			t.FailNow()
		}
	})

	t.Run("when first page has nothing", func(t *testing.T) {
		var (
			expectedTotalEntries int64 = 0
			expectedPageSize     int64 = 10
			expectedLenData      int64 = 0
			expectedTotalPages   int64 = 1
			expectedCurrentPage  int64 = 1
			expectedNextPage     int64 = 1
			expectedHasNextPage  bool  = false
			expectedPage               = gormpager.Page[TestingModel]{
				PageSize:    expectedPageSize,
				CurrentPage: expectedCurrentPage,
			}
		)

		if err := expectedPage.SelectPages(pager, db.Where("user_id = 4")); err != nil {
			t.Error(err)
			t.FailNow()
		}

		if expectedPage.CurrentPage != expectedCurrentPage {
			t.Log("unexpected current page", expectedPage.CurrentPage)
			t.FailNow()
		}
		if expectedPage.TotalEntries != expectedTotalEntries {
			t.Log("unexpected total entries", expectedPage.TotalEntries)
			t.FailNow()
		}
		if expectedPage.PageSize != expectedPageSize {
			t.Log("unexpected page size", expectedPage.PageSize)
			t.FailNow()
		}
		if len(expectedPage.Data) != int(expectedLenData) {
			t.Log("unexpected data len", len(expectedPage.Data))
			t.FailNow()
		}
		if expectedPage.NextPage != expectedNextPage {
			t.Log("unexpected next page", expectedPage.NextPage)
			t.FailNow()
		}
		if expectedPage.HasNextPage() != expectedHasNextPage {
			t.Log("unexpected has next page", expectedPage.HasNextPage())
			t.FailNow()
		}
		if expectedPage.TotalPages != expectedTotalPages {
			t.Log("unexpected total pages", expectedPage.TotalPages)
			t.FailNow()
		}
		if expectedPage.EntriesCount != expectedLenData {
			t.Log("unexpected entries count", expectedPage.EntriesCount)
			t.FailNow()
		}
	})
}
