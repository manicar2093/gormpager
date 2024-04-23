package gormpager_test

import (
	"testing"

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
			expectedPage = createAndGenerateTestingModel(t, db, testExpectations)
		)

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
			expectedPage = createAndGenerateTestingModel(t, db, testExpectations)
		)

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
			expectedPage = createAndGenerateTestingModel(t, db, testExpectations)
		)

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
			expectedPage = createAndGenerateTestingModel(t, db, testExpectations)
		)

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
			expectedPage = createAndGenerateTestingModel(t, db, testExpectations)
		)

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
			expectedPage = createAndGenerateTestingModel(t, db, testExpectations)
		)

		if err := expectedPage.SelectPages(pager, db.Where("user_id = ?", testExpectations.expectedUserId)); err != nil {
			t.Error(err)
			t.FailNow()
		}
	})

	t.Run("creates a connection with deafult options when options is empty", func(t *testing.T) {
		var (
			pagerWithOptions = gormpager.WrapGormDBWithOptions(
				db,
				gormpager.Options{},
			)
			testExpectations = testingExpects[TestingModel]{
				expectedTotalEntries: 44,
				expectedPageSize:     10,
				expectedLenData:      10,
				expectedTotalPages:   5,
				expectedCurrentPage:  1,
				expectedNextPage:     2,
				expectedHasNextPage:  true,
				expectedUserId:       8,
			}
			expectedPage = createAndGenerateTestingModel(t, db, testExpectations)
		)

		if err := expectedPage.SelectPages(pagerWithOptions, db.Where("user_id = ?", testExpectations.expectedUserId)); err != nil {
			t.Error(err)
			t.FailNow()
		}

		validator(t, expectedPage, testExpectations)
	})

	t.Run("creates a new page with options", func(t *testing.T) {
		var (
			pagerWithOptions = gormpager.WrapGormDBWithOptions(
				db,
				gormpager.Options{
					PageSizeLowerLimit: 5,
				},
			)
			testExpectations = testingExpects[TestingModel]{
				expectedTotalEntries: 44,
				expectedPageSize:     5,
				expectedLenData:      5,
				expectedTotalPages:   9,
				expectedCurrentPage:  1,
				expectedNextPage:     2,
				expectedHasNextPage:  true,
				expectedUserId:       9,
			}
			gotPage = createAndGenerateTestingModel(t, db, testExpectations)
		)

		if err := gotPage.SelectPages(pagerWithOptions, db.Where("user_id = ?", testExpectations.expectedUserId)); err != nil {
			t.Error(err)
			t.FailNow()
		}

		validator(t, gotPage, testExpectations)
	})
}
