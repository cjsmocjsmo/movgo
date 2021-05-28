package movgo

import (
	"fmt"
	//because I want it
	"github.com/disintegration/imaging"
	"github.com/globalsign/mgo/bson"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"path"
	"strconv"
	"time"
)

//ThumbInFo struct exported to setup
type ThumbInFo struct {
	ID            bson.ObjectId `bson:"_id,omitempty"`
	MovName       string        `bson:"movname"`
	BasePath      string        `bson:"baspath"`
	DirPATH       string        `bson:"dirpath"`
	ThumbPath     string        `bson:"thumbpath"`
	ThumbID       string        `bson:"thumbid"`
	HTTPThumbPath string 		`bson:"httpthumbpath"`
}

//UUID holds the unique identifier for the file
// a comment
func UUID() (UUID string) {
	aSeed := time.Now()
	aseed := aSeed.UnixNano()
	rand.Seed(aseed)
	u := rand.Int63n(aseed)
	UUID = strconv.FormatInt(u, 10)
	return
}

func myPathSplit(myPath string) (DirPath string, BaseNAme string, MOvName string, Ext string) {
	DirPath, BaseNAme = path.Split(myPath)
	Ext = BaseNAme[3:]
	factor := len(BaseNAme) - 4
	MOvName = BaseNAme[:factor]
	return
}

func getServerAddr() (addr string) {
	addr = os.Getenv("MEDIACENTER_SERVER_ADDRESS")
	return
}

func getServerPort() (port string) {
	port = os.Getenv("MEDIACENTER_SERVER_PORT")
	return
}

func getThumbPath() (tpath string) {
	tpath = os.Getenv("MEDIACENTER_THUMBNAIL_PIC_PATH")
	return
}

//CreateMoviesThumbnail exported to setup
func CreateMoviesThumbnail(p string) (ThumbINFO ThumbInFo) {
	dirpath, basepath, movname, ext := myPathSplit(p)
	MSA := getServerAddr()
	MSP := getServerPort()
	MTPP := getThumbPath()
	// BP := "/" + url.QueryEscape(basepath)
	// thumbpathtwo := MSA + ":" + MSP + MTPP + BP
	// thumbpathone := "./static/" + basepath
	var BP string = "/" + basepath
	var thumbpathtwo string = MSA + ":" + MSP + MTPP + BP
	var thumbpathone string = "static/" + basepath



	ThumbINFO.ID = bson.NewObjectId()
	ThumbINFO.MovName = movname
	ThumbINFO.BasePath = basepath
	ThumbINFO.DirPATH = dirpath
	ThumbINFO.HTTPThumbPath = thumbpathtwo
	ThumbINFO.ThumbPath = thumbpathone
	ThumbINFO.ThumbID = UUID()

	if ext == ".txt" {
		fmt.Print("what the fuck a text file remove it")
		os.Remove(p)
	} else if ext == ".srt" {
		os.Remove(p)
	} else {

		_, err := os.Stat(thumbpathone)
		if err == nil {
			log.Printf("FILE %s EXISTS", thumbpathone)
		} else if os.IsNotExist(err) {
			pic, err := imaging.Open(p)
			if err != nil {
				log.Printf("\n this is file Open error jpgthumb %v \n", p)
			}
			thumbImage := imaging.Resize(pic, 0, 250, imaging.Lanczos)
			err = imaging.Save(thumbImage, thumbpathone)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			log.Printf("file %s stat error: %v", thumbpathone, err)
		}
		cmtses := MovDBcon()
		defer cmtses.Close()
		CMTc := cmtses.DB("movbsthumb").C("movbsthumb")
		err = CMTc.Insert(&ThumbINFO)
		if err != nil {
			fmt.Println(err)
		}
	}
	return
}

//NoArtList exported to setup
var NoArtList []string
//a comment
//FindPicPaths exported to setup
func FindPicPaths(mpath string, noartpicpath string) (result, result2 string) {
	_, _, movename, _ := myPathSplit(mpath)
	Tses := MovDBcon()
	defer Tses.Close()
	MTc := Tses.DB("movbsthumb").C("movbsthumb")
	b1 := bson.M{"movname": movename}
	b2 := bson.M{"_Id": 0}
	var ThumbI []map[string]string
	err := MTc.Find(b1).Select(b2).All(&ThumbI)
	if err != nil {
		log.Println(err)
	}
	LenI := len(ThumbI)
	// fmt.Printf("THIS IS THUMBI %s \n", ThumbI)
	// var result string
	// var result2 string
	if LenI == 0 {
		NoArtList = append(NoArtList, mpath)
		result = noartpicpath
		result2 = noartpicpath
	} else {
		result = ThumbI[0]["thumbpath"]
		result2 = ThumbI[0]["httpthumbpath"]
	}
	fmt.Printf("this is result %s", result)
	return
}

