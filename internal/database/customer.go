package database

import (
	"html"
	"log"
	"unify/internal/models"
	"unify/validation"
)

func SignUpCustomer(req models.CustomerRequestPayload, SessionID string) (res models.Customer) {
	log.Printf("SignUpCustomer Called")
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの準備
	//UID,Name,Address,Email,PhoneNumber,Register,CreatedDate,ModifiedDate,RegisteredDate,LastLogInDate

	ins, err := db.Prepare("INSERT INTO user VALUES(?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(req.UID, "default", "default", req.Email, "00000000000", false, GetDate(), 20000101, 20000101, SessionID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(req.Email)

	return GetCustomer(req.UID)
}

func RegisterCustomer(usr validation.User, customer models.CustomerRegisterPayload) (res models.Customer) {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの準備
	//UID,Name,Address,Email,PhoneNumber,Register,CreatedDate,ModifiedDate,RegisteredDate,LastLogInDate
	ins, err := db.Prepare(
		`UPDATE user SET 
		Name = ?,
		Address = ?,
		PhoneNumber = ?,
		Register = ?,
		RegisteredDate = ?,
		WHERE uid = ?`)
	if err != nil {
		log.Fatal(err)
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(html.EscapeString(customer.Name), html.EscapeString(customer.Address), customer.PhoneNumber, true, GetDate(), usr.Userdata.UID)
	if err != nil {
		log.Fatal(err)
	}
	return GetCustomer(usr.Userdata.UID)
}

func ModifyCustomer(usr validation.User, customer models.CustomerRegisterPayload) (res models.Customer) {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの準備
	//UID,Name,Address,Email,PhoneNumber,Register,CreatedDate,ModifiedDate,RegisteredDate,LastLogInDate
	ins, err := db.Prepare(
		`UPDATE user SET 
		Name = ?,
		Address = ?,
		PhoneNumber = ?,
		Register = ?,
		ModifiedDate = ?,
		WHERE uid = ?`)
	if err != nil {
		log.Fatal(err)
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(html.EscapeString(customer.Name), html.EscapeString(customer.Address), customer.PhoneNumber, true, GetDate(), usr.Userdata.UID)
	if err != nil {
		log.Fatal(err)
	}
	return GetCustomer(usr.Userdata.UID)
}

func StoredLogInCustomer(uid string, NewSessionKey string, OldSessionKey string) {
	LogInLog(uid, NewSessionKey)
	Invalid(OldSessionKey)
	db := ConnectSQL()
	// SQLの実行
	rows, err := db.Query("SELECT * FROM user WHERE uid = ?", uid)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var Customer models.Customer
	// SQLの実行
	for rows.Next() {
		//err := rows.Scan(&Customer)
		err := rows.Scan(&Customer.UID, &Customer.Name, &Customer.Address, &Customer.Email, &Customer.PhoneNumber, &Customer.Register, &Customer.CreatedDate, &Customer.ModifiedDate, &Customer.RegisteredDate, &Customer.LastSessionId)
		if err != nil {
			panic(err.Error())
		}
	}
}
func VerifyCustomer(OldSessionKey string) bool {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの実行
	rows, err := db.Query("SELECT Available FROM loginlog WHERE SessionKey = ?", OldSessionKey)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var loginlog bool
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&loginlog)

		if err != nil {
			panic(err.Error())
		}
	}
	return loginlog
}
func NewLogInCustomer(uid string, NewSessionKey string) {
	LogInLog(uid, NewSessionKey)
	db := ConnectSQL()
	// SQLの実行
	rows, err := db.Query("SELECT * FROM user WHERE uid = ?", uid)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var Customer models.Customer
	// SQLの実行
	for rows.Next() {
		//err := rows.Scan(&Customer)
		err := rows.Scan(&Customer.UID, &Customer.Name, &Customer.Address, &Customer.Email, &Customer.PhoneNumber, &Customer.Register, &Customer.CreatedDate, &Customer.ModifiedDate, &Customer.RegisteredDate, &Customer.LastSessionId)
		if err != nil {
			panic(err.Error())
		}
	}
}

func LogInLog(uid string, NewSessionKey string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの準備
	//UID SessionKey LoginedDate Available
	ins, err := db.Prepare("INSERT INTO loginlog VALUES(?,?,?,true)")
	if err != nil {
		log.Fatal(err)
	}

	// SQLの実行
	_, err = ins.Exec(uid, NewSessionKey, GetDate())
	if err != nil {
		log.Fatal(err)
	}
	defer ins.Close()

}
func Invalid(SessionKey string) {
	log.Println("Invalid called")
	// データベースのハンドルを取得する
	db := ConnectSQL()
	ins, err := db.Prepare("UPDATE loginlog SET Available = false WHERE SessionKey = ?")
	if err != nil {
		log.Fatal(err)
	}
	// SQLの実行
	_, err = ins.Exec(SessionKey)
	defer ins.Close()
}

func GetCustomer(uid string) (res models.Customer) {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの実行
	rows, err := db.Query("SELECT * FROM user WHERE uid = ?", uid)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var Customer models.Customer
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&Customer.UID, &Customer.Name, &Customer.Address, &Customer.Email, &Customer.PhoneNumber, &Customer.Register, &Customer.CreatedDate, &Customer.ModifiedDate, &Customer.RegisteredDate, &Customer.LastSessionId)

		if err != nil {
			panic(err.Error())
		}
	}
	return Customer
}
