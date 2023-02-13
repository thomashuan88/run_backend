package gamebo

import (
	"fmt"
	"rabbit-mq-consumer/model"
)

type GameBoRequestApiLog struct {
	Request string
}

func (req GameBoRequestApiLog) InsertLog() {
	stmt, err := model.MiscLogsDB.Prepare("INSERT INTO game_bo_request_api_log (request) VALUES (?)")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(req.Request)
	if err != nil {
		fmt.Println(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Game Request API Log with ID %d saved successfully", id)
}
