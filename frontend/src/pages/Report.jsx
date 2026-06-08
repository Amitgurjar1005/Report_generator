import { useLocation, useNavigate } from "react-router-dom";
import "./Report.css";
import jsPDF from "jspdf";
import html2canvas from "html2canvas";

export default function Report() {
  const { state: report } = useLocation();
  const navigate = useNavigate();

  if (!report) {
    return (
      <div className="reportPage">
        <h2>No Report Found</h2>
        <button onClick={() => navigate("/")}>Go Back</button>
      </div>
    );
  }

  const downloadPDF = async () => {
    console.log("📄 PDF START");

    const element = document.querySelector(".a4");

    if (!element) {
      console.error("❌ .a4 element not found");
      return;
    }

    const clone = element.cloneNode(true);

    clone.style.position = "absolute";
    clone.style.left = "-99999px";
    clone.style.top = "0";
    clone.style.background = "#fff";
    clone.style.width = "210mm";

    document.body.appendChild(clone);

    const images = clone.querySelectorAll("img");

    await Promise.all(
      Array.from(images).map(async (img) => {
        try {
          const res = await fetch(img.src, { mode: "cors" });
          const blob = await res.blob();

          const base64 = await new Promise((resolve) => {
            const reader = new FileReader();
            reader.onloadend = () => resolve(reader.result);
            reader.readAsDataURL(blob);
          });

          img.src = base64;
        } catch (err) {
          img.style.display = "none";
        }
      })
    );

    await new Promise((r) => setTimeout(r, 500));

    const canvas = await html2canvas(clone, {
      scale: 2,
      useCORS: true,
      backgroundColor: "#ffffff",
    });

    document.body.removeChild(clone);

    const pdf = new jsPDF("p", "mm", "a4");

    const pdfWidth = 210;
    const pdfHeight = 297;

    const imgWidth = pdfWidth;
    const imgHeight = (canvas.height * pdfWidth) / canvas.width;

    let heightLeft = imgHeight;
    let position = 0;

    pdf.addImage(
      canvas.toDataURL("image/png"),
      "PNG",
      0,
      position,
      imgWidth,
      imgHeight
    );

    heightLeft -= pdfHeight;

    while (heightLeft > 0) {
      position -= pdfHeight;

      pdf.addPage();
      pdf.addImage(
        canvas.toDataURL("image/png"),
        "PNG",
        0,
        position,
        imgWidth,
        imgHeight
      );

      heightLeft -= pdfHeight;
    }

    pdf.save("website-report.pdf");

    console.log("✅ PDF DONE");
  };

  return (
    <div className="reportPage">
      <div className="topBar no-print">
        <div>
          <h1>Website Intelligence Report</h1>
          <p className="sub">Auto Generated Website Audit</p>
        </div>

        <div className="btnGroup">
          <button onClick={downloadPDF}>⬇ Download PDF</button>
          <button onClick={() => navigate("/")}>🏠 Home</button>
        </div>
      </div>

      <div className="a4">
        <h1 className="title">WEBSITE ANALYSIS REPORT</h1>

        <Section title="CONTENT INFO" data={report.content_info} />
        <Section title="PERFORMANCE" data={report.page_performance} />
        <Section title="SECURITY RISK" data={report.risk_infor} />
        <Section title="SEO ANALYSIS" data={report.seo_performance} />
        <Section title="NETWORK INFO" data={report.network_info} />
        <Section title="SOCIAL LINKS" data={report.sociallink_info} />
        <Section title="COOKIE INFO" data={report.cookie_info} />
      </div>
    </div>
  );
}
function Section({ title, data }) {
  return (
    <div className="section">
      <h2>{title}</h2>

      <table>
        <tbody>
          {Object.entries(data || {}).map(([key, value]) => (
            <tr key={key}>
              <td className="key">{key}</td>

              <td className="value">
                {value === "" || value === null || value === undefined ? (
                  "-"
                ) : key === "broken_images_details" &&
                  Array.isArray(value) ? (
                  value.length === 0 ? (
                    "No Broken Images"
                  ) : (
                    value.map((item, index) => (
                      <div key={index} className="broken-image-card">
                        <p><strong>Page URL:</strong></p>
                        <a href={item.page_url} target="_blank" rel="noreferrer">
                          {item.page_url}
                        </a>

                        <p style={{ marginTop: "10px" }}>
                          <strong>Image URL:</strong>
                        </p>
                        <a href={item.image_url} target="_blank" rel="noreferrer">
                          {item.image_url}
                        </a>

                        <div style={{ marginTop: "10px" }}>
                          <img
                            src={`http://localhost:8086/${item.screenshot_path}`}
                            alt="Broken Screenshot"
                            className="broken-image-preview"
                          />
                        </div>
                      </div>
                    ))
                  )
                ) : key === "broken_link_details" &&
                  Array.isArray(value) ? (
                  value.length === 0 ? (
                    "No Broken Links"
                  ) : (
                    value.map((item, index) => (
                      <div key={index} className="broken-link-card">
                        <p><strong>Page URL:</strong></p>
                        <a href={item.page_url} target="_blank" rel="noreferrer">
                          {item.page_url}
                        </a>

                        <p style={{ marginTop: "10px" }}>
                          <strong>Broken URL:</strong>
                        </p>
                        <a href={item.url} target="_blank" rel="noreferrer">
                          {item.url}
                        </a>
                      </div>
                    ))
                  )
                ) : typeof value === "object" ? (
                  Object.keys(value).length === 0 ? (
                    "No Data"
                  ) : (
                    <pre className="json">
                      {JSON.stringify(value, null, 2)}
                    </pre>
                  )
                ) : (
                  String(value)
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}