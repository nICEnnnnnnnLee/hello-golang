package test

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptrace"
	"net/url"
	"testing"
)

func TestCookieJar(t *testing.T) {
	const (
		dstUrl         = "https://www.baidu.com/s?wd=golang"
		cookieDomain   = ".baidu.com"
		queryCookieUrl = "https://www.baidu.com/"
	)
	// 通过jar来管理cookie
	jar, _ := cookiejar.New(nil)
	// 人工添加一个cookie(假设是持久化中读取而来)
	var cookies []*http.Cookie
	firstCookie := &http.Cookie{
		Name:   "testName",
		Value:  "testValue",
		Domain: cookieDomain,
		Path:   "/",
	}
	cookies = append(cookies, firstCookie)
	cookieURL, _ := url.Parse("https://" + cookieDomain)
	jar.SetCookies(cookieURL, cookies)
	client := &http.Client{
		Jar: jar,
	}
	//提交请求
	request, err := http.NewRequest("GET", dstUrl, nil)
	if err != nil {
		panic(err)
	}

	//增加header选项
	request.Header.Add("Cookie", "xxxxxx")
	request.Header.Add("User-Agent", "xxx")
	request.Header.Add("X-Requested-With", "xxxx")
	//处理返回结果
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.Status)
	for key, value := range response.Header {
		fmt.Println(key, value)
	}

	domain, _ := url.Parse(queryCookieUrl)

	for _, cookie := range jar.Cookies(domain) {
		fmt.Println(cookie.Name, cookie.Value)
	}

	// 读取body
	// body, _ := io.ReadAll(response.Body)
	// fmt.Println(string(body))
	defer response.Body.Close()
}

func TestNormalHTTP(t *testing.T) {
	res, err := http.Get("https://baidu.com")
	if err != nil {
		panic(err)
	}
	fmt.Println(res.Status)
}

func TestClientTrace(t *testing.T) {
	trace := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) {},
		DNSDone:  func(_ httptrace.DNSDoneInfo) {},
		ConnectStart: func(net, addr string) {
			fmt.Printf("ConnectStart addr=%s\n", addr)
		},
		ConnectDone: func(net, addr string, err error) {
			fmt.Printf("ConnectDone addr=%s\n", addr)
		},
		GotConn:              func(_ httptrace.GotConnInfo) {},
		GotFirstResponseByte: func() {},
		TLSHandshakeStart:    func() {},
		TLSHandshakeDone:     func(_ tls.ConnectionState, _ error) {},
	}

	req, err := http.NewRequest(http.MethodGet, "http://www.baidu.com", nil)
	if err != nil {
		log.Fatal(err)
	}

	req = req.WithContext(httptrace.WithClientTrace(context.Background(), trace))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}

func TestSpecialHTTPConfig(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "https://baidu.com", nil)
	if err != nil {
		log.Fatal(err)
	}

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			fmt.Println(addr)
			// conn, err := net.Dial(network, "127.0.0.1:443")
			conn, err := net.Dial(network, addr)
			fmt.Println(conn.RemoteAddr())
			req.RemoteAddr = conn.RemoteAddr().String()
			return conn, err
		},
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         "baidu.com",
			// VerifyConnection: func(connState tls.ConnectionState) error {
			// 	return connState.PeerCertificates[0].VerifyHostname(connState.ServerName)
			// },
		},
		ForceAttemptHTTP2: true,
	}
	client := http.Client{
		Transport: transport,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	fmt.Println("Status:", resp.Status)
	fmt.Println("RemoteAddr:", req.RemoteAddr)
}

// var SpecifiDnsServerTransport http.RoundTripper = &http.Transport{
//     Proxy: http.ProxyFromEnvironment,
//     DialContext: (&net.Dialer{
//         Timeout:   30 * time.Second,
//         KeepAlive: 30 * time.Second,
//         Resolver: &net.Resolver{PreferGo: true, Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
//             d := net.Dialer{}
//             address = "8.8.8.8:53"
//             return d.DialContext(ctx, network, address)
//         }},
//     }).DialContext,
//     ForceAttemptHTTP2:     true,
//     MaxIdleConns:          100,
//     IdleConnTimeout:       90 * time.Second,
//     TLSHandshakeTimeout:   10 * time.Second,
//     ExpectContinueTimeout: 1 * time.Second,
// }
// addr, err := net.LookupHost("www.baidu.com")
// if err != nil {
//     panic(err)
// }
// fmt.Println(addr)
