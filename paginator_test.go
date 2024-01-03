package gormpager_test

import (
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/manicar2093/gormpager"

	"gorm.io/gorm"
)

type (
	TestingModel struct {
		gorm.Model
		Name   string `json:"name,omitempty"`
		Age    uint   `json:"age,omitempty"`
		Hobbie string `json:"hobbies,omitempty"`
		UserID uint
	}

	testingExpects[T any] struct {
		expectedTotalEntries int64
		expectedPageSize     int64
		expectedLenData      int64
		expectedTotalPages   int64
		expectedCurrentPage  int64
		expectedNextPage     int64
		expectedHasNextPage  bool
		expectedUserId       int64
	}
)

func TestWrapGormDB(t *testing.T) {

	var (
		db    = getDbConnection(t)
		pager = gormpager.WrapGormDB(db)
	)

	t.Cleanup(func() {
		db.Exec("TRUNCATE TABLE testing_models")
	})

	t.Run("creates a new page", func(t *testing.T) {
		var (
			testExpectations = testingExpects[TestingModel]{
				expectedTotalEntries: 44,
				expectedPageSize:     15,
				expectedLenData:      15,
				expectedTotalPages:   3,
				expectedCurrentPage:  1,
				expectedNextPage:     2,
				expectedHasNextPage:  true,
				expectedUserId:       1,
			}
			savedData = sliceGenerator(int(testExpectations.expectedTotalEntries), func() *TestingModel {
				return &TestingModel{
					Name:   faker.Name(),
					Age:    uint(faker.RandomUnixTime()),
					Hobbie: faker.Name(),
					UserID: uint(testExpectations.expectedUserId),
				}
			})
			expectedPage = gormpager.Page[TestingModel]{
				PageSize:    testExpectations.expectedPageSize,
				CurrentPage: testExpectations.expectedCurrentPage,
			}
		)
		db.Create(&savedData)

		if err := expectedPage.SelectPages(pager, db.Where("user_id = ?", testExpectations.expectedUserId)); err != nil {
			t.Error(err)
			t.FailNow()
		}

		validator(t, expectedPage, testExpectations)
	})

	t.Run("when there is a single page with few entries", func(t *testing.T) {
		var (
			testExpectations = testingExpects[TestingModel]{
				expectedTotalEntries: 2,
				expectedPageSize:     15,
				expectedLenData:      2,
				expectedTotalPages:   1,
				expectedCurrentPage:  1,
				expectedNextPage:     1,
				expectedHasNextPage:  false,
				expectedUserId:       2,
			}
			savedData = sliceGenerator(int(testExpectations.expectedTotalEntries), func() *TestingModel {
				return &TestingModel{
					Name:   faker.Name(),
					Age:    uint(faker.RandomUnixTime()),
					Hobbie: faker.Name(),
					UserID: uint(testExpectations.expectedUserId),
				}
			})

			expectedPage = gormpager.Page[TestingModel]{
				PageSize:    testExpectations.expectedPageSize,
				CurrentPage: testExpectations.expectedCurrentPage,
			}
		)
		db.Create(&savedData)

		if err := expectedPage.SelectPages(pager, db.Where("user_id = ?", testExpectations.expectedUserId)); err != nil {
			t.Error(err)
			t.FailNow()
		}

		validator(t, expectedPage, testExpectations)
	})

	t.Run("when first page has nothing", func(t *testing.T) {
		var (
			testExpectations = testingExpects[TestingModel]{
				expectedTotalEntries: 0,
				expectedPageSize:     10,
				expectedLenData:      0,
				expectedTotalPages:   1,
				expectedCurrentPage:  1,
				expectedNextPage:     1,
				expectedHasNextPage:  false,
				expectedUserId:       3,
			}
			expectedPage = gormpager.Page[TestingModel]{
				PageSize:    testExpectations.expectedPageSize,
				CurrentPage: testExpectations.expectedCurrentPage,
			}
		)

		if err := expectedPage.SelectPages(pager, db.Where("user_id = ?", testExpectations.expectedUserId)); err != nil {
			t.Error(err)
			t.FailNow()
		}

	})

	t.Run("when there is a single page with few entries", func(t *testing.T) {
		var (
			testExpectations = testingExpects[TestingModel]{
				expectedTotalEntries: 2,
				expectedPageSize:     15,
				expectedLenData:      2,
				expectedTotalPages:   1,
				expectedCurrentPage:  1,
				expectedNextPage:     1,
				expectedHasNextPage:  false,
				expectedUserId:       4,
			}
			savedData = sliceGenerator(int(testExpectations.expectedTotalEntries), func() *TestingModel {
				return &TestingModel{
					Name:   faker.Name(),
					Age:    uint(faker.RandomUnixTime()),
					Hobbie: faker.Name(),
					UserID: 4,
				}
			})

			expectedPage = gormpager.Page[TestingModel]{
				PageSize:    testExpectations.expectedPageSize,
				CurrentPage: testExpectations.expectedCurrentPage,
			}
		)
		db.Create(&savedData)

		if err := expectedPage.SelectPages(pager, db.Where("user_id = ?", testExpectations.expectedUserId)); err != nil {
			t.Error(err)
			t.FailNow()
		}

		validator(t, expectedPage, testExpectations)
	})

	t.Run("gets page size from page when into limits", func(t *testing.T) {
		var (
			testExpectations = testingExpects[TestingModel]{
				expectedTotalEntries: 55,
				expectedPageSize:     40,
				expectedLenData:      40,
				expectedTotalPages:   2,
				expectedCurrentPage:  1,
				expectedNextPage:     2,
				expectedHasNextPage:  true,
				expectedUserId:       5,
			}
			savedData = sliceGenerator(int(testExpectations.expectedTotalEntries), func() *TestingModel {
				return &TestingModel{
					Name:   faker.Name(),
					Age:    uint(faker.RandomUnixTime()),
					Hobbie: faker.Name(),
					UserID: uint(testExpectations.expectedUserId),
				}
			})
			expectedPage = gormpager.Page[TestingModel]{
				PageSize:    testExpectations.expectedPageSize,
				CurrentPage: testExpectations.expectedCurrentPage,
			}
		)
		db.Create(&savedData)

		if err := expectedPage.SelectPages(pager, db.Where("user_id = ?", testExpectations.expectedUserId)); err != nil {
			t.Error(err)
			t.FailNow()
		}

	})

	t.Run("restart page lower limit", func(t *testing.T) {
		var (
			testExpectations = testingExpects[TestingModel]{
				expectedTotalEntries: 55,
				expectedPageSize:     10,
				expectedLenData:      40,
				expectedTotalPages:   2,
				expectedCurrentPage:  1,
				expectedNextPage:     2,
				expectedHasNextPage:  true,
				expectedUserId:       6,
			}
			savedData = sliceGenerator(int(testExpectations.expectedTotalEntries), func() *TestingModel {
				return &TestingModel{
					Name:   faker.Name(),
					Age:    uint(faker.RandomUnixTime()),
					Hobbie: faker.Name(),
					UserID: uint(testExpectations.expectedUserId),
				}
			})
			expectedPage = gormpager.Page[TestingModel]{
				PageSize:    5,
				CurrentPage: testExpectations.expectedCurrentPage,
			}
		)
		db.Create(&savedData)

		if err := expectedPage.SelectPages(pager, db.Where("user_id = ?", testExpectations.expectedUserId)); err != nil {
			t.Error(err)
			t.FailNow()
		}

	})

	t.Run("restart page upper limit", func(t *testing.T) {
		var (
			testExpectations = testingExpects[TestingModel]{
				expectedTotalEntries: 55,
				expectedPageSize:     100,
				expectedLenData:      40,
				expectedTotalPages:   2,
				expectedCurrentPage:  1,
				expectedNextPage:     2,
				expectedHasNextPage:  true,
				expectedUserId:       7,
			}
			savedData = sliceGenerator(int(testExpectations.expectedTotalEntries), func() *TestingModel {
				return &TestingModel{
					Name:   faker.Name(),
					Age:    uint(faker.RandomUnixTime()),
					Hobbie: faker.Name(),
					UserID: uint(testExpectations.expectedUserId),
				}
			})
			expectedPage = gormpager.Page[TestingModel]{
				PageSize:    1000,
				CurrentPage: testExpectations.expectedCurrentPage,
			}
		)
		db.Create(&savedData)

		if err := expectedPage.SelectPages(pager, db.Where("user_id = ?", testExpectations.expectedUserId)); err != nil {
			t.Error(err)
			t.FailNow()
		}

	})
}
