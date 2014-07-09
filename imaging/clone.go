package imaging

import (
	"fmt"
	"image"
	"image/color"
	"reflect"
)

// Clone returns a copy of the given image in NRGBA kind.
func Clone(img image.Image) *image.NRGBA {
	srcBounds := img.Bounds()
	dstBounds := srcBounds.Sub(srcBounds.Min)

	dst := image.NewNRGBA(dstBounds)

	dstMinX := dstBounds.Min.X
	dstMinY := dstBounds.Min.Y

	srcMinX := srcBounds.Min.X
	srcMinY := srcBounds.Min.Y
	srcMaxX := srcBounds.Max.X
	srcMaxY := srcBounds.Max.Y

	switch src0 := img.(type) {

	case *image.NRGBA:
		rowSize := srcBounds.Dx() * 4
		numRows := srcBounds.Dy()

		i0 := dst.PixOffset(dstMinX, dstMinY)
		j0 := src0.PixOffset(srcMinX, srcMinY)

		di := dst.Stride
		dj := src0.Stride

		for row := 0; row < numRows; row++ {
			copy(dst.Pix[i0:i0+rowSize], src0.Pix[j0:j0+rowSize])
			i0 += di
			j0 += dj
		}

	case *image.NRGBA64:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)

				dst.Pix[i+0] = src0.Pix[j+0]
				dst.Pix[i+1] = src0.Pix[j+2]
				dst.Pix[i+2] = src0.Pix[j+4]
				dst.Pix[i+3] = src0.Pix[j+6]

			}
		}

	case *image.RGBA:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)
				a := src0.Pix[j+3]
				dst.Pix[i+3] = a

				switch a {
				case 0:
					dst.Pix[i+0] = 0
					dst.Pix[i+1] = 0
					dst.Pix[i+2] = 0
				case 0xff:
					dst.Pix[i+0] = src0.Pix[j+0]
					dst.Pix[i+1] = src0.Pix[j+1]
					dst.Pix[i+2] = src0.Pix[j+2]
				default:
					dst.Pix[i+0] = uint8(uint16(src0.Pix[j+0]) * 0xff / uint16(a))
					dst.Pix[i+1] = uint8(uint16(src0.Pix[j+1]) * 0xff / uint16(a))
					dst.Pix[i+2] = uint8(uint16(src0.Pix[j+2]) * 0xff / uint16(a))
				}
			}
		}

	case *image.RGBA64:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)
				a := src0.Pix[j+6]
				dst.Pix[i+3] = a

				switch a {
				case 0:
					dst.Pix[i+0] = 0
					dst.Pix[i+1] = 0
					dst.Pix[i+2] = 0
				case 0xff:
					dst.Pix[i+0] = src0.Pix[j+0]
					dst.Pix[i+1] = src0.Pix[j+2]
					dst.Pix[i+2] = src0.Pix[j+4]
				default:
					dst.Pix[i+0] = uint8(uint16(src0.Pix[j+0]) * 0xff / uint16(a))
					dst.Pix[i+1] = uint8(uint16(src0.Pix[j+2]) * 0xff / uint16(a))
					dst.Pix[i+2] = uint8(uint16(src0.Pix[j+4]) * 0xff / uint16(a))
				}
			}
		}

	case *image.Gray:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)
				c := src0.Pix[j]
				dst.Pix[i+0] = c
				dst.Pix[i+1] = c
				dst.Pix[i+2] = c
				dst.Pix[i+3] = 0xff

			}
		}

	case *image.Gray16:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)
				c := src0.Pix[j]
				dst.Pix[i+0] = c
				dst.Pix[i+1] = c
				dst.Pix[i+2] = c
				dst.Pix[i+3] = 0xff

			}
		}

	case *image.YCbCr:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				yj := src0.YOffset(x, y)
				cj := src0.COffset(x, y)
				r, g, b := color.YCbCrToRGB(src0.Y[yj], src0.Cb[cj], src0.Cr[cj])

				dst.Pix[i+0] = r
				dst.Pix[i+1] = g
				dst.Pix[i+2] = b
				dst.Pix[i+3] = 0xff

			}
		}

	case *image.Paletted:
		plen := len(src0.Palette)
		pnew := make([]color.NRGBA, plen)
		for i := 0; i < plen; i++ {
			pnew[i] = color.NRGBAModel.Convert(src0.Palette[i]).(color.NRGBA)
		}

		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)
				c := pnew[src0.Pix[j]]

				dst.Pix[i+0] = c.R
				dst.Pix[i+1] = c.G
				dst.Pix[i+2] = c.B
				dst.Pix[i+3] = c.A

			}
		}

	default:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				c := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)

				dst.Pix[i+0] = c.R
				dst.Pix[i+1] = c.G
				dst.Pix[i+2] = c.B
				dst.Pix[i+3] = c.A

			}
		}
	}

	return dst
}

// Copies an image to a new image of the same kind:
func CloneKind(src image.Image) image.Image {
	srcBounds := src.Bounds().Canon()
	zeroedBounds := srcBounds.Sub(srcBounds.Min)

	switch si := src.(type) {
	case *image.RGBA:
		out := image.NewRGBA(zeroedBounds)
		for y := srcBounds.Min.Y; y < srcBounds.Max.Y; y++ {
			for x := srcBounds.Min.X; x < srcBounds.Max.X; x++ {
				out.SetRGBA(x-srcBounds.Min.X, y-srcBounds.Min.Y, si.At(x, y).(color.RGBA))
			}
		}
		return out
	case *image.YCbCr:
		out := image.NewYCbCr(zeroedBounds, si.SubsampleRatio)
		for y := srcBounds.Min.Y; y < srcBounds.Max.Y; y++ {
			for x := srcBounds.Min.X; x < srcBounds.Max.X; x++ {
				ycbcr := si.At(x, y).(color.YCbCr)
				yoffs := out.YOffset(x-srcBounds.Min.X, y-srcBounds.Min.Y)
				coffs := out.COffset(x-srcBounds.Min.X, y-srcBounds.Min.Y)
				out.Y[yoffs] = ycbcr.Y
				out.Cb[coffs] = ycbcr.Cb
				out.Cr[coffs] = ycbcr.Cr
			}
		}
		return out
	case *image.Paletted:
		out := image.NewPaletted(zeroedBounds, si.Palette)
		for y := srcBounds.Min.Y; y < srcBounds.Max.Y; y++ {
			for x := srcBounds.Min.X; x < srcBounds.Max.X; x++ {
				out.SetColorIndex(x-srcBounds.Min.X, y-srcBounds.Min.Y, si.ColorIndexAt(x, y))
			}
		}
		return out
	case *image.RGBA64:
		out := image.NewRGBA64(zeroedBounds)
		for y := srcBounds.Min.Y; y < srcBounds.Max.Y; y++ {
			for x := srcBounds.Min.X; x < srcBounds.Max.X; x++ {
				out.SetRGBA64(x-srcBounds.Min.X, y-srcBounds.Min.Y, si.At(x, y).(color.RGBA64))
			}
		}
		return out
	case *image.NRGBA:
		out := image.NewNRGBA(zeroedBounds)
		for y := srcBounds.Min.Y; y < srcBounds.Max.Y; y++ {
			for x := srcBounds.Min.X; x < srcBounds.Max.X; x++ {
				out.SetNRGBA(x-srcBounds.Min.X, y-srcBounds.Min.Y, si.At(x, y).(color.NRGBA))
			}
		}
		return out
	case *image.NRGBA64:
		out := image.NewNRGBA64(zeroedBounds)
		for y := srcBounds.Min.Y; y < srcBounds.Max.Y; y++ {
			for x := srcBounds.Min.X; x < srcBounds.Max.X; x++ {
				out.SetNRGBA64(x-srcBounds.Min.X, y-srcBounds.Min.Y, si.At(x, y).(color.NRGBA64))
			}
		}
		return out
	case *image.Alpha:
		out := image.NewAlpha(zeroedBounds)
		for y := srcBounds.Min.Y; y < srcBounds.Max.Y; y++ {
			for x := srcBounds.Min.X; x < srcBounds.Max.X; x++ {
				out.SetAlpha(x-srcBounds.Min.X, y-srcBounds.Min.Y, si.At(x, y).(color.Alpha))
			}
		}
		return out
	case *image.Alpha16:
		out := image.NewAlpha16(zeroedBounds)
		for y := srcBounds.Min.Y; y < srcBounds.Max.Y; y++ {
			for x := srcBounds.Min.X; x < srcBounds.Max.X; x++ {
				out.SetAlpha16(x-srcBounds.Min.X, y-srcBounds.Min.Y, si.At(x, y).(color.Alpha16))
			}
		}
		return out
	case *image.Gray:
		out := image.NewGray(zeroedBounds)
		for y := srcBounds.Min.Y; y < srcBounds.Max.Y; y++ {
			for x := srcBounds.Min.X; x < srcBounds.Max.X; x++ {
				out.SetGray(x-srcBounds.Min.X, y-srcBounds.Min.Y, si.At(x, y).(color.Gray))
			}
		}
		return out
	case *image.Gray16:
		out := image.NewGray16(zeroedBounds)
		for y := srcBounds.Min.Y; y < srcBounds.Max.Y; y++ {
			for x := srcBounds.Min.X; x < srcBounds.Max.X; x++ {
				out.SetGray16(x-srcBounds.Min.X, y-srcBounds.Min.Y, si.At(x, y).(color.Gray16))
			}
		}
		return out
	default:
		panic(fmt.Errorf("Unhandled image format type: %s", reflect.TypeOf(src).Name()))
	}
}
