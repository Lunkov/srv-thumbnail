package main

import (
  "os"
  "io"
  "path/filepath"
  "strings"
  "fmt"
  
  "image"
  "image/jpeg"
	"image/png"
  
  "github.com/nfnt/resize"

  "github.com/golang/glog"
)

// Format is an image file format.
type Format int

// Image file formats.
const (
	JPEG Format = iota
	PNG
)

var formatExts = map[string]Format{
	"jpg":  JPEG,
	"jpeg": JPEG,
	"png":  PNG,
}

type SrcImage struct {
  Name         string
  Type         Format
  Img          image.Image
}

func newFilename(filename string, width uint, height uint) string {
  dir, file := filepath.Split(filename)
  newfile := filepath.Clean(fmt.Sprintf("%s/i/%dx%d/%s", dir, width, height, file))
  os.MkdirAll(filepath.Dir(newfile), os.ModePerm)
  return newfile
}

// FormatFromExtension parses image format from formatFromFilename:
// "jpg" (or "jpeg"), "png", "gif", "tif" (or "tiff") and "bmp" are supported.
func formatFromFilename(filename string) Format {
	if f, ok := formatExts[strings.ToLower(strings.TrimPrefix(filepath.Ext(filename), "."))]; ok {
		return f
	}
  glog.Errorf("ERR: ErrUnsupportedFormat: filename='%s'", filename)
	return -1
}

func (i *SrcImage) load(filename string) bool {
  i.Type = formatFromFilename(filename)
  if i.Type == -1 {
    glog.Errorf("ERR: Image Format(%s): ErrUnsupportedFormat", filename)
    return false
  }

  handle, err := os.Open(filename)
  if err != nil {
    glog.Errorf("ERR: Image Open(%s): %v", filename, err  )
    return false
  }
  defer handle.Close()
  
  var errImg error
  i.Img, _, errImg = image.Decode(handle)
  if errImg != nil {
    glog.Errorf("ERR: Image Decode(%s): %v", filename, errImg)
    return false
  }
  
  i.Name = filename
  return true
}

func (i *SrcImage) resizeImage(width uint, height uint) bool {
  m := resize.Resize(width, height, i.Img, resize.Lanczos3)
  newName := newFilename(i.Name, width, height)
	file, err := fs.Create(newName)
	if err != nil {
    glog.Errorf("ERR: Create(%s): %v", newName, err)
		return false
	}
	ok := encodeImage(file, m, i.Type)
  file.Close()
	return ok
}

// Encode writes the image img to w in the specified format (JPEG, PNG, GIF, TIFF or BMP).
func encodeImage(w io.Writer, img image.Image, format Format) bool {
	switch format {
	case JPEG:
		return jpeg.Encode(w, img, &jpeg.Options{Quality: 80}) != nil

	case PNG:
		encoder := png.Encoder{CompressionLevel: 90 }
		return encoder.Encode(w, img) != nil

	}
  glog.Errorf("ERR: ErrUnsupportedFormat")
	return false
}

