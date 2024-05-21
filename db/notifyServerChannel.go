package db

import "fmt"

func checkExistingServerData(serverID int64) (int64, error) {
	db := ConnectDB()
	if db == nil {
		err := fmt.Errorf("%s", "Cannot connect to database")
		db.Close()
		return -1, err
	}
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf(`
		SELECT channel_id FROM notify_servers WHERE server_id = %d`, serverID))
	if err != nil {
		return -1, err
	}

	var channelID int64

	for rows.Next() {
		err = rows.Scan(&channelID)
		if err != nil {
			return -1, err
		}
	}

	return channelID, nil
}

func GetNotifiedChannellsID() ([]int64, error) {
	var ids []int64
	db := ConnectDB()
	if db == nil {
		err := fmt.Errorf("%s", "Cannot connect to database")
		db.Close()
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT channel_id FROM notify_servers")
	if err != nil {
		return nil, err
	}

	var id int64

	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func NotifyServerChannel(serverID int64, channelID int64) (string, error) {
	var message string
	db := ConnectDB()
	if db == nil {
		err := fmt.Errorf("%s", "Cannot connect to database")
		db.Close()
		return message, err
	}
	defer db.Close()

	dbChannelID, err := checkExistingServerData(serverID)

	if err != nil {
		return message, err
	}

	if channelID == dbChannelID {
		return message, nil
	}

	if dbChannelID == 0 {
		stmt, err := db.Prepare(`
			INSERT INTO notify_servers (server_id, channel_id) VALUES ($1, $2)`)
		if err != nil {
			return message, err
		}
		_, err = stmt.Exec(serverID, channelID)
		if err != nil {
			return message, err
		}
		message = "From now new forum posts will be notified on this channel :tada:"
	} else {
		stmt, err := db.Prepare("UPDATE notify_servers SET channel_id = $1")
		if err != nil {
			return message, err
		}
		_, err = stmt.Exec(channelID)
		if err != nil {
			return message, err
		}
		message = "Channel for notifying of new forum posts was changed :arrows_clockwise:"
	}

	return message, nil
}
