package main

import (
	"LLM-Chat/config"
	"LLM-Chat/routers"
)

func main() {
	config.InitConfig()

	r := routers.InitRouter()
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
