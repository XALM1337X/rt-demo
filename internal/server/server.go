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
	"time"
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
	start := time.Now()
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
	//connect to DB
	db, err2 := DbConnect()
	if err2 != nil {
		w.Write([]byte(fmt.Sprintf("Error: %s", err2.Error())))
	}
	defer db.Close()
	//Ping it (officially open connection)
	err3 := db.Ping()
	if err3 != nil {
		w.Write([]byte(fmt.Sprintf("Error: %s", err3.Error())))
		return
	} 

	//Check if key(n'th fib number) exists.
	val, exists, err_cache := CheckCache(req["lookup"].(string), db)
	if err_cache != nil && !re.MatchString(err_cache.Error()){
		w.Write([]byte(err_cache.Error()))
		return
	}

	if exists {
		//if it does return it.
		duration := time.Since(start)
		w.Write([]byte(fmt.Sprintf("Successfully retrieved from database: Key: %v , Value: %v, Elapsed Time: %v", val.Nth_fib, val.Nth_fib_result, duration)))
		return
	} 
	
	//Calculate it
	//Cheese the off by 1 thats not a real off by 1
	//Typically fibbonaci starts at 1, mine starts at 0
	cheese_req, cheese_err := strconv.Atoi(req["lookup"].(string))
	if cheese_err != nil {
		w.Write([]byte(cheese_err.Error()))
	}
	cheese_format := cheese_req + 1 
	cheese := strconv.Itoa(cheese_format)
	res, fib_err := FibGenerate(cheese)
	if fib_err != nil {
		w.Write([]byte(fib_err.Error()))
		return
	}	
	
	//Store it
	insert_err := DbInsert(req["lookup"].(string), res, db)
	if insert_err != nil {
		w.Write([]byte(insert_err.Error()))
	} else {
		duration := time.Since(start)
		w.Write([]byte(fmt.Sprintf("Successfully stored in database: Key: %v , Value: %v, Elapsed Time: %v", req["lookup"].(string), res, duration)))
	}
	
}


func CheckCache(lookup string, db *sql.DB) (TFibParse, bool, error) {
	var fib TFibParse
	query := "SELECT fib_num, result FROM event WHERE fib_num='"+lookup+"';"
	err := db.QueryRow(query).Scan(&fib.Nth_fib, &fib.Nth_fib_result)
	if err != nil {
		return fib, false, err
	}
	fib.CacheStatus = "Entry found"
	return fib, true, nil
}


func FibGenerate(nth_fib_str string) (string,error) {
	n, n_err := strconv.Atoi(nth_fib_str)
	store := "0"
	current := "0"
	previous := "0"
	var err error
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
				store = current
				current, err = FibCrunchStrings(current, previous)
				if err != nil {
					return "", err
				}
				previous = store
			}			
		}
	}
	return current, nil
}

func FibCrunchStrings(current string, previous string) (string, error) {

	carry := 0
	crunch1 := 0 
	crunch2 := 0
	result := 0 
	remainder := 0
	var new_str_rev string = ""
	var err_conv1, err_conv2, err_conv3 error

	for i:=len(current)-1; i >= 0; i-- {
		if len(previous) > 0 {
			crunch1, err_conv1 = strconv.Atoi(string(current[i]))
			if err_conv1 != nil {
				return "",err_conv1
			}
			crunch2, err_conv2 = strconv.Atoi(string(previous[len(previous)-1]))
			if err_conv2 != nil {
				return "",err_conv2
			}
			result = crunch1+crunch2+carry
			carry = 0
			remainder = 0

			if result >= 10 {
				remainder = result % 10
				carry++
				ascii := strconv.Itoa(remainder)				
				new_str_rev += ascii
			} else {
				new_str_rev += strconv.Itoa(result)
			}
			if len(previous) > 0 {
				previous = previous[:len(previous)-1]
			}			
			
		} else {
			if carry > 0 {
				crunch1, err_conv3 = strconv.Atoi(string(current[i]))
				if err_conv3 != nil {
					return "",err_conv3
				}
				result = crunch1 + carry
				carry = 0
				remainder = 0

				if result >= 10 {
					remainder = result % 10
					carry++
					ascii := strconv.Itoa(remainder)				
					new_str_rev += ascii
				} else {
					new_str_rev += strconv.Itoa(result)
				}
			} else {
				new_str_rev += string(current[i])
			}
		}

		if i == 0 && carry > 0 && len(previous) < 1{
			new_str_rev += "1"			
		}

	}
	ret_string :=""
	for i:=len(new_str_rev)-1; i>=0; i-- {
		ret_string += string(new_str_rev[i])
	}
	return ret_string, nil

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