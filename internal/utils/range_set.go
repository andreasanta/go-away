package utils

// Courtesy of https://github.com/rolyatmax/go-rangeset/blob/master/range_set.go
// With some performance hacks by Andrea Santa <andrea.santa@gmail.com>

import (
	"math"
)

type Range struct {
	Low, High uint32
	Metadata  map[string]string
}

type RangeSet struct {
	Ranges []Range
}

func (rs *RangeSet) AddInts(nums []uint32) {
	for _, num := range nums {
		if len(rs.Ranges) == 0 {
			rs.Ranges = append(rs.Ranges, Range{num, num, nil})
			continue
		}

		for j, curRange := range rs.Ranges {
			low := curRange.Low
			high := curRange.High
			isLastLoop := len(rs.Ranges)-1 == j

			if contains(curRange, num) {
				break
			}

			if low-1 == num {
				rs.Ranges[j].Low = num
				break
			}

			if high+1 == num {
				rs.Ranges[j].High = num
				if !isLastLoop {
					nextRange := rs.Ranges[j+1]
					if nextRange.Low-1 == num {
						// closes a gap
						rs.Ranges = splice(rs.Ranges, j, 2, Range{low, nextRange.High, nil})
					}
				}
				break
			}

			if num < low {
				rs.Ranges = splice(rs.Ranges, j, 0, Range{num, num, nil})
				break
			}

			// if none of the previous ranges or gaps contain the num
			if isLastLoop {
				rs.Ranges = append(rs.Ranges, Range{num, num, nil})
			}
		}
	}
}

func (rs *RangeSet) RemoveInts(nums []uint32) {
	for _, num := range nums {
		for j, curRange := range rs.Ranges {
			if !contains(curRange, num) {
				continue
			}

			low := curRange.Low
			high := curRange.High

			if low == num && high == num {
				rs.Ranges = remove(rs.Ranges, j, 1)
			} else if low == num {
				rs.Ranges[j].Low = low + 1
			} else if high == num {
				rs.Ranges[j].High = high - 1
			} else {
				rs.Ranges = splice(rs.Ranges, j, 1, Range{low, num - 1, curRange.Metadata})
				rs.Ranges = splice(rs.Ranges, j+1, 0, Range{num + 1, high, curRange.Metadata})
			}
			break
		}
	}
}

func (rs *RangeSet) AddRange(r Range) {
	if r.Low > r.High {
		// throw an error
		//log.Panic("Low above high")
	}

	if len(rs.Ranges) == 0 {
		rs.Ranges = append(rs.Ranges, r)
		//log.Printf("Appending first range")
		return
	}

	var overlapStart uint32
	overlapStartIdx := -1
	for i, curRange := range rs.Ranges {
		// if the range comes before all the other ranges with no overlap
		if r.High < curRange.Low-1 {
			//log.Printf("Simple splice range")
			rs.Ranges = splice(rs.Ranges, i, 0, r)
			return
		}

		if overlapStartIdx == -1 && hasOverlap(curRange, r) {
			overlapStartIdx = i
			overlapStart = curRange.Low
		}

		isLastLoop := len(rs.Ranges)-1 == i
		if overlapStartIdx == -1 && isLastLoop {
			// last loop and no overlapStart found
			// it must come after all the other ranges
			//log.Printf("Simple append range")
			rs.Ranges = append(rs.Ranges, r)
			return
		}

		isLastOverlap := isLastLoop || !hasOverlap(r, rs.Ranges[i+1])
		if overlapStartIdx != -1 && isLastOverlap {
			// curRange is the last overlapping range
			low := math.Min(float64(overlapStart), float64(r.Low))
			high := math.Max(float64(curRange.High), float64(r.High))
			overlappingRangeCount := i - overlapStartIdx + 1
			newRange := Range{uint32(low), uint32(high), r.Metadata}
			//log.Printf("Advanced splice range")
			rs.Ranges = splice(rs.Ranges, overlapStartIdx, overlappingRangeCount, newRange)
			return
		}
	}
}

func (rs *RangeSet) RemoveRange(r Range) {
	if r.Low > r.High {
		// throw an error
	}

	var rangesToRemove []int
	for i, curRange := range rs.Ranges {
		if r.High < curRange.Low {
			break
		}

		if r.Low > curRange.High {
			continue
		}

		if r.Low <= curRange.Low {
			if r.High < curRange.High {
				rs.Ranges[i].Low = r.High + 1
			} else {
				rangesToRemove = append(rangesToRemove, i)
			}
		} else {
			if r.High >= curRange.High {
				rs.Ranges[i].High = r.Low - 1
			} else {
				rs.Ranges = splice(rs.Ranges, i, 1, Range{curRange.Low, r.Low - 1, curRange.Metadata})
				rs.Ranges = splice(rs.Ranges, i+1, 0, Range{r.High + 1, curRange.High, curRange.Metadata})
				return
			}
		}
	}
	if len(rangesToRemove) != 0 {
		rs.Ranges = remove(rs.Ranges, rangesToRemove[0], len(rangesToRemove))
	}
}

func (rs *RangeSet) Contains(num uint32) *Range {

	/*
	 * We perform binary chop here, instead of complete linear lookup.
	 *
	 * This has been modified by Andrea to improve performance when numerous
	 * ranges are introduced.
	 */
	begin := 0
	end := len(rs.Ranges)
	var midpoint int
	var curRange Range

	for end > begin {

		midpoint = (end + begin) / 2

		curRange = rs.Ranges[midpoint]
		if contains(curRange, num) {
			return &curRange
		}

		if num < curRange.Low {
			end = midpoint - 1
		} else {
			begin = midpoint + 1
		}

	}

	return nil
}

// helpers

func contains(r Range, num uint32) bool {
	return num >= r.Low && num <= r.High
}

func splice(ranges []Range, startIdx int, elCount int, toInsert Range) []Range {
	temp := make([]Range, startIdx)
	copy(temp, ranges)
	temp = append(temp, toInsert)
	return append(temp, ranges[startIdx+elCount:]...)
}

func remove(ranges []Range, startIdx int, elCount int) []Range {
	return append(ranges[:startIdx], ranges[startIdx+elCount:]...)
}

func hasOverlap(rangeOne, rangeTwo Range) bool {
	var lowest, highest Range
	if rangeOne.Low <= rangeTwo.Low {
		lowest = rangeOne
		highest = rangeTwo
	} else {
		lowest = rangeTwo
		highest = rangeOne
	}
	return lowest.High >= highest.Low-1
}
