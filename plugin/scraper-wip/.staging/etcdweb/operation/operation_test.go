package operation

import "testing"
import "fmt"

var etcdClient = EtcdClient{
	Endpoint: "http://192.168.14.166:32379",
}

func TestGetValue(t *testing.T) {
	etcdClient.InitClient("http://192.168.14.166:32379")
	_, err := etcdClient.GetKeyValue("ufleet")
	if err != nil {
		t.Error("false")
	}
	// fmt.Println(value)
}

func TestGetDirKeys(t *testing.T) {
	etcdClient.InitClient("http://192.168.14.166:32379")
	keys, err := etcdClient.GetDirKeys("/ufleet/user/detail")
	if err != nil {
		t.Error("false")
	}
	fmt.Println(keys)
}
