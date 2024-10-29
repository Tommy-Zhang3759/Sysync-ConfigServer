package clientManage

import (
	"ConfigServer/APIGateway"
	DataFrame "ConfigServer/utils/dataFrame"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net"
	"sort"
	"time"
)

var container *CliContainer = nil

func Init(dbPath string) {
	container = NewCliContainer(dbPath)
	container.Init(dbPath)

	CliUdpApiGateway = APIGateway.NewUDPAPIGateway(UdpClientPort)
	CliUdpApiGateway.Init()
}

type Client struct {
	HostName   string           `json:"host_name"`
	IpAddr     net.Addr         `json:"ip_addr"`
	MacAddr    net.HardwareAddr `json:"mac_addr"`
	StatusCode int              `json:"status_code"`
	OsVersion  string           `json:"os_version"`
	ProductId  string           `json:"product_id"`
	SysyncId   [32]byte         `json:"sysync_id"`

	conn   *net.Conn
	caught bool
}

func NewClient(
	hostName string,
	ipAddr net.Addr,
	macAddr net.HardwareAddr,
	statusCode int,
	osVersion string,
	productId string,
) *Client {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand := make([]byte, 40)

	// 将字节映射到 letters 字符串中
	for i := range rand {
		rand[i] = letters[rand[i]%byte(len(letters))]
	}

	sysyncId := sha256.Sum256(append([]byte(macAddr.String()+productId), rand...))

	return &Client{
		HostName:   hostName,
		IpAddr:     ipAddr,
		MacAddr:    macAddr,
		StatusCode: statusCode,
		OsVersion:  osVersion,
		ProductId:  productId,
		SysyncId:   sysyncId,
	}
}

func (c *Client) logIn() {
	c.updateStatusCode(500)
}

func (c *Client) logOut() {
	c.updateStatusCode(200)
}

func (c *Client) updateStatusCode(a int) {
	c.StatusCode = a
}

type FriendlyClient struct {
	HostName   string `json:"host_name"`
	IpAddr     string `json:"ip_addr"`
	MacAddr    string `json:"mac_addr"`
	StatusCode int    `json:"status_code"`
	OsVersion  string `json:"os_version"`
	ProductId  string `json:"product_id"`
	SysyncId   string `json:"sysync_id"`
}

// HumanFriendly converts a Client struct to a human-friendly JSON format.
func (c *Client) HumanFriendly() (FriendlyClient, error) {
	friendly := FriendlyClient{
		HostName:   c.HostName,
		IpAddr:     c.IpAddr.String(),
		MacAddr:    c.MacAddr.String(),
		StatusCode: c.StatusCode,
		OsVersion:  c.OsVersion,
		ProductId:  c.ProductId,
		SysyncId:   hex.EncodeToString(c.SysyncId[:]),
	}

	return friendly, nil
}

type CliContainer struct {
	container map[string]*Client //id as key
	dbPath    string
	db        *DataFrame.SQLite

	initiated bool
}

func NewCliContainer(dbPath string) *CliContainer {
	return &CliContainer{
		container: make(map[string]*Client),
		initiated: false,
		dbPath:    dbPath,
	}
}

func (c *CliContainer) Init(dbPath string) error {
	db := &DataFrame.SQLite{}
	err := db.Connect(dbPath)
	if err != nil {
		return fmt.Errorf("failed to connect to database")
	}
	//defer func(db *DataFrame.SQLite) {
	//	_ = db.Close()
	//}(db)

	err = c.loadClientsFromDB(db)
	if err != nil {
		return fmt.Errorf("failed to load from database for initiation: %e", err)
	}
	c.initiated = true
	return nil
}

func (c *CliContainer) loadClientsFromDB(db *DataFrame.SQLite) error {
	query := `SELECT host_name, IP_address, MAC_address, status_code, OS_version, product_ID, sysync_ID FROM win_cli`
	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query database: %v", err)
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var hostName, ipAddrStr, macAddrStr, osVersion, productId string
		var statusCode int
		var sysyncId []byte

		if err = rows.Scan(&hostName, &ipAddrStr, &macAddrStr, &statusCode, &osVersion, &productId, &sysyncId); err != nil {
			return fmt.Errorf("failed to scan row: %v", err)
		}

		// 解析 IP 和 MAC 地址
		ipAddr := net.ParseIP(ipAddrStr)
		macAddr, err := net.ParseMAC(macAddrStr)
		if err != nil {
			return fmt.Errorf("invalid MAC address: %v", err)
		}

		// 创建 Client 实例并添加到容器
		client := &Client{
			HostName:   hostName,
			IpAddr:     &net.IPAddr{IP: ipAddr},
			MacAddr:    macAddr,
			StatusCode: statusCode,
			OsVersion:  osVersion,
			ProductId:  productId,
			SysyncId:   [32]byte(sysyncId),
		}

		c.Push(client)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error reading rows: %v", err)
	}

	return nil
}

func (c *CliContainer) Push(cli *Client) error {
	c.container[cli.HostName] = cli
	return nil
}

func (c *CliContainer) Delete(key string) error {
	if c.Exists(key) {
		c.container[key] = nil
		delete(c.container, key)
		return nil
	} else {
		return fmt.Errorf("key not found")
	}
}

func (c *CliContainer) Pop(key string) (*Client, error) {
	if c.Exists(key) {
		cli := c.container[key]
		delete(c.container, key)
		return cli, nil
	} else {
		return nil, fmt.Errorf("host name dose not exists")
	}
}

func (c *CliContainer) Get(key string) (*Client, error) {
	if c.Exists(key) {
		cli := c.container[key]
		return cli, nil
	} else {
		return nil, fmt.Errorf("host name dose not exists")
	}
}

func (c *CliContainer) Exists(key string) bool {
	if _, ok := c.container[key]; ok {
		return true
	} else {
		return false
	}
}

func (c *CliContainer) AllHostName() []string {
	keys := make([]string, 0, len(c.container))
	for k := range c.container {
		keys = append(keys, k)
	}
	return keys
}

func DiscoverClient(container *CliContainer, port int) {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.IPv4(0, 0, 0, 0),
	}

	// 监听 UDP 广播
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer func(conn *net.UDPConn) {
		_ = conn.Close()
	}(conn)

	fmt.Println("Listening for UDP broadcast on port", addr.Port)

	buf := make([]byte, 1024)

	for {
		n, srcAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}
		fmt.Printf("Received UDP broadcast from %s: %s\n", srcAddr, string(buf[:n]))
	}
}

func AllHostName() []string {
	keys := make([]string, 0, len(container.container))
	for k := range container.container {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func DetailedInfo(keys []string) ([]Client, time.Time, error) {
	cliList := make([]Client, len(keys), len(keys))

	t := time.Now()

	for i, name := range keys {
		if container.Exists(name) {
			cliList[i] = *container.container[name]
		}
	}

	return cliList, t, nil
}
