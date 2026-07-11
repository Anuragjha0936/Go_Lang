# AI Resume Parser

An AI-powered Resume Parser built with **Go**, **Gin**, and **NVIDIA Llama 3.1**. Upload a PDF resume and receive structured candidate information including personal details, skills, work experience, education, and certifications.

## Features

- 📄 Upload PDF resumes
- 🤖 AI-powered resume parsing using NVIDIA Llama 3.1
- 📑 Extract text from PDF documents
- 📋 Structured JSON output
- 🎨 Simple web interface
- 🌐 REST API
- ⚡ Fast and lightweight Go backend

---

## Tech Stack

### Backend
- Go
- Gin Web Framework
- NVIDIA AI API (Llama 3.1 8B Instruct)

### PDF Processing
- ledongthuc/pdf

### Frontend
- HTML
- CSS
- JavaScript

---

## Project Structure

```
resume-parser/
│
├── extractor/
│   └── extractor.go
├── parser/
│   └── parser.go
├── models/
│   └── resume.go
├── main.go
├── index.html
├── go.mod
├── go.sum
└── README.md
```

---

## How It Works

1. Upload a PDF resume.
2. The backend extracts text from the PDF.
3. The extracted text is sent to NVIDIA's Llama model.
4. The AI converts the resume into structured JSON.
5. The parsed information is displayed in the web interface.

---

## Extracted Information

The parser extracts:

- Name
- Email
- Phone Number
- Location
- Professional Summary
- Skills
- Work Experience
- Education
- Certifications

---

## Installation

### Clone the repository

```bash
git clone https://github.com/<your-username>/resume-parser.git

cd resume-parser
```

### Install dependencies

```bash
go mod tidy
```

### Configure Environment Variables

Create a `.env` file or export the following environment variable:

```env
NVIDIA_API_KEY=your_api_key
```

---

## Run the Application

```bash
go run main.go
```

The server starts on:

```
http://localhost:8080
```

---

## API Endpoints

### Health Check

```http
GET /health
```

Response

```json
{
  "status": "ok"
}
```

---

### Parse Resume

```http
POST /parse
```

Request

```
multipart/form-data

file=<resume.pdf>
```

Example Response

```json
{
  "status": "success",
  "data": {
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "+91 9876543210",
    "location": "Bangalore",
    "summary": "...",
    "skills": [
      "Go",
      "Docker",
      "MySQL"
    ],
    "experience": [
      {
        "company": "ABC Pvt Ltd",
        "title": "Backend Developer",
        "dates": "2023-Present"
      }
    ],
    "education": [
      {
        "institution": "XYZ University",
        "degree": "B.Tech",
        "year": "2025"
      }
    ],
    "certifications": [
      "AWS Cloud Practitioner"
    ]
  }
}
```

---

## Screenshots

Add screenshots of:

- Upload Page
- Parsed Resume Output

---

## Future Improvements

- Support DOCX resumes
- OCR support for scanned PDFs
- Resume scoring
- ATS compatibility score
- Export parsed data as JSON or CSV
- Authentication for API
- Docker support
- Batch resume parsing

---

## Author

**Anurag Jha**

- GitHub: https://github.com/<your-username>
- LinkedIn: https://linkedin.com/in/<your-linkedin>

---

## License

This project is licensed under the MIT License.
