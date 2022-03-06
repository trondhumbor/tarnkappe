package tarnkappe

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"math"
	"os"
)

const useBits = 1

func Hide(inPath, outPath, hidePath string) (int, error) {
	hideFile, hideFileErr := os.Open(hidePath)
	if hideFileErr != nil {
		panic(hideFileErr)
	}
	defer hideFile.Close()

	stat, statErr := hideFile.Stat()
	if statErr != nil {
		panic(statErr)
	}

	content := make([]byte, stat.Size())
	binary.Read(hideFile, binary.LittleEndian, &content)

	inFile, inFileErr := os.Open(inPath)
	if inFileErr != nil {
		panic(inFileErr)
	}
	defer inFile.Close()

	inImg, _, imageErr := image.Decode(inFile)
	if imageErr != nil {
		panic(imageErr)
	}

	rect := inImg.Bounds()
	xDim, yDim := rect.Max.X, rect.Max.Y
	storableBits := (xDim * yDim) * (3 * useBits) // number of bits we can store in the image when using three colors (rgb)
	neededBits := len(content) * 8                // number of bits we need to store the content

	if neededBits > storableBits {
		return 0, fmt.Errorf("cannot fit content into image. image can hold %d bits but %d were given", storableBits, neededBits)
	}

	// convert image to NRGBA
	img := image.NewNRGBA(image.Rect(0, 0, rect.Dx(), rect.Dy()))
	draw.Draw(img, img.Bounds(), inImg, rect.Min, draw.Src)

	var chunks []byte
	for _, b := range content {
		numbytes := math.Ceil(float64(8) / useBits) // how many bytes of size 'useBits' we have to use for storing each 8 bit byte
		for i := 0; i < int(numbytes); i++ {
			v := b >> (useBits * i) & ((1 << useBits) - 1)
			chunks = append(chunks, v)
		}
	}

	i := 0
	for j := 0; j < len(img.Pix); j += 4 { // 4 because r, g, b, a
		if i > len(chunks) {
			break
		}
		r, g, b, a := img.Pix[j], img.Pix[j+1], img.Pix[j+2], img.Pix[j+3]

		var w int
		if i < len(chunks) {
			r = ((r >> useBits) << useBits) | chunks[i]
			w = 1
		}

		if i+1 < len(chunks) {
			g = ((g >> useBits) << useBits) | chunks[i+1]
			w = 2
		}

		if i+2 < len(chunks) {
			b = ((b >> useBits) << useBits) | chunks[i+2]
			w = 3
		}

		img.Pix[j], img.Pix[j+1], img.Pix[j+2], img.Pix[j+3] = r, g, b, a
		i += w
	}

	outFile, outFileErr := os.Create(outPath)
	if outFileErr != nil {
		panic(outFileErr)
	}
	defer outFile.Close()

	png.Encode(outFile, img)
	return len(chunks), nil
}

func Reveal(inPath string, outPath string, length int) {
	inFile, inFileErr := os.Open(inPath)
	if inFileErr != nil {
		panic(inFileErr)
	}
	defer inFile.Close()

	inImg, _, imageErr := image.Decode(inFile)
	if imageErr != nil {
		panic(imageErr)
	}
	rect := inImg.Bounds()

	img := image.NewNRGBA(image.Rect(0, 0, rect.Dx(), rect.Dy()))
	draw.Draw(img, img.Bounds(), inImg, rect.Min, draw.Src)

	var chunks []byte
	i := 0
	for j := 0; j < len(img.Pix); j += 4 { // 4 because r, g, b, a
		if i > length {
			break
		}

		r, g, b := img.Pix[j], img.Pix[j+1], img.Pix[j+2]

		var w int
		if i < length {
			r &= (1 << useBits) - 1
			chunks = append(chunks, r)
			w = 1
		}

		if i+1 < length {
			g &= (1 << useBits) - 1
			chunks = append(chunks, g)
			w = 2
		}

		if i+2 < length {
			b &= (1 << useBits) - 1
			chunks = append(chunks, b)
			w = 3
		}

		i += w
	}

	var content []byte
	numbytes := int(math.Ceil(float64(8) / useBits)) // how many bytes of size 'useBits' each 8 bit byte is split over
	for i := 0; i < length; i += numbytes {
		b := chunks[i : i+numbytes]
		var tmp byte
		for j := 0; j < len(b); j++ {
			tmp |= (b[j] << (useBits * j))
		}
		content = append(content, tmp)
	}

	outFile, outFileErr := os.Create(outPath)
	if outFileErr != nil {
		panic(outFileErr)
	}
	defer outFile.Close()

	binary.Write(outFile, binary.LittleEndian, &content)
}
