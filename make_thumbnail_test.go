package main

import (
  "flag"
  "testing"
  "github.com/stretchr/testify/assert"
)


/////////////////////////
// TESTS
/////////////////////////
func TestMakeThumbnail(t *testing.T) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", ".")
	flag.Set("v", "9")
	flag.Parse()
  
  cfg := ConfigInfo{Thumbnail: ThumbnailInfo{Sizes: []ResizeInfo{{Width:100, Height: 0}, {Width:100, Height: 100}}}}
  
  prop := map[string]string(nil)
  res, ok := makeThumbnail(&cfg, &prop)
  assert.Equal(t, false, ok)
  
  res_need := map[string]string{}
  assert.Equal(t, res_need, res)

  prop = map[string]string{"SRC_FILE": "storage/smart-city.png",
             }
  res, ok = makeThumbnail(&cfg, &prop)
  assert.Equal(t, true, ok)

  prop = map[string]string{"FILE": "storage/smart-city.png", "THUMBNAIL_LINK": "FILE"}
  res, ok = makeThumbnail(&cfg, &prop)
  assert.Equal(t, true, ok)

}
