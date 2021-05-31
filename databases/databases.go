package databases

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
)

var App *Application

type Config viper.Viper

type Application struct {
	Name        string        `json:"name"`
	Version     string        `json:"version"`
	ENV         string        `json:"env"`
	AppConfig   Config        `json:"application_config"`
	DBConfig    *gorm.DB      `json:"database_config"`
	RedisConfig *redis.Client `json:"redis_config"`
}

type Database struct {
	Driver            string
	Host              string
	User              string
	Password          string
	DBName            string
	DBNumber          int
	Port              int
	API_URL           string
	ReconnectRetry    int
	ReconnectInterval int64
	DebugMode         bool
	Pool              Pool
}

type Pool struct {
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  int
}

func init() {
	App = &Application{}
	App.Name = "APP_NAME"
	App.Version = "APP_VERSION"
	App.loadENV()
	App.loadAppConfig()
	App.loadDBConfig()
	App.loadRedisConfig()
}

// loadAppConfig: read application config and build viper object
func (app *Application) loadAppConfig() {
	var (
		appConfig *viper.Viper
		err       error
	)
	appConfig = viper.New()
	appConfig.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	appConfig.SetEnvPrefix("APP_")
	appConfig.AutomaticEnv()
	appConfig.SetConfigName("config")
	appConfig.AddConfigPath(".")
	appConfig.SetConfigType("json")
	if err = appConfig.ReadInConfig(); err != nil {
		panic(err)
	}
	appConfig.WatchConfig()
	appConfig.OnConfigChange(func(e fsnotify.Event) {
		//	glog.Info("App Config file changed %s:", e.Name)
	})

	app.AppConfig = Config(*appConfig)
}

// loadDBConfig: read application config and build viper object
func (app *Application) loadDBConfig() {
	dbConfig := viper.New()
	dbConfig.SetConfigType("json")
	dbConfig.AddConfigPath(".")
	dbConfig.SetConfigName("config")

	if err := dbConfig.ReadInConfig(); err != nil {
		panic(err)
	}

	dbConfig.WatchConfig()
	dbConfig.OnConfigChange(func(e fsnotify.Event) {
		// log.Info("App Config file changed %s:", e.Name)
	})

	viperConfig := viper.Viper(Config(*dbConfig))

	db := viperConfig.Sub("database.mysql")
	conf := Database{
		Driver:    db.GetString("driver"),
		Host:      db.GetString("host"),
		User:      db.GetString("user"),
		Password:  db.GetString("password"),
		DBName:    db.GetString("db_name"),
		DBNumber:  db.GetInt("db_number"),
		Port:      db.GetInt("port"),
		DebugMode: db.GetBool("debug"),
		Pool: Pool{
			MaxOpenConns: db.GetInt("maxOpenConns"),
			MaxIdleConns: db.GetInt("maxIdleConns"),
			MaxLifetime:  db.GetInt("maxLifetime"),
		},
	}

	app.DBConfig = mysqlConnect(conf)
}

func mysqlConnect(config Database) *gorm.DB {
	var (
		connectionString string
		err              error
		db               *gorm.DB
	)

	connectionString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", config.User, config.Password, config.Host, config.Port, config.DBName)

	if db, err = gorm.Open("mysql", connectionString); err != nil {
		panic(err)
	}
	if err = db.DB().Ping(); err != nil {
		panic(err)
	}

	db.LogMode(true)
	db.DB().SetMaxOpenConns(config.Pool.MaxOpenConns)
	db.DB().SetMaxIdleConns(config.Pool.MaxIdleConns)

	return db
}

// loadDBConfig: read application config and build viper object
func (app *Application) loadRedisConfig() {
	redisConfig := viper.New()
	redisConfig.SetConfigType("json")
	redisConfig.AddConfigPath(".")
	redisConfig.SetConfigName("config")

	if err := redisConfig.ReadInConfig(); err != nil {
		panic(err)
	}

	redisConfig.WatchConfig()
	redisConfig.OnConfigChange(func(e fsnotify.Event) {
		//	glog.Info("App Config file changed %s:", e.Name)
	})

	viperConfig := viper.Viper(Config(*redisConfig))

	db := viperConfig.Sub("database.redis")
	conf := Database{
		Driver:   db.GetString("driver"),
		Host:     db.GetString("host"),
		User:     db.GetString("user"),
		Password: db.GetString("password"),
		DBNumber: db.GetInt("db_number"),
		Port:     db.GetInt("port"),
	}

	app.RedisConfig = RedisConnect(conf)
}

func RedisConnect(config Database) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + strconv.Itoa(config.Port),
		Password: config.Password,
		DB:       config.DBNumber,
	})

	return client
}

// loadENV
func (app *Application) loadENV() {
	var APPENV string
	var appConfig viper.Viper
	appConfig = viper.Viper(app.AppConfig)
	APPENV = appConfig.GetString("env")
	switch APPENV {
	case "dev":
		app.ENV = "dev"
		break
	case "staging":
		app.ENV = "staging"
		break
	case "production":
		app.ENV = "production"
		break
	default:
		app.ENV = "dev"
		break
	}
}

// String: read string value from viper.Viper
func (config *Config) String(key string) string {
	var viperConfig viper.Viper
	viperConfig = viper.Viper(*config)
	return viperConfig.GetString(fmt.Sprintf("%s.%s", App.ENV, key))
}

// Int: read int value from viper.Viper
func (config *Config) Int(key string) int {
	var viperConfig viper.Viper
	viperConfig = viper.Viper(*config)
	return viperConfig.GetInt(fmt.Sprintf("%s.%s", App.ENV, key))
}

// Boolean: read boolean value from viper.Viper
func (config *Config) Boolean(key string) bool {
	var viperConfig viper.Viper
	viperConfig = viper.Viper(*config)
	return viperConfig.GetBool(fmt.Sprintf("%s.%s", App.ENV, key))
}
