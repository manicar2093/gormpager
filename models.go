package gormpager

type (
	Options struct {
		PageSizeUpperLimit uint
		PageSizeLowerLimit uint
	}

	Page[T any] struct {
		// CurrentPage is which page I am in
		CurrentPage int64 `json:"current_page"`
		// PageSize is how many items can I expect in Data
		PageSize int64 `json:"page_size"`
		// TotalEntries indicates how many data exists
		TotalEntries int64 `json:"total_entries"`
		// TotalPages is how many pages are in total
		TotalPages int64 `json:"total_pages"`
		// NextPage says what page is next from current page. If current page is the
		// last, this data will return to first page
		NextPage int64 `json:"next_page"`
		// EntriesCount indicates how many data the current page contains
		EntriesCount int64 `json:"entries_count"`
		// Data is what was found in db
		Data []T `json:"data"`
	}
)
