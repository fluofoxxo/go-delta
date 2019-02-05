// -----------------------------------------------------------------------------
// (c) balarabe@protonmail.com                                      License: MIT
// :v: 2019-02-05 16:22:32 8E63AB                  go-delta/[experiment_test.go]
// -----------------------------------------------------------------------------

package delta

/*
to generate a test coverage report for the whole module use:
    go test -coverprofile cover.out
    go tool cover -html=cover.out
*/

import (
	"testing"
)

const RunExperiments = false

// -----------------------------------------------------------------------------
// # Experimental / Auxiliary Tests

// go test --run Test01
func Test01(t *testing.T) {
	if !RunExperiments {
		return
	}
	if PrintTestNames {
		printTestName()
	}
	PL("Test01 " + Line)
	//
	var m1 = newChunkMap(readData("test1.zip"))
	PL("Created m1. len(m1):", len(m1))
	//
	var m2 = newChunkMap(readData("test2.zip"))
	PL("Created m2. len(m2):", len(m2))
	//
	if false {
		const MaxLines = 0
		var i = 1
		for k, v := range m1 {
			PL("key:", k, "val:", v)
			i++
			if i > MaxLines {
				break
			}
		}
	}
	if true {
		for k, v := range m2 {
			_, exist := m1[k]
			PL("key:", k, "val:", v, "exist:", exist)
		}
	}
} //                                                                      Test01

// go test --run Test02
func Test02(t *testing.T) {
	if !RunExperiments {
		return
	}
	if PrintTestNames {
		printTestName()
	}
	var a, b []byte
	switch 5 {
	case 1:
		a = ab(AtoM + " " + AtoS + " " + AtoZ)
		b = ab("0x0x0x" + AtoZ + " " + AtoZ + " " + AtoZ + " " + Nums)
	case 2:
		a = ab(AtoM + " " + AtoS + " " + AtoZ)
		b = ab(atoz + " " + atoz + " " + atoz + " " + Nums)
	case 3:
		/*
			Target array's size: 16,994,304 bytes
			-
			Before optimizing newChunkMap():
			--------------------------------------------------------------
			uncompressed delta length: 1,855,440 bytes
			compressed delta length:     704,583 (4.15% of target's size)
			elapsed time:              171.4 seconds
			--------------------------------------------------------------
			171.25880: delta.Make
			  0.16411: hashOfBytes
			  3.78551: newChunkMap
			165.82172: longestMatch
			  0.09878: write
			  0.13109: compressBytes
			-
			After optimizing newChunkMap():
			--------------------------------------------------------------
			uncompressed delta length: 1,952,772 bytes
			compressed delta length:     729,574 (4.29% of target's size)
			elapsed time:                2.4 seconds
			--------------------------------------------------------------
			  2.40135: delta.Make
			  0.11608: hashOfBytes
			  1.28985: newChunkMap
			  0.14999: longestMatch
			  0.07882: write
			  0.09806: compressBytes
			-
			After adding backward-scanning in longestMatch()
			--------------------------------------------------------------
			uncompressed delta length: 1,675,811 bytes
			compressed delta length:     666,880 (3.92% of target's size)
			elapsed time:                    2.4 seconds
			--------------------------------------------------------------
			  2.45898: delta.Make
			  0.15910: hashOfBytes
			  1.49399: newChunkMap
			  0.16595: longestMatch
			  0.07311: write
			  0.12408: compressBytes
		*/
		a = readData("test1.file")
		b = readData("test2.file")
		PL("loaded data")
	case 4:
		/*
			target size:        10,356,821
			uncompressed delta:  5,414,754
			compressed delta:    5,258,684 (50.7% of file size)
			elapsed time:              6.2 seconds
		*/
		a = readData("test1.zip")
		b = readData("test2.zip")
		PL("loaded data")
	case 5:
		/*
				target size:        17,096,704 bytes
				uncompressed delta:     64,081 bytes
				compressed delta:       25,967 (50.7% of file size)
				elapsed time:             2.06 seconds
				--------------------------------------------------------------
			  	  2.06019: delta.Make
				  0.11507: hashOfBytes
				  1.44146: newChunkMap
				  0.05109: longestMatch
				  0.00349: write
				  0.00600: compressBytes
				  3.67731
		*/
		a = readData("day1.data")
		b = readData("day2.data")
		PL("loaded data")
	}
	if DebugTiming {
		tmr.Start("delta.Make")
	}
	{
		var d = Make(a, b)
		d.Bytes()
	}
	if DebugTiming {
		tmr.Stop("delta.Make")
		tmr.Print()
	}
} //                                                                      Test02

// go test --run Test03
func Test03(t *testing.T) {
	if !RunExperiments {
		return
	}
	if PrintTestNames {
		printTestName()
	}
	var a, b []byte
	switch 1 {
	case 1:
		a = ab(AtoM + " " + AtoS + " " + AtoZ)
		b = ab("000" + AtoZ + " " + AtoZ + " " + AtoZ + " " + Nums)
	}
	// -------------------------------------------------------------------------
	PL("\n" + Line)
	var d1 = Make(a, b)
	PL("CREATED d1:")
	d1.Dump()
	//
	var dbytes = d1.Bytes()
	PL("got 'dbytes'")
	// -------------------------------------------------------------------------
	PL("\n" + Line)
	if DebugTiming {
		tmr.Start("Load")
	}
	var d2, err = Load(dbytes)
	PL("CREATED d2: err:", err)
	d2.Dump()
	if DebugTiming {
		tmr.Stop("Load")
		tmr.Print()
	}
} //                                                                      Test03

// go test --run Test04
func Test04(t *testing.T) {
	if !RunExperiments {
		return
	}
	if PrintTestNames {
		printTestName()
	}
	var d = Delta{
		sourceSize: 111,
		sourceHash: []byte("SOURCE"),
		targetSize: 222,
		targetHash: []byte("TARGET"),
		newCount:   333,
		oldCount:   444,
		parts: []deltaPart{
			{},
			{},
		},
	}
	PL(d.GoString())
} //                                                                      Test04

//end
