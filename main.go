package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

//go:embed templates/*.html
var templateFS embed.FS

type CertificateData struct {
	EventName       string
	StartDate       string
	EndDate         string
	ParticipantName string
	LogoURL         string
	GeneratedDate   string
	CertificateID   string
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	templates := getAvailableTemplates()

	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Certificate Generator</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 600px; margin: 50px auto; padding: 20px; }
        .form-group { margin: 15px 0; }
        label { display: block; margin-bottom: 5px; font-weight: bold; }
        input, select { width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; }
        button { background: #007cba; color: white; padding: 10px 20px; border: none; border-radius: 4px; cursor: pointer; }
        button:hover { background: #005a87; }
    </style>
</head>
<body>
    <h1>Certificate Generator</h1>
    <form action="/generate" method="post" target="_blank">
        <div class="form-group">
            <label>Event Name:</label>
            <input type="text" name="event_name" required>
        </div>
        <div class="form-group">
            <label>Start Date:</label>
            <input type="date" name="start_date" required>
        </div>
        <div class="form-group">
            <label>End Date:</label>
            <input type="date" name="end_date" required>
        </div>
        <div class="form-group">
            <label>Participant Name:</label>
            <input type="text" name="participant_name" required>
        </div>
        <div class="form-group">
            <label>Logo URL (optional):</label>
            <input type="url" name="logo_url" placeholder="https://example.com/logo.png">
        </div>
        <div class="form-group">
            <label>Template:</label>
            <select name="template">`

	for _, tmpl := range templates {
		html += `<option value="` + tmpl + `">` + strings.Title(tmpl) + `</option>`
	}

	html += `</select>
        </div>
        <button type="submit">Generate Certificate</button>
    </form>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data := CertificateData{
		EventName:       r.FormValue("event_name"),
		StartDate:       formatDate(r.FormValue("start_date")),
		EndDate:         formatDate(r.FormValue("end_date")),
		ParticipantName: r.FormValue("participant_name"),
		LogoURL:         r.FormValue("logo_url"),
		GeneratedDate:   time.Now().Format("January 2, 2006"),
		CertificateID:   generateCertificateID(),
	}

	templateName := r.FormValue("template")
	if templateName == "" {
		templateName = "classic"
	}

	tmpl := getTemplate(templateName)

	// Generate the certificate HTML
	var htmlContent strings.Builder
	tmpl.Execute(&htmlContent, data)
	certificate := htmlContent.String()

	// Add download buttons to the certificate
	downloadHTML := `
<!DOCTYPE html>
<html>
<head>
    <title>Certificate Generated</title>
    <style>
        .download-bar {
            position: fixed;
            top: 0;
            left: 0;
            right: 0;
            background: #333;
            color: white;
            padding: 10px;
            text-align: center;
            z-index: 1000;
        }
        .download-btn {
            background: #007cba;
            color: white;
            padding: 8px 16px;
            border: none;
            border-radius: 4px;
            margin: 0 5px;
            cursor: pointer;
            text-decoration: none;
            display: inline-block;
        }
        .download-btn:hover { background: #005a87; }
        .certificate-container { margin-top: 60px; }
        @media print { .download-bar { display: none; } .certificate-container { margin-top: 0; } }
    </style>
</head>
<body>
    <div class="download-bar">
        <button class="download-btn" onclick="window.print()">Print Certificate</button>
    </div>
    <div class="certificate-container">
` + certificate + `
    </div>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(downloadHTML))
}

func formatDate(dateStr string) string {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return dateStr
	}
	return t.Format("January 2, 2006")
}

func getTemplate(name string) *template.Template {
	templatePath := "templates/" + name + ".html"

	content, err := templateFS.ReadFile(templatePath)
	if err != nil {
		log.Printf("Template %s not found, falling back to classic: %v", name, err)
		// Fall back to classic template if file doesn't exist
		content, err = templateFS.ReadFile("templates/classic.html")
		if err != nil {
			log.Fatalf("Critical error: classic template not found: %v", err)
		}
	}

	return template.Must(template.New(name).Parse(string(content)))
}

func getAvailableTemplates() []string {
	var templates []string

	entries, err := templateFS.ReadDir("templates")
	if err != nil {
		log.Printf("Error reading embedded templates: %v", err)
		return []string{"classic"}
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".html") {
			name := strings.TrimSuffix(entry.Name(), ".html")
			templates = append(templates, name)
		}
	}

	if len(templates) == 0 {
		return []string{"classic"}
	}

	return templates
}



func downloadHandler(w http.ResponseWriter, r *http.Request) {
	data := CertificateData{
		EventName:       r.URL.Query().Get("event_name"),
		StartDate:       formatDate(r.URL.Query().Get("start_date")),
		EndDate:         formatDate(r.URL.Query().Get("end_date")),
		ParticipantName: r.URL.Query().Get("participant_name"),
		LogoURL:         r.URL.Query().Get("logo_url"),
		GeneratedDate:   time.Now().Format("January 2, 2006"),
		CertificateID:   generateCertificateID(),
	}

	templateName := r.URL.Query().Get("template")
	if templateName == "" {
		templateName = "classic"
	}

	tmpl := getTemplate(templateName)

	// Set download headers
	filename := strings.ReplaceAll(data.ParticipantName, " ", "_") + "_certificate.html"
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "text/html")

	tmpl.Execute(w, data)
}

func generateCertificateID() string {
	return "CERT-" + time.Now().Format("20060102") + "-" + time.Now().Format("150405")
}


func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/generate", generateHandler)
	http.HandleFunc("/download", downloadHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
