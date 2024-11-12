package clientManage

import (
	"ConfigServer/APIGateway"
	"ConfigServer/utils/database"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"path/filepath"
	"sort"
)

var Container *CliContainer = nil

func Init(dbPath string) {
	Container = NewCliContainer(dbPath)
	Container.Init(dbPath)
	// TODO: support identify the server ip that is under the same net range as clients
	CliUdpApiGateway = APIGateway.NewUDPAPIGateway(UdpHostPort, "0.0.0.0")
	CliUdpApiGateway.Init()
}

type Client struct {
	HostName   string           `json:"host_name"`
	IpAddr     net.IP           `json:"ip_addr"`
	MacAddr    net.HardwareAddr `json:"mac_addr"`
	StatusCode int              `json:"status_code"`
	OsVersion  string           `json:"os_version"`
	ProductId  string           `json:"product_id"`
	SysyncId   [32]byte         `json:"sysync_id"`

	conn   *net.Conn
	caught bool
}

func CreateNewClientInfo(
	hostName string,
	ipAddr net.IP,
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
func (c *Client) HumanFriendly() FriendlyClient {
	friendly := FriendlyClient{
		HostName:   c.HostName,
		IpAddr:     c.IpAddr.String(),
		MacAddr:    c.MacAddr.String(),
		StatusCode: c.StatusCode,
		OsVersion:  c.OsVersion,
		ProductId:  c.ProductId,
		SysyncId:   hex.EncodeToString(c.SysyncId[:]),
	}

	return friendly
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
	c.db = &DataFrame.SQLite{}
	err := c.db.Connect(dbPath)

	if err != nil {
		absDbPath, err := filepath.Abs(dbPath)
		if err != nil {
			return fmt.Errorf("failed to get absolute path of database: %v", err)
		}
		return fmt.Errorf("failed to connect to database at %s (absolute path: %s)", dbPath, absDbPath)
	}

	//err = c.loadClientsFromDB(db)
	//if err != nil {
	//	return fmt.Errorf("failed to load from database for initiation: %e", err)
	//}
	c.initiated = true
	return nil
}

func (c *CliContainer) DataFrameConn() *DataFrame.SQLite {
	return c.db
}

func DataFrameConn() *DataFrame.SQLite {
	return Container.DataFrameConn()
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
			IpAddr:     ipAddr,
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
	query := "INSERT INTO win_cli (host_name, IP_address, MAC_address, status_code, OS_version, product_ID, sysync_ID) VALUES (?, ?, ?, ?, ?, ?, ?)"

	sysyncIdStr := hex.EncodeToString(cli.SysyncId[:]) // 将字节数组转换为十六进制字符串

	_, err := c.db.Insert(query, cli.HostName, cli.IpAddr.String(), cli.MacAddr.String(), 000, cli.OsVersion, cli.ProductId, sysyncIdStr)
	if err != nil {
		return err
	}
	c.container[cli.HostName] = cli
	return nil
}

func Push(cli *Client) error {
	return Container.Push(cli)
}

func (c *CliContainer) Delete(sysyncID string) error {
	query := "DELETE FROM win_cli WHERE sysync_ID = ?"
	_, err := c.db.Query(query, sysyncID)
	if err != nil {
		return fmt.Errorf("failed to delete rows: %v", err)
	}
	return nil
}

func (c *CliContainer) Pop(key string) (*Client, error) {

	if existed, err := c.Exists(key); existed && err == nil {
		cli := c.container[key]
		delete(c.container, key)
		return cli, nil
	} else if err != nil {
		return nil, err
	} else {
		return nil, fmt.Errorf("host name dose not exists")
	}
}

func (c *CliContainer) Get(sysyncID string) (*Client, error) {
	if !c.initiated {
		return nil, fmt.Errorf("call the container before initiate")
	}

	query := "SELECT host_name, IP_address, MAC_address, status_code, OS_version, product_ID FROM win_cli WHERE sysync_ID = ? "
	rows, err := c.db.Query(query, sysyncID)
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	if err != nil {
		log.Fatalf("Failed to query win_cli: %v", err)
	}

	if rows.Next() {
		var hostName, ipAddrStr, macAddrStr, statCode, osVersion string
		var productId []byte
		if err := rows.Scan(&hostName, &ipAddrStr, &macAddrStr, &statCode, &osVersion, &productId); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		if rows.Next() {
			return nil, fmt.Errorf("sysync ID is duplicated")
		}

		macAddr, err := net.ParseMAC(macAddrStr)
		if err != nil {
			return nil, fmt.Errorf("invalid MAC address: %v", err)
		}
		return &Client{
			HostName:   hostName,
			IpAddr:     net.ParseIP(ipAddrStr),
			MacAddr:    macAddr,
			StatusCode: 0,
			OsVersion:  "",
			ProductId:  "",
			SysyncId:   [32]byte{},
			conn:       nil,
			caught:     false,
		}, nil
	} else {
		return nil, fmt.Errorf("sysync ID: %s dose not exists", sysyncID)
	}
}

func Get(sysyncID string) (*Client, error) {
	return Container.Get(sysyncID)
}

func (c *CliContainer) MacExists(mac string) (bool, error) {
	if !c.initiated {
		return false, fmt.Errorf("call the container before initiate")
	}

	query := "SELECT host_name, IP_address, sysync_ID FROM win_cli WHERE MAC_address = ? "
	rows, err := c.db.Query(query, mac)
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	if err != nil {
		log.Fatalf("Failed to query win_cli: %v", err)
	}

	if rows.Next() {
		return true, nil
	}
	return false, nil
}

func MacExists(mac string) (bool, error) {
	return Container.MacExists(mac)
}

func (c *CliContainer) Exists(sysyncID string) (bool, error) {
	if !c.initiated {
		return false, fmt.Errorf("call the container before initiate")
	}

	query := "SELECT host_name, IP_address, MAC_address FROM win_cli WHERE sysync_ID = ? "
	rows, err := c.db.Query(query, sysyncID)
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	if err != nil {
		log.Fatalf("Failed to query win_cli: %v", err)
	}

	if rows.Next() {
		return true, nil
	}
	return false, nil
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
	keys := make([]string, 0, len(Container.container))
	for k := range Container.container {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
