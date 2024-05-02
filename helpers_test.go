package gormpager_test

import (
	"os"
	"testing"

	"github.com/go-faker/faker/v4"
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

func createAndGenerateTestingModel(t *testing.T, db *gorm.DB, testExpectations testingExpects[TestingModel]) gormpager.Page[TestingModel] {
	savedData := sliceGenerator(int(testExpectations.expectedTotalEntries), func() *TestingModel {
		return &TestingModel{
			Name:   faker.Name(),
			Age:    uint(faker.RandomUnixTime()),
			Hobbie: faker.Name(),
			UserID: uint(testExpectations.expectedUserId),
		}
	})

	if res := db.Create(&savedData); res.Error != nil {
		t.Error(res.Error)
		t.FailNow()
	}

	return gormpager.Page[TestingModel]{
		PageSize:    testExpectations.expectedPageSize,
		CurrentPage: testExpectations.expectedCurrentPage,
	}
}

func validator[T any](t *testing.T, gotPage gormpager.Page[T], testExpectations testingExpects[T]) {
	failer := func(message string, expectedData, got interface{}) {
		t.Logf("%s \n\tgot: %v \texpected: %v", message, got, expectedData)
		t.FailNow()
	}

	switch {
	case gotPage.CurrentPage != testExpectations.expectedCurrentPage:
		failer("unexpected current page", testExpectations.expectedCurrentPage, gotPage.CurrentPage)
		fallthrough
	case gotPage.TotalEntries != testExpectations.expectedTotalEntries:
		failer("unexpected total entries", testExpectations.expectedTotalEntries, gotPage.TotalEntries)
	case gotPage.PageSize != testExpectations.expectedPageSize:
		failer("unexpected page size", testExpectations.expectedPageSize, gotPage.PageSize)
	case len(gotPage.Data) != int(testExpectations.expectedEntriesCount):
		failer("unexpected data len", testExpectations.expectedEntriesCount, len(gotPage.Data))
	case gotPage.NextPage != testExpectations.expectedNextPage:
		failer("unexpected next page", testExpectations.expectedNextPage, gotPage.NextPage)
	case gotPage.HasNextPage() != testExpectations.expectedHasNextPage:
		failer("unexpected has next page", testExpectations.expectedHasNextPage, gotPage.HasNextPage())
	case gotPage.TotalPages != testExpectations.expectedTotalPages:
		failer("unexpected total pages", testExpectations.expectedTotalPages, gotPage.TotalPages)
	case gotPage.EntriesCount != testExpectations.expectedEntriesCount:
		failer("unexpected entries count", testExpectations.expectedEntriesCount, gotPage.EntriesCount)
	}

}
