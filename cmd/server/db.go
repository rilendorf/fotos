package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3" // MIT licensed.

	"log"
)

var (
	db *sql.DB
)

var dblog *log.Logger

var (
	addImg          *sql.Stmt
	readImgs        *sql.Stmt
	readImgsDeleted *sql.Stmt
	readImgPath     *sql.Stmt

	delImg     *sql.Stmt
	delImgHard *sql.Stmt

	addAcc   *sql.Stmt
	delAcc   *sql.Stmt
	getTkn   *sql.Stmt
	getAcc   *sql.Stmt
	rmImgAcc *sql.Stmt

	markImgPrint   *sql.Stmt
	getImg2Print   *sql.Stmt
	markImgPrinted *sql.Stmt
)

func openDB(path string) {
	var err error

	dblog = log.New(log.Writer(), "[DB] ", 0)
	db, err = sql.Open("sqlite3", path)
	if err != nil {
		dblog.Fatalf("open file %s: %s", path, err)
	}

	// create tables
	// images: name, blob, account (shared secret)
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `images` (`name` TEXT PRIMARY KEY, `blob` BLOB, `account` TEXT, `deleted` BOOLEAN);")
	if err != nil {
		dblog.Fatalf("Can't prepare statement %s", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `printimages` (`name` TEXT PRIMARY KEY, `account` TEXT, `printed` BOOLEAN);")
	if err != nil {
		dblog.Fatalf("Can't prepare statement %s", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `tokens` (`account` TEXT PRIMARY KEY, `token` VARCHAR[25] KEY, `viewtoken` VARCHAR[25]);")
	if err != nil {
		dblog.Fatalf("Can't prepare statement %s", err)
	}

	// marks image as deleted and sets blob to 0
	delImg, err = db.Prepare("UPDATE images SET deleted = true, blob = 0 WHERE name = ? AND account = ?;")
	if err != nil {
		dblog.Fatalf("Can't prepare statement %s", err)
	}

	delImgHard, err = db.Prepare("DELETE FROM images WHERE name = ? AND account = ?;")
	if err != nil {
		dblog.Fatalf("Can't prepare statement %s", err)
	}

	rmImgAcc, err = db.Prepare("DELETE FROM images WHERE account = ?")
	if err != nil {
		dblog.Fatalf("Can't prepare statement %s", err)
	}

	addImg, err = db.Prepare("INSERT OR REPLACE INTO images (name, account, blob, deleted) VALUES (?, ?, ?, false)")
	if err != nil {
		dblog.Fatalf("Can't prepare statement %s", err)
	}

	readImgs, err = db.Prepare("SELECT name FROM images WHERE account = ? AND deleted != true")
	if err != nil {
		dblog.Fatalf("Can't prepare statement %s", err)
	}

	readImgsDeleted, err = db.Prepare("SELECT name FROM images WHERE account = ? AND deleted = true")
	if err != nil {
		dblog.Fatalf("Can't prepare statement %s", err)
	}

	readImgPath, err = db.Prepare("SELECT blob FROM images WHERE name = ? AND account = ?")
	if err != nil {
		dblog.Fatalf("Can't prepare statement %s", err)
	}

	addAcc, err = db.Prepare("INSERT OR REPLACE INTO tokens (account, token, viewtoken) VALUES (?, ?, ?)")
	if err != nil {
		dblog.Fatalf("Can't prepare statement %s", err)
	}

	delAcc, err = db.Prepare("DELETE FROM tokens WHERE account = ?")
	if err != nil {
		dblog.Fatalf("Can't prepare statement %s", err)
	}

	getTkn, err = db.Prepare("SELECT token, viewtoken FROM tokens WHERE account = ?")
	if err != nil {
		dblog.Fatalf("Can't prepare statement %s", err)
	}

	getAcc, err = db.Prepare("SELECT account FROM tokens WHERE token = ?")
	if err != nil {
		dblog.Fatalf("Can't prepare statement %s", err)
	}

	// 	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `printimages` (`name` TEXT PRIMARY KEY, `account` TEXT, printed` BOOLEAN);")
	markImgPrint, err = db.Prepare("INSERT OR REPLACE INTO printimages (name, account, printed) VALUES (?, ?, false)")
	if err != nil {
		dblog.Fatalf("Can't prepare statement %s", err)
	}

	getImg2Print, err = db.Prepare("SELECT name FROM printimages WHERE account = ? AND printed = FALSE")
	if err != nil {
		dblog.Fatalf("Can't prepare statement %s", err)
	}

	markImgPrinted, err = db.Prepare("UPDATE printimages SET printed = true WHERE name = ? AND account = ?")
	if err != nil {
		dblog.Fatalf("Can't prepare statement %s", err)
	}

	return
}

func MarkImagePrinted(name, account string) error {
	_, err := markImgPrinted.Exec(name, account)
	return err
}

func MarkImagePrint(name, account string) error {
	_, err := markImgPrint.Exec(name, account)
	return err
}

func GetImages2Print(account string) ([]string, error) {
	r, err := getImg2Print.Query(account)
	if err != nil {
		return nil, err
	}

	s := make([]string, 0)
	var buf string

	for r.Next() {
		err = r.Scan(&buf)
		if err != nil {
			return nil, err
		}

		s = append(s, buf)
	}

	return s, nil
}

func addImage(name, account string, data []byte) error {
	_, err := addImg.Exec(name, account, data)
	if err != nil {
		dblog.Printf("Error adding Image, '%s' len(%d): %s\n", name, len(data), err)
	}

	return err
}

func readImages(account string) (s []string) {
	r, err := readImgs.Query(account)
	if err != nil {
		dblog.Printf("Error reading Images: %s\n", err)
		return
	}

	defer r.Close()

	for r.Next() {
		var name string
		if err := r.Scan(&name); err != nil {
			if err != nil {
				dblog.Printf("Error reading rows: %s\n", err)
				return
			}
		}
		s = append(s, name)
	}
	return
}

func readImage(account, name string) (b []byte) {
	r, err := readImgPath.Query(name, account)
	if err != nil {
		dblog.Printf("Error reading image at path '%s': %s\n", name, err)
		return
	}

	defer r.Close()

	if r.Next() {
		err = r.Scan(&b)
		if err != nil {
			dblog.Printf("Error scanning image data: %s\n", err)
		}
	} else {
		dblog.Printf("Tried to read nonexistent image %s/%s", account, name)
	}

	return
}

func haveImage(name, account string) bool {
	return len(readImage(name, account)) == 0
}

// sets blob to 0 and marks image as deleted
func deleteImage(name, account string) error {
	_, err := delImg.Exec(name, account)
	if err != nil {
		dblog.Printf("Error deleting image '%s'\n", name)
	}

	return err
}

// also deletes entry
func removeImage(name, account string) error {
	_, err := delImg.Exec(name, account)
	if err != nil {
		dblog.Printf("Error removing image '%s'\n", name)
	}

	return err
}

func addAccount(name, token, viewtoken string) error {
	_, err := addAcc.Exec(name, token, viewtoken)
	if err != nil {
		dblog.Printf("Error adding Account, '%s': %s\n", name, err)
	}

	return err
}

func removeAccount(name string) {
	_, err := delAcc.Exec(name)
	if err != nil {
		dblog.Printf("Error removing account '%s'\n", name)
	}
}

func getTokenByName(acc string) (upload, view string) {
	r, err := getTkn.Query(acc)
	if err != nil {
		dblog.Printf("Error reading account '%s's token: %s\n", acc, err)
		return
	}

	defer r.Close()

	if r.Next() {
		err = r.Scan(&upload, &view)
		if err != nil {
			dblog.Printf("Error scanning account token: %s\n", err)
		}
	}

	return
}

func getNameByToken(tkn string) (s string) {
	r, err := getAcc.Query(tkn)
	if err != nil {
		dblog.Printf("Error reading account by token '%s': %s\n", tkn, err)
		return
	}

	defer r.Close()

	if r.Next() {
		err = r.Scan(&s)
		if err != nil {
			dblog.Printf("Error scanning account name: %s\n", err)
		}
	}

	return

}

func removeAccountImages(acc string) {
	_, err := rmImgAcc.Exec(acc)
	if err != nil {
		log.Printf("Error removing account '%s's images: %s", acc, err)
	}
}
