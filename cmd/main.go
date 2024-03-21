package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"log"
	todo "todo-app"
	"todo-app/pkg/handler"
	"todo-app/pkg/repository"
	"todo-app/pkg/repository/database/mysql"
	"todo-app/pkg/service"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	db, err := mysql.NewMySqlDB(&mysql.Config{
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		Port:     viper.GetString("db.port"),
		Host:     viper.GetString("db.host"),
		DBName:   viper.GetString("db.dbname"),
	})
	if err != nil {
		log.Fatalf("error occured while initializating database: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRouter()); err != nil {
		log.Fatalf("error occured while running server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
