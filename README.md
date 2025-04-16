# Bidirectional ClickHouse & Flat File Data Ingestion Tool

A web-based tool for bidirectional data ingestion between ClickHouse and CSV files.

## Setup

1. Install Go, Node.js, Docker.
2. Run ClickHouse:
    ```bash
    docker run -d -p 8123:8123 -p 9000:9000 clickhouse/clickhouse-server
    ```
