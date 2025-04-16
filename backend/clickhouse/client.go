package clickhouse

import (
    "context"
    "fmt"
    "io"
    "os"
    "strings"
    "github.com/ClickHouse/clickhouse-go/v2"
    "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type Client struct {
    conn driver.Conn
}

func NewClient(host, port, user, database, jwt string) (*Client, error) {
    conn, err := clickhouse.Open(&clickhouse.Options{
        Addr: []string{fmt.Sprintf("%s:%s", host, port)},
        Auth: clickhouse.Auth{
            Database: database,
            Username: user,
            Password: jwt,
        },
        Protocol: clickhouse.HTTP,
    })
    if err != nil {
        return nil, err
    }
    return &Client{conn: conn}, nil
}

func (c *Client) GetTables(ctx context.Context) ([]string, error) {
    rows, err := c.conn.Query(ctx, "SHOW TABLES")
    if err != nil {
        return nil, err
    }
    var tables []string
    for rows.Next() {
        var table string
        if err := rows.Scan(&table); err != nil {
            return nil, err
        }
        tables = append(tables, table)
    }
    return tables, nil
}

func (c *Client) GetColumns(ctx context.Context, table string) ([]string, error) {
    rows, err := c.conn.Query(ctx, fmt.Sprintf("DESCRIBE TABLE %s", table))
    if err != nil {
        return nil, err
    }
    var columns []string
    for rows.Next() {
        var name, typ string
        if err := rows.Scan(&name, &typ); err != nil {
            return nil, err
        }
        columns = append(columns, name)
    }
    return columns, nil
}

func (c *Client) IngestToFlatFile(ctx context.Context, table string, columns []string, outputPath, delimiter string) (int64, error) {
    query := fmt.Sprintf("SELECT %s FROM %s", strings.Join(columns, ","), table)
    rows, err := c.conn.Query(ctx, query)
    if err != nil {
        return 0, err
    }
    defer rows.Close()

    file, err := os.Create(outputPath)
    if err != nil {
        return 0, err
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    writer.Comma = rune(delimiter[0])
    if err := writer.Write(columns); err != nil {
        return 0, err
    }

    var count int64
    for rows.Next() {
        values := make([]string, len(columns))
        if err := rows.Scan(&values); err != nil {
            return count, err
        }
        if err := writer.Write(values); err != nil {
            return count, err
        }
        count++
    }
    writer.Flush()
    return count, writer.Error()
}

func (c *Client) IngestFromFlatFile(ctx context.Context, inputPath, delimiter, table string, columns []string) (int64, error) {
    file, err := os.Open(inputPath)
    if err != nil {
        return 0, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.Comma = rune(delimiter[0])
    _, err = reader.Read() // Skip headers
    if err != nil {
        return 0, err
    }

    var count int64
    batchSize := 1000
    var batch []interface{}
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return count, err
        }
        batch = append(batch, record)
        if len(batch) >= batchSize {
            if err := c.conn.Exec(ctx, fmt.Sprintf("INSERT INTO %s (%s) VALUES", table, strings.Join(columns, ",")), batch...); err != nil {
                return count, err
            }
            count += int64(len(batch))
            batch = nil
        }
    }
    if len(batch) > 0 {
        if err := c.conn.Exec(ctx, fmt.Sprintf("INSERT INTO %s (%s) VALUES", table, strings.Join(columns, ",")), batch...); err != nil {
            return count, err
        }
        count += int64(len(batch))
    }
    return count, nil
}

func (c *Client) IngestJoinedToFlatFile(ctx context.Context, tables []string, columns []string, joinCondition string, outputPath, delimiter string) (int64, error) {
    query := fmt.Sprintf("SELECT %s FROM %s %s", strings.Join(columns, ","), tables[0], joinCondition)
    rows, err := c.conn.Query(ctx, query)
    if err != nil {
        return 0, err
    }
    defer rows.Close()

    file, err := os.Create(outputPath)
    if err != nil {
        return 0, err
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    writer.Comma = rune(delimiter[0])
    if err := writer.Write(columns); err != nil {
        return 0, err
    }

    var count int64
    for rows.Next() {
        values := make([]string, len(columns))
        if err := rows.Scan(&values); err != nil {
            return count, err
        }
        if err := writer.Write(values); err != nil {
            return count, err
        }
        count++
    }
    writer.Flush()
    return count, writer.Error()
}

func (c *Client) PreviewData(ctx context.Context, table string, columns []string) ([][]string, error) {
    query := fmt.Sprintf("SELECT %s FROM %s LIMIT 100", strings.Join(columns, ","), table)
    rows, err := c.conn.Query(ctx, query)
    if err != nil {
        return nil, err
    }
    var result [][]string
    for rows.Next() {
        values := make([]string, len(columns))
        if err := rows.Scan(&values); err != nil {
            return nil, err
        }
        result = append(result, values)
    }
    return result, nil
}