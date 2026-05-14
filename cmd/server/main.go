package main

import (
	"gin-quickstart/config"
	"gin-quickstart/routes"
)

func main() {

	{
		_, err := config.InitDB()

		if err != nil {
			panic("failed to connect to database: " + err.Error())
		}
	}

	{
		redis, err := config.InitRedis()

		if err != nil {
			panic("failed to connect to Redis: " + err.Error())
		}

		defer redis.Close()

		config.RedisClient = redis
	}

	r := routes.SetupRouter()
	r.Run(":8080")
}
