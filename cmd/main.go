package main

import (
	"fmt"
	"time"

	"os"

	"tax-management/external/exkafka"
	"tax-management/external/exkafka/messages"
	"tax-management/external/gateway"
	taxorganization "tax-management/external/gateway/tax_organization"
	"tax-management/external/pg"
	"tax-management/pkg"

	"github.com/gin-gonic/gin"

	"github.com/go-redis/redis/v8"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"

	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	setUpViper()
	db := getGormDb()
	runDbMigrations()
	repository := pg.RepositoryImpl{
		DB: db,
	}

	client := gateway.ClientLoggerExtensionImpl{
		GatewayRepository: repository,
	}
	service := pg.ServiceImp{
		Repository: repository,
		Client:     client,
	}
	rdb := getActualRedisClient()
	redisClient := getRedisClient(rdb)
	taxClient := taxClientConfiguration(repository, client)
	kafkaService := kafkaConfiguration(repository, redisClient, taxClient)
	controller := pkg.Controller{
		Service:      service,
		KafkaService: kafkaService,
	}
	router := gin.New()
	controller.SetRoutes(router)

	router.Run(viper.GetString("serverPort"))
}

func taxClientConfiguration(repository pg.RepositoryImpl, client gateway.ClientLoggerExtensionImpl) pkg.TaxClient {
	url := viper.GetString("taxOrg.url")
	serverUrl := viper.GetString("taxOrg.serverInformationUrl")
	taxClient := taxorganization.ClientImpl{
		Repository:           repository,
		HttpClient:           client,
		Url:                  url,
		ServerInformationUrl: serverUrl,
	}
	return taxClient
}
func kafkaConfiguration(repository pg.RepositoryImpl, redis pg.RedisServiceImpl, client pkg.TaxClient) exkafka.KafkaServiceImpl {
	topic := viper.GetString("kafka.topic")
	bs := viper.GetString("kafka.urls")
	writer := &kafka.Writer{
		Addr:     kafka.TCP(bs),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}
	rawConsumer := messages.Consumer[*messages.RawTransaction]{
		Dialer:  dialer,
		Topic:   topic,
		Brokers: []string{bs},
	}
	rawConsumer.CreateConnection()

	kafkaConf := exkafka.KafkaServiceImpl{

		Writer:        writer,
		Url:           viper.GetString("taxOrg.url"),
		TokenUrl:      viper.GetString("taxOrg.tokenUrl"),
		ServerInfoUrl: viper.GetString("taxOrg.serverInformationUrl"),
		Redis:         redis,
		Repository:    repository,
		TaxClient:     client,
	}

	go rawConsumer.Read(&messages.RawTransaction{}, func(rt *messages.RawTransaction, err error) {
		if err != nil {
			fmt.Println(err.Error())
		}
		kafkaConf.Consumer(rt, err)
	})
	return kafkaConf

}

func getActualRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.url"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})
}
func setUpViper() {
	viper.SetConfigName(getEnv("CONFIG_NAME", "dev-conf"))
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./conf")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %+v \n", err)
	}
}
func getRedisClient(rdb *redis.Client) pg.RedisServiceImpl {
	return pg.RedisServiceImpl{Rdb: rdb}
}
func getGormDb() *gorm.DB {
	connection := viper.GetString("postgresSource")
	db, err := gorm.Open(postgres.Open(connection), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to initial gorm DB")
	}

	return db
}
func runDbMigrations() {
	fmt.Print("start migrations for second time")
	m, err := migrate.New(
		"file://db/migrations",
		viper.GetString("pgMigrationSource"))
	if err != nil {
		log.Fatalf("failed to find migration files")
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to apply migrations%s", err.Error())
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
