package main

import (
	"fmt"
	"godice/roller"
	"html/template"
	"net/http"
	"strings"

	_ "github.com/lib/pq"

	"github.com/google/uuid"
)

func main() {
	// Initialize the database connection
	db := connectDB()
	defer db.Close()

	// Ensure the default profile exists
	ensureDefaultProfileExists(db)

	// Serve CSS files
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("public/css"))))

	// Serve images
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("public/images"))))

	// Serve HTML files
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		profiles, profileErr := getProfiles(db)
		if profileErr != nil {
			http.Error(w, profileErr.Error(), http.StatusInternalServerError)
			return
		}

		rollResults, rollErr := getRollResultsForProfile(db, "Default")
		if rollErr != nil {
			http.Error(w, rollErr.Error(), http.StatusInternalServerError)
			return
		}

		rollContent := convertResultsToHTML(rollResults)
		profileContent := convertProfilesToHTML(profiles)

		// Parse the HTML template
		tmpl, tmplErr := template.ParseFiles("public/html/index.html")
		if tmplErr != nil {
			http.Error(w, tmplErr.Error(), http.StatusInternalServerError)
			return
		}

		// Execute the template with htmlContent
		data := struct {
			HtmlContent template.HTML
			Profiles    template.HTML
		}{
			HtmlContent: template.HTML(rollContent),
			Profiles:    template.HTML(profileContent),
		}

		err := tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/roll", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		diceString := r.FormValue("diceString")
		result, err := roller.RollDiceString(diceString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Call the SaveRollResult function
		saveerr := SaveRollResult(db, "Default", result)
		if saveerr != nil {
			http.Error(w, "Error saving roll result", http.StatusInternalServerError)
			return
		}

		// Build HTML response
		var sb strings.Builder
		sb.WriteString("<li class='flex flex gap-1 text-xl font-bold font-mono'>")

		for i, set := range result.Sets {
			newUUID := uuid.New().ID()
			sb.WriteString(fmt.Sprintf("<div class='cursor-pointer' id='%d' onclick='expandDice(this)'>%d</div>", newUUID, set.Total))

			sb.WriteString(fmt.Sprintf("<ul class='flex hidden cursor-pointer' onclick='collapseDice(this, %d)'>", newUUID))
			for _, roll := range set.Rolls {
				sb.WriteString(fmt.Sprintf("<li class='%s'>%d</li>", roll.RollType, roll.Value))
			}
			sb.WriteString("</ul>")

			if i < len(result.Operands) && result.Operands[i] != "" {
				sb.WriteString(fmt.Sprintf("<div>%s</div>", result.Operands[i]))
			}

			if i == len(result.Sets)-1 {
				sb.WriteString(fmt.Sprintf("<div>%s</div>", "="))
			}
		}

		// Append grand total
		sb.WriteString(fmt.Sprintf("%d", result.GrandTotal))

		sb.WriteString("</li>")
		fmt.Fprint(w, sb.String())
	})

	http.ListenAndServe(":8080", nil)
}
