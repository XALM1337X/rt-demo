package server

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"database/sql"
	"github.com/gorilla/mux"
	"errors"
	"fmt"
	"os"
	"regexp"
)

var (
	re = regexp.MustCompile("sql: no rows")
)
type TFibParse struct {
	
	Nth_fib string `json:"nth_fib"`
	Nth_fib_result string  `json:"nth_fib_result"`
	CacheStatus string `json:"cache_status"`
}


func NewHTTPServer(addr string) *http.Server {
	r := mux.NewRouter()

	r.HandleFunc("/", EntryHandler).Methods("GET")
	r.HandleFunc("/fib_check", FibHandler).Methods("POST")
	r.PathPrefix("/www/").Handler(http.StripPrefix("/www/", http.FileServer(http.Dir("www"))))

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

func EntryHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("www/index.gohtml")
	if err != nil {
		w.Write([]byte("Error index.gohtml not found"))
		return
	}
	tmp := struct {
		Display string
	}{
		"",
	}
	w.Header().Set("Content-Type", "text/html")
	t.ExecuteTemplate(w, "index", tmp)
}

func FibHandler(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	w.Header().Set("Content-Type", "text/html")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	err_unmarsh := json.Unmarshal(body, &req)
	if err_unmarsh != nil {
		w.Write([]byte(err_unmarsh.Error()))
		return
	}
	if _, ok := req["lookup"]; !ok {
		w.Write([]byte("Error key not found in req map."))
		return
	}
	db, err2 := DbConnect()
	if err2 != nil {
		w.Write([]byte(fmt.Sprintf("Error: %s", err2.Error())))
	}
	defer db.Close()

	err3 := db.Ping()
	if err3 != nil {
		w.Write([]byte(fmt.Sprintf("Error: %s", err3.Error())))
		return
	} 

	//DB connected


	//Check if key(n'th fib number) exists.
	val, exists, err_cache := CheckCache(req["lookup"].(string), db)
	if err_cache != nil && !re.MatchString(err_cache.Error()){
		w.Write([]byte(err_cache.Error()))
		return
	}

	if exists {
		//if it does return it.
		w.Header().Set("Content-Type","application/json")
		w.Write(val)
	} else {
		//Calculate it
		//TODO

		//Store it
		insert_err := DbInsert(req["lookup"].(string), "1337", db)
		if insert_err != nil {
			w.Write([]byte(insert_err.Error()))
		} else {
			w.Write([]byte("Successfully added entry to db"))
		}
	}	
}


func CheckCache(lookup string, db *sql.DB) ([]byte, bool, error) {
	var fib TFibParse
	query := "SELECT fib_num, result FROM event WHERE fib_num='"+lookup+"';"
	err := db.QueryRow(query).Scan(&fib.Nth_fib, &fib.Nth_fib_result)
	if err != nil {
		return []byte{}, false, err
	}
	fib.CacheStatus = "Entry found"
	bytes, marsh_err := json.Marshal(fib)
	if marsh_err != nil {
		return []byte{}, false, errors.New(fmt.Sprintf("Error: %s", marsh_err.Error()))
	}
	return bytes, true, nil
}


func FibGenerate(nth_fib_str string) {

}


func DbInsert(key string, value string, db *sql.DB) error {
	sqlStatement := `INSERT INTO event (fib_num, result)
						   VALUES ($1, $2)`

	_, err := db.Exec(sqlStatement, key, value)
	if err != nil {
		return errors.New(fmt.Sprintf("Error:DbInsert %s", err.Error()))
	} 
	return nil
}



func DbConnect() (*sql.DB, error) {
	var (
		host     = "postgres"
  		port     = 5432
  		user     = "docker"
	    password = "docker"
  		dbname   = "docker"
	)
	psqlInfo := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable", user, password, host, port, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error: %s", err.Error()))	
	}
	return db, nil
}