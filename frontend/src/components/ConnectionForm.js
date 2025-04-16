import React, { useState } from "react";
import axios from "axios";

const ConnectionForm = ({ onConnect }) => {
    const [config, setConfig] = useState({
        host: "",
        port: "",
        user: "",
        database: "",
        jwt: "",
    });
    const [error, setError] = useState("");

    const handleConnect = async () => {
        try {
            const response = await axios.post("http://localhost:8080/connect/clickhouse", config);
            onConnect(response.data);
            setError("");
        } catch (err) {
            setError(err.response?.data?.error || "Connection failed");
        }
    };

    return (
        <div className="section">
            <h2 className="section-title">ClickHouse Connection</h2>
            <div className="input-group">
                <input
                    type="text"
                    placeholder="Host"
                    value={config.host}
                    onChange={(e) => setConfig({ ...config, host: e.target.value })}
                />
            </div>
            <div className="input-group">
                <input
                    type="text"
                    placeholder="Port"
                    value={config.port}
                    onChange={(e) => setConfig({ ...config, port: e.target.value })}
                />
            </div>
            <div className="input-group">
                <input
                    type="text"
                    placeholder="User"
                    value={config.user}
                    onChange={(e) => setConfig({ ...config, user: e.target.value })}
                />
            </div>
            <div className="input-group">
                <input
                    type="text"
                    placeholder="Database"
                    value={config.database}
                    onChange={(e) => setConfig({ ...config, database: e.target.value })}
                />
            </div>
            <div className="input-group">
                <input
                    type="password"
                    placeholder="JWT Token"
                    value={config.jwt}
                    onChange={(e) => setConfig({ ...config, jwt: e.target.value })}
                />
            </div>
            <button className="button-primary" onClick={handleConnect}>
                Connect
            </button>
            {error && <p className="error">{error}</p>}
        </div>
    );
};

export default ConnectionForm;
