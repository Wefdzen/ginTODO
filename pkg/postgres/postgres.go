package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"wefdzen/cmd/postes"
	"wefdzen/cmd/users"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func RegistrationUser(newUser users.User) bool {
	// Connect
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	urlToDataBase := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", Cfg.PGuser, Cfg.PGpassword, Cfg.PGaddress, Cfg.PGPort, Cfg.PGdbname)
	conn, err := pgx.Connect(context.Background(), urlToDataBase)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer conn.Close(context.Background())

	//check in db availible login
	var not_exists bool = false //  чтобы не записывать в бд
	query := fmt.Sprintf("SELECT NOT EXISTS(SELECT login FROM %s WHERE login=$1)", "users")
	err = conn.QueryRow(context.Background(), query, newUser.Login).Scan(&not_exists)
	if err != nil {
		fmt.Println(err.Error())
	}

	//add new user
	if not_exists { // если такой логин не занят то добавляем
		_, err = conn.Exec(context.Background(), `INSERT INTO users (login, email, password) VALUES ($1, $2, $3)`, newUser.Login, newUser.Email, newUser.Password)
		if err != nil {
			fmt.Println(err.Error())
		}
		return true // регистрация успешна
	}
	fmt.Println("login уже занят")
	return false // login занят
}

func CheckDataForLogin(login, password string) bool {
	// Connect
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	urlToDataBase := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", Cfg.PGuser, Cfg.PGpassword, Cfg.PGaddress, Cfg.PGPort, Cfg.PGdbname)
	conn, err := pgx.Connect(context.Background(), urlToDataBase)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer conn.Close(context.Background())

	passFromDB := ""
	err = conn.QueryRow(context.Background(), `SELECT password FROM users WHERE login=$1`, login).Scan(&passFromDB)
	if err != nil {
		fmt.Println(err.Error())
	}
	if success := bcrypt.CompareHashAndPassword([]byte(passFromDB), []byte(password)); success == nil {
		return true //login success
	}
	return false //login not success
}

func InsertNewPost(title string, text string) error {
	// Connect
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	urlToDataBase := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", Cfg.PGuser, Cfg.PGpassword, Cfg.PGaddress, Cfg.PGPort, Cfg.PGdbname)
	conn, err := pgx.Connect(context.Background(), urlToDataBase)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	command := fmt.Sprintf(`INSERT INTO %s (title, source_text) VALUES ($1, $2)`, Cfg.PGnameTable)
	_, err = conn.Exec(context.Background(), command, title, text)
	if err != nil {
		return err
	}
	return nil
}

func GetAllPost() []postes.PostUser {
	tmp := postes.New()

	// Connect
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	urlToDataBase := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", Cfg.PGuser, Cfg.PGpassword, Cfg.PGaddress, Cfg.PGPort, Cfg.PGdbname)
	conn, err := pgx.Connect(context.Background(), urlToDataBase)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer conn.Close(context.Background())

	command := fmt.Sprintf(`SELECT title, source_text FROM %s`, Cfg.PGnameTable)
	rows, _ := conn.Query(context.Background(), command)
	defer rows.Close()
	for rows.Next() {
		var tempTitle, tempText string
		rows.Scan(&tempTitle, &tempText)
		tmp.Add(postes.PostUser{Title: tempTitle, Post: tempText})
	}
	return tmp.Items
}

func DeletePostByID(id int) {
	// Connect
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	urlToDataBase := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", Cfg.PGuser, Cfg.PGpassword, Cfg.PGaddress, Cfg.PGPort, Cfg.PGdbname)
	conn, err := pgx.Connect(context.Background(), urlToDataBase)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer conn.Close(context.Background())

	command := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, Cfg.PGnameTable)
	_, err = conn.Exec(context.Background(), command, strconv.Itoa(id))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func FullEditPostByID(id int, newTitle, newText string) {
	// Connect
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	urlToDataBase := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", Cfg.PGuser, Cfg.PGpassword, Cfg.PGaddress, Cfg.PGPort, Cfg.PGdbname)
	conn, err := pgx.Connect(context.Background(), urlToDataBase)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer conn.Close(context.Background())

	command := fmt.Sprintf(`UPDATE %s SET title=$1, source_text=$2  WHERE id=$3`, Cfg.PGnameTable)
	_, err = conn.Exec(context.Background(), command, newTitle, newText, strconv.Itoa(id))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func WatchPostByID(id int) postes.PostUser {
	// Connect
	tmp := postes.PostUser{Title: "", Post: ""}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	urlToDataBase := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", Cfg.PGuser, Cfg.PGpassword, Cfg.PGaddress, Cfg.PGPort, Cfg.PGdbname)
	conn, err := pgx.Connect(context.Background(), urlToDataBase)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer conn.Close(context.Background())

	command := fmt.Sprintf(`SELECT title, source_text FROM %s WHERE id=$1`, Cfg.PGnameTable)
	err = conn.QueryRow(context.Background(), command, id).Scan(&tmp.Title, &tmp.Post)
	if err != nil {
		fmt.Println(err.Error())
	}
	return tmp
}

func init() {
	file, err := os.Open("config.cfg")
	if err != nil {
		fmt.Println("Error open .cfg", err)
		panic("Can't open the file \"setting.cfg\"")
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	readSetting := make([]byte, fileInfo.Size())
	_, err = file.Read(readSetting)
	if err != nil {
		panic("can't read file")
	}

	err = json.Unmarshal(readSetting, &Cfg)
	if err != nil {
		panic("json err")
	}
}

type setting struct {
	PGaddress   string
	PGpassword  string
	PGuser      string
	PGdbname    string
	PGPort      string
	PGnameTable string
}

var (
	Cfg setting
)
