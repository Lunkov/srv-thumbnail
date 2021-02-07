package main

import (
  "flag"
  "net"
  "io/ioutil"
  "path/filepath"
  "gopkg.in/yaml.v2"
  "github.com/golang/glog"

  "golang.org/x/net/context"
  "google.golang.org/grpc"
  "google.golang.org/grpc/reflection"

  "github.com/Lunkov/lib-env"
  "github.com/Lunkov/grpc-bpmn"
)

type ResizeInfo struct {
  Width     uint      `yaml:"width"`
  Height    uint      `yaml:"height"`
}

type ThumbnailInfo struct {
  StoragePath     string         `yaml:"storage"`
  Sizes           []ResizeInfo   `yaml:"sizes"`
}

type BPMNInfo struct {
  ConnectStr      string  `yaml:"connect"`
}

type ConfigInfo struct {
  ConfigPath      string
  Thumbnail       ThumbnailInfo    `yaml:"thumbnail"`
  BPMN            BPMNInfo         `yaml:"bpmn"`
}

var globConf ConfigInfo

type BPMNJobService struct{}

func (s *BPMNJobService) CallFunction(ctx context.Context, in *srv_bpmn.RPCBPMNJob) (*srv_bpmn.RPCBPMNJobResponse, error) {
  results, ok := makeThumbnail(&globConf, &in.Parameters)
	return &srv_bpmn.RPCBPMNJobResponse{BpmnProcessId: in.BpmnProcessId, Ok: ok, Results: results}, nil
}

func loadConfig(filename string) ConfigInfo {
  var err error
  cfg := ConfigInfo{}

  if !env.WaitFile(filename, 300) {
    return cfg
  }

  yamlFile, err := ioutil.ReadFile(filename)
  if err != nil {
    glog.Errorf("ERR: yamlFile(%s)  #%v ", filename, err)
    return cfg
  }
  err = yaml.Unmarshal(yamlFile, &cfg)
  if err != nil {
    glog.Errorf("ERR: yamlFile(%s): YAML: %v", filename, err)
  }
  if cfg.ConfigPath == "" {
    cfg.ConfigPath = filepath.Dir(filename)
  }
  return cfg
}

func main() {
  flag.Set("alsologtostderr", "true")
  flag.Set("log_dir", ".")
  // flag.Set("v", "9")
  configPath := flag.String("config_path", "./etc/", "Config path")
  flag.Parse()

  globConf = loadConfig(*configPath + "config.yaml")

  glog.Infof("LOG: Start gRPC Maker Report")

  lis, err := net.Listen("tcp", "0.0.0.0:3000")
  if err != nil {
    glog.Fatalf("Can not listen the port：%v", err)
  }

  s := grpc.NewServer()

  srv_bpmn.RegisterBPMNJobServer(s, &BPMNJobService{})

  reflection.Register(s)
  if err := s.Serve(lis); err != nil {
    glog.Fatalf("ERR: Can not provide service：%v", err)
  }
}
