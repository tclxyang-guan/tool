package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"transfDoc/conf"
	"transfDoc/endpoint"
	"transfDoc/pkg/logging"
)

func init() {
	conf.Setup()
	logging.Setup()
}

func main() {

	cfg := conf.GetConfig()
	gin.SetMode(cfg.RunMode)
	routersInit := endpoint.InitRouter()
	fmt.Println("|-----------------------------------|")
	fmt.Println("|   dhcc-go transDoc service  |")
	fmt.Println("|-----------------------------------|")
	fmt.Println("|  Go Http Server Start Successful  |")
	fmt.Println("|   Port" + cfg.PrefixUrl + "     Pid:" + fmt.Sprintf("%d", os.Getpid()) + "   |")
	fmt.Println("|-----------------------------------|")
	fmt.Println("")
	log.Fatal(routersInit.Run(cfg.PrefixUrl))
}
