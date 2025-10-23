package virtualhere

import (
	"encoding/xml"
	"errors"
	"time"
)

// CommandResult represents the result of a VirtualHere command execution
type CommandResult struct {
	Success bool   `json:"success"`
	Output  string `json:"output"`
	Error   error  `json:"error,omitempty"`
}

// Device represents a USB device connected to a VirtualHere hub
type Device struct {
	Address  string `json:"address"`  // e.g., "raspberrypi.114"
	Name     string `json:"name"`     // e.g., "Ultra USB 3.0"
	AutoUse  bool   `json:"auto_use"` // Whether auto-use is enabled for this device
	InUse    bool   `json:"in_use"`   // Whether the device is currently in use
	Nickname string `json:"nickname"` // Custom nickname if set
}

// Hub represents a VirtualHere USB server/hub
type Hub struct {
	Name    string   `json:"name"`    // e.g., "Raspberry Hub"
	Address string   `json:"address"` // e.g., "raspberrypi:7575"
	Devices []Device `json:"devices"` // Devices connected to this hub
}

// ClientState represents the overall state of the VirtualHere client
type ClientState struct {
	Hubs              []Hub `json:"hubs"`
	AutoFindEnabled   bool  `json:"auto_find_enabled"`
	AutoUseAllEnabled bool  `json:"auto_use_all_enabled"`
	ReverseLookup     bool  `json:"reverse_lookup"`
	RunningAsService  bool  `json:"running_as_service"`
}

// DeviceInfo represents detailed information about a device
type DeviceInfo struct {
	Address   string `json:"address"`
	Vendor    string `json:"vendor"`
	VendorID  string `json:"vendor_id"`
	Product   string `json:"product"`
	ProductID string `json:"product_id"`
	Serial    string `json:"serial"`
	InUseBy   string `json:"in_use_by"` // "NO ONE" if not in use, otherwise client hostname
}

// ServerInfo represents detailed information about a server/hub
type ServerInfo struct {
	Name           string `json:"name"`
	Version        string `json:"version"`
	State          string `json:"state"`
	Address        string `json:"address"`
	Port           string `json:"port"`
	ConnectedFor   string `json:"connected_for"` // e.g., "9265 sec"
	MaxDevices     string `json:"max_devices"`
	ConnectionID   string `json:"connection_id"`
	Interface      string `json:"interface"`
	SerialNumber   string `json:"serial_number"`
	EasyFind       string `json:"easy_find"`
}

// ReverseClient represents a reverse client connection
type ReverseClient struct {
	ServerSerial  string `json:"server_serial"`
	ClientAddress string `json:"client_address"`
}

// License represents a VirtualHere license
type License struct {
	Key    string `json:"key"`
	Status string `json:"status"`
}

// XMLClientState represents the root XML structure from GET CLIENT STATE
type XMLClientState struct {
	XMLName xml.Name    `xml:"state" json:"-"`
	Servers []XMLServer `xml:"server" json:"servers"`
}

// XMLServer represents a server connection in the XML state
type XMLServer struct {
	XMLName    xml.Name    `xml:"server" json:"-"`
	Connection XMLServerConnection `xml:"connection" json:"connection"`
	Devices    []XMLDevice `xml:"device" json:"devices"`
}

// XMLServerConnection represents server connection details
type XMLServerConnection struct {
	ConnectionID        int       `xml:"connectionId,attr" json:"connection_id"`
	Secure              bool      `xml:"secure,attr" json:"secure"`
	ServerMajor         int       `xml:"serverMajor,attr" json:"server_major"`
	ServerMinor         int       `xml:"serverMinor,attr" json:"server_minor"`
	ServerRevision      int       `xml:"serverRevision,attr" json:"server_revision"`
	RemoteAdmin         bool      `xml:"remoteAdmin,attr" json:"remote_admin"`
	ServerName          string    `xml:"serverName,attr" json:"server_name"`
	InterfaceName       string    `xml:"interfaceName,attr" json:"interface_name"`
	Hostname            string    `xml:"hostname,attr" json:"hostname"`
	ServerSerial        string    `xml:"serverSerial,attr" json:"server_serial"`
	LicenseMaxDevices   int       `xml:"license_max_devices,attr" json:"license_max_devices"`
	State               int       `xml:"state,attr" json:"state"`
	ConnectedTime       time.Time `xml:"connectedTime,attr" json:"connected_time"`
	Host                string    `xml:"host,attr" json:"host"`
	Port                int       `xml:"port,attr" json:"port"`
	Error               bool      `xml:"error,attr" json:"error"`
	UUID                string    `xml:"uuid,attr" json:"uuid"`
	TransportID         string    `xml:"transportId,attr" json:"transport_id"`
	EasyFindEnabled     bool      `xml:"easyFindEnabled,attr" json:"easy_find_enabled"`
	EasyFindAvailable   bool      `xml:"easyFindAvailable,attr" json:"easy_find_available"`
	EasyFindID          string    `xml:"easyFindId,attr" json:"easy_find_id"`
	EasyFindPin         string    `xml:"easyFindPin,attr" json:"easy_find_pin"`
	EasyFindAuthorized  int       `xml:"easyFindAuthorized,attr" json:"easy_find_authorized"`
	IP                  string    `xml:"ip,attr" json:"ip"`
}

// XMLDevice represents a device in the XML state
type XMLDevice struct {
	Vendor                         string `xml:"vendor,attr" json:"vendor"`
	Product                        string `xml:"product,attr" json:"product"`
	IDVendor                       int    `xml:"idVendor,attr" json:"id_vendor"`
	IDProduct                      int    `xml:"idProduct,attr" json:"id_product"`
	Address                        int    `xml:"address,attr" json:"address"`
	ConnectionID                   int    `xml:"connectionId,attr" json:"connection_id"`
	State                          int    `xml:"state,attr" json:"state"`
	ServerSerial                   string `xml:"serverSerial,attr" json:"server_serial"`
	ServerName                     string `xml:"serverName,attr" json:"server_name"`
	ServerInterfaceName            string `xml:"serverInterfaceName,attr" json:"server_interface_name"`
	DeviceSerial                   string `xml:"deviceSerial,attr" json:"device_serial"`
	ConnectionUUID                 string `xml:"connectionUUID,attr" json:"connection_uuid"`
	BoundConnectionUUID            string `xml:"boundConnectionUUID,attr" json:"bound_connection_uuid"`
	BoundConnectionIP              string `xml:"boundConnectionIp,attr" json:"bound_connection_ip"`
	BoundConnectionIP6             string `xml:"boundConnectionIp6,attr" json:"bound_connection_ip6"`
	BoundClientHostname            string `xml:"boundClientHostname,attr" json:"bound_client_hostname"`
	Nickname                       string `xml:"nickname,attr" json:"nickname"`
	ClientID                       string `xml:"clientId,attr" json:"client_id"`
	NumConfigurations              int    `xml:"numConfigurations,attr" json:"num_configurations"`
	NumInterfacesInFirstConfiguration int `xml:"numInterfacesInFirstConfiguration,attr" json:"num_interfaces_in_first_configuration"`
	FirstInterfaceClass            int    `xml:"firstInterfaceClass,attr" json:"first_interface_class"`
	FirstInterfaceSubClass         int    `xml:"firstInterfaceSubClass,attr" json:"first_interface_sub_class"`
	FirstInterfaceProtocol         int    `xml:"firstInterfaceProtocol,attr" json:"first_interface_protocol"`
	HideClientInfo                 bool   `xml:"hideClientInfo,attr" json:"hide_client_info"`
	BadSerial                      bool   `xml:"badSerial,attr" json:"bad_serial"`
	ParentHubPort                  int    `xml:"parentHubPort,attr" json:"parent_hub_port"`
	ParentHubAddress               int    `xml:"parentHubAddress,attr" json:"parent_hub_address"`
	ParentHubContainerID           string `xml:"parentHubContainerID,attr" json:"parent_hub_container_id"`
	ParentHubContainerIDPrefix     int    `xml:"parentHubContainerIDPrefix,attr" json:"parent_hub_container_id_prefix"`
	ContainerID                    string `xml:"containerID,attr" json:"container_id"`
	ContainerIDPrefix              int    `xml:"containerIDPrefix,attr" json:"container_id_prefix"`
	NumPorts                       int    `xml:"numPorts,attr" json:"num_ports"`
	AutoUse                        string `xml:"autoUse,attr" json:"auto_use"` // "not-set", "on", "off", etc.
}

// Common errors
var (
	ErrCommandFailed    = errors.New("command failed")
	ErrCommandTimeout   = errors.New("command timeout (>5 seconds)")
	ErrInvalidAddress   = errors.New("invalid address")
	ErrServerNotFound   = errors.New("server not found")
	ErrDeviceNotFound   = errors.New("device not found")
	ErrDeviceInUse      = errors.New("device already in use")
	ErrBinaryNotFound   = errors.New("virtualhere binary not found")
	ErrInvalidResponse  = errors.New("invalid response from client")
)
