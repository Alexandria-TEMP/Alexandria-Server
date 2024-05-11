package main

import (
	_ "github.com/swaggo/gin-swagger"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/routers"
)

func main() {
	routers.Init()
}
