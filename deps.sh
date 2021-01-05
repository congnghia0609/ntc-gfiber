#!/bin/bash
# Author:       nghiatc
# Email:        congnghia0609@gmail.com

# source /etc/profile

echo "Install library dependencies..."
go get -u github.com/tools/godep
go get -u github.com/congnghia0609/ntc-gconf
go get -u github.com/spf13/viper
go get -u github.com/gorilla/mux
go get -u github.com/sirupsen/logrus
go get -u github.com/natefinch/lumberjack
go get -u github.com/satori/go.uuid
go get -u github.com/valyala/fasthttp
go get -u github.com/valyala/tcplisten
go get -u github.com/gofiber/fiber/v2
go get -u go.mongodb.org/mongo-driver
echo "Install dependencies complete..."
