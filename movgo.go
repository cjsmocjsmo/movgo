package movgo

import (
	"github.com/globalsign/mgo"
	"io"
	"io/ioutil"
	"log"
	// "log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//MovDBcon is exported for all our db connection objects
func MovDBcon() *mgo.Session {
	log.Println("Starting DB session")
	s, err := mgo.Dial(os.Getenv("MEDIACENTER_MONGODB_ADDRESS"))
	if err != nil {
		log.Println("this is dial err")
		panic(err)
	}
	return s
}

func isDirEmpty(name string) (bool, error) {
	log.Printf("\n this is name from isDirEmpty %s \n", name)
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// read in ONLY one file
	_, err = f.Readdir(1)

	// and if the file is EOF... well, the dir is empty.
	if err == io.EOF {
		return true, nil
	}
	log.Println("isDirEmpty is complete")
	return false, err
}

//ProcessMovs is needed in update
func ProcessMovs(pAth string) {
	log.Println("Process_Movs has started")
	// var movpicPath string
	// var httppicPath string

	movpicPath, httppicPath := FindPicPaths(pAth, os.Getenv("MOVIEGOBS_NO_ART_PIC_PATH"))
	log.Printf("\n\n THIS IS MOVPICPATH %s", movpicPath)
	var MovI MOVI
	MovI = GetMovieInfo(pAth, movpicPath, httppicPath)
	ses := MovDBcon()
	defer ses.Close()
	MTc := ses.DB("moviegobs").C("moviegobs")
	err := MTc.Insert(&MovI)
	if err != nil {
		log.Println(err)
	}
}

// func processTVShow(pAth string) {
// 	log.Println("Process_Movs has started")
// 	var tvpicpath string
// 	tvpicpath = FindPicPaths(pAth, os.Getenv("MOVIEGOBS_NO_ART_PIC_PATH"))

// 	var TVShowI tVShowInfoS
// 	TVShowI = getTvShowInfo(pAth, tvpicpath)

// 	ses := MovDBcon()
// 	defer ses.Close()
// 	MTc := ses.DB("moviegobs").C("tvshows")
// 	err := MTc.Insert(&TVShowI)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	return
// }

func PosterDirVisit(posterpath string, f os.FileInfo, err error) error {
	log.Printf("\n\n this is posterpath from PosterDirVisit %s\n\n", posterpath)
	ext := filepath.Ext(posterpath)
	log.Printf("\n\n this is ext from posterdirvistit %s \n\n", ext)
	if err != nil {
		log.Println(err) // can't walk here,
		return nil       // but continue walking elsewhere
	}
	if f.IsDir() {
		log.Println("fi its a dir")
		log.Println(posterpath)
	} else if ext == ".txt" {
		log.Printf("\n\n its a txt file %s", f)
	} else if strings.Contains(posterpath, "TVShows") {
		log.Println("\nstarting createtvshowthumbnail")
		// CreateTVShowsThumbnail(posterpath)
	} else {
		log.Println("\n starting createmoviesthumbnail this is posterpath")
		log.Println(posterpath)
		CreateMoviesThumbnail(posterpath)
	}
	return nil
}

func genMatch(patH string, mtv bool) {
	if mtv {
		log.Println(patH)
		// processTVShow(patH)
	} else {
		ProcessMovs(patH)
	}
}

func myDirVisit(pAth string, f os.FileInfo, err error) error {
	log.Printf("this is path: %s", pAth)
	if err != nil {
		log.Println(err) // can't walk here,
		return nil       // not a file.  ignore.
	}
	if f.IsDir() {
		return nil
	}
	ext := filepath.Ext(pAth)
	if ext == "" {
		return nil //not a pic or movie
	}
	matchedTV, err := filepath.Match("*TVShows", f.Name())
	if err != nil {
		log.Println(err) // malformed pattern
		return err       // this is fatal.
	}
	switch {
	case ext == ".mp4":
		genMatch(pAth, matchedTV)
	case ext == ".mkv":
		genMatch(pAth, matchedTV)
	case ext == ".avi":
		genMatch(pAth, matchedTV)
	case ext == ".m4v":
		genMatch(pAth, matchedTV)
	}
	return nil
}

func removeFiles() {
	dir, _ := ioutil.ReadDir("/root/static")
	for _, d := range dir {
		os.RemoveAll(path.Join([]string{"tmp", d.Name()}...))
	}
}

func posterTotal() int {
	posters, _ := filepath.Glob("/root/fsData/Posters2/*.*")
	posttotal := len(posters)
	return posttotal
}

func thumbTotal() int {
	thumb, _ := filepath.Glob("/root/fsData/Thumbnails/*.*")
	thumbtotal := len(thumb)
	return thumbtotal
}

func picUpdateStatus() (updateStat bool) {
	pt := posterTotal()
	tt := thumbTotal()

	lpp := strconv.Itoa(pt)
	ltt := strconv.Itoa(tt)
	log.Printf("this is lp %s", lpp)
	log.Printf("this is lt %s", ltt)

	if pt != tt {
		updateStat = true
	} else {
		updateStat = false
	}
	return
}

func setupLogging() {
	logfile := os.Getenv("MEDIACENTER_LOG_BASE_PATH") + "/moviegobsMOV.log"
	// If the file doesn't exist, create it or append to the file
	file, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	log.Println("Mov logging started")
}

//MovSetUp is exported to main
func MovSetUp() (ExStat int) {
	//Start the timer
	starttime := time.Now().Unix()
	startTime2 := strconv.FormatInt(starttime, 10)
	// starttime := strconv.Itoa(s)

	log.Printf("setup function has started at: %s", startTime2)
	//Connect to the DB
	sess := MovDBcon()
	err := sess.DB("moviegobs").DropDatabase()
	if err != nil {
		log.Println(err)
	}
	err = sess.DB("movbsthumb").DropDatabase()
	if err != nil {
		log.Println(err)
	}
	sess.Close()
	log.Println("moviegobs and movbsthumb dbs have been dropped")

	//Check thumbnail dir create thumbs if empty
	empty, err := isDirEmpty("/root/static")
	if err != nil {
		log.Println(err)
	}
	if empty {
		filepath.Walk("/root/fsData/Posters2", PosterDirVisit)
	} else {
		if picUpdateStatus() {
			removeFiles()
			filepath.Walk("/root/fsData/Posters2", PosterDirVisit)
		}
	}

	err = filepath.Walk(os.Getenv("MEDIACENTER_MOVIES_PATH"), myDirVisit)
	if err != nil {
		log.Println(err)
	}

	os.Setenv("MEDIACENTER_SETUP", "0")
	log.Printf("this is noartlist :: %s", NoArtList)
	log.Println(startTime2)
	stopTime := time.Now().Unix()
	log.Println(stopTime)
	etime := stopTime - starttime
	log.Println(etime)
	log.Println("SETUP IS COMPLETE")
	ExStat = 0
	return
}
