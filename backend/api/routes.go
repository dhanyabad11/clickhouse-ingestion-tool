
package api

import (
    "github.com/dhanyabad11/clickhouse-ingestion-tool/clickhouse"
    "github.com/dhanyabad11/clickhouse-ingestion-tool/flatfile"
    "github.com/gin-gonic/gin"
    "net/http"
)

func SetupRoutes(r *gin.Engine, client *clickhouse.Client) {
    r.POST("/connect/clickhouse", func(c *gin.Context) {
        var config struct {
            Host     string `json:"host"`
            Port     string `json:"port"`
            User     string `json:"user"`
            Database string `json:"database"`
            JWT      string `json:"jwt"`
        }
        if err := c.BindJSON(&config); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        newClient := clickhouse.NewClient(config.Host, config.Port, config.User, config.Database, config.JWT)
        if err := newClient.Connect(); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        *client = *newClient
        c.JSON(http.StatusOK, gin.H{"status": "connected"})
    })

    r.GET("/tables", func(c *gin.Context) {
        tables, err := client.GetTables()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"tables": tables})
    })

    r.GET("/columns", func(c *gin.Context) {
        table := c.Query("table")
        columns, err := client.GetColumns(table)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"columns": columns})
    })

    r.POST("/ingest", func(c *gin.Context) {
        var req struct {
            Source    string   `json:"source"`
            Table     string   `json:"table"`
            Columns   []string `json:"columns"`
            OutputPath string  `json:"outputPath"`
            Delimiter string  `json:"delimiter"`
            FilePath  string  `json:"filePath"`
        }
        if err := c.BindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        if req.Source == "clickhouse" {
            count, err := client.Ingest(req.Table, req.Columns, req.OutputPath, req.Delimiter)
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            c.JSON(http.StatusOK, gin.H{"record_count": count})
        } else {
            count, err := flatfile.Ingest(req.FilePath, req.Table, req.Delimiter, client)
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            c.JSON(http.StatusOK, gin.H{"record_count": count})
        }
    })

    r.POST("/ingest/join", func(c *gin.Context) {
        var req struct {
            Tables       []string `json:"tables"`
            Columns      []string `json:"columns"`
            JoinCondition string  `json:"joinCondition"`
            OutputPath   string  `json:"outputPath"`
            Delimiter    string  `json:"delimiter"`
        }
        if err := c.BindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        count, err := client.IngestJoin(req.Tables, req.Columns, req.JoinCondition, req.OutputPath, req.Delimiter)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"record_count": count})
    })

    r.POST("/preview", func(c *gin.Context) {
        var req struct {
            Source  string   `json:"source"`
            Table   string   `json:"table"`
            Columns []string `json:"columns"`
        }
        if err := c.BindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        data, err := client.Preview(req.Table, req.Columns)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"data": data})
    })
}
