package imaging

import (
	"fmt"
	"image"
	"reflect"
)

// Calls SubImage on `img` and returns a new image of the same kind:
func SubImageKind(img image.Image, srcBounds image.Rectangle) image.Image {
	switch si := img.(type) {
	case *image.RGBA:
		return si.SubImage(srcBounds)
	case *image.YCbCr:
		return si.SubImage(srcBounds)
	case *image.Paletted:
		return si.SubImage(srcBounds)
	case *image.RGBA64:
		return si.SubImage(srcBounds)
	case *image.NRGBA:
		return si.SubImage(srcBounds)
	case *image.NRGBA64:
		return si.SubImage(srcBounds)
	case *image.Alpha:
		return si.SubImage(srcBounds)
	case *image.Alpha16:
		return si.SubImage(srcBounds)
	case *image.Gray:
		return si.SubImage(srcBounds)
	case *image.Gray16:
		return si.SubImage(srcBounds)
	default:
		panic(fmt.Errorf("Unhandled image format type: %s", reflect.TypeOf(img).Name()))
	}
}
