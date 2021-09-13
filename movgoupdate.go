package movgo

import (
	"fmt"
	// "io"
	"log"
	"os"
	"path/filepath"
	// "time"
	// "strings"
	// "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

var finished bool = false

func scanFileNames() {
	err := filepath.Walk(os.Getenv("MOVIEGOBS_MOVIES_PATH"), updateDirVisit)
	if err != nil {
		fmt.Println(err)
	}
}

func updateDirVisit(pAth string, f os.FileInfo, err error) error {
	log.Printf("this is path: %s", pAth)
	if err != nil {
		fmt.Println(err) // can't walk here,
		return nil       // not a file.  ignore.
	}
	if f.IsDir() {
		return nil
	}
	ext := filepath.Ext(pAth)
	if ext == "" {
		return nil //not a pic or movie
	} else if ext == ".mp4" {
		if !movNameInDbCheck(pAth) {
			ProcessMovs(pAth)
		}
	} else if ext == ".avi" {
		if !movNameInDbCheck(pAth) {
			ProcessMovs(pAth)
		}
	} else if ext == ".mkv" {
		if !movNameInDbCheck(pAth) {
			ProcessMovs(pAth)
		}
	} else if ext == ".m4v" {
		if !movNameInDbCheck(pAth) {
			ProcessMovs(pAth)
		}
	}
	return nil
}

func movNameInDbCheck(fn string) (result bool) {
	sess := MovDBcon()
	defer sess.Close()
	MTc := sess.DB("moviegobs").C("moviegobs")
	b1 := bson.M{"filepath": fn}
	b2 := bson.M{"_Id": 0}
	var PMedia []map[string]string
	err := MTc.Find(b1).Select(b2).All(&PMedia)
	if err != nil {
		log.Println(err)
	}
	num := len(PMedia)
	if num != 0 {
		result = true
	} else {
		result = false
	}
	return
}

// func ProcessMovs(fn string) {
// 	isMovNameInDB := movNameInDbCheck(fn)
// 	if !isMovNameInDB {
		
// 	}

// }

// MovUpdate needs to be exported
func MovUpdate() (finished bool) {
	scanFileNames()
	finished = true
	return
}
