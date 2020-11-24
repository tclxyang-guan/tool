/**
* @Auther:gy
* @Date:2020/11/23 16:23
 */
package conf

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

// 初始化所有项目配置
func Setup() {
	filename := "/data/service/transfDoc/configs/app.yaml"
	exist, err := isFileExist(filename)
	if err != nil {
		log.Fatal(fmt.Errorf("check config file error, err: %v, filename: %s", err, filename))
	}
	if !exist {
		filename = "conf/app.yaml"
		exist, err = isFileExist(filename)
		if err != nil {
			log.Fatal(fmt.Errorf("check config file error, err: %v, filename: %s", err, filename))
		}
	}

	confValue = &Conf{}
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(fmt.Errorf("read config file error, err: %v, filename: %s", err, filename))
	}

	err = yaml.Unmarshal(yamlFile, confValue)
	if err != nil {
		log.Fatal(fmt.Errorf("unmarshal config file error, err: %v, filename: %s", err, filename))
	}
}

var confValue *Conf

// 获取配置对象
func GetConfig() *Conf {
	return confValue
}

// 获取mysql配置项
func (c *Conf) GetMysqlConf() *MySQLConf {
	return &c.Env.MySQL
}

//Conf config of different environments
type Conf struct {
	ApiKey          string `yaml:"apikey"`
	ApiToken        string `yaml:"apitoken"`
	ShowDocOpen     bool   `yaml:"showdocopen"`
	ShowDocUrl      string `yaml:"showdocurl"`
	PrefixUrl       string `yaml:"prefixurl"`
	RunMode         string `yaml:"runmode"`
	RuntimeRootPath string `yaml:"runtimerootpath"`
	LogSavePath     string `yaml:"logsavepath"`
	LogSaveName     string `yaml:"logsavename"`
	LogFileExt      string `yaml:"logfileext"`
	TimeFormat      string `yaml:"timeformat"`
	Env             EnvConf
}

// EnvConf envirnoment config
type EnvConf struct {
	MySQL MySQLConf `yaml:"mysql"`
}

type MySQLConf struct {
	Address         string `yaml:"address"`
	DbName          string `yaml:"dbname"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	MaxOpenConns    int    `yaml:"maxconns"`
	MaxIdleConns    int    `yaml:"maxidleconns"`
	ConnMaxLifetime int64  `yaml:"connmaxlifetime"`
}

func isFileExist(fliename string) (bool, error) {
	_, err := os.Stat(fliename)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Logger :
type Logger interface {
	Infof(format string, a ...interface{})
	Warnf(format string, a ...interface{})
	Errorf(format string, a ...interface{})
}
