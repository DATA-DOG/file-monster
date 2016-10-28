package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"gopkg.in/jmcvetta/napping.v1"
	"code.cloudfoundry.org/bytefmt"
)

type fileList struct {
	Ok    bool   `json:"ok"`
	Files []file `json:"files"`
}

type file struct {
	ID         string `json:"id"`
	IsExternal bool   `json:"is_external"`
	Size			 uint64 `json:"size"`
}

type filesDelete struct {
	Ok bool `json:"ok"`
}

func eatPublicFiles(token string) {
	timeUntil := strconv.FormatInt(time.Now().Add(-time.Hour*24*30).Unix(), 10)
	log.Println("Eating public files until", timeUntil)
	params := &napping.Params{"token": token, "ts_to": timeUntil, "count": "1000"}
	list := fileList{}
	var err interface{}

	napping.Get("https://slack.com/api/files.list", params, &list, err)

	if err != nil {
		log.Panic(err)
	}

	deleteList(token, list)
}

func eatUserFiles(user string, token string, userName string) {
	timeUntil := strconv.FormatInt(time.Now().Add(-time.Hour*24*30).Unix(), 10)
	log.Println("Eating user", user, "files until", timeUntil)
	params := &napping.Params{"token": token, "user": user, "ts_to": timeUntil, "count": "1000"}
	list := fileList{}
	var err interface{}

	napping.Get("https://slack.com/api/files.list", params, &list, err)

	if err != nil {
		log.Panic(err)
	}

	count, totalSize := deleteList(token, list)
	log.Println("Done for user", user, userName)

	if count > 0 {
		notifyUser(userName, fmt.Sprint("Your ", count, " files (", bytefmt.ByteSize(totalSize), ") were delicious. Remember to feed me soon!"))
	} else {
		notifyUser(userName, "I did not find anything to eat.")
	}
}

func deleteList(token string, list fileList) (count int, totalSize uint64) {
	if !list.Ok {
		log.Print(list)
	}

	count = 0
  totalSize = 0

	for _, file := range list.Files {
		if file.IsExternal {
			continue
		}

		delParams := &napping.Params{"token": token, "file": file.ID}
		response := filesDelete{}
		var err interface{}
		napping.Get("https://slack.com/api/files.delete", delParams, &response, err)

		if response.Ok {
			log.Print("Deleted file #", file.ID)
		}

		time.Sleep(time.Second)
		count++
		totalSize += file.Size
	}

	return
}
