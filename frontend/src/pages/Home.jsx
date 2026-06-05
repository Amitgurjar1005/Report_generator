import { useState } from "react";
import { useNavigate } from "react-router-dom";
import "./Home.css";

export default function Home() {
  const [url, setUrl] = useState("");
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const generateReport = async () => {
    if (!url) return;
    try {
      setLoading(true);

      const res = await fetch("http://localhost:8085/report", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ url }),
      });

      const data = await res.json();
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
        <div className="badge">✨ Professional Audit</div>
        <h1>Website Intelligence Tool</h1>
        <p>
          Analyze your website's <strong>SEO</strong>, <strong>Performance</strong>, 
          <strong>Security</strong>, and <strong>Network</strong> metrics instantly.
        </p>

        <div className="searchContainer">
          <div className="searchBox">
            <span className="searchIcon" aria-hidden="true">🌐</span>
            <input
              type="url"
              placeholder="Enter website URL (e.g., https://example.com)"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              disabled={loading}
              aria-label="Website URL"
            />
          </div>

          <button
            onClick={generateReport}
            disabled={loading || !url}
            className="generateBtn"
          >
            {loading ? (
              <>
                <span className="spinner"></span> Analyzing...
              </>
            ) : (
              "Analyze Website"
            )}
          </button>
        </div>

        {loading && (
          <div className="loadingCard">
            <div className="progressBar">
              <div className="progressLine"></div>
            </div>
            <p>Scanning security protocols, performance budgets, and SEO data...</p>
          </div>
        )}
      </div>
    </div>
  );
}