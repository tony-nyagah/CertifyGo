# CertifyGo

A simple web-based application that generates participation certificates with customizable templates.

## What it does

- Creates professional-looking participation certificates
- Supports multiple template styles (Classic, Modern, Elegant)
- Generates printable HTML certificates
- Web form interface for easy certificate creation

## Tech Stack

- **Backend**: Go 1.23.2
- **Frontend**: HTML/CSS (embedded templates)
- **Templates**: Go html/template package
- **Server**: Go net/http

## How to run

1. Make sure you have Go installed (version 1.23.2 or later)

2. Clone and navigate to the project:
   ```bash
   git clone <repo-url>
   cd certificate-generator
   ```

3. Run the application:
   ```bash
   go run main.go
   ```

4. Open your browser and go to: http://localhost:8080

5. Fill in the form and generate your certificate!

## Usage

1. Enter event name, dates, and participant name
2. Choose a template style
3. Click "Generate Certificate"
4. Print or save the certificate using your browser

That's it! CertifyGo will automatically create template files on first run.

## Adding New Templates

To add a custom template:

1. Create a new `.html` file in the `templates/` directory (e.g., `templates/corporate.html`)
2. Use these template variables in your HTML:
   - `{{.EventName}}` - Event or course name
   - `{{.ParticipantName}}` - Participant's name  
   - `{{.StartDate}}` - Formatted start date
   - `{{.EndDate}}` - Formatted end date
   - `{{.LogoURL}}` - Logo image URL (optional)
   - `{{.GeneratedDate}}` - Current date when certificate was generated
3. Rebuild the application: `go run main.go`
4. Your new template will appear in the dropdown automatically

Example template structure:
```html
<!DOCTYPE html>
<html>
<head>
    <title>Certificate</title>
    <style>
        /* Your custom styles here */
    </style>
</head>
<body>
    <div class="certificate">
        {{if .LogoURL}}<img src="{{.LogoURL}}" alt="Logo" class="logo">{{end}}
        <h1>{{.EventName}}</h1>
        <p>{{.ParticipantName}}</p>
        <p>{{.StartDate}} - {{.EndDate}}</p>
    </div>
</body>
</html>
```

## Generate Templates with AI

You can use any LLM (ChatGPT, Claude, Gemini, etc.) to generate custom certificate templates. Just copy-paste this prompt:

### LLM Prompt for Template Generation:

```
Create an HTML certificate template for [YOUR COMPANY NAME] with the following requirements:

TEMPLATE VARIABLES (must include these exact placeholders):
- {{.EventName}} - Event or course name
- {{.ParticipantName}} - Participant's name  
- {{.StartDate}} - Formatted start date
- {{.EndDate}} - Formatted end date
- {{.LogoURL}} - Logo image URL (optional, use conditional: {{if .LogoURL}}<img src="{{.LogoURL}}" alt="Logo">{{end}})
- {{.GeneratedDate}} - Current date when certificate was generated

DESIGN REQUIREMENTS:
- Complete HTML document with embedded CSS (no external stylesheets)
- Certificate size: 800px width x 600px height
- Print-friendly design with good contrast
- Professional appearance suitable for [SPECIFY YOUR THEME: corporate/academic/creative/etc.]
- Responsive layout that looks good on screen and when printed
- Include proper margins and padding for printing

STYLE INSPIRATION:
Create a [SPECIFY STYLE: modern/classic/elegant/minimalist/colorful] certificate with [SPECIFY ELEMENTS: borders/gradients/ornaments/etc.]

EXAMPLE STRUCTURE (modify the styling but keep the template variables):
<!DOCTYPE html>
<html>
<head>
    <title>Certificate</title>
    <style>
        /* Your CSS here */
        .certificate { width: 800px; height: 600px; /* other styles */ }
    </style>
</head>
<body>
    <div class="certificate">
        {{if .LogoURL}}<img src="{{.LogoURL}}" alt="Logo" class="logo">{{end}}
        <h1>{{.EventName}}</h1>
        <p>{{.ParticipantName}}</p>
        <p>{{.StartDate}} - {{.EndDate}}</p>
        <small>{{.GeneratedDate}}</small>
    </div>
</body>
</html>

Please generate a complete, ready-to-use HTML template.
```

### How to use:
1. Copy the prompt above
2. Customize the bracketed sections [SPECIFY YOUR THEME], [SPECIFY STYLE], etc.
3. Paste into your preferred LLM
4. Save the generated HTML as `templates/yourname.html`
5. Rebuild CertifyGo: `go run main.go`

### Logo Support:
- Add a "Logo URL (optional)" field to include organization logos
- Logos work with any publicly accessible image URL (PNG, JPG, SVG)
- Logos are automatically positioned and sized in each template
- Use conditional rendering: `{{if .LogoURL}}<img src="{{.LogoURL}}" alt="Logo" class="logo">{{end}}`

### Pro tip: 
Include one of the existing templates (classic/modern/elegant) as reference by adding "Here's an example template for reference:" and pasting the HTML.
