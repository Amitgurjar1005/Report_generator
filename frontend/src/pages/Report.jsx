import { useLocation, useNavigate } from "react-router-dom";
import "./Report.css";
import jsPDF from "jspdf";
import autoTable from "jspdf-autotable";
import { useState, useEffect } from "react";
export default function Report() {
  const { state: report } = useLocation();
  const navigate = useNavigate();

  const [selectedImage, setSelectedImage] = useState(null);

  if (!report) {
    return (
      <div className="reportPage">
        <h2>No Report Found</h2>
        <button onClick={() => navigate("/")}>Go Back</button>
      </div>
    );
  }

  // const downloadPDF = async () => {
  //   console.log("📄 PDF START");

  //   const element = document.querySelector(".a4");

  //   if (!element) {
  //     console.error("❌ .a4 element not found");
  //     return;
  //   }

  //   const clone = element.cloneNode(true);

  //   clone.style.position = "absolute";
  //   clone.style.left = "-99999px";
  //   clone.style.top = "0";
  //   clone.style.background = "#fff";
  //   clone.style.width = "210mm";

  //   document.body.appendChild(clone);

  //   const images = clone.querySelectorAll("img");

  //   await Promise.all(
  //     Array.from(images).map(async (img) => {
  //       try {
  //         const res = await fetch(img.src, { mode: "cors" });
  //         const blob = await res.blob();

  //         const base64 = await new Promise((resolve) => {
  //           const reader = new FileReader();
  //           reader.onloadend = () => resolve(reader.result);
  //           reader.readAsDataURL(blob);
  //         });

  //         img.src = base64;
  //       } catch (err) {
  //         img.style.display = "none";
  //       }
  //     })
  //   );

  //   await new Promise((r) => setTimeout(r, 500));

  //   const canvas = await html2canvas(clone, {
  //     scale: 2,
  //     useCORS: true,
  //     backgroundColor: "#ffffff",
  //   });

  //   document.body.removeChild(clone);

  //   const pdf = new jsPDF("p", "mm", "a4");

  //   const pdfWidth = 210;
  //   const pdfHeight = 297;

  //   const imgWidth = pdfWidth;
  //   const imgHeight = (canvas.height * pdfWidth) / canvas.width;

  //   let heightLeft = imgHeight;
  //   let position = 0;

  //   pdf.addImage(
  //     canvas.toDataURL("image/png"),
  //     "PNG",
  //     0,
  //     position,
  //     imgWidth,
  //     imgHeight
  //   );

  //   heightLeft -= pdfHeight;

  //   while (heightLeft > 0) {
  //     position -= pdfHeight;

  //     pdf.addPage();
  //     pdf.addImage(
  //       canvas.toDataURL("image/png"),
  //       "PNG",
  //       0,
  //       position,
  //       imgWidth,
  //       imgHeight
  //     );

  //     heightLeft -= pdfHeight;
  //   }

  //   pdf.save("website-report.pdf");

  //   console.log("✅ PDF DONE");
  // };

  
useEffect(() => {
  const autoSend = async () => {
    const blob = await generatepdfblob();
    await sendEmailWithPDF(blob);
  };

  autoSend();
}, []);

  const downloadPDF = async () => {
  const blob = await generatepdfblob();

  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = "website-report.pdf";
  a.click();
  URL.revokeObjectURL(url);
};
const generatepdfblob= async () => {
  const pdf = new jsPDF("p", "mm", "a4");

  let currentY = 15;
  const pageHeight = 297;

  const checkPageBreak = (spaceNeeded = 20) => {
    if (currentY + spaceNeeded > pageHeight - 10) {
      pdf.addPage();
      currentY = 15;
    }
  };

  pdf.setFontSize(20);
  pdf.setFont(undefined, "bold");
  pdf.text("WEBSITE ANALYSIS REPORT", 14, currentY);

  currentY += 12;

  const addSection = async (title, data) => {
    if (!data) return;

    checkPageBreak(30);

    pdf.setFontSize(14);
    pdf.setFont(undefined, "bold");
    pdf.text(title, 14, currentY);

    currentY += 5;

    const rows = [];

    Object.entries(data).forEach(([key, value]) => {
      if (
        key !== "broken_images_details" &&
        key !== "broken_link_details"
      ) {
        rows.push([
          key,
          typeof value === "object"
            ? JSON.stringify(value, null, 2)
            : String(value ?? "-"),
        ]);
      }
    });

    if (rows.length > 0) {
      autoTable(pdf, {
        startY: currentY,
        head: [["Metric", "Value"]],
        body: rows,
        theme: "grid",
        styles: {
          fontSize: 9,
          cellPadding: 3,
          overflow: "linebreak",
        },
        headStyles: {
          fillColor: [17, 17, 17],
          textColor: [255, 255, 255],
        },
        columnStyles: {
          0: { cellWidth: 55, fontStyle: "bold" },
          1: { cellWidth: 120 },
        },
      });

      currentY = pdf.lastAutoTable.finalY + 8;
    }

    // ================= BROKEN LINKS =================
    if (data.broken_link_details?.length > 0) {
      checkPageBreak(20);

      pdf.setFontSize(12);
      pdf.setFont(undefined, "bold");
      pdf.text("Broken Links", 14, currentY);
      currentY += 8;

      data.broken_link_details.forEach((item, index) => {
        checkPageBreak(40);

        pdf.rect(10, currentY - 5, 190, 30);

        pdf.setFont(undefined, "bold");
        pdf.text(`Broken Link #${index + 1}`, 15, currentY);
        currentY += 10;

        pdf.setFont(undefined, "bold");
        pdf.text("Page URL:", 15, currentY);

        pdf.setFont(undefined, "normal");
        pdf.textWithLink(item.page_url, 40, currentY, { url: item.page_url });

        currentY += 8;

        pdf.setFont(undefined, "bold");
        pdf.text("Broken URL:", 15, currentY);

        pdf.setFont(undefined, "normal");
        pdf.textWithLink(item.url, 40, currentY, { url: item.url });

        currentY += 15;
      });
    }

    // ================= BROKEN IMAGES =================
    if (data.broken_images_details?.length > 0) {
      checkPageBreak(20);

      pdf.setFontSize(12);
      pdf.setFont(undefined, "bold");
      pdf.text("Broken Images", 14, currentY);
      currentY += 8;

      for (let i = 0; i < data.broken_images_details.length; i++) {
        const item = data.broken_images_details[i];

        try {
          checkPageBreak(110);

          pdf.rect(10, currentY - 5, 190, 95);

          pdf.setFont(undefined, "bold");
          pdf.text(`Broken Image #${i + 1}`, 15, currentY);
          currentY += 10;

          pdf.setFont(undefined, "bold");
          pdf.text("Page URL:", 15, currentY);

          pdf.setFont(undefined, "normal");
          pdf.textWithLink(item.page_url, 40, currentY, { url: item.page_url });

          currentY += 8;

          pdf.setFont(undefined, "bold");
          pdf.text("Image URL:", 15, currentY);

          pdf.setFont(undefined, "normal");
          pdf.textWithLink(item.image_url, 40, currentY, { url: item.image_url });

          currentY += 10;

          const imageUrl = `http://localhost:8086/${item.screenshot_path}`;
          const response = await fetch(imageUrl);
          const blob = await response.blob();

          const base64 = await new Promise((resolve) => {
            const reader = new FileReader();
            reader.onloadend = () => resolve(reader.result);
            reader.readAsDataURL(blob);
          });

          pdf.addImage(base64, "PNG", 15, currentY, 100, 55);
          currentY += 65;
        } catch (err) {
          console.error("Image PDF Error:", err);
        }
      }
    }
  };

  // ================= SECTIONS =================
  await addSection("CONTENT INFO", report.content_info);
  await addSection("PERFORMANCE", report.page_performance);
  await addSection("SECURITY RISK", report.risk_infor);
  await addSection("SEO ANALYSIS", report.seo_performance);
  await addSection("NETWORK INFO", report.network_info);
  await addSection("SOCIAL LINKS", report.sociallink_info);
  await addSection("COOKIE INFO", report.cookie_info);

  return pdf.output("blob");
};
 const sendEmailWithPDF = async (pdfBlob) => {
  try {
    const formData = new FormData();

    formData.append(
      "pdf",
      pdfBlob,
      "website-report.pdf"
    );

    formData.append("email", "amitpunase4@gmail.com");

    await fetch("http://localhost:8086/send-report", {
      method: "POST",
      body: formData,
    });

    console.log("✅ Email sent automatically");
  } catch (err) {
    console.error("❌ Email error:", err);
  }
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

        <Section
          title="CONTENT INFO"
          data={report.content_info}
          setSelectedImage={setSelectedImage}
        />
        <Section
          title="PERFORMANCE"
          data={report.page_performance}
          setSelectedImage={setSelectedImage}
        />
        <Section
          title="SECURITY RISK"
          data={report.risk_infor}
          setSelectedImage={setSelectedImage}
        />
        <Section
          title="SEO ANALYSIS"
          data={report.seo_performance}
          setSelectedImage={setSelectedImage}
        />
        <Section
          title="NETWORK INFO"
          data={report.network_info}
          setSelectedImage={setSelectedImage}
        />
        <Section
          title="SOCIAL LINKS"
          data={report.sociallink_info}
          setSelectedImage={setSelectedImage}
        />
        <Section
          title="COOKIE INFO"
          data={report.cookie_info}
          setSelectedImage={setSelectedImage}
        />
      </div>

      {/* ✅ IMAGE MODAL */}
      {selectedImage && (
        <div
          className="image-modal"
          onClick={() => setSelectedImage(null)}
        >
          <img
            src={selectedImage}
            alt="Preview"
            className="modal-image"
            onClick={(e) => e.stopPropagation()}
          />
        </div>
      )}
    </div>
  );
}

/* ================= SECTION ================= */

function Section({ title, data, setSelectedImage }) {
  return (
    <div className="section">
      <h2>{title}</h2>

      <table>
        <tbody>
          {Object.entries(data || {}).map(([key, value]) => (
            <tr key={key}>
              <td className="key">{key}</td>

              <td className="value">
                {value === "" ||
                value === null ||
                value === undefined ? (
                  "-"
                ) : key === "broken_images_details" &&
                  Array.isArray(value) ? (
                  value.length === 0 ? (
                    "No Broken Images"
                  ) : (
                    value.map((item, index) => (
                      <div
                        key={index}
                        className="broken-image-card"
                      >
                        <p>
                          <strong>Page URL:</strong>
                        </p>
                        <a
                          href={item.page_url}
                          target="_blank"
                          rel="noreferrer"
                        >
                          {item.page_url}
                        </a>

                        <p style={{ marginTop: "10px" }}>
                          <strong>Image URL:</strong>
                        </p>
                        <a
                          href={item.image_url}
                          target="_blank"
                          rel="noreferrer"
                        >
                          {item.image_url}
                        </a>

                        <div style={{ marginTop: "10px" }}>
                          <img
                            src={`http://localhost:8086/${item.screenshot_path}`}
                            alt="Broken Screenshot"
                            className="broken-image-preview"
                            onClick={() =>
                              setSelectedImage(
                                `http://localhost:8086/${item.screenshot_path}`
                              )
                            }
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
                      <div
                        key={index}
                        className="broken-link-card"
                      >
                        <p>
                          <strong>Page URL:</strong>
                        </p>
                        <a
                          href={item.page_url}
                          target="_blank"
                          rel="noreferrer"
                        >
                          {item.page_url}
                        </a>

                        <p style={{ marginTop: "10px" }}>
                          <strong>Broken URL:</strong>
                        </p>
                        <a
                          href={item.url}
                          target="_blank"
                          rel="noreferrer"
                        >
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