// -----------------------------------------------------------------------------
// (c) balarabe@protonmail.com                                      License: MIT
// :v: 2019-01-15 20:06:41 D6BF78                             go-delta/[diff.go]
// -----------------------------------------------------------------------------

package bdelta

// Diff __
type Diff struct {
	sourceHash []byte
	targetHash []byte
	parts      []diffPart
} //                                                                        Diff

// diffPart __
type diffPart struct {
	sourceLoc int
	size      int
	data      []byte
} //                                                                    diffPart

// writePart __
func (ob *Diff) writePart(sourceLoc, size int, data []byte) {
	PL("WP",
		"sourceLoc:", sourceLoc,
		"size:", size,
		"data:", data, string(data))
} //                                                                   writePart

//end
