package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Emmanuel-MacAnThony/snippetbox/internal/models"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *models.SnippetModel
	users *models.UserModel
	templateCache map[string]*template.Template
	formDecoder *form.Decoder
	sessionManager *scs.SessionManager
	
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil

}

func main() {

	addr := flag.String("addr", ":4000", "http network address")
	dsn := flag.String("sql", "root:feb0699@/snippetbox?parseTime=true", "MySql data source name")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// create a db connection pool
	db, err := openDB(*dsn)

	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	// initialize template cache
	templateCache, err := newTemplateCache()

	if err != nil{
		errorLog.Fatal(err)
	}



	// Initialize a new instance of our application struct, containing the
	// dependencies.

	// initialize a new form decoder
	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &models.SnippetModel{DB: db},
		users: &models.UserModel{DB: db},
		templateCache: templateCache,
		formDecoder: formDecoder,
		sessionManager: sessionManager,
		
	}

	srv := &http.Server{
		ErrorLog: errorLog,
		Addr:     *addr,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)

	err = srv.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}
}
