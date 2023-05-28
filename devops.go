package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"unicode"
)

var db *sql.DB

// validateUser validates that a user provided only contains letters.
func validateUser(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// convertStringToDate converts a date provided as string to a valid date format using date format YYYY-MM-DD.
//
// It has one parameter: the date in a string format and it returns the equivalent in a parameter of type time.Time along with an error (if any).
func convertStringToDate(date string) (time.Time, error) {
	dateFormat := "2006-01-02"
	d, err := time.Parse(dateFormat, date)
	if err != nil {
		return time.Time{}, err
	}
	return d, nil
}

// validateDate validates that the date of birth provided is not after today.
//
// It has one parameter: a string representing the date obtained after unmarshalling the JSON object provided by client along with an error (if any).
func validateDate(date string) error {
	d, err := convertStringToDate(date)
	if err != nil {
		return err
	}
	today := time.Now().UTC()
	if d.Before(today) {
		return nil
	} else {
		return errors.New("The date is not before today")
	}
}

// getUser extracts the username from the URI
//
// It has one parameter: an http.Request object from which the username will be obtained from.
// It returns the username after being extracted along with an error (if any)
func getUser(r *http.Request) (string, error) {
	person := r.URL.Path[len("/user/"):]
	person = person[1:]
	if validateUser(person) {
		return person, nil
	} else {
		return "", errors.New("Username provided does not contain only letters")
	}
}

// User represents the combination of username and its date of birth provided through a PUT request to the API
type User struct {
	ID          int64
	Username    string `json:"username"`
	DateofBirth string `json:"dateOfBirth"`
}

// helloBirthday is the HTTP Handler for the /hello endpoint
func helloBirthday(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		var user User
		var ErrBadUserName, ErrBadBirthdayDate = false, false
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// validating the date:
		err = validateDate(user.DateofBirth)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			ErrBadBirthdayDate = true
		}

		u, err := getUser(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			ErrBadUserName = true
		}
		if (!ErrBadUserName) && (!ErrBadBirthdayDate) {
			err = addUser(u, user.DateofBirth)
			if err != nil {
				fmt.Println("User data could not be saved")
				fmt.Println("Error: ", err)
				os.Exit(2)
			}
			w.WriteHeader(http.StatusNoContent)
		}

		// If everything goes fine, then return an HTTP code 204 No Content
	} else if r.Method == http.MethodGet {
		username, err := getUser(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Check if username exists in the DB and if so retrieve the row with its data
		user, err := retrieveUser(username)
		if err != nil {
			fmt.Printf("%+v", err)
			return
		}
		today := time.Now().UTC().Truncate(24 * time.Hour)

		// Check if user's birthday is today
		birthDate, err := convertStringToDate(user.DateofBirth)
		if err != nil {
			fmt.Printf("%+v", err)
			return
		}

		birthDate = birthDate.UTC().Truncate(24 * time.Hour)
		nextBirthday := time.Date(today.Year(), birthDate.Month(), birthDate.Day(), 0, 0, 0, 0, time.UTC)
		var message string

		if birthDate.Day() == today.Day() && birthDate.Month() == today.Month() {
			message = "Hello, " + username + "! Happy birthday!"
			sendJSONResponse(message, w)
			return
		}
		if today.After(nextBirthday) {
			nextBirthday = nextBirthday.AddDate(1, 0, 0)

			// Otherwise let user know that his/her birthday is in N days
		}
		days := int(nextBirthday.Sub(today).Hours() / 24)
		message = "Hello, " + username + "! Your birhtday is in " + strconv.Itoa(days) + " day(s)"

		sendJSONResponse(message, w)
	}
}

// sendJSONResponse sends the reply to the client in JSON format.
//
// It has two parameters: the message string that will be sent as a response to the client, and an instance of a http.ResponseWriter for being able to write back to the client's request.
func sendJSONResponse(message string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)

}

// retrieveUser checks if username provided already exists when user hits the /hello/<user> endpoint.
//
// It has one parameter: the username that will be checked against the database.
// It has two return values: a User instance which will contain the data fetched from the database results and an error instance in case there is any.
func retrieveUser(username string) (User, error) {
	var user User
	row := db.QueryRow("SELECT * FROM users WHERE username = ?", username)
	switch err := row.Scan(&user.ID, &user.Username, &user.DateofBirth); err {
	case sql.ErrNoRows:
		fmt.Println("No user ", username, "found")
		return user, fmt.Errorf("No user %s, found %w", username, sql.ErrNoRows)
	case nil:
		return user, nil
	default:
		panic(err)
	}
}

// addUser adds the username to the database.
//
// It has two parameters: the username as a string variable and dateOfBirth also as a string value which both will be inserted into the database.
// It returns an error (if any).
func addUser(user string, dateOfBirth string) error {
	result, err := db.Exec("INSERT INTO users (username, dateOfBirth) VALUES (?, ?)", user, dateOfBirth)
	if err != nil {
		return fmt.Errorf("addUser: %v", err)
	}
	_, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("addUser: %v", err)
	}
	return nil
}

func main() {
	// Capture connection properties.
	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 os.Getenv("DBHOST"),
		DBName:               "users",
		AllowNativePasswords: true,
	}
	// Get a database handle.
	fmt.Println("DB HOST: ", os.Getenv("DBHOST"), "DB user: ", os.Getenv("DBUSER"), "DB PASS: ", os.Getenv("DBPASS"))
	fmt.Println(cfg.FormatDSN())
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	// If database connection is successful, then print to the stdout "Connected".
	fmt.Println("Connected!")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/hello/", helloBirthday)
	http.ListenAndServe(":"+port, mux)
}
