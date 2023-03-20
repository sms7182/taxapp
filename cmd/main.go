package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
	"os"
	kafka2 "tax-management/external/kafka"
	"tax-management/external/pg"
	redis2 "tax-management/external/redis"
	"tax-management/pkg"
	terminal "tax-management/taxDep"
	"tax-management/taxDep/types"

	"github.com/go-redis/redis/v8"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"

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
	repository := pg.RepositoryImpl{DB: db}

	//rdb := getActualRedisClient()
	//redisClient := getRedisClient(rdb)

	kafkaConn := NewConsumer()
	syncConsumer := kafka2.SyncConsumer{Conn: kafkaConn}

	araJahanUsername := viper.GetString("araJahanUsername")
	delijanUsername := viper.GetString("delijanUsername")

	arajahanTerminal, err := terminal.New(
		types.TerminalOptions{
			TripPrivatePemPath: "./sign_ara.key",
			ClientID:           araJahanUsername,
			TerminalBaseURl:    viper.GetString("taxOrg.url"),
		},
	)
	if err != nil {
		panic(fmt.Sprintf("failed to create client for %s, err: %+v", araJahanUsername, err))
	}

	delijanTerminal, err := terminal.New(
		types.TerminalOptions{
			TripPrivatePemPath: "sign_ara.key",
			ClientID:           delijanUsername,
			TerminalBaseURl:    viper.GetString("taxOrg.url"),
		},
	)
	if err != nil {
		panic(fmt.Sprintf("failed to create client for %s, err: %+v", delijanUsername, err))
	}

	taxClientMap := make(map[string]*terminal.Terminal)
	taxClientMap[araJahanUsername] = arajahanTerminal
	taxClientMap[delijanUsername] = delijanTerminal
	service := pkg.Service{Repository: repository, TaxClient: taxClientMap}
	syncConsumer.StartConsuming(viper.GetStringSlice("kafka.consumerTopics"), service.ProcessKafkaMessage)
	controller := pkg.Controller{}
	router := gin.New()
	controller.SetRoutes(router)
	router.Run(viper.GetString("serverPort"))
}

func NewConsumer() *kafka.Consumer {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": viper.GetString("kafka.urls"),
		"group.id":          "tax-management2",
		"auto.offset.reset": "smallest"})
	if err != nil {
		panic("failed to create consumer")
	}
	return consumer
}

func getPrivateKey(pvPath string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	prvPemBytes, err := os.ReadFile(pvPath)
	if err != nil {
		return nil, nil, err
	}

	prvBlock, _ := pem.Decode(prvPemBytes)
	if prvBlock == nil {
		return nil, nil, fmt.Errorf("invalid trip private key")
	}

	prv, err := x509.ParsePKCS8PrivateKey(prvBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}

	privateKey := prv.(*rsa.PrivateKey)

	return privateKey, &privateKey.PublicKey, nil
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
func getRedisClient(rdb *redis.Client) redis2.ServiceImpl {
	return redis2.ServiceImpl{Rdb: rdb}
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
