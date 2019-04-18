package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/rms1000watt/golang-integration-test/person"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	Port string
	Addr string

	PGHost       string
	PGPort       string
	PGUser       string
	PGPass       string
	PGDB         string
	PGConnString string
}

type Server struct {
	DB *sqlx.DB
}

func (s *Server) HandlerPerson(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.HandlerPersonGET(w, r)
	case http.MethodPost:
		s.HandlerPersonPOST(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (s *Server) HandlerPersonGET(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if strings.TrimSpace(name) == "" {
		fmt.Println("ERROR: No name provided")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var p person.Person
	if err := s.DB.Get(&p, "SELECT * FROM person WHERE name=$1 LIMIT 1;", name); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Person not found: " + name)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		fmt.Println("Failed running query:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.Marshal(p)
	if err != nil {
		fmt.Println("Failed json marshal:", err)
		fmt.Println("GET:", p)
		fmt.Fprintln(w, p)
		return
	}

	jsonString := string(jsonBytes)
	fmt.Println("GET:", jsonString)
	fmt.Fprintln(w, jsonString)
}

func (s *Server) HandlerPersonPOST(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if strings.TrimSpace(name) == "" {
		fmt.Println("ERROR: No name provided")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	age, err := strconv.Atoi(r.URL.Query().Get("age"))
	if err != nil {
		fmt.Println("Failed converting Atoi:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	p := &person.Person{
		Name: name,
		Age:  age,
	}
	res, err := s.DB.NamedExec("INSERT INTO person (name, age) VALUES (:name, :age);", p)
	if err != nil {
		fmt.Println("Failed NamedExec:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	rowsCnt, err := res.RowsAffected()
	if err != nil {
		fmt.Println("Failed inserting:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if rowsCnt == 0 {
		fmt.Println("ERROR: zero rows inserted")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fmt.Println("POST:", p)
	fmt.Fprintln(w, http.StatusText(http.StatusOK))
}

func getenvDefault(key, defaultValue string) (value string) {
	value = os.Getenv(key)

	if strings.TrimSpace(value) == "" {
		value = defaultValue
	}

	return
}

func getConfig() (config Config, err error) {
	config.PGHost = getenvDefault("PERSON_SVC_PG_HOST", "postgres")
	config.PGPort = getenvDefault("PERSON_SVC_PG_PORT", "5432")
	config.PGUser = getenvDefault("PERSON_SVC_PG_USER", "postgres")
	config.PGPass = getenvDefault("PERSON_SVC_PG_PASS", "")
	config.PGDB = getenvDefault("PERSON_SVC_PG_DB", "person")

	passwordStr := "password=" + config.PGPass
	if strings.TrimSpace(config.PGPass) == "" {
		passwordStr = ""
	}

	config.PGConnString = fmt.Sprintf(
		"host=%s port=%s user=%s %s dbname=%s sslmode=disable",
		config.PGHost,
		config.PGPort,
		config.PGUser,
		passwordStr,
		config.PGDB)

	config.Port = getenvDefault("PERSON_SVC_PORT", "9999")
	config.Addr = ":" + config.Port

	return
}

func main() {
	config, err := getConfig()
	if err != nil {
		fmt.Println("Failed getting config:", err)
		return
	}

	fmt.Println("Connecting to Postgres:", config.PGConnString)
	db, err := sqlx.Connect("postgres", config.PGConnString)
	if err != nil {
		fmt.Println("Unable to connect to db:", err)
		return
	}

	s := Server{
		DB: db,
	}

	fmt.Println("Starting server on", config.Addr)

	http.HandleFunc("/person", s.HandlerPerson)
	log.Fatalln(http.ListenAndServe(config.Addr, nil))
}
