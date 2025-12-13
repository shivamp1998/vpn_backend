package wireguard

import (
	"fmt"
	"strings"
)

func GenerateClientConfig(
	clientPrivateKey string,
	serverPublicKey string,
	serverEndpoint string,
	clientIp string,
	dns string,
) string {
	var config strings.Builder

	// [Interface] section - Client configuration
	config.WriteString("[Interface]\n")
	config.WriteString(fmt.Sprintf("PrivateKey = %s\n", clientPrivateKey))
	config.WriteString(fmt.Sprintf("Address = %s\n", clientIp))

	if dns != "" {
		config.WriteString(fmt.Sprintf("DNS = %s\n", dns))
	}

	config.WriteString("\n")

	// [Peer] section - Server Configuration

	config.WriteString("[Peer]\n")
	config.WriteString(fmt.Sprintf("PublicKey = %s\n", serverPublicKey))
	config.WriteString(fmt.Sprintf("Endpoint = %s\n", serverEndpoint))
	config.WriteString("AllowedIps = 0.0.0.0/0\n")
	config.WriteString("PersistentKeepalive = 25\n")

	return config.String()
}

type PeerConfig struct {
	PublicKey  string
	AllowedIps string
}

func GenerateServerConfig(
	serverPrivatekey string,
	serverIp string,
	port int,
	clients []PeerConfig,
) string {
	var config strings.Builder

	config.WriteString("[Interface]\n")
	config.WriteString(fmt.Sprintf("PrivateKey = %s\n", serverPrivatekey))
	config.WriteString(fmt.Sprintf("Address = %s\n", serverIp))
	config.WriteString(fmt.Sprintf("ListenPort = %d\n", port))
	config.WriteString("PostUp = iptables -A FORWARD -i wg0 -j ACCEPT; iptables -A FORWARD -o wg0 -j ACCEPT; iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE\n")
	config.WriteString("PostDown = iptables -D FORWARD -i wg0 -j ACCEPT; iptables -D FORWARD -o wg0 -j ACCEPT; iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE\n")
	config.WriteString("\n")

	for _, client := range clients {
		config.WriteString("[Peer]\n")
		config.WriteString(fmt.Sprintf("Publickey = %s\n", client.PublicKey))
		config.WriteString(fmt.Sprintf("AllowedIps = %s\n", client.AllowedIps))
		config.WriteString("\n")
	}

	return config.String()
}
