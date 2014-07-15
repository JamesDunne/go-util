package base

import (
	"fmt"
	"io"
	"log"
)

func HexDumpToWriter(b []byte, w io.Writer) {
	const digits = "0123456789abcdef"
	line := make([]byte, 3*16+16)

	lines := len(b) >> 4
	remainder := len(b) & 15

	t := 0
	for i := 0; i < lines; i++ {
		for j := 0; j < 16; j++ {
			d := b[t]
			t++
			line[j*3+0] = digits[(d>>4)&15]
			line[j*3+1] = digits[d&15]
			line[j*3+2] = ' '

			if d >= 32 && d <= 127 {
				line[16*3+j] = d
			} else {
				line[16*3+j] = '.'
			}
		}

		fmt.Fprintln(w, string(line))
	}

	if remainder > 0 {
		for j := 0; j < remainder; j++ {
			d := b[t]
			t++
			line[j*3+0] = digits[(d>>4)&15]
			line[j*3+1] = digits[d&15]
			line[j*3+2] = ' '

			if d >= 32 && d <= 127 {
				line[16*3+j] = d
			} else {
				line[16*3+j] = '.'
			}
		}

		for j := remainder; j < 16; j++ {
			line[j*3+0] = ' '
			line[j*3+1] = ' '
			line[j*3+2] = ' '
			line[16*3+j] = ' '
		}

		fmt.Fprintln(w, string(line))
	}
}

func HexDumpToLogger(b []byte, o *log.Logger) {
	const digits = "0123456789abcdef"
	line := make([]byte, 3*16+16)

	lines := len(b) >> 4
	remainder := len(b) & 15

	t := 0
	for i := 0; i < lines; i++ {
		for j := 0; j < 16; j++ {
			d := b[t]
			t++
			line[j*3+0] = digits[(d>>4)&15]
			line[j*3+1] = digits[d&15]
			line[j*3+2] = ' '

			if d >= 32 && d <= 127 {
				line[16*3+j] = d
			} else {
				line[16*3+j] = '.'
			}
		}

		o.Println(string(line))
	}

	if remainder > 0 {
		for j := 0; j < remainder; j++ {
			d := b[t]
			t++
			line[j*3+0] = digits[(d>>4)&15]
			line[j*3+1] = digits[d&15]
			line[j*3+2] = ' '

			if d >= 32 && d <= 127 {
				line[16*3+j] = d
			} else {
				line[16*3+j] = '.'
			}
		}

		for j := remainder; j < 16; j++ {
			line[j*3+0] = ' '
			line[j*3+1] = ' '
			line[j*3+2] = ' '
			line[16*3+j] = ' '
		}

		o.Println(string(line))
	}
}
