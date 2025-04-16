package main

import (
    "github.com/dhanyabad11/clickhouse-ingestion-tool/api"
    "github.com/dhanyabad11/clickhouse-ingestion-tool/clickhouse"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    client := &clickhouse.Client{}
    api.SetupRoutes(r, client)
    r.Run(":8080")
}
