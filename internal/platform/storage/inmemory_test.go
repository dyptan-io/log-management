package storage

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type testItem struct {
	Id string
}

func (i testItem) ID() ID {
	return ID(i.Id)
}

func TestInMemory_Insert(t *testing.T) {
	tests := map[string]struct {
		giveItems []testItem
		wantItems []testItem
		wantErr   error
	}{
		"one record": {
			giveItems: []testItem{{Id: "test"}},
			wantItems: []testItem{{Id: "test"}},
		},
		"many records": {
			giveItems: []testItem{{Id: "test-1"}, {Id: "test-2"}},
			wantItems: []testItem{{Id: "test-1"}, {Id: "test-2"}},
		},
		"error - missing ID": {
			giveItems: []testItem{{Id: ""}},
			wantErr:   ErrMissingID,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			store := NewInMemory[testItem]()

			for _, item := range test.giveItems {
				err := store.Insert(item)

				if test.wantErr != nil {
					require.Equal(t, test.wantErr, err)
				} else {
					require.NoError(t, err)
				}
			}

			items, err := store.Find(func(_ testItem) bool {
				// Return everything.
				return true
			})

			require.NoError(t, err)
			require.ElementsMatch(t, test.wantItems, items)
		})
	}
}

func TestInMemory_Get(t *testing.T) {
	tests := map[string]struct {
		giveItems []testItem
		giveID    ID
		wantItem  testItem
		wantErr   error
	}{
		"one record": {
			giveItems: []testItem{{Id: "test"}},
			giveID:    "test",
			wantItem:  testItem{Id: "test"},
		},
		"many records": {
			giveItems: []testItem{{Id: "test-1"}, {Id: "test-2"}},
			giveID:    "test-2",
			wantItem:  testItem{Id: "test-2"},
		},
		"error - not found": {
			giveItems: []testItem{{Id: "test-1"}, {Id: "test-2"}},
			giveID:    "does-not-exist",
			wantErr:   ErrNotFound,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			store := NewInMemory[testItem]()

			for _, item := range test.giveItems {
				require.NoError(t, store.Insert(item))
			}

			item, err := store.Get(test.giveID)

			if test.wantErr != nil {
				require.Equal(t, test.wantErr, err)
			} else {
				require.NoError(t, err)
			}

			require.EqualValues(t, test.wantItem, item)
		})
	}
}

func TestInMemory_Find(t *testing.T) {
	tests := map[string]struct {
		giveItems   []testItem
		giveMatcher Matcher[testItem]
		wantItems   []testItem
		wantErr     error
	}{
		"get all": {
			giveItems:   []testItem{{Id: "test-1"}, {Id: "test-2"}},
			giveMatcher: func(_ testItem) bool { return true },
			wantItems:   []testItem{{Id: "test-1"}, {Id: "test-2"}},
		},
		"get one with predicate": {
			giveItems:   []testItem{{Id: "test-1"}, {Id: "test-2"}, {Id: "test-3"}},
			giveMatcher: func(v testItem) bool { return strings.Contains(v.Id, "2") },
			wantItems:   []testItem{{Id: "test-2"}},
		},
		"not found with predicate": {
			giveItems:   []testItem{{Id: "test-1"}, {Id: "test-2"}, {Id: "test-2"}},
			giveMatcher: func(v testItem) bool { return strings.Contains(v.Id, "something-wierd") },
			wantItems:   nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			store := NewInMemory[testItem]()

			for _, item := range test.giveItems {
				require.NoError(t, store.Insert(item))
			}

			items, err := store.Find(test.giveMatcher)

			if test.wantErr != nil {
				require.Equal(t, test.wantErr, err)
			} else {
				require.NoError(t, err)
			}

			require.ElementsMatch(t, test.wantItems, items)
		})
	}
}
