package database

import "fmt"

func CreateTables() {
	createUserTable := `
	CREATE TABLE IF NOT EXISTS users (
		user_id INT AUTO_INCREMENT PRIMARY KEY,
		word_used INT,
		words_left INT default 1000000,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
	`
	_, err := DB.Exec(createUserTable)
	if err != nil {
		fmt.Println("Error creating user table:", err)
	}

	createRequestTable := `
	CREATE TABLE IF NOT EXISTS requests (
		request_id INT AUTO_INCREMENT PRIMARY KEY,
		user_id INT NOT NULL,
		data TEXT,
		duration INT,  
		FOREIGN KEY (user_id) REFERENCES users(user_id)
	)
	`
	_, err = DB.Exec(createRequestTable)
	if err != nil {
		fmt.Println("Error creating request table:", err)
	}

	fmt.Println("tables are created successfully")
}
