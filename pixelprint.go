// Package pixelprint fingerprints decoded images for duplicate detection.
package pixelprint

import (
	"crypto/sha256"
	"encoding/binary"
	"image"
)

// Fingerprint hashes a non-nil image's dimensions and decoded pixel samples.
// It writes the width and height as big-endian uint64 values, then each
// pixel's four premultiplied RGBA values as big-endian uint16 values, scanning
// from top left to bottom right. The bounds origin is ignored.
func Fingerprint(img image.Image) [sha256.Size]byte {
	bounds := img.Bounds()
	digest := sha256.New()
	var dimensions [16]byte
	binary.BigEndian.PutUint64(dimensions[:8], uint64(bounds.Dx()))
	binary.BigEndian.PutUint64(dimensions[8:], uint64(bounds.Dy()))
	digest.Write(dimensions[:])

	var pixels [4096]byte
	offset := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			binary.BigEndian.PutUint16(pixels[offset:offset+2], uint16(r))
			binary.BigEndian.PutUint16(pixels[offset+2:offset+4], uint16(g))
			binary.BigEndian.PutUint16(pixels[offset+4:offset+6], uint16(b))
			binary.BigEndian.PutUint16(pixels[offset+6:offset+8], uint16(a))
			offset += 8
			if offset == len(pixels) {
				digest.Write(pixels[:])
				offset = 0
			}
		}
	}
	digest.Write(pixels[:offset])

	var fingerprint [sha256.Size]byte
	digest.Sum(fingerprint[:0])
	return fingerprint
}

// DuplicateGroups returns groups of input indexes with identical fingerprints.
// Each image must be non-nil. Groups follow the order of their first member,
// and indexes within each group follow input order. Unique images are omitted.
func DuplicateGroups(images []image.Image) [][]int {
	groupByFingerprint := make(map[[sha256.Size]byte]int, len(images))
	var groups [][]int
	for index, img := range images {
		fingerprint := Fingerprint(img)
		groupIndex, found := groupByFingerprint[fingerprint]
		if !found {
			groupByFingerprint[fingerprint] = len(groups)
			groups = append(groups, []int{index})
			continue
		}
		groups[groupIndex] = append(groups[groupIndex], index)
	}

	var duplicates [][]int
	for _, group := range groups {
		if len(group) > 1 {
			duplicates = append(duplicates, group)
		}
	}
	return duplicates
}
