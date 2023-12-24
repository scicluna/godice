package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"godice/roller"
	"log"
	"strings"

	"github.com/google/uuid"
)

const (
	host     = "db"       // This should match the service name in docker-compose.yml
	port     = 5432       // Default PostgreSQL port
	user     = "user"     // Your DB username
	password = "password" // Your DB password
	dbname   = "godice"   // Your DB name
)

func connectDB() *sql.DB {
	// Construct the connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open a connection to the database
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}

	// Check the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v\n", err)
	}

	fmt.Println("Successfully connected to database!")
	return db
}

func ensureDefaultProfileExists(db *sql.DB) {
	// SQL query to insert the Default profile if it doesn't exist
	query := `
        INSERT INTO profiles (name)
        SELECT 'Default'
        WHERE NOT EXISTS (
            SELECT id FROM profiles WHERE name = 'Default'
        );
    `

	// Execute the query
	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create default profile: %v", err)
	}
}

func getProfiles(db *sql.DB) ([]string, error) {

	var profiles []string
	// SQL query to get all current profiles
	query := `
            SELECT * FROM profiles;
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var profile string
		if err := rows.Scan(&profile); err != nil {
			return nil, err
		}

		profiles = append(profiles, profile)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return profiles, nil
}

func getDefaultProfileID(db *sql.DB) (int, error) {
	var id int
	query := `SELECT id FROM profiles WHERE name = 'Default';`
	err := db.QueryRow(query).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// SaveRollResult saves a dice roll result to the database under a given profile.
func SaveRollResult(db *sql.DB, profileName string, result *roller.RollResult) error {
	// Serialize the result into JSON
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return err
	}

	// Get the ID of the specified profile
	profileID, err := getProfileID(db, profileName)
	if err != nil {
		return err
	}

	// Prepare the SQL INSERT statement
	insertQuery := `INSERT INTO roll_results (profile_id, roll_data, created_at) VALUES ($1, $2, NOW())`

	// Execute the insert query
	_, err = db.Exec(insertQuery, profileID, resultJSON)
	if err != nil {
		return err
	}

	return nil
}

// getProfileID retrieves the ID for a given profile name.
func getProfileID(db *sql.DB, profileName string) (int, error) {
	var id int
	query := `SELECT id FROM profiles WHERE name = $1;`
	err := db.QueryRow(query, profileName).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func getRollResultsForProfile(db *sql.DB, profileName string) ([]roller.RollResult, error) {
	var results []roller.RollResult

	// SQL query to select roll results for a given profile
	query := `SELECT roll_data FROM roll_results WHERE profile_id = (SELECT id FROM profiles WHERE name = $1);`

	rows, err := db.Query(query, profileName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rollData string
		if err := rows.Scan(&rollData); err != nil {
			return nil, err
		}

		var result roller.RollResult
		if err := json.Unmarshal([]byte(rollData), &result); err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func convertResultsToHTML(results []roller.RollResult) string {
	var htmlBuilder strings.Builder

	for _, result := range results {
		// Build HTML response for each result
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
		sb.WriteString(fmt.Sprintf("<div>Total: %d</div>", result.GrandTotal))

		sb.WriteString("</li>")

		// Append the result's HTML to the main HTML builder
		htmlBuilder.WriteString(sb.String())
	}

	return htmlBuilder.String()
}

func convertProfilesToHTML(profiles []string) string {
	var htmlBuilder strings.Builder

	for _, profile := range profiles {
		// Build HTML response for each result
		var sb strings.Builder
		sb.WriteString("<option class='flex flex gap-1 text-xl font-bold font-mono'>")

		sb.WriteString(fmt.Sprintf("%s", profile))

		sb.WriteString("</option>")

		// Append the result's HTML to the main HTML builder
		htmlBuilder.WriteString(sb.String())

	}

	return htmlBuilder.String()
}
