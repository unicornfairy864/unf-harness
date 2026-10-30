package main

import (
	"github.com/unicornfairy864/unf-harness/client"
	"github.com/unicornfairy864/unf-harness/config"
)

func main() {
	config.Init(&config.UH_Config)

	client.Init(&client.UH_Client)

}
