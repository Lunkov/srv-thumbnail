package main

import (
  "flag"
  "testing"
  "github.com/stretchr/testify/assert"

  "net"
  "context"
  
  "github.com/golang/glog"
  
  "google.golang.org/grpc"
  "google.golang.org/grpc/test/bufconn"
  
  "github.com/Lunkov/grpc-bpmn"
)

/////////////////////////
// TESTS
/////////////////////////

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
    lis = bufconn.Listen(bufSize)
    s := grpc.NewServer()
    srv_bpmn.RegisterBPMNJobServer(s, &BPMNJobService{})
    go func() {
        if err := s.Serve(lis); err != nil {
          glog.Errorf("ERR: Server exited with error: %v", err)
        }
    }()
}

func bufDialer(context.Context, string) (net.Conn, error) {
    return lis.Dial()
}

/////
func TestGRPCMakeThumbnail(t *testing.T) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", ".")
	flag.Set("v", "9")
	flag.Parse()
  
  globConf = loadConfig("./etc4test/config.yaml")
  
  ctx := context.Background()
  conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
  if err != nil {
      t.Fatalf("Failed to dial bufnet: %v", err)
  }
  defer conn.Close()
  client := srv_bpmn.NewBPMNJobClient(conn)
  
  
  assert.Equal(t, ConfigInfo(ConfigInfo{ConfigPath:"etc4test", Thumbnail:ThumbnailInfo{StoragePath:"storage/", Sizes:[]ResizeInfo{ResizeInfo{Width:0x32, Height:0x0}, ResizeInfo{Width:0x0, Height:0x32}, ResizeInfo{Width:0x96, Height:0x96}}}, BPMN:BPMNInfo{ConnectStr:"localhost"}}), globConf)
  
  prop := map[string]string{"SRC_FILE": "storage/oblojka1-1024x576.jpg",
             }
  
  resp, err := client.CallFunction(ctx, &srv_bpmn.RPCBPMNJob{Parameters: prop})
  if err != nil {
    t.Fatalf("SayHello failed: %v", err)
  }
  assert.Equal(t, true, resp.Ok)
  res_need := map[string]string(nil)
  assert.Equal(t, res_need, resp.Results)
}
