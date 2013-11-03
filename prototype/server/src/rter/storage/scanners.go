package storage

import (
	"database/sql"
	"log"
	"rter/data"
	"time"
	"net/http"
	"fmt"
	"net"
)

var timeout = time.Duration(10 * time.Millisecond)

func dialTimeout(network, addr string) (net.Conn, error) {
    return net.DialTimeout(network, addr, timeout)
}

func scanItemComment(comment *data.ItemComment, rows *sql.Rows) error {
	var updateTimeString string

	err := rows.Scan(
		&comment.ID,
		&comment.ItemID,
		&comment.Author,
		&comment.Body,
		&updateTimeString,
	)

	if err != nil {
		return err
	}

	// TODO: this is a hacky fix for null times
	if updateTimeString == "0000-00-00 00:00:00" {
		updateTimeString = "0001-01-01 00:00:00"
	}
	updateTime, err := time.Parse("2006-01-02 15:04:05", updateTimeString) // this assumes UTC as timezone

	if err != nil {
		log.Println("ItemComment scanner failed to parse time. " + updateTimeString)
		return err
	}

	comment.UpdateTime = updateTime

	return nil
}

func scanItem(item *data.Item, rows *sql.Rows) error {
	var startTimeString, stopTimeString string

	err := rows.Scan(
		&item.ID,
		&item.Type,
		&item.Author,
		&item.ThumbnailURI,
		&item.ContentURI,
		&item.UploadURI,
		&item.ContentToken,
		&item.HasHeading,
		&item.Heading,
		&item.HasGeo,
		&item.Lat,
		&item.Lng,
		&item.Radius,
		&item.Live,
		&startTimeString,
		&stopTimeString,
	)

	if err != nil {
		return err
	}

	// TODO: this is a hacky fix for null times
	if startTimeString == "0000-00-00 00:00:00" {
		startTimeString = "0001-01-01 00:00:00"
	}
	startTime, err := time.Parse("2006-01-02 15:04:05", startTimeString) // this assumes UTC as timezone

	if err != nil {
		log.Println("Item scanner failed to parse time. " + startTimeString)
		return err
	}

	item.StartTime = startTime

	// TODO: this is a hacky fix for null times
	if stopTimeString == "0000-00-00 00:00:00" {
		stopTimeString = "0001-01-01 00:00:00"
	}
	stopTime, err := time.Parse("2006-01-02 15:04:05", stopTimeString) // this assumes UTC as timezone

	if err != nil {
		log.Println("Item scanner failed to parse time. " + stopTimeString)
		return err
	}

	item.StopTime = stopTime

	transport := http.Transport{
        Dial: dialTimeout,
    }

    client := http.Client{
        Transport: &transport,
    }

	for thumbnailID := 1; ; thumbnailID += 1 {
		thumbnailURI := fmt.Sprintf(item.ContentURI + "/thumb/%09d.jpg", thumbnailID)
		resp, err := client.Get(thumbnailURI)
		if err != nil {
			break
		}
		if resp.StatusCode == 404 {
			if thumbnailID > 1 {
				item.ThumbnailURI = fmt.Sprintf(item.ContentURI + "/thumb/%09d.jpg", thumbnailID - 1)
			}
			break
		}
		resp.Body.Close()
	}
	

	return nil
}

func scanTerm(term *data.Term, rows *sql.Rows) error {
	var updateTimeString string

	cols, err := rows.Columns()

	if err != nil {
		return err
	}

	if len(cols) < 5 {
		err = rows.Scan(
			&term.Term,
			&term.Automated,
			&term.Author,
			&updateTimeString,
		)
	} else {
		err = rows.Scan(
			&term.Term,
			&term.Automated,
			&term.Author,
			&updateTimeString,
			&term.Count,
		)
	}

	if err != nil {
		return err
	}

	// TODO: this is a hacky fix for null times
	if updateTimeString == "0000-00-00 00:00:00" {
		updateTimeString = "0001-01-01 00:00:00"
	}
	updateTime, err := time.Parse("2006-01-02 15:04:05", updateTimeString) // this assumes UTC as timezone

	if err != nil {
		log.Println("Term scanner failed to parse time.")
		return err
	}

	term.UpdateTime = updateTime

	return nil
}

func scanTermRelationship(relationship *data.TermRelationship, rows *sql.Rows) error {
	err := rows.Scan(
		&relationship.Term,
		&relationship.ItemID,
	)

	return err
}

func scanTermRanking(ranking *data.TermRanking, rows *sql.Rows) error {
	var updateTimeString string

	err := rows.Scan(
		&ranking.Term,
		&ranking.Ranking,
		&updateTimeString,
	)

	if err != nil {
		return err
	}

	// TODO: this is a hacky fix for null times
	if updateTimeString == "0000-00-00 00:00:00" {
		updateTimeString = "0001-01-01 00:00:00"
	}
	updateTime, err := time.Parse("2006-01-02 15:04:05", updateTimeString) // this assumes UTC as timezone

	if err != nil {
		log.Println("TermRanking scanner failed to parse time.")
		return err
	}

	ranking.UpdateTime = updateTime

	return nil
}

func scanRole(role *data.Role, rows *sql.Rows) error {
	err := rows.Scan(
		&role.Title,
		&role.Permissions,
	)

	return err
}

func scanUser(user *data.User, rows *sql.Rows) error {
	var createTimeString string
	var updateTimeString string
	var statusTimeString string

	err := rows.Scan(
		&user.Username,
		&user.Password,
		&user.Salt,
		&user.Role,
		&user.TrustLevel,
		&createTimeString,
		&user.Heading,
		&user.Lat,
		&user.Lng,
		&updateTimeString,
		&user.Status,
		&statusTimeString,
	)

	if err != nil {
		return err
	}

	// TODO: this is a hacky fix for null times
	if createTimeString == "0000-00-00 00:00:00" {
		createTimeString = "0001-01-01 00:00:00"
	}
	createTime, err := time.Parse("2006-01-02 15:04:05", createTimeString) // this assumes UTC as timezone

	if err != nil {
		log.Println("User scanner failed to parse create time.")
		return err
	}

	user.CreateTime = createTime

	// TODO: this is a hacky fix for null times
	if updateTimeString == "0000-00-00 00:00:00" {
		updateTimeString = "0001-01-01 00:00:00"
	}
	updateTime, err := time.Parse("2006-01-02 15:04:05", updateTimeString) // this assumes UTC as timezone

	if err != nil {
		log.Println("User scanner failed to parse update time.")
		return err
	}

	user.UpdateTime = updateTime

	// TODO: this is a hacky fix for null times
	if statusTimeString == "0000-00-00 00:00:00" {
		statusTimeString = "0001-01-01 00:00:00"
	}
	statusTime, err := time.Parse("2006-01-02 15:04:05", statusTimeString) // this assumes UTC as timezone

	if err != nil {
		log.Println("User scanner failed to parse time.")
		return err
	}

	user.StatusTime = statusTime

	return nil
}

func scanUserDirection(direction *data.UserDirection, rows *sql.Rows) error {
	var updateTimeString string

	err := rows.Scan(
		&direction.Username,
		&direction.LockUsername,
		&direction.Command,
		&direction.Heading,
		&direction.Lat,
		&direction.Lng,
		&updateTimeString,
	)

	if err != nil {
		return err
	}

	// TODO: this is a hacky fix for null times
	if updateTimeString == "0000-00-00 00:00:00" {
		updateTimeString = "0001-01-01 00:00:00"
	}
	updateTime, err := time.Parse("2006-01-02 15:04:05", updateTimeString) // this assumes UTC as timezone

	if err != nil {
		log.Println("UserDirection scanner failed to parse time.")
		return err
	}

	direction.UpdateTime = updateTime

	return nil
}
