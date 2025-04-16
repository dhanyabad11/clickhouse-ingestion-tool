import React, { useState } from "react";
import SourceSelection from "./components/SourceSelection";
import ConnectionForm from "./components/ConnectionForm";
import ColumnSelection from "./components/ColumnSelection";
import axios from "axios";
import "./App.css";

const App = () => {
    const [source, setSource] = useState("");
    const [tables, setTables] = useState([]);
    const [selectedTable, setSelectedTable] = useState("");
    const [selectedTables, setSelectedTables] = useState([]);
    const [selectedColumns, setSelectedColumns] = useState([]);
    const [joinCondition, setJoinCondition] = useState("");
    const [status, setStatus] = useState("");
    const [recordCount, setRecordCount] = useState(null);
    const [progress, setProgress] = useState(0);
    const [previewData, setPreviewData] = useState(null);

    const handleConnect = async () => {
        setStatus("Connecting...");
        try {
            const res = await axios.get("http://localhost:8080/tables");
            setTables(res.data.tables);
            setStatus("Connected");
        } catch (err) {
            setStatus("Error: " + err.response?.data?.error);
        }
    };

    const handleTableSelect = (e) => {
        const table = e.target.value;
        if (table && !selectedTables.includes(table)) {
            setSelectedTables([...selectedTables, table]);
        }
    };

    const handleIngest = async () => {
        setStatus("Ingesting...");
        setProgress(0);
        try {
            const res = await axios.post("http://localhost:8080/ingest", {
                source,
                table: selectedTable,
                columns: selectedColumns,
                outputPath: "output.csv",
                delimiter: ",",
            });
            setRecordCount(res.data.record_count);
            setProgress(100);
            setStatus("Completed");
        } catch (err) {
            setStatus("Error: " + err.response?.data?.error);
        }
    };

    const handleJoinIngest = async () => {
        setStatus("Ingesting...");
        setProgress(0);
        try {
            const res = await axios.post("http://localhost:8080/ingest/join", {
                tables: selectedTables,
                columns: selectedColumns,
                joinCondition,
                outputPath: "output_joined.csv",
                delimiter: ",",
            });
            setRecordCount(res.data.record_count);
            setProgress(100);
            setStatus("Completed");
        } catch (err) {
            setStatus("Error: " + err.response?.data?.error);
        }
    };

    const handlePreview = async () => {
        try {
            const res = await axios.post("http://localhost:8080/preview", {
                source,
                table: selectedTable,
                columns: selectedColumns,
            });
            setPreviewData(res.data.data);
        } catch (err) {
            setStatus("Error: " + err.response?.data?.error);
        }
    };

    return (
        <div className="container">
            <h1 className="header">Data Ingestion Tool</h1>
            <SourceSelection onSelect={setSource} />
            {source === "clickhouse" && <ConnectionForm onConnect={handleConnect} />}
            {tables.length > 0 && (
                <div className="section">
                    <label className="section-title">Select Table(s):</label>
                    <select onChange={(e) => setSelectedTable(e.target.value)}>
                        <option value="">Select</option>
                        {tables.map((t) => (
                            <option key={t} value={t}>
                                {t}
                            </option>
                        ))}
                    </select>
                    <select onChange={handleTableSelect}>
                        <option value="">Select for Join</option>
                        {tables.map((t) => (
                            <option key={t} value={t}>
                                {t}
                            </option>
                        ))}
                    </select>
                    <p>Selected for Join: {selectedTables.join(", ")}</p>
                    {selectedTables.length > 1 && (
                        <input
                            type="text"
                            placeholder="JOIN Condition (e.g., JOIN table2 ON table1.id = table2.id)"
                            value={joinCondition}
                            onChange={(e) => setJoinCondition(e.target.value)}
                            className="input-group"
                        />
                    )}
                </div>
            )}
            {selectedTable && (
                <ColumnSelection
                    source={source}
                    table={selectedTable}
                    onSelect={setSelectedColumns}
                />
            )}
            {selectedColumns.length > 0 && (
                <div className="section">
                    <button className="button-green" onClick={handleIngest}>
                        Start Ingestion
                    </button>
                    <button className="button-yellow" onClick={handlePreview}>
                        Preview Data
                    </button>
                    {selectedTables.length > 1 && (
                        <button className="button-purple" onClick={handleJoinIngest}>
                            Start Joined Ingestion
                        </button>
                    )}
                </div>
            )}
            <div className="section">
                <p className="status">Status: {status}</p>
                {recordCount !== null && <p>Records Processed: {recordCount}</p>}
                <div className="progress-bar">
                    <div className="progress-fill" style={{ width: `${progress}%` }}></div>
                </div>
            </div>
            {previewData && (
                <table>
                    <thead>
                        <tr>
                            {selectedColumns.map((col) => (
                                <th key={col}>{col}</th>
                            ))}
                        </tr>
                    </thead>
                    <tbody>
                        {previewData.map((row, i) => (
                            <tr key={i}>
                                {row.map((val, j) => (
                                    <td key={j}>{val}</td>
                                ))}
                            </tr>
                        ))}
                    </tbody>
                </table>
            )}
        </div>
    );
};

export default App;
