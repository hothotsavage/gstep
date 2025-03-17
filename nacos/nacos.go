package nacos

import (
	"fmt"
	"github.com/hothotsavage/gstep/config"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"net"
	"strings"
)

var NamingClient naming_client.INamingClient
var ConfigClient config_client.IConfigClient

// 初始化nacos
func Setup() {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(config.Config.Nacos.Host, config.Config.Nacos.Port, constant.WithContextPath("/nacos")),
	}

	cc := *constant.NewClientConfig(
		constant.WithNamespaceId(config.Config.Nacos.Namespace),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		//constant.WithLogDir("/tmp/nacos/log"),
		//constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithUsername("nacos"),
		constant.WithPassword("nacos"),
	)

	client, _ := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	configClient, _ := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	ConfigClient = configClient
	NamingClient = client

	// 注册服务
	registerServiceInstance(client, vo.RegisterInstanceParam{
		Ip:          config.Config.Nacos.ServiceIP,
		Port:        config.Config.Port,
		ServiceName: config.Config.Nacos.ServiceName,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	})
}

// 注册服务
func registerServiceInstance(nacosClient naming_client.INamingClient, param vo.RegisterInstanceParam) {
	success, err := nacosClient.RegisterInstance(param)
	if !success || err != nil {
		panic("register Service Instance failed!")
	}
}

// 获取本机ip地址
func getHostIp() string {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		fmt.Println("get current host ip err: ", err)
		return ""
	}
	addr := conn.LocalAddr().(*net.UDPAddr)
	ip := strings.Split(addr.String(), ":")[0]
	return ip
}
