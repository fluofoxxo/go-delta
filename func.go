// -----------------------------------------------------------------------------
// (c) balarabe@protonmail.com                                      License: MIT
// :v: 2019-01-16 15:07:20 554E1D                             go-delta/[func.go]
// -----------------------------------------------------------------------------

package bdelta

import (
	"bytes"
	"compress/zlib"
	"crypto/sha512"
)

// ApplyDiff __
func ApplyDiff(source []byte, diff Diff) []byte {
	if DebugTiming {
		tmr.Start("ApplyDiff")
		defer tmr.Stop("ApplyDiff")
	}
	return []byte{}
} //                                                                   ApplyDiff

// MakeDiff given two byte arrays 'a' and 'b', calculates the binary
// delta difference between the two arrays and returns it as a Diff.
// You can then use ApplyDiff() to generate 'b' from 'a' the Diff.
func MakeDiff(a, b []byte) Diff {
	if DebugTiming {
		tmr.Start("MakeDiff")
		defer tmr.Stop("MakeDiff")
	}
	var ret = Diff{targetHash: makeHash(b)}
	var lenB = len(b)
	if lenB < ChunkSize {
		ret.parts = []diffPart{{sourceLoc: -1, size: lenB, data: b}}
		return ret
	}
	var m = makeMap(a)
	var chunk [ChunkSize]byte
	for i := 0; i < lenB; {
		if lenB-i < ChunkSize {
			ret.appendPart(-1, lenB-i, b[i:])
			ret.newCount++
			break
		}
		var locs []int
		var found = false
		if lenB-i >= ChunkSize {
			copy(chunk[:], b[i:])
			locs, found = m[chunk]
		}
		if found {
			var at, size = longestMatch(a, locs, b, i)
			ret.appendPart(at, size, nil)
			i += size
			ret.oldCount++
			continue
		}
		ret.appendPart(-1, ChunkSize, chunk[:])
		i += ChunkSize
		ret.newCount++
	}
	ret.sourceHash = makeHash(a)
	return ret
} //                                                                    MakeDiff

// -----------------------------------------------------------------------------
// # Helper Functions

// compressBytes compresses an array of bytes and
// returns the ZLIB-compressed array of bytes.
func compressBytes(data []byte) []byte {
	if DebugTiming {
		tmr.Start("compressBytes")
		defer tmr.Stop("compressBytes")
	}
	if len(data) == 0 {
		return []byte{}
	}
	// zip data in standard manner
	var b bytes.Buffer
	var w = zlib.NewWriter(&b)
	var _, err = w.Write(data)
	w.Close()
	//
	// log any problem
	const ERRM = "Failed compressing data with zlib:"
	if err != nil {
		mod.Error(ERRM, err)
		return []byte{}
	}
	var ret = b.Bytes()
	if len(ret) < 3 {
		mod.Error(ERRM, "length < 3")
		return []byte{}
	}
	return ret
} //                                                               compressBytes

// longestMatch __
func longestMatch(a []byte, aLocs []int, b []byte, bLoc int) (loc, size int) {
	if DebugTiming {
		tmr.Start("longestMatch")
		defer tmr.Stop("longestMatch")
	}
	if len(aLocs) < 1 {
		mod.Error("aLocs is empty")
		return -1, -1
	}
	var bEnd = len(b) - 1
	if bLoc < 0 || bLoc > bEnd {
		mod.Error("bLoc", bLoc, "out of range [0 -", len(b), "]")
		return -1, -1
	}
	var aEnd = len(a) - 1
	var retLoc = -1
	var retSize = -1
	for _, aLoc := range aLocs {
		var n = ChunkSize
		if !bytes.Equal(a[aLoc:aLoc+n], b[bLoc:bLoc+n]) {
			mod.Error("mismatch at aLoc:", aLoc, "bLoc:", bLoc)
			continue
		}
		// extend match forward
		for aLoc+n <= aEnd && bLoc+n <= bEnd && a[aLoc+n] == b[bLoc+n] {
			n++
		}
		if n > retSize {
			retLoc = aLoc
			retSize = n
		}
	}
	return retLoc, retSize
} //                                                                longestMatch

// makeHash returns the SHA-512 hash of byte slice 'data'.
func makeHash(data []byte) []byte {
	if DebugTiming {
		tmr.Start("makeHash")
		defer tmr.Stop("makeHash")
	}
	var ret = sha512.Sum512(data)
	return ret[:]
} //                                                                    makeHash

// makeMap creates a map of unique chunks in 'data'.
// The key specifies the unique chunk of bytes, while the
// values array returns the positions of the chunk in 'data'.
func makeMap(data []byte) (ret map[[ChunkSize]byte][]int) {
	if DebugTiming {
		tmr.Start("makeMap")
		defer tmr.Stop("makeMap")
	}
	ret = make(map[[ChunkSize]byte][]int, 0)
	if len(data) < ChunkSize {
		return ret
	}
	var key [ChunkSize]byte
	for i := 0; i < len(data)-ChunkSize; i++ {
		copy(key[:], data[i:])
		var _, found = ret[key]
		if found {
			ret[key] = append(ret[key], i)
		} else {
			ret[key] = []int{i}
		}
	}
	return ret
} //                                                                     makeMap

//end
