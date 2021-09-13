package movgo

import (
	// "fmt"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"path"
	"strconv"
	"strings"
	"time"
)

// func getmovName(movname string) (movName string) {
// 	_, fname := path.Split(movname)
// 	if strings.Contains(fname, "(") {
// 		fi := strings.Index(fname, "(")
// 		fdex := fi - 1
// 		movName = fname[:fdex]
// 	} else {
// 		ddex := len(fname) - 11
// 		movName = fname[ddex:]
// 	}
// 	return
// }

func getMovieYear(apath string) (movyr string) {
	_, filename := path.Split(apath)
	fi := strings.Index(filename, "(")
	fdex := fi + 1
	ldex := strings.LastIndex(filename, ")")
	movyr = filename[fdex:ldex]
	return
}

func moviesUUID() (UUID string) {
	aseed := time.Now()
	aSeed := aseed.UnixNano()
	rand.Seed(aSeed)
	u := rand.Int63n(aSeed)
	UUID = strconv.FormatInt(u, 10)
	// p := strconv.FormatInt(u, 10)
	return
}

//MOVI is exported because I want it so
type MOVI struct {
	ID             bson.ObjectId `bson:"_id,omitempty"`
	DirPath        string        `bson:"dirpath"`
	Filepath       string        `bson:"filepath"`
	MediaID        string        `bson:"mediaid"`
	Movname        string        `bson:"movname"`
	Genre          string        `bson:"genre"`
	Catagory       string        `bson:"catagory"`
	MovFSPath      string        `bson:"movfspath"`
	ThumbPath      string        `bson:"thumbpath"`
	HTTPThumbPath  string		 `bson:"httpthumbpath"`
	MovYear        string        `bson:"movyear"`
}

// 
// ../pkg/mod/github.com/cjsmocjsmo/movgo@v0.1.8/movgo.go:52:14: undefined: httppicpath
// ../pkg/mod/github.com/cjsmocjsmo/movgo@v0.1.8/movgo.go:55:40: undefined: httppicpath
// ../pkg/mod/github.com/cjsmocjsmo/movgo@v0.1.8/movies.go:71:26: undefined: httppicInfo
// 

// GetMovieInfo comment
func GetMovieInfo(apath string, movpicInfo string, httppicInfo string) (MovInfo MOVI) {
	filesystempath := "/media/pi/PiTB/media/" + apath[13:]
	dirp, _ := path.Split(apath)
	MovInfo.ID = bson.NewObjectId()
	MovInfo.DirPath = dirp
	MovInfo.Filepath = apath
	MovInfo.MediaID = moviesUUID()
	MovInfo.Genre = "movies"
	MovInfo.MovFSPath = filesystempath
	MovInfo.ThumbPath = movpicInfo
	MovInfo.HTTPThumbPath = httppicInfo
	
	MovInfo.MovYear = getMovieYear(apath)
	switch {
	case strings.Contains(apath, "SciFi"):
		MovInfo.Catagory = "SciFi"
	case strings.Contains(apath, "Cartoons"):
		MovInfo.Catagory = "Cartoons"
	case strings.Contains(apath, "Godzilla"):
		MovInfo.Catagory = "Godzilla"
	case strings.Contains(apath, "Kingsman"):
		MovInfo.Catagory = "Kingsman"
	case strings.Contains(apath, "StarTrek") && !strings.Contains(apath, " STTV "):
		MovInfo.Catagory = "StarTrek"
	case strings.Contains(apath, "StarWars"):
		MovInfo.Catagory = "StarWars"
	case strings.Contains(apath, "SuperHeros"):
		MovInfo.Catagory = "SuperHeros"
	case strings.Contains(apath, "IndianaJones"):
		MovInfo.Catagory = "IndianaJones"
	case strings.Contains(apath, "Action"):
		MovInfo.Catagory = "Action"
	case strings.Contains(apath, "Comedy"):
		MovInfo.Catagory = "Comedy"
	case strings.Contains(apath, "Drama"):
		MovInfo.Catagory = "Drama"
	case strings.Contains(apath, "JurassicPark"):
		MovInfo.Catagory = "JurassicPark"
	case strings.Contains(apath, "JohnWayne"):
		MovInfo.Catagory = "JohnWayne"
	case strings.Contains(apath, "JohnWick"):
		MovInfo.Catagory = "JohnWick"
	case strings.Contains(apath, "MenInBlack"):
		MovInfo.Catagory = "MenInBlack"
	case strings.Contains(apath, "HarryPotter"):
		MovInfo.Catagory = "HarryPotter"
	case strings.Contains(apath, "Tremors"):
		MovInfo.Catagory = "Tremors"
	case strings.Contains(apath, "Misc"):
		MovInfo.Catagory = "Misc"
	case strings.Contains(apath, "BruceWillis"):
		MovInfo.Catagory = "BruceWillis"
	case strings.Contains(apath, "Pirates"):
		MovInfo.Catagory = "Pirates"
	case strings.Contains(apath, "Fantasy"):
		MovInfo.Catagory = "Fantasy"
	case strings.Contains(apath, "XMen"):
		MovInfo.Catagory = "XMen"
	case strings.Contains(apath, "TomCruize"):
		MovInfo.Catagory = "TomCruize"
	case strings.Contains(apath, "Riddick"):
		MovInfo.Catagory = "Riddick"
	case strings.Contains(apath, "Documentary"):
		MovInfo.Catagory = "Documentary"
	case strings.Contains(apath, "TheRock"):
		MovInfo.Catagory = "TheRock"
	}
	return
}
