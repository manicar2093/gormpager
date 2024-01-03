package gormpager_test

import (
	"os"
	"testing"

	"github.com/manicar2093/gormpager"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func sliceGenerator[T any](count int, generator func() T) []T {
	res := []T{}
	for i := 0; i < count; i++ {
		res = append(res, generator())
	}
	return res
}

func getDbConnection(t *testing.T) *gorm.DB {
	db, err := gorm.Open(postgres.Open(os.Getenv("TEST_DB_DNS")))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := db.AutoMigrate(&TestingModel{}); err != nil {
		t.Error(err)
		t.FailNow()
	}
	return db
}

func validator[T any](t *testing.T, expectedPage gormpager.Page[T], testExpectations testingExpects[T]) {
	failer := func(message string, expectedData interface{}) {
		t.Log(message, expectedData)
		t.FailNow()
	}

	if expectedPage.CurrentPage != testExpectations.expectedCurrentPage {
		failer("unexpected current page", expectedPage.CurrentPage)
	}
	if expectedPage.TotalEntries != testExpectations.expectedTotalEntries {
		failer("unexpected total entries", expectedPage.TotalEntries)
	}
	if expectedPage.PageSize != testExpectations.expectedPageSize {
		failer("unexpected page size", expectedPage.PageSize)
	}
	if len(expectedPage.Data) != int(testExpectations.expectedLenData) {
		failer("unexpected data len", len(expectedPage.Data))
	}
	if expectedPage.NextPage != testExpectations.expectedNextPage {
		failer("unexpected next page", expectedPage.NextPage)
	}
	if expectedPage.HasNextPage() != testExpectations.expectedHasNextPage {
		failer("unexpected has next page", expectedPage.HasNextPage())
	}
	if expectedPage.TotalPages != testExpectations.expectedTotalPages {
		failer("unexpected total pages", expectedPage.TotalPages)
	}
	if expectedPage.EntriesCount != testExpectations.expectedLenData {
		failer("unexpected entries count", expectedPage.EntriesCount)
	}
}
