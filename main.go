package main

import (
	"fmt"
	"net/http"

	"godice/roller"
)

func main() {
	// Serve CSS files
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("public/css"))))

	// Serve HTML files
	http.Handle("/", http.FileServer(http.Dir("public/html")))

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

		fmt.Fprintf(w, "<div class='p-4 mt-4 bg-gray-100 rounded'>Grand Total: %d</div>", result.GrandTotal)
	})

	http.ListenAndServe(":8080", nil)
}
