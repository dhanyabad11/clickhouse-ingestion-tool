package models

type ConnectionConfig struct {
    Host     string `json:"host"`
    Port     string `json:"port"`
    User     string `json:"user"`
    Database string `json:"database"`
    JWT      string `json:"jwt"`
}

type IngestionRequest struct {
    Source       string   `json:"source"`
    Table        string   `json:"table"`
    Columns      []string `json:"columns"`
    OutputPath   string   `json:"outputPath"`
    Delimiter    string   `json:"delimiter"`
}

type JoinRequest struct {
    Tables       []string `json:"tables"`
    Columns      []string `json:"columns"`
    JoinCondition string   `json:"joinCondition"`
    OutputPath   string   `json:"outputPath"`
    Delimiter    string   `json:"delimiter"`
}