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
	EventName            string
	StartDate            string
	EndDate              string
	ParticipantName      string
	LogoURL              string
	GeneratedDate        string
	CertificateID        string
	CertificateType      string // e.g. Presenter, Sponsor, Participant
	RecognitionStatement string
	ReceiverRole         string
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
            <label>Certificate Type:</label>
            <select name="participant_type">
                <optgroup label="Participation">
                    <option value="attendee">Attendee</option>
                    <option value="delegate">Delegate</option>
                    <option value="participant" selected>Participant</option>
                </optgroup>
                <optgroup label="Contribution">
                    <option value="speaker">Speaker</option>
                    <option value="presenter">Presenter</option>
                    <option value="panelist">Panelist</option>
                    <option value="moderator">Moderator</option>
                    <option value="instructor">Instructor</option>
                    <option value="facilitator">Workshop Facilitator</option>
                </optgroup>
                <optgroup label="Support">
                    <option value="sponsor">Sponsor</option>
                    <option value="partner">Partner</option>
                    <option value="exhibitor">Exhibitor</option>
                    <option value="patron">Patron</option>
                </optgroup>
                <optgroup label="Organization">
                    <option value="organizer">Organizer</option>
                    <option value="volunteer">Volunteer</option>
                    <option value="committee">Committee Member</option>
                    <option value="judge">Judge</option>
                    <option value="chair">Chair</option>
                </optgroup>
            </select>
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

func getReceiverRole(certificateType string) string {
	switch strings.ToLower(certificateType) {
	// Contribution categories
	case "speaker":
		return "Speaker"
	case "presenter":
		return "Presenter"
	case "panelist":
		return "Panelist"
	case "moderator":
		return "Moderator"
	case "instructor":
		return "Instructor"
	case "facilitator":
		return "Workshop Facilitator"

	// Support categories
	case "sponsor":
		return "Sponsor"
	case "partner":
		return "Partner Organization"
	case "exhibitor":
		return "Exhibitor"
	case "patron":
		return "Patron"

	// Organizing categories
	case "organizer":
		return "Event Organizer"
	case "volunteer":
		return "Volunteer"
	case "committee":
		return "Committee Member"
	case "judge":
		return "Judge"
	case "chair":
		return "Chair"

	// Participation categories
	case "attendee":
		return "Attendee"
	case "delegate":
		return "Delegate"
	case "participant":
		return "Participant"

	default:
		return "Participant"
	}
}

func getRecognitionStatement(certificateType string) string {
	switch strings.ToLower(certificateType) {
	// Contribution categories
	case "speaker":
		return "for delivering an exceptional presentation and sharing valuable insights in the"
	case "presenter":
		return "for presenting important work and contributing to knowledge exchange in the"
	case "panelist":
		return "for providing expert perspectives as a discussion panelist in the"
	case "moderator":
		return "for skillfully moderating discussions and facilitating engagement in the"
	case "instructor":
		return "for instructing an educational session with expertise and dedication in the"
	case "facilitator":
		return "for facilitating an engaging and productive workshop in the"

	// Support categories
	case "sponsor":
		return "for generously supporting and making this event possible"
	case "partner":
		return "for valuable partnership and collaboration in delivering this event"
	case "exhibitor":
		return "for participating as an exhibitor and enriching the event experience"
	case "patron":
		return "for distinguished patronage and support of this event"

	// Organizing categories
	case "organizer":
		return "for outstanding contributions to organizing and executing this event"
	case "volunteer":
		return "for dedicated volunteer service that helped make this event successful"
	case "committee":
		return "for diligent service as a committee member"
	case "judge":
		return "for providing expert evaluation and judgment"
	case "chair":
		return "for leadership and guidance as a chair"

	// Participation categories
	case "attendee":
		return "for active attendance and participation"
	case "delegate":
		return "for representing their organization as an official delegate"
	case "participant":
		return "for meaningful participation and engagement"

	default:
		return "for participating in"
	}
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data := CertificateData{
		EventName:            r.FormValue("event_name"),
		StartDate:            formatDate(r.FormValue("start_date")),
		EndDate:              formatDate(r.FormValue("end_date")),
		ParticipantName:      r.FormValue("participant_name"),
		LogoURL:              r.FormValue("logo_url"),
		GeneratedDate:        time.Now().Format("January 2, 2006"),
		CertificateID:        generateCertificateID(),
		CertificateType:      r.FormValue("participant_type"),
		ReceiverRole:         getReceiverRole(r.FormValue("participant_type")),
		RecognitionStatement: getRecognitionStatement(r.FormValue("participant_type")),
	}

	templateName := r.FormValue("template")
	if templateName == "" {
		templateName = "classic"
	}

	// Ensure default value if not specified
	if data.CertificateType == "" {
		data.CertificateType = "participant"
		data.ReceiverRole = "Participant"
		data.RecognitionStatement = getRecognitionStatement("participant")
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
` + certificate + `
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
