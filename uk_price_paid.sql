CREATE TABLE uk_price_paid (
    price UInt32,
    date Date,
    postcode1 String,
    postcode2 String,
    type Enum8('terraced' = 1, 'semi-detached' = 2, 'detached' = 3, 'flat' = 4, 'other' = 0),
    is_new UInt8,
    duration Enum8('freehold' = 1, 'leasehold' = 2, 'unknown' = 0),
    addr1 String,
    addr2 String,
    street String,
    locality String,
    town String,
    district String,
    county String
)
ENGINE = MergeTree()
ORDER BY (date, price);

INSERT INTO uk_price_paid
SELECT *
FROM s3('https://datasets.clickhouse.com/uk_price_paid/rows.csv', 'CSVWithNames')
SETTINGS input_format_csv_skip_first_lines = 1;
