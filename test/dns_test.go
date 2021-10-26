package test

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"
)

//"golang.org/x/net/dns/dnsmessage"

func TestDefalutDNS(t *testing.T) {

	ns, err := net.LookupHost("www.baidu.com")
	if err != nil {
		t.Fail()
	}
	for _, n := range ns {
		fmt.Printf("--%s\n", n)
	}

}

//https://pkg.go.dev/net#hdr-Name_Resolution
// Windows下总是会调用C库，所以无法指定
// https://github.com/golang/go/issues/22846
// 容器里需要注意新建 nsswitch.conf文件
//		hosts:	files dns
func TestSpecificDNS(t *testing.T) {
	fmt.Println("TestSpecificDNS~...")
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Second * 10,
			}
			fmt.Println("\n\n\nQuery dns for: ", address)
			return d.DialContext(ctx, "udp", "127.0.0.1:53")
		},
	}
	ns, err := r.LookupHost(context.Background(), "www.baidu.com")
	if err != nil {
		t.Fail()
	}
	for _, n := range ns {
		fmt.Printf("--%s\n", n)
	}

}
