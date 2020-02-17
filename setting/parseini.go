package setting

import (
	"flag"
	"github.com/go-ini/ini"
	"log"
	"time"
)

type App struct {
	AllowIps           string
	AwsCredential      string
	AwsRegion          string
	AwsAccessKeyId     string
	AwsSecretAccessKey string
	TokenSalt          string
	ScriptsPath        string
	LogsPath           string
	AwsSesSenderEmail  string
}

var AppSetting = &App{}

type Server struct {
	RunMode      string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

type Mailgroups struct {
	Manager 	string
	Admin   	string
	Devops  	string
}

var MailgroupsSetting = &Mailgroups{}

var cfg *ini.File

func Setup() {
	var err error

	file := confFile()

	cfg, err = ini.Load(file)
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
	}

	mapTo("app", AppSetting)
	mapTo("mailgroups", MailgroupsSetting)
	mapTo("server", ServerSetting)

	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.ReadTimeout * time.Second

}

func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo err: %v", err)
	}
}

func confFile() string {

	confPtr := flag.String("conf", "", "set private app.ini file")
	flag.Parse()

	if *confPtr != "" {
		return *confPtr
	} else {
		return "conf/app.ini"
	}
}
