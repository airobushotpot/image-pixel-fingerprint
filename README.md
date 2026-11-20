# Image Pixel Fingerprint

Find duplicates among downloaded renders by comparing their
decoded pixels. The Go package gives you input indexes so you can decide which
copy to keep in your own asset workflow.

## Usage

Pass images decoded with Go's `image.Decode` or created in memory. The example
checks three generated color swatches, two of which share their pixel data.

```go
package main

import (
	"fmt"
	"image"
	"image/color"

	pixelprint "github.com/airobushotpot/image-pixel-fingerprint"
)

func main() {
	original := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	crop := image.NewNRGBA(image.Rect(5, 5, 6, 6))
	other := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	red := color.NRGBA{R: 255, A: 255}
	original.SetNRGBA(0, 0, red)
	crop.SetNRGBA(5, 5, red)
	other.SetNRGBA(0, 0, color.NRGBA{B: 255, A: 255})
	fmt.Println(pixelprint.DuplicateGroups([]image.Image{original, crop, other}))
	// [[0 1]]
}
```

## Features

- SHA-256 fingerprints include image dimensions and premultiplied RGBA samples.
- File metadata and the pixel coordinate origin do not affect the fingerprint.
- Duplicate groups retain input order so the first copy is easy to identify.
- Uses only the Go standard library. Callers provide decoded images.

This package compares exact decoded samples. Resizing, recompression that changes
pixel values, and visually similar renders produce different fingerprints.

## Related workflow

To create images from a prompt and download candidates for review, see
[AI H3 Image's MiniMax H3 Text to Image workspace](https://aih3image.com/minimax-h3-text-to-image).
