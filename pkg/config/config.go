package config

import (
	"flag"
	"fmt"
	"github.com/jinzhu/gorm"
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

	testdbUser string
	testdbPswd string
	testdbHost string
	testdbPort string
	testdbName string

	apiPort string
	migrate string

	redisHost string
	redisPort string
}

func Get() *Config {
	conf := &Config{}

	flag.StringVar(&conf.appEnv, "appenv", os.Getenv("APP_ENV"), "Application Environment")

	flag.StringVar(&conf.dbUser, "dbuser", os.Getenv("DB_USER"), "DB user name")
	flag.StringVar(&conf.dbPswd, "dbpswd", os.Getenv("DB_PASSWORD"), "DB pass")
	flag.StringVar(&conf.dbPort, "dbport", os.Getenv("DB_PORT"), "DB port")
	flag.StringVar(&conf.dbHost, "dbhost", os.Getenv("DB_HOST"), "DB host")
	flag.StringVar(&conf.dbName, "dbname", os.Getenv("DB_NAME"), "DB name")

	flag.StringVar(&conf.testdbUser, "testdbuser", os.Getenv("TESTDB_USER"), "DB user name")
	flag.StringVar(&conf.testdbPswd, "testdbpswd", os.Getenv("TESTDB_PASSWORD"), "DB pass")
	flag.StringVar(&conf.testdbHost, "testdbhost", os.Getenv("TESTDB_PORT"), "DB port")
	flag.StringVar(&conf.testdbPort, "testdbport", os.Getenv("TESTDB_HOST"), "DB host")
	flag.StringVar(&conf.testdbName, "testdbname", os.Getenv("TESTDB_NAME"), "DB name")

	flag.StringVar(&conf.apiPort, "apiPort", os.Getenv("API_PORT"), "API Port")
	flag.StringVar(&conf.migrate, "migrate", "up", "specify if we should be migrating DB 'up' or 'down'")

	flag.StringVar(&conf.redisHost, "redisHost", os.Getenv("REDIS_HOST"), "Redis Host")
	flag.StringVar(&conf.redisPort, "redisPort", os.Getenv("REDIS_PORT"), "Redis Port")

	flag.Parse()

	return conf
}

func (c *Config) GetAppEnv() string {
	return c.appEnv
}

func (c *Config) GetDBConnStr() string {
	return c.getDBConnStr(c.dbHost, c.dbName)
}

func (c *Config) GetTestDBConnStr() string {
	return c.getDBConnStr(c.testdbHost, c.testdbName)
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
	fmt.Printf("Connecting to database @ %s:%s\n", c.dbHost, c.dbPort)
	db, err := gorm.Open("mysql", c.GetDBConnStr())
	if err != nil {
		panic(err)
	}

	return db
}
