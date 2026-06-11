
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import "./Home.css";
import User from "./User"
export default function Home() {
  const [url, setUrl] = useState("");
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
const [showModal, setShowModal] = useState(false);
const user = JSON.parse(
    localStorage.getItem("user")
  );

  const generateReport = async () => {
    if (!url) return;

    try {
      setLoading(true);

      const res = await fetch("http://localhost:8086/report", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ url }),
      });

      const data = await res.json();

      navigate("/report", {
        state: data,
      });
    } catch (err) {
      console.error("Error generating report:", err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="homePage">
      <div className="bgGridPattern"></div>
      <div className="ambientGlow"></div>

      <div className="heroContainer">

        <div className="featureBadge">
          <span className="badgePulse"></span>
          PROFESSIONAL WEBSITE AUDIT
        </div>

        <h1 className="mainHeading">
          Website <span className="gradientText">Intelligence</span>
          <br />
          Platform
        </h1>

        <p className="subHeading">
          Analyze your website's <strong>SEO</strong>,
          <strong> Security</strong>,
          <strong> Performance</strong>,
          <strong> Content</strong>,
          <strong> Cookies</strong>,
          <strong> Social Media</strong> and
          <strong> Network</strong> metrics instantly.
        </p>

        <div className="searchWrapper">
          <div className="inputContainer">
            <span className="inputIcon">🌐</span>

            <input
              type="url"
              placeholder="https://example.com"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              disabled={loading}
            />
          </div>

          <button
            className="actionButton"
            onClick={generateReport}
            disabled={!url || loading}
          >
            {loading ? (
              <>
                <span className="btnSpinner"></span>
                Analyzing...
              </>
            ) : (
              <>
                Analyze Website
                <span className="arrowIcon">→</span>
              </>
            )}
          </button>
        </div>

        {loading && (
          <div className="loadingContainer">
            <div className="skeletonTrack">
              <div className="skeletonBar"></div>
            </div>

            <p className="loadingText">
              Scanning SEO, Security, Content and Performance metrics...
            </p>
          </div>
        )}
   <div className="featureCard">
  <div className="featureEmoji">👤</div>
  <h3>User Registration</h3>
 {user?.email ? (
    <>
      <p className="successText">
        ✅ Successfully Registered
      </p>

      <p>
        Reports will be automatically sent to:
        <br />
        <strong>{user.email}</strong>
      </p>
        <button
        className="registerBtn"
        onClick={() => setShowModal(true)}
      >
        Change Email
      </button>
    </>
  ) : (
    <>
      <p>
        Register yourself so reports can be sent automatically to your email.
      </p>

      <button
        className="registerBtn"
        onClick={() => setShowModal(true)}
      >
        Register
      </button>
    </>
  )}
  
</div>
        <div className="featureCards">

          <div className="featureCard">
            <div className="featureEmoji">🔍</div>
            <h3>SEO Analysis</h3>
            <p>
              Check meta tags, headings, page structure and indexing readiness.
            </p>
          </div>

          <div className="featureCard">
            <div className="featureEmoji">🔒</div>
            <h3>Security Audit</h3>
            <p>
              Analyze cookies, HTTPS configuration and security headers.
            </p>
          </div>

          <div className="featureCard">
            <div className="featureEmoji">⚡</div>
            <h3>Performance</h3>
            <p>
              Review page assets, load behavior and optimization opportunities.
            </p>
          </div>

          <div className="featureCard">
            <div className="featureEmoji">📄</div>
            <h3>Content Insights</h3>
            <p>
              Discover content statistics, images, links and page elements.
            </p>
          </div>

          <div className="featureCard">
            <div className="featureEmoji">🔗</div>
            <h3>Broken Links</h3>
            <p>
              Detect broken links and missing resources across your website.
            </p>
          </div>

          <div className="featureCard">
            <div className="featureEmoji">📱</div>
            <h3>Social Media</h3>
            <p>
              Extract and verify connected social media profiles and links.
            </p>
          </div>

       

        </div>
      </div>
      {showModal && (
  <div
    className="modalOverlay"
    onClick={() => setShowModal(false)}
  >
    <div
      className="modalContent"
      onClick={(e) => e.stopPropagation()}
    >
      <User setShowModal={setShowModal} />
    </div>
  </div>
)}
    </div>
  );
}
