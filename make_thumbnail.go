package main

import (
  "github.com/golang/glog"
)

func makeThumbnail(cfg *ConfigInfo, prop *map[string]string) (map[string]string, bool) {
  if glog.V(2) {
    glog.Infof("LOG: THUMBNAIL: run prepare...\n")
  }
  result := make(map[string]string)
  var ok bool
  var link string
  var srcFile string
  
  link, ok = (*prop)["THUMBNAIL_LINK"]
  if !ok {
    srcFile, ok = (*prop)["SRC_FILE"]
    if !ok {
      glog.Errorf("ERR: THUMBNAIL: 'SRC_FILE' don`t set")
      return result, false
    }
  } else {
    srcFile, ok = (*prop)[link]
    if !ok {
      glog.Errorf("ERR: THUMBNAIL: '%s' don`t set", link)
      return result, false
    }
  }
   
  var img SrcImage
  ok = img.load(srcFile)
  if ok {
    for _, v := range cfg.Thumbnail.Sizes {
      ok = img.resizeImage(v.Width, v.Height)
      glog.Errorf("ERR: THUMBNAIL: resizeImage(%d, %d)", v.Width, v.Height)
    }
  }
  
  if glog.V(2) {
    glog.Infof("LOG: THUMBNAIL(%s): DONE", srcFile)
  }
  return result, true
}

