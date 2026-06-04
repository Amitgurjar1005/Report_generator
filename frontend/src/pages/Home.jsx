import { useState } from "react";
import { useNavigate } from "react-router-dom";
import "./Home.css";

export default function Home() {
  const [url, setUrl] = useState("");
  const [report, setReport] = useState(null);
  const [loading, setLoading] = useState(false);

  const navigate = useNavigate();


  const generateReport = async () => {
    try {
      setLoading(true);
      setReport(null);

      const res = await fetch("http://localhost:8084/report", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ url }),
      });

      const data = await res.json();
      setReport(data);

  
      navigate("/report", { state: data });

    } catch (err) {
      console.error("Error:", err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="homePage">

      <div className="hero">
        <h1>🚀 Website Intelligence Tool</h1>
        <p>Analyze SEO, Performance, Security & Network in seconds</p>

    
        <div className="searchBox">
          <input
            type="text"
            placeholder="Enter website URL (https://example.com)"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
          />
        </div>

        <div className="buttonGroup">

          <button
            onClick={generateReport}
            disabled={loading || !url}
            className="generateBtn"
          >
            {loading ? "Analyzing..." : "Generate Report"}
          </button>

        </div>
      </div>

      {loading && (
        <div className="loadingCard">
          🔍 Scanning website... please wait
        </div>
      )}

    </div>
  );
}