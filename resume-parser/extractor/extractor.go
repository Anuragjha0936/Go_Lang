package extractor



import (
    "bytes"
    "fmt"
    "strings"
    "github.com/ledongthuc/pdf"
)

func Extract(filename string, data []byte) (string, error) {
    if !strings.HasSuffix(filename, ".pdf") {
        return "", fmt.Errorf("unsupported file type: only PDF is supported")
    }
    return extractPDF(data)
}

func extractPDF(data []byte) (string, error) {
    r, err := pdf.NewReader(bytes.NewReader(data), int64(len(data)))
    if err != nil {
        return "", fmt.Errorf("failed to open PDF: %w", err)
    }

    var sb strings.Builder
    for i := 1; i <= r.NumPage(); i++ {
        page := r.Page(i)
        if page.V.IsNull() {
            continue
        }
        text, err := page.GetPlainText(nil)
        if err == nil {
            sb.WriteString(text)
            sb.WriteString("\n")
        }
    }

    if sb.Len() == 0 {
        return "", fmt.Errorf("no text found in PDF (possibly scanned/image-based)")
    }

    return sb.String(), nil
}