package models


type Experience struct {
    Company    string   `json:"company"`
    Title      string   `json:"title"`
    Dates      string   `json:"dates"`
    Location   string   `json:"location"`
    Highlights []string `json:"highlights"`
}

type Education struct {
    Institution string `json:"institution"`
    Degree      string `json:"degree"`
    Field       string `json:"field"`
    Year        string `json:"year"`
    GPA         string `json:"gpa"`
}

type Project struct {
    Name         string   `json:"name"`
    Description  string   `json:"description"`
    Technologies []string `json:"technologies"`
    URL          string   `json:"url"`
}

type Resume struct {
    Name           string       `json:"name"`
    Email          string       `json:"email"`
    Phone          string       `json:"phone"`
    Location       string       `json:"location"`
    LinkedIn       string       `json:"linkedin"`
    GitHub         string       `json:"github"`
    Summary        string       `json:"summary"`
    Skills         []string     `json:"skills"`
    Experience     []Experience `json:"experience"`
    Education      []Education  `json:"education"`
    Certifications []string     `json:"certifications"`
    Projects       []Project    `json:"projects"`
}