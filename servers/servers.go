package servers

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"

	"github.com/legenove/cocore"
	"github.com/legenove/utils"
)

const (
	SERVER_TYPE_SERVER = "server" // rest and rpc and tcp
	SERVER_TYPE_ASYNC  = "async"  //
	SERVER_TYPE_CRON   = "cron"   // 定时任务发布
)

const (
	SecretNormalType = "normal" // 普通校验对比
	SecretMD5Type    = "md5"    // md5加密
	SecretBase64Type = "base64" // base64加密
)

var Server = &ServerConf{}

type ServerConf struct {
	// for cocore Config
	Debug       bool
	AppENV      string
	ConfigDir   string
	AppConfName string
	// for db or redis
	DBDebugLog bool
	// for servers Config
	Doc           bool
	MaxProcs      int
	DocDir        string
	Title         string         `json:"server_title" mapstructure:"server_title"`
	Group         string         `json:"server_group" mapstructure:"server_group"`
	Name          string         `json:"server_name" mapstructure:"server_name"`
	Type          string         `json:"server_type"  mapstructure:"server_type"`
	Host          string         `json:"host" mapstructure:"host"`
	Secrets       []ServerSecret `json:"secrets" mapstructure:"secrets"`
	IPStrategy    string         `json:"ip_strategy" mapstructure:"ip_strategy"` // ip策略
	stringSecrets []string
}

func InitServer(secretKey, secretType string) {
	secret := ServerSecret{
		Secret: secretKey,
		Type:   secretType,
	}
	Server.Secrets = []ServerSecret{secret}

	if strings.HasPrefix(Server.DocDir, "$GOPATH") {
		Server.DocDir = filepath.Join(os.Getenv("GOPATH"), Server.DocDir[7:])
	}
	cocore.InitApp(Server.Debug, Server.AppENV, Server.ConfigDir, Server.AppConfName)
	InitServerLog()
	cocore.RegisterInitFunc("serverLog", InitServerLog)
}

func (s *ServerConf) GetServerGroup() string {
	return s.Group
}

func (s *ServerConf) GetServerName() string {
	return s.Name
}

func (s *ServerConf) GetServerTitle() string {
	return s.Title
}

func (s *ServerConf) Validator(value string) bool {
	for _, v := range s.stringSecrets {
		if v == value {
			return true
		}
	}
	return false
}

func (s *ServerConf) SetStringSecret() {
	res := make([]string, len(s.Secrets))
	for i, v := range s.Secrets {
		res[i] = v.getSecret()
	}
	s.stringSecrets = res
}

type ServerSecret struct {
	Type   string `json:"type" mapstructure:"type"`
	Secret string `json:"secret" mapstructure:"secret"`
}

func (as *ServerSecret) getSecret() string {
	switch as.Type {
	case SecretNormalType:
		return as.Secret
	case SecretBase64Type:
		return base64.RawURLEncoding.EncodeToString([]byte(as.Secret))
	case SecretMD5Type:
		return utils.GetMD5Hash(as.Secret)
	}
	return as.Secret
}
