package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"lab-platform/core/loadconfig"
	"lab-platform/lib/auth"
	"lab-platform/lib/config"
	"lab-platform/lib/logger"
	"lab-platform/router"
)

func Init() (*core, error) {

	// _init ConfigData properties
	if err := initViper(); err != nil {
		return nil, err
	}
	loadConfig()

	// init log package
	if err := logger.Init(); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return nil, err
	}


	//users, authConf, err := initAuth()
	//if err != nil {
	//	return nil, err
	//}

	// init gin
	engine := initGin()

	return NewCore(engine), nil
}

func initGin() *gin.Engine {
	// Set gin mode.
	gin.SetMode(loadconfig.Config.Mode)

	// Create the Gin engine.
	g := gin.New()

	// Routes.
	router.Load(
		// Cores.
		g,
		// Middleware's.
		//gin.BasicAuth(users),
		// todo: 认证模块
		//authz.NewAuthorizer(authConf),
		logger.GinLogger(),
		logger.GinRecovery(true),
	)
	return g
}

func initViper() error {
	if config.Flag.Cfg != "" {
		viper.SetConfigFile(config.Flag.Cfg)
	} else {
		viper.AddConfigPath("configs")
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

func loadConfig() {
	loadconfig.BasicConfig()
	auth.LoadConfig()
	logger.LoadConfig()
}

//func initAuth() (gin.Accounts, *casbin.Enforcer, error) {
//	users := gin.Accounts{}
//	for username, password := range auth.Config.Users {
//		users[username] = password.(string)
//	}
//
//	confPath := auth.Config.Path
//	authConf, err := casbin.NewEnforcer(confPath+"authz_model.conf", confPath+"authz_policy.csv")
//	return users, authConf, err
//}