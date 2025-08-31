package request

import (
	"encoding/binary"
	"errors"
	"github.com/patrickmn/go-cache"
	"math/rand"
	"net"
	"strings"
	"time"
)

// GetUrlIpv4 获取Url的IPv4地址
func GetUrlIpv4(domain string) string {
	if val, found := DnsCache.Get(domain); found {
		if ips, ok := val.([]string); ok {
			return ips[rand.Intn(len(ips))]
		}
	}

	if res, _ := DNSQuery(domain); res != nil && len(res) > 0 {
		return res[rand.Intn(len(res))]
	}
	if res := LoopIp(domain); res != "" {
		return res
	}
	DnsCache.Set(domain, []string{""}, cache.DefaultExpiration)
	return ""
}

// DnsServer Dns Server默认列表
var DnsServer = []string{
	"223.5.5.5",
	"223.6.6.6",
	"119.29.29.29",
	"182.254.116.116",
	"180.76.76.76",
	"114.114.114.114",
	"114.114.115.115",
	"8.8.8.8",
}
var DnsCache = cache.New(60*time.Minute, 60*time.Second)

// LoopIp nsLookup获取IP
func LoopIp(domain string) string {
	if ips, err := net.LookupIP(domain); err == nil { // A记录
		for _, ip := range ips {
			if ip4 := ip.To4(); ip4 != nil {
				DnsCache.Set(domain, []string{ip4.String()}, cache.DefaultExpiration)
				return ip4.String()
			}
		}
	}
	return ""
}

// DNSHeader 定义 DNS 包头结构
type DNSHeader struct {
	ID      uint16
	Flags   uint16
	QDCount uint16
	ANCount uint16
	NSCount uint16
	ARCount uint16
}

// DNSQuery 发送 DNS 查询包并获取响应
func DNSQuery(domain string) ([]string, error) {
	rand.Seed(time.Now().UnixNano())

	var lastErr error
	for _, dnsServer := range shuffleServers(DnsServer) {
		for i := 0; i < 3; i++ {
			ips, err := queryDNSOnce(dnsServer, domain)
			if err == nil && len(ips) > 0 {
				DnsCache.Set(domain, ips, cache.DefaultExpiration)
				return ips, nil
			}
			lastErr = err
			time.Sleep(200 * time.Millisecond)
		}
	}
	return nil, lastErr
}

// 单次请求逻辑
func queryDNSOnce(dnsServer, domain string) ([]string, error) {
	timeout := 5 * time.Second
	conn, err := net.DialTimeout("udp", net.JoinHostPort(dnsServer, "53"), timeout)
	if err != nil {
		return nil, err
	}
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	query := buildDNSQuery(domain)
	if _, err = conn.Write(query); err != nil {
		return nil, err
	}

	response := make([]byte, 512)
	if err = conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return nil, err
	}
	n, err := conn.Read(response)
	if err != nil {
		return nil, err
	}

	if n < 12 {
		return nil, errors.New("invalid DNS response length")
	}

	// 解析 Flags 和 RCode
	flags := binary.BigEndian.Uint16(response[2:4])
	rcode := flags & 0x000F
	if flags&0x0200 != 0 {
		return nil, errors.New("server truncation (TC bit set)")
	}
	if rcode != 0 {
		return nil, errors.New("DNS response error, RCode=" + string(rune(rcode)))
	}

	ips := parseDNSResponse(response[:n])
	return ips, nil
}

// buildDNSQuery 构建 DNS 查询包
func buildDNSQuery(domain string) []byte {
	buf := make([]byte, 0)
	header := DNSHeader{ID: uint16(rand.Intn(65535)), Flags: 0x0100, QDCount: 1}
	hdrBuf := make([]byte, 12)
	binary.BigEndian.PutUint16(hdrBuf[0:], header.ID)
	binary.BigEndian.PutUint16(hdrBuf[2:], header.Flags)
	binary.BigEndian.PutUint16(hdrBuf[4:], header.QDCount)
	buf = append(buf, hdrBuf...)

	qName := encodeDomainName(domain)
	buf = append(buf, qName...)

	qBuf := make([]byte, 4)
	binary.BigEndian.PutUint16(qBuf[0:], 1)
	binary.BigEndian.PutUint16(qBuf[2:], 1)
	buf = append(buf, qBuf...)
	return buf
}

// encodeDomainName 将域名转为 DNS 查询格式
func encodeDomainName(domain string) []byte {
	parts := strings.Split(domain, ".")
	buf := make([]byte, 0)
	for _, part := range parts {
		buf = append(buf, byte(len(part)))
		buf = append(buf, []byte(part)...)
	}
	buf = append(buf, 0)
	return buf
}

// parseDNSResponse 解析响应提取 IP
func parseDNSResponse(response []byte) []string {
	var ips []string
	answerCount := binary.BigEndian.Uint16(response[6:8])
	offset := 12
	for response[offset] != 0 {
		offset++
	}
	offset += 5 // QName + QType + QClass

	for i := 0; i < int(answerCount); i++ {
		offset += 2 // Name
		qType := binary.BigEndian.Uint16(response[offset : offset+2])
		offset += 8 // Type, Class, TTL
		dataLen := binary.BigEndian.Uint16(response[offset : offset+2])
		offset += 2

		if qType == 1 && dataLen == 4 {
			ip := net.IP(response[offset : offset+4])
			ips = append(ips, ip.String())
		}
		offset += int(dataLen)
	}
	return ips
}

// shuffleServers 随机打乱 DNS Server 顺序
func shuffleServers(servers []string) []string {
	shuffled := make([]string, len(servers))
	copy(shuffled, servers)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return shuffled
}
