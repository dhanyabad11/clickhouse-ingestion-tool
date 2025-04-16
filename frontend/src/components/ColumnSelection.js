import React, { useState, useEffect } from "react";
import axios from "axios";

const ColumnSelection = ({ source, table, onSelect }) => {
    const [columns, setColumns] = useState([]);
    const [selectedColumns, setSelectedColumns] = useState([]);

    useEffect(() => {
        if (source && table) {
            axios
                .get(`http://localhost:8080/columns?table=${table}`)
                .then((res) => setColumns(res.data.columns))
                .catch((err) => console.error(err));
        }
    }, [source, table]);

    const handleSelect = (column) => {
        const updated = selectedColumns.includes(column)
            ? selectedColumns.filter((c) => c !== column)
            : [...selectedColumns, column];
        setSelectedColumns(updated);
        onSelect(updated);
    };

    return (
        <div className="section">
            <h2 className="section-title">Select Columns</h2>
            {columns.map((col) => (
                <div key={col} className="checkbox-group">
                    <input
                        type="checkbox"
                        checked={selectedColumns.includes(col)}
                        onChange={() => handleSelect(col)}
                    />
                    <label>{col}</label>
                </div>
            ))}
        </div>
    );
};

export default ColumnSelection;
