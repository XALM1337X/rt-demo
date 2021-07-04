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
	"regexp"
	"strconv"
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
		return
	} 
	
	//Calculate it
	//TODO
	res, fib_err := FibGenerate(req["lookup"].(string))
	if fib_err != nil {
		w.Write([]byte(fib_err.Error()))
		return
	}
	fmt.Println(fmt.Sprintf("%v", res))
	//Store it
	/*insert_err := DbInsert(req["lookup"].(string), res, db)
	if insert_err != nil {
		w.Write([]byte(insert_err.Error()))
	} else {
		w.Write([]byte("Successfully added entry to db"))
	}*/
	
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


func FibGenerate(nth_fib_str string) (string,error) {
	n, n_err := strconv.Atoi(nth_fib_str)
	store := "0"
	current := "0"
	previous := "0"
	if n_err != nil {
		return "", errors.New(fmt.Sprintf("Error:FibGenerate: %s", n_err.Error()))
	}
	if n < 1 {
		return "", errors.New(fmt.Sprintf("Error:FibGenerate: Value must be greater than 0."))
	} else if n == 1 {
		return "0", nil
	} else {
		for i:=0; i<n; i++ {
			if i == 0 {
				continue
			} else if i == 1 {
				current = "1"
			} else {
				current, previous, store = FibCrunchStrings(current, previous, store)
			}			
		}
	}
	return strconv.Itoa(current), nil
}

func FibCrunchStrings(current string, previous string, store string) (string, string, string, error) {
/*
	store = current
	current += previous
	previous = store
*/

	carry := 0
	crunch1 := 0 
	crunch2 := 0
	result := 0 
	remainder := 0
	var new_str string = ""
	for i:=len(current); i > 0; i-- {
		if len(previous) > 0 {
			crunch1 = strconv.Atoi(current[i-1])
			crunch2 = strconv.Atoi(previous[len(previous)-1])
			result = crunch1+crunch2+carry
			carry = 0
			remainder = 0

			if result >= 10 {
				remainder = result % 10
				carry++
				ascii, err := strconv.Atoi(remainder)
				if err != nil {
					return "","","",err.Error()
				}
				new_str = append(new_str, ascii)
			}



		} else {

		}
		
	}

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