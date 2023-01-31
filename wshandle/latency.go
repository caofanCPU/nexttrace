package wshandle

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	result  string
	results = make(chan string)
)

func GetFastIP(domain string, port string) string {

	ips, err := net.LookupIP(domain)
	if err != nil {
		log.Fatal("DNS 解析失败，请检查您的系统 DNS 设置")
		return ""
	}

	for _, ip := range ips {
		go checkLatency(ip.String(), port)
	}

	select {
	case result = <-results:
	case <-time.After(1 * time.Second):

	}
	if result == "" {
		log.Fatal("IP 连接均超时，请检查您的网络")
	}
	res := strings.Split(result, "-")

	if len(ips) > 1 {
		fmt.Fprintf(color.Output, "%s 已为您优选最近的节点 %s - %s\n",
			color.New(color.FgWhite, color.Bold).Sprintf("[NextTrace API]"),
			color.New(color.FgGreen, color.Bold).Sprintf("%s", res[0]),
			color.New(color.FgCyan, color.Bold).Sprintf("%sms", res[1]),
		)
	}

	return res[0]
}

func checkLatency(ip string, port string) {
	start := time.Now()
	if !strings.Contains(ip, ".") {
		ip = "[" + ip + "]"
	}
	conn, err := net.DialTimeout("tcp", ip+":"+port, time.Second*1)
	if err != nil {
		return
	}
	defer conn.Close()
	if result == "" {
		result = fmt.Sprintf("%s-%.2f", ip, float64(time.Since(start))/float64(time.Millisecond))
		results <- result
		return
	}
}