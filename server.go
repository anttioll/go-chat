package main

import (
	//"database/sql"
	"encoding/json"
	"fmt"
	"net"
	//"os"
	"time"
	//"github.com/go-sql-driver/mysql"
)

type MessageRow struct {
	UUID     string `json:"uuid"`
	User     string `json:"user"`
	Datetime string `json:"datetime"`
	Text     string `json:"text"`
}

//type Connections struct {
//  id int
//  connection []net.Conn
//}

var connections []net.Conn

func main() {
	server, err := net.Listen("tcp", ":42069")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("%s Connection accepted from: %v\n", time.Now().UTC().Format(time.DateTime), conn.RemoteAddr())
		connections = append(connections, conn)
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	var message MessageRow
	var bufferLen int
	var err error
	buffer := make([]byte, 1024)
	for {
		bufferLen, err = conn.Read(buffer)
		if err != nil {
			fmt.Println(err)
			removeConnection(conn)
			break
		}
		err = json.Unmarshal(buffer[:bufferLen], &message)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s Received: %v\n", time.Now().UTC().Format(time.DateTime), string(message.Text))
		//if _, err = addMessageToDB(message); err != nil {
		//	panic(err)
		//}
		sendToAllClients(buffer, bufferLen)
	}
}

func removeConnection(conn net.Conn) {
	for i, v := range connections {
		if v == conn {
			connections = append(connections[:i], connections[i+1:]...)
		}
	}
}

func sendToAllClients(buffer []byte, bufferLen int) {
	for _, conn := range connections {
		if _, err := conn.Write([]byte(buffer[:bufferLen])); err != nil {
			fmt.Printf("%s Couldn't send message to all clients: %v\n", time.Now().UTC().Format(time.DateTime), err)
		}
		fmt.Printf("%s Sent %v message to %s\n", time.Now().UTC().Format(time.DateTime), string(buffer[:bufferLen]), conn.RemoteAddr())
	}
}

//func connectToDB() *sql.DB {
//	var err error
//	var db *sql.DB
//	cfg := mysql.Config{
//		User:                 os.Getenv("DBUSER"),
//		Passwd:               os.Getenv("DBPASS"),
//		Net:                  "tcp",
//		Addr:                 os.Getenv("DBADDR"),
//		DBName:               os.Getenv("DBNAME"),
//		AllowNativePasswords: true,
//		ParseTime:            true,
//	}
//	db, err = sql.Open("mysql", cfg.FormatDSN())
//	if err != nil {
//		panic(err)
//	}
//	pingErr := db.Ping()
//	if pingErr != nil {
//		panic(pingErr)
//	}
//	return db
//}

//func addMessageToDB(msg MessageRow) (int64, error) {
//	db := connectToDB()
//	defer db.Close()
//	result, err := db.Exec("INSERT INTO messages (uuid, user, datetime, text) VALUES (?, ?, ?, ?)", msg.UUID, msg.User, msg.Datetime, msg.Text)
//	if err != nil {
//		return 0, fmt.Errorf("addMessageToDB: %v", err)
//	}
//	id, err := result.LastInsertId()
//	if err != nil {
//		return 0, fmt.Errorf("addMessageToDB: %v", err)
//	}
//	return id, nil
//}

//func readMessagesByUsername(username string) ([]messageRow, error) {
//	var messages []messageRow
//
//	rows, err := db.Query("SELECT * FROM messages WHERE user = ?", username)
//	if err != nil {
//		return nil, fmt.Errorf("readMessagesByUsername %q: %v", username, err)
//	}
//	defer rows.Close()
//	for rows.Next() {
//		var msg messageRow
//		if err := rows.Scan(&msg.ID, &msg.UUID, &msg.user, &msg.text, &msg.datetime); err != nil {
//			return nil, fmt.Errorf("readMessagesByUsername %q: %v", username, err)
//		}
//		messages = append(messages, msg)
//	}
//	if err := rows.Err(); err != nil {
//		return nil, fmt.Errorf("readMessagesByUsername %q: %v", username, err)
//	}
//	return messages, nil
//}
