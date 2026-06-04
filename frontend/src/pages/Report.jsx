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
  const element = document.querySelector(".a4");

  const canvas = await html2canvas(element, {
    scale: 2,
    useCORS: true,
    scrollY: -window.scrollY,
  });

  const imgData = canvas.toDataURL("image/png");

  const pdf = new jsPDF("p", "mm", "a4");

  const pdfWidth = 210;
  const pdfHeight = 297;

  const imgWidth = pdfWidth;
  const imgHeight = (canvas.height * imgWidth) / canvas.width;

  let heightLeft = imgHeight;
  let position = 0;

  pdf.addImage(imgData, "PNG", 0, position, imgWidth, imgHeight);
  heightLeft -= pdfHeight;

  while (heightLeft > 0) {
    position = heightLeft - imgHeight;
    pdf.addPage();
    pdf.addImage(imgData, "PNG", 0, position, imgWidth, imgHeight);
    heightLeft -= pdfHeight;
  }

  pdf.save("website-report.pdf");
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
) : typeof value === "object" ? (
  Object.keys(value).length === 0 ? (
    "No Data"
  ) : (
    <pre className="json">{JSON.stringify(value, null, 2)}</pre>
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