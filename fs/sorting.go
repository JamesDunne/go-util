package fs

import "os"

// For directory entry sorting:

type Entries []os.FileInfo

func (s Entries) Len() int      { return len(s) }
func (s Entries) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type SortBy int

const (
	SortByName SortBy = iota
	SortByDate
	SortBySize
)

type SortDirection int

const (
	SortAscending SortDirection = iota
	SortDescending
)

// Sort by last modified time:
type ByDate struct {
	Entries
	dir SortDirection
}

func (s ByDate) Less(i, j int) bool {
	if s.Entries[i].IsDir() && !s.Entries[j].IsDir() {
		return true
	}
	if !s.Entries[i].IsDir() && s.Entries[j].IsDir() {
		return false
	}

	if s.dir == SortAscending {
		if s.Entries[i].ModTime().Equal(s.Entries[j].ModTime()) {
			return s.Entries[i].Name() > s.Entries[j].Name()
		} else {
			return s.Entries[i].ModTime().Before(s.Entries[j].ModTime())
		}
	} else {
		if s.Entries[i].ModTime().Equal(s.Entries[j].ModTime()) {
			return s.Entries[i].Name() > s.Entries[j].Name()
		} else {
			return s.Entries[i].ModTime().After(s.Entries[j].ModTime())
		}
	}
}
