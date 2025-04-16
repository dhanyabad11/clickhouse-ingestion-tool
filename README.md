Bidirectional ClickHouse & Flat File Data Ingestion Tool
This tool enables bidirectional data transfer between ClickHouse databases and CSV files, with a web interface for selecting tables, columns, and configuring ingestion. It uses a Go/Gin backend and a React frontend (vanilla CSS). The project is hosted at https://github.com/dhanyabad11/clickhouse-ingestion-tool.
This README.md provides step-by-step instructions to set up and run the tool locally, including prerequisites, environment configuration, dataset loading, and starting the backend and frontend.
Prerequisites
Before setting up, ensure you have the following installed:

Git: To clone the repository.
Verify: git --version

Go: Version 1.23.6 or later (for backend).
Verify: go version

Node.js and npm: Node 18+ (for frontend).
Verify: node --version && npm --version

Docker: To run ClickHouse server.
Verify: docker --version

macOS/Linux/Windows: Instructions are macOS-focused (Apple Silicon tested) but adaptable.
Internet Access: For downloading dependencies and datasets.

Setup Instructions

1. Clone the Repository
   Clone the project to your local machine:
   git clone https://github.com/dhanyabad11/clickhouse-ingestion-tool.git
   cd clickhouse-ingestion-tool

2. Install Dependencies
   Backend (Go)

Navigate to the backend:cd backend

Install Go dependencies:go mod tidy

Installs github.com/gin-gonic/gin, github.com/ClickHouse/clickhouse-go/v2, and others per go.mod.

Verify:ls go.mod go.sum

Should list go.mod and go.sum.

Frontend (React)

Navigate to the frontend:cd ../frontend

Install Node.js dependencies:npm install

Installs React, Axios, and other packages per package.json.

Verify:ls node_modules

Should show node_modules directory.

3. Set Up ClickHouse
   The tool uses ClickHouse to store and query datasets (uk_price_paid, ontime).

Start ClickHouse Server:Run ClickHouse in Docker:
docker run -d -p 8123:8123 -p 9000:9000 --name clickhouse clickhouse/clickhouse-server

Exposes HTTP port 8123 and TCP port 9000.

Verify ClickHouse:
curl http://localhost:8123

Should return Ok.

docker ps

Should list clickhouse/clickhouse-server.

Load Datasets:

Create SQL files for datasets:cd ..
cat << 'EOF' > uk_price_paid.sql
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
SELECT \*
FROM s3('https://datasets.clickhouse.com/uk_price_paid/rows.csv', 'CSVWithNames')
SETTINGS input_format_csv_skip_first_lines = 1;
EOF
cat << 'EOF' > ontime.sql
CREATE TABLE ontime (
Year UInt16,
Quarter UInt8,
Month UInt8,
DayofMonth UInt8,
DayOfWeek UInt8,
FlightDate Date,
UniqueCarrier FixedString(7),
AirlineID Int32,
Carrier FixedString(2),
TailNum String,
FlightNum String,
OriginAirportID Int32,
OriginAirportSeqID Int32,
OriginCityMarketID Int32,
Origin String,
OriginCityName String,
OriginState String,
OriginStateFips String,
OriginStateName String,
OriginWac Int32,
DestAirportID Int32,
DestAirportSeqID Int32,
DestCityMarketID Int32,
Dest String,
DestCityName String,
DestState String,
DestStateFips String,
DestStateName String,
DestWac Int32,
CRSDepTime Int32,
DepTime Int32,
DepDelay Int32,
DepDelayMinutes Int32,
DepDel15 Int32,
DepartureDelayGroups Int32,
DepTimeBlk String,
TaxiOut Int32,
WheelsOff Int32,
WheelsOn Int32,
TaxiIn Int32,
CRSArrTime Int32,
ArrTime Int32,
ArrDelay Int32,
ArrDelayMinutes Int32,
ArrDel15 Int32,
ArrivalDelayGroups Int32,
ArrTimeBlk String,
Cancelled UInt8,
CancellationCode FixedString(1),
Diverted UInt8,
CRSElapsedTime Int32,
ActualElapsedTime Int32,
AirTime Int32,
Flights Int32,
Distance Int32,
DistanceGroup UInt8,
CarrierDelay Int32,
WeatherDelay Int32,
NASDelay Int32,
SecurityDelay Int32,
LateAircraftDelay Int32,
FirstDepTime String,
TotalAddGTime String,
LongestAddGTime String,
DivAirportLandings String,
DivReachedDest String,
DivActualElapsedTime String,
DivArrDelay String,
DivDistance String,
Div1Airport String,
Div1AirportID Int32,
Div1AirportSeqID Int32,
Div1WheelsOn String,
Div1Total
