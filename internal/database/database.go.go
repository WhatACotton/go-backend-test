package database

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func GetDate() string {
	t := time.Now()
	ft := t.Format("2006-01-02 15:04:05")
	return ft
}

// []uint8型の値をtime.Time型に変換する
func ConvertBytesToTime(b []uint8) time.Time {
	str := string(b)
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return time.Time{}
	}
	return time.Unix(i, 0)
}

type PatchRequestPayload struct {
	ID        string `json:"id"`
	Attribute string `json:"attribute"`
	Value     string `json:"value"`
}

func (patchItem PatchRequestPayload) Patch(table string, where string) {
	db := ConnectSQL()
	defer db.Close()

	// SQLの準備
	upd, err := db.Prepare("UPDATE ? SET ? = ? WHERE ? = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer upd.Close()

	if http.DetectContentType([]byte(patchItem.Value)) == "int" {
		value, err := strconv.Atoi(patchItem.Value)
		if err != nil {
			log.Fatal(err)
		}
		// SQLの実行
		_, err = upd.Exec(table, patchItem.Attribute, value, where, patchItem.ID)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// SQLの実行
		_, err = upd.Exec(table, patchItem.Attribute, patchItem.Attribute, where, patchItem.ID)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func Delete(table string, where string, id string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()

	// SQLの実行
	del, err := db.Prepare("DELETE FROM ? WHERE ? = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer del.Close()

	// SQLの実行
	_, err = del.Exec(table, where, id)
	if err != nil {
		log.Fatal(err)
	}
}

func ConnectSQL() (db *sql.DB) {
	// データベースのハンドルを取得する
	db, err := sql.Open("mysql", os.Getenv("MYSQL_USER")+":"+os.Getenv("MYSQL_PASS")+"@tcp(localhost:3306)/go_test")

	if err != nil {
		// ここではエラーを返さない
		log.Fatal(err)
	}
	return db
}

func TestSQL() {
	// データベースのハンドルを取得する
	mysql := os.Getenv("MYSQL_USER") + ":" + os.Getenv("MYSQL_PASS") + "@tcp(localhost:3306)/go_test"
	log.Println(mysql)
	db, err := sql.Open("mysql", mysql)

	// 実際に接続する
	err = db.Ping()
	if err != nil {
		log.Println("データベースに接続できません。MySQLが起動しているか、環境変数が設定されているか確認してください。")
		log.Fatal(err)
		return
	} else {
		log.Println("データベース接続確認")
	}

}
