package tcp

import (
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
)

// commonServices 常见端口到服务名称的映射
// 基于IANA标准端口分配和常见服务端口
// 作为包级变量避免每次调用重复分配内存
var commonServices = map[int]string{
	20:    "ftp-data",
	21:    "ftp",
	22:    "ssh",
	23:    "telnet",
	25:    "smtp",
	53:    "dns",
	67:    "dhcp",
	68:    "dhcp",
	80:    "http",
	110:   "pop3",
	111:   "rpcbind",
	113:   "ident",
	119:   "nntp",
	123:   "ntp",
	135:   "msrpc",
	139:   "netbios-ssn",
	143:   "imap",
	161:   "snmp",
	162:   "snmptrap",
	179:   "bgp",
	194:   "irc",
	389:   "ldap",
	443:   "https",
	445:   "smb",
	465:   "smtps",
	514:   "syslog",
	515:   "printer",
	587:   "submission",
	631:   "ipp",
	636:   "ldaps",
	873:   "rsync",
	989:   "ftps-data",
	990:   "ftps",
	993:   "imaps",
	995:   "pop3s",
	1080:  "socks",
	1194:  "openvpn",
	1433:  "mssql",
	1434:  "ms-sql-m",
	1521:  "oracle",
	1701:  "l2tp",
	1723:  "pptp",
	1883:  "mqtt",
	2049:  "nfs",
	2082:  "cpanel",
	2083:  "cpanel-ssl",
	2086:  "whm",
	2087:  "whm-ssl",
	2222:  "directadmin",
	2375:  "docker",
	2376:  "docker-ssl",
	3000:  "grafana",
	3306:  "mysql",
	3389:  "rdp",
	3690:  "svn",
	4333:  "mimer",
	4444:  "metasploit",
	4500:  "ipsec-nat-t",
	5000:  "upnp",
	5001:  "iperf",
	5060:  "sip",
	5061:  "sips",
	5432:  "postgresql",
	5900:  "vnc",
	5984:  "couchdb",
	5985:  "winrm",
	5986:  "winrm-ssl",
	6379:  "redis",
	6443:  "kubernetes",
	7001:  "weblogic",
	8000:  "http-alt",
	8008:  "http",
	8080:  "http-proxy",
	8086:  "influxdb",
	8443:  "https-alt",
	8888:  "http-alt",
	9000:  "sonarqube",
	9090:  "prometheus",
	9092:  "kafka",
	9200:  "elasticsearch",
	9300:  "elasticsearch-transport",
	9418:  "git",
	9999:  "abyss",
	10000: "webmin",
	11211: "memcached",
	27017: "mongodb",
	27018: "mongodb-shard",
	27019: "mongodb-config",
	28017: "mongodb-web",
	50000: "sap",
	50070: "hadoop-namenode",
}

// guessService 根据端口号猜测服务名称
// 基于IANA标准端口分配和常见服务端口
//
// 参数:
//   - port: 端口号
//
// 返回值:
//   - string: 服务名称,如"http"、"ssh"等;未知端口返回"unknown"
//
// 说明:
//   - 包含常见服务的知名端口映射
//   - 仅供参考,实际服务可能不同
//   - 涵盖Web、数据库、消息队列等常见服务
func guessService(port int) string {
	if service, ok := commonServices[port]; ok {
		return service
	}
	return "unknown"
}

// resolveHost 解析主机地址
// 支持IP地址和域名,优先返回IPv4地址
//
// 参数:
//   - host: 主机名或IP地址
//
// 返回值:
//   - net.IP: 解析后的IP地址
//   - error: 解析错误,如DNS解析失败等
//
// 解析逻辑:
//   - 如果是有效IP地址,直接返回
//   - 否则进行DNS解析
//   - 优先返回IPv4地址
//   - 如无IPv4,返回第一个可用地址
func resolveHost(host string) (net.IP, error) {
	// 尝试解析为 IP
	ip := net.ParseIP(host)
	if ip != nil {
		return ip, nil
	}

	// DNS 解析
	addrs, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}

	// 优先使用 IPv4
	for _, addr := range addrs {
		if addr.To4() != nil {
			return addr, nil
		}
	}

	// 如果没有 IPv4，使用第一个
	if len(addrs) > 0 {
		return addrs[0], nil
	}

	return nil, fmt.Errorf("no IP address found for host: %s", host)
}

// formatErrorState 格式化错误状态为简洁的字符串
//
// 参数:
//   - err: 连接错误
//
// 返回值:
//   - string: 状态字符串 (timeout/refused/filtered/error)
func formatErrorState(err error) string {
	errStr := err.Error()
	if strings.Contains(errStr, "timeout") || strings.Contains(errStr, "timed out") {
		return "timeout"
	}
	if strings.Contains(errStr, "refused") || strings.Contains(errStr, "connection refused") {
		return "refused"
	}
	if strings.Contains(errStr, "no such host") || strings.Contains(errStr, "i/o timeout") {
		return "timeout"
	}
	if strings.Contains(errStr, "network is unreachable") || strings.Contains(errStr, "no route to host") {
		return "unreachable"
	}
	return "error"
}

// parsePortRange 解析端口范围字符串
// 支持多种格式,自动去重和排序
//
// 参数:
//   - rangeStr: 端口范围字符串
//
// 返回值:
//   - []int: 解析后的端口号列表(已排序)
//   - error: 解析错误,如格式无效、端口越界等
//
// 支持的格式:
//   - 单个端口: "80"
//   - 范围: "80-100"
//   - 多个端口: "80,443,8080"
//   - 混合: "1-100,443,8080-8090"
//
// 注意事项:
//   - 端口范围: 1-65535
//   - 自动去重
//   - 结果按升序排列
func parsePortRange(rangeStr string) ([]int, error) {
	if rangeStr == "" {
		return nil, fmt.Errorf("empty port range")
	}

	ports := make(map[int]bool)
	parts := strings.Split(rangeStr, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if strings.Contains(part, "-") {
			// 范围格式: 80-100
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid port range: %s", part)
			}

			start, err := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
			if err != nil {
				return nil, fmt.Errorf("invalid port number: %s", rangeParts[0])
			}

			end, err := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
			if err != nil {
				return nil, fmt.Errorf("invalid port number: %s", rangeParts[1])
			}

			if start < 1 || end > 65535 || start > end {
				return nil, fmt.Errorf("invalid port range: %d-%d", start, end)
			}

			for i := start; i <= end; i++ {
				ports[i] = true
			}
		} else {
			// 单个端口
			port, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid port number: %s", part)
			}
			if port < 1 || port > 65535 {
				return nil, fmt.Errorf("port out of range: %d", port)
			}
			ports[port] = true
		}
	}

	// 转换为切片
	result := make([]int, 0, len(ports))
	for port := range ports {
		result = append(result, port)
	}

	// 排序
	sort.Ints(result)

	return result, nil
}
