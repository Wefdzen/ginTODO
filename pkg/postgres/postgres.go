package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

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
