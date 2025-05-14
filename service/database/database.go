/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetMyUser(cs int) (string, []byte, int, error)
	SetMyPhoto([]byte, int) error
	SetMyUserName(string, int) error

	DoLogin(string) (int, error)

	CreateChat(int, int) (int, error)
	CheckNames(int, string) ([]UtenteDb, error)
	CheckChatNames(int, string) ([]ChatUtenteDb, error)
	LeaveGroup(int, int) error
	GetMyConversations(int) ([]ChatUtenteDb, error)
	GetConversation(int, int) (bool, string, []MessDb, error)
	SendMessage(int, int, string, []byte, int) error
	DeleteMessage(int, int, int, int) error
	ForwardMessage(int, []int, int, int) error
	CreateGroup(int, string, []byte, []int) (int, error)
	SetGroupName(int, int, string) error
	SetGroupPhoto(int, int, []byte) error
	AddToGroup(int, int, []int) error
	CommentMessage(int, int, string, int) error
	UncommentMessage(int, int, int) error
	GetGroupUsers(int, int) ([]UtenteDb, error)
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='example_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {

		sqlStmt := `CREATE TABLE if not exists utenti (id INTEGER NOT NULL PRIMARY KEY, username TEXT, propic BLOB);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure utenti: %w", err)
		}

		sqlStmt = `CREATE TABLE if not exists chat (id INTEGER NOT NULL PRIMARY KEY, nome TEXT, propic BLOB,gruppo BOOL);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure chat: %w", err)
		}

		sqlStmt = `CREATE TABLE if not exists membri (
			id_utenti INTEGER NOT NULL,
			id_chat INTEGER NOT NULL,
			FOREIGN KEY (id_utenti) REFERENCES utenti(id) ON DELETE CASCADE, 
			FOREIGN KEY (id_chat) REFERENCES chat(id) ON DELETE CASCADE
			); `
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure membri: %w", err)
		}

		sqlStmt = `CREATE TABLE if not exists messaggi (id INTEGER NOT NULL PRIMARY KEY, testo TEXT, image BLOB, data TIMESTAMP DEFAULT CURRENT_TIMESTAMP, mittente INTEGER NOT NULL, FOREIGN KEY(mittente) REFERENCES utenti(id));`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure messaggi: %w", err)
		}

		sqlStmt = `CREATE TABLE if not exists messaggi_di_chat (id INTEGER NOT NULL PRIMARY KEY, id_chat INTEGER NOT NULL, id_messaggio INTEGER NOT NULL,id_forward INTEGER,id_forw_mit INTEGER, forward_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP, id_reply INTEGER, visual BOOLEAN ,FOREIGN KEY(id_forward) REFERENCES utenti(id) , FOREIGN KEY(id_chat) REFERENCES chat(id), FOREIGN KEY(id_messaggio) REFERENCES messaggi(id), FOREIGN KEY (id_reply) REFERENCES messaggi_di_chat(id));`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure messaggi_di_chat: %w", err)
		}

		sqlStmt = `CREATE TABLE if not exists emoticon (id_utente INTEGER NOT NULL, id_messaggio INTEGER NOT NULL,emoji VARCHAR(3) NOT NULL, FOREIGN KEY(id_utente) REFERENCES utenti(id), FOREIGN KEY(id_messaggio) REFERENCES messaggi(id) ON DELETE CASCADE);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure emoticon: %w", err)
		}

		sqlStmt = `CREATE TABLE if not exists accessi_chat (id_utente INTEGER NOT NULL, id_chat INTEGER NOT NULL, data TIMESTAMP DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (id_utente, id_chat), FOREIGN KEY(id_utente) REFERENCES utenti(id) ON DELETE CASCADE, FOREIGN KEY(id_chat) REFERENCES chat(id) ON DELETE CASCADE);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure accessi_chat: %w", err)
		}
	}
	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
