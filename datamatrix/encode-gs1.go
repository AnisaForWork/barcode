// Package datamatrix can create GS1 Datamatrix barcodes -- y
package datamatrix

import (
	"errors"

	"github.com/boombuler/barcode"
)

// Encode returns a Datamatrix barcode for the given content and color scheme
func EncodeWithColorGS1(content string, color barcode.ColorScheme) (barcode.Barcode, error) {
	data := encodeTextGS1(content)

	var size *dmCodeSize
	for _, s := range codeSizes {
		if s.DataCodewords() >= len(data) {
			size = s
			break
		}
	}
	if size == nil {
		return nil, errors.New("to much data to encode")
	}
	data = addPadding(data, size.DataCodewords())
	data = ec.calcECC(data, size)
	code := render(data, size, color)
	if code != nil {
		code.content = content
		return code, nil
	}
	return nil, errors.New("unable to render barcode")
}

// Encode returns a Datamatrix barcode for the given content
func EncodeGS1(content string) (barcode.Barcode, error) {
	return EncodeWithColor(content, barcode.ColorScheme16)
}

func encodeTextGS1(content string) []byte {
	var result []byte
	input := []byte(content)

	for i := 0; i < len(input); {
		c := input[i]
		i++

		// ensure that we have FNC1 (ascii 232) at the very beginning (GS1 DMX) -- added by y, 2020/04/26
		if i == 1 && c != 232 {
			result = append(result, 232)
		}

		if c >= '0' && c <= '9' && i < len(input) && input[i] >= '0' && input[i] <= '9' {
			// two numbers...
			c2 := input[i]
			i++
			cw := byte(((c-'0')*10 + (c2 - '0')) + 130)
			result = append(result, cw)
		} else if c > 127 {
			// not correct... needs to be redone later...
			result = append(result, 235, c-127)
		} else {
			result = append(result, c+1)
		}
	}
	return result
}
