package config

import (
	"flag"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/shadowbane/golang-microservice-sekolah/pkg/logger"
	"go.uber.org/zap"
	"os"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Config struct {
	appEnv string
	dbUser string
	dbPswd string
	dbHost string
	dbPort string
	dbName string

	apiPort string
	migrate string

	redisHost string
	redisPort string

	LogConfig *logger.LogConfig
}

func Get() *Config {
	conf := &Config{}

	flag.StringVar(&conf.appEnv, "appenv", getenv("APP_ENV", "production"), "Application Environment")

	flag.StringVar(&conf.dbUser, "dbuser", getenv("DB_USERNAME", "root"), "DB user name")
	flag.StringVar(&conf.dbPswd, "dbpswd", getenv("DB_PASSWORD", "password"), "DB pass")
	flag.StringVar(&conf.dbPort, "dbport", getenv("DB_PORT", "3306"), "DB port")
	flag.StringVar(&conf.dbHost, "dbhost", getenv("DB_HOST", "localhost"), "DB host")
	flag.StringVar(&conf.dbName, "dbname", getenv("DB_DATABASE", "microservice_sekolah"), "DB name")

	flag.StringVar(&conf.apiPort, "apiPort", getenv("API_PORT", "8080"), "API Port")

	flag.StringVar(&conf.redisHost, "redisHost", getenv("REDIS_HOST", "localhost"), "Redis Host")
	flag.StringVar(&conf.redisPort, "redisPort", getenv("REDIS_PORT", "6379"), "Redis Port")

	conf.LogConfig = logger.LoadEnvForLogger()

	flag.Parse()

	return conf
}

// getenv get environment variable or fallback to default value if not set
func getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func (c *Config) GetAppEnv() string {
	return c.appEnv
}

func (c *Config) GetDBConnStr() string {
	return c.getDBConnStr(c.dbHost, c.dbName)
}

func (c *Config) GetDBConnStrForMigration() string {
	return fmt.Sprintf(
		"%s://%s",
		"mysql",
		c.GetDBConnStr(),
	)
}

func (c *Config) getDBConnStr(dbhost, dbname string) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?multiStatements=true&charset=utf8mb4&parseTime=True&loc=Local",
		c.dbUser,
		c.dbPswd,
		dbhost,
		c.dbPort,
		dbname,
	)
}

func (c *Config) GetAPIPort() string {
	return ":" + c.apiPort
}

func (c *Config) GetMigration() string {
	return c.migrate
}

func (c *Config) GetRedisHost() string {
	return c.redisHost
}

func (c *Config) GetRedisPort() string {
	return c.redisPort
}

func (c *Config) GetRedisConnStr() string {
	return c.redisHost + ":" + c.redisPort
}

func (c *Config) ConnectToDatabase() *gorm.DB {
	if c.appEnv != "production" {
		zap.S().Debugf("Connecting to database @ %s:%s\n", c.dbHost, c.dbPort)
	}

	db, err := gorm.Open("mysql", c.GetDBConnStr())
	if err != nil {
		panic(err)
	}

	return db
}
