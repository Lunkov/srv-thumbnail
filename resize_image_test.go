package main

import (
  "flag"
  "testing"
  "github.com/stretchr/testify/assert"
)


/////////////////////////
// TESTS
/////////////////////////
func TestImageThumbnail(t *testing.T) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", ".")
	flag.Set("v", "9")
	flag.Parse()
  
  assert.Equal(t, "/folder1/folder2/i/300x200/image.jpeg", newFilename("/folder1/folder2/image.jpeg", 300, 200))
  assert.Equal(t, "/folder1/i/250x0/image.png", newFilename("/folder1/image.png", 250, 0))


}
