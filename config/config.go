package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var (
	// Env environment
	Env  Environment
	once sync.Once
)

// InitEnvironment init env
func InitEnvironment(confPath string) error {
	var err error
	once.Do(func() {
		Env, err = loadEnvironment(confPath)
	})
	return err
}

type Environment struct {
	Core         sectionCore
	Log          sectionLog
	MySQL        sectionMySQL
	Redis        sectionRedis
	SessionAuth  sectionSessionAuth
	SectionImage sectionImage
}

type sectionCore struct {
	Mode             string
	Port             string
	SkipRateLimitKey string
}
type sectionLog struct {
	Format string
	Output string
	Level  string
}

// SectionMySQL is sub section of config.
type sectionMySQL struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type sectionRedis struct {
	Host     string
	Port     string
	Password string
}
type sectionSessionAuth struct {
	Name        string
	MaxLifeTime uint
}

type sectionImage struct {
	Size int64
}

func loadEnvironment(path string) (Environment, error) {
	var env Environment

	viper.AutomaticEnv()
	viper.SetConfigType("env")

	if path != "" {
		content, err := os.ReadFile(filepath.Clean(path))
		if err != nil {
			return Environment{}, err
		}
		if err := viper.ReadConfig(bytes.NewBuffer(content)); err != nil {
			return Environment{}, err
		}
	} else {
		// look for config in the working directory
		viper.AddConfigPath(".")
		viper.SetConfigFile(".env")

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// go service
	env.Core.Mode = viper.GetString("core_bk_mode")
	env.Core.Port = viper.GetString("core_bk_port")
	if len(env.Core.Port) == 0 {
		env.Core.Port = "8020"
	}
	env.Core.SkipRateLimitKey = viper.GetString("core_skip_rate_limit_key")

	// log
	env.Log.Format = viper.GetString("log_format")
	env.Log.Level = viper.GetString("log_level")
	env.Log.Output = viper.GetString("log_output")

	// postgres
	env.MySQL.Host = viper.GetString("mysql_host")
	env.MySQL.Port = viper.GetString("mysql_port")
	env.MySQL.User = viper.GetString("mysql_user")
	env.MySQL.Password = viper.GetString("mysql_password")
	env.MySQL.DBName = viper.GetString("mysql_db_name")

	// redis
	env.Redis.Host = viper.GetString("redis_host")
	env.Redis.Port = viper.GetString("redis_port")
	env.Redis.Password = viper.GetString("redis_password")

	// session auth
	env.SessionAuth.Name = viper.GetString("session_auth_name")
	env.SessionAuth.MaxLifeTime = viper.GetUint("session_auth_max_life_time")

	// image
	env.SectionImage.Size = viper.GetInt64("image_size")

	return env, nil
}
