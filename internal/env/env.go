package env

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	APP_MODE         = "PROD"
	LOG_LEVEL        = "INFO"
	PORT             = "3000"
	DEFAULT_PROJECT  = "auto"
	DISCORD_TOKEN    = ""
	DISCORD_APP_ID   = ""
	DISCORD_GUILD_ID = ""
)

func PopulateEnvironment() bool {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		} else {
			fmt.Printf("Error reading configuration file: %s\n", err)
			return false
		}
	}

	if viper.IsSet("APP_MODE") {
		APP_MODE = viper.GetString("APP_MODE")
		log.Printf("[ENV] Application Mode: %s", APP_MODE)
	}

	if APP_MODE == "PROD" {
		gin.SetMode(gin.ReleaseMode)
	}

	if viper.IsSet("LOG_LEVEL") {
		LOG_LEVEL = viper.GetString("LOG_LEVEL")
		log.Printf("[ENV] Log Level: %s", LOG_LEVEL)
	}

	if viper.IsSet("DEFAULT_PROJECT") {
		DEFAULT_PROJECT = viper.GetString("DEFAULT_PROJECT")
		log.Printf("[ENV] Default Project: %s", DEFAULT_PROJECT)
	}

	if viper.IsSet("PORT") {
		PORT = viper.GetString("PORT")
		log.Printf("[ENV] Listen Port: %s", PORT)
	}

	if viper.IsSet("DISCORD_TOKEN") {
		DISCORD_TOKEN = viper.GetString("DISCORD_TOKEN")
		log.Print("[ENV] Discord token configured")
	} else {
		log.Print("[ENV] no Discord token configured")
		return false
	}

	if viper.IsSet("DISCORD_APP_ID") {
		DISCORD_APP_ID = viper.GetString("DISCORD_APP_ID")
		log.Print("[ENV] Discord App ID configured")
	} else {
		log.Print("[ENV] no Discord App ID configured")
		return false
	}

	if viper.IsSet("DISCORD_GUILD_ID") {
		DISCORD_GUILD_ID = viper.GetString("DISCORD_GUILD_ID")
		log.Print("[ENV] Discord Guild configured")
	} else {
		log.Print("[ENV] no Discord Guild configured")
		return false
	}

	return true
}
