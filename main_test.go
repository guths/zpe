package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/guths/zpe/config"
	"github.com/guths/zpe/router"
)

var r *gin.Engine

func TestMain(m *testing.M) {
	r = router.InitializeRouter()

	config.InitializeAppConfig()

	fmt.Printf("%+v", config.AppConfig)

	code := m.Run()

	os.Exit(code)
}
