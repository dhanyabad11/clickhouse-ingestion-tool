import React, { useState } from "react";

const SourceSelection = ({ onSelect }) => {
    const [source, setSource] = useState("");

    const handleSelect = (e) => {
        setSource(e.target.value);
        onSelect(e.target.value);
    };

    return (
        <div className="section">
            <label className="section-title">Select Source:</label>
            <select value={source} onChange={handleSelect}>
                <option value="">Select</option>
                <option value="clickhouse">ClickHouse</option>
                <option value="flatfile">Flat File</option>
            </select>
        </div>
    );
};

export default SourceSelection;
