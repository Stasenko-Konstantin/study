package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/fs"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type gameDir = []fs.FileInfo

type game struct {
	name        string
	description string
	number      string
	image1      string
	image2      string
	file        string
	comments    string
}

func newGame(d fs.FileInfo) (*game, error) {
	var r game
	r.number = d.Name()
	file, err := os.ReadFile("assets/games/" + r.number + "/text")
	if err != nil {
		return nil, err
	}
	text := strings.Split(string(file), "\n")
	r.name, r.description = text[0], sliceToStr(text[1:])
	r.image1 = "assets/games/" + r.number + "/image1.jpg"
	r.image2 = "assets/games/" + r.number + "/image2.jpg"
	r.file = "assets/games/" + r.number + "/game.nes"
	r.comments = "assets/games/" + r.number + "/comments"
	return &r, nil
}

var (
	gameDirs gameDir
	games    []game
	slide    = []string{
		`<div class="nk-image-slider-item">
                <img src="assets/games/`, // № игры
		`/image2.jpg" alt="" class="nk-image-slider-img"
                     data-thumb="assets/games/`, // № игры
		`/image1.jpg">
                <div class="nk-image-slider-content">
                    <h3 class="h4">`, // название
		`</h3>
                    <p class="text-white">`, // описание
		`</p>
                    <a href="#" class="btn btn-slide btn-rounded btn-hover-color-main-1 ">Далее</a>
                </div>
            </div>`,
	}
	mainList = []string{
		`<div class="col-md-6">
                            <div class="blog-post">
                                <a href="article" class="post-img post-img-main">
                                    <img src="assets/games/`, // № игры
		`/image1.jpg" alt="">
                                </a>
                                <div class="gap"></div>
                                <h2 class="post-title h4"><a href="article">`, // название
		`</a></h2>
                                <div class="gap"></div>
                                <div class="post-text">
                                    <p>`, // описание
		`</p>
                                </div>
                                <div class="gap"></div>
                                <a href="article"
                                   class="btn btn-main btn-rounded btn-hover-color-main-1">Далее</a>
                            </div>
                        </div>`,
	}
	gameList = []string{
		`<div class="blog-post">
                        <div class="row vertical-gap">
                            <div class="col-md-5 col-lg-6">
                                <a href="article" class="post-img">
                                    <img src="assets/games/`, // № игры
		`/image1.jpg" alt="">
								</a>
                            </div>
                            <div class="col-md-7 col-lg-6">
                                <h2 class="post-title h4"><a href="article">`, // название
		`</a></h2>
                                <div class="post-text">
                                    <p>`, // описание
		`</p>
                                </div>
                            </div>
                        </div>
                    </div>`,
	}
)

func getRandGame(g *gameDir) fs.FileInfo {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(*g))
	r := (*g)[n]
	tmp := append((*g)[:n], (*g)[n+1:]...)
	*g = tmp
	return r
}

func sliceToStr(s []string) string {
	var r string
	for _, i := range s {
		r += i
	}
	return r
}

func splitSentences(d string) string {
	var (
		r     string
		count = 0
	)
	for n, i := range d {
		if i == '.' || i == '?' || i == '!' {
			count += 1
		}
		if count == 3 {
			r = d[:n] + ".."
			break
		}
	}
	return r
}

func handler(h string, w http.ResponseWriter, r *http.Request) {
	file, err := os.ReadFile("assets/" + h + ".html")
	if err != nil {
		fmt.Fprintf(w, "Server error: "+err.Error())
	}
	fmt.Fprintf(w, string(file))
}

func gamesHandler(w http.ResponseWriter, r *http.Request) {
	var gamesList []string
	file, err := os.ReadFile("assets/games.html")
	if err != nil {
		fmt.Fprintf(w, "Server error: "+err.Error())
		os.Exit(1)
	}
	gamesList = strings.Split(string(file), "Список игр")
	var tmp []string
	for i := 0; i < len(gameDirs); i++ {
		g, err := newGame(gameDirs[i])
		if err != nil {
			fmt.Fprintf(w, "cant make game struct: "+err.Error())
			os.Exit(1)
		}
		tmp = append(tmp, gameList[0]+g.number+gameList[1]+g.name+gameList[2]+splitSentences(g.description)+gameList[3]+"ㅤ")
	}
	gamesList = append([]string{}, gamesList[0]+sliceToStr(tmp)+gamesList[1])
	fmt.Fprintf(w, sliceToStr(gamesList))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var index []string
	file, err := os.ReadFile("assets/index.html")
	if err != nil {
		fmt.Fprintf(w, "Server error: "+err.Error())
		os.Exit(1)
	}
	index = strings.Split(string(file), "Слайдер")
	thisGames := append(gameDir{}, gameDirs...)
	var tmp []string
	for i := 0; i < 4; i++ {
		g, err := newGame(getRandGame(&thisGames))
		if err != nil {
			fmt.Fprintf(os.Stderr, "cant make game struct: "+err.Error())
			os.Exit(1)
		}
		tmp = append(tmp, slide[0]+g.number+slide[1]+g.number+slide[2]+g.name+slide[3]+g.description+slide[4])
	}
	index = append([]string{}, index[0]+sliceToStr(tmp)+index[1])
	index = strings.Split(index[0], "Главный список")
	tmp = []string{}
	for i := 0; i < 6; i++ {
		g, err := newGame(getRandGame(&thisGames))
		if err != nil {
			fmt.Fprintf(os.Stderr, "cant make game struct: "+err.Error())
			os.Exit(1)
		}
		tmp = append(tmp, mainList[0]+g.number+mainList[1]+g.name+mainList[2]+splitSentences(g.description)+mainList[3])
	}
	index = append([]string{}, index[0]+sliceToStr(tmp)+index[1])
	fmt.Fprintf(w, sliceToStr(index))
}

func getGames() (gameDir, error) {
	dirs, err := ioutil.ReadDir("./assets/games")
	if err != nil {
		return gameDir{}, err
	}
	return dirs, nil
}

func main() {
	fmt.Println("server start")
	dirs, err := getGames()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cant read gameDirs directories: "+err.Error())
		os.Exit(1)
	}
	gameDirs = dirs
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/games", gamesHandler).Methods("GET")
	r.HandleFunc("/article",
		func(w http.ResponseWriter, r *http.Request) { handler("article", w, r) }).Methods("GET")
	staticFileDir := http.Dir("./assets/")
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDir))
	r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")
	http.ListenAndServe(":8080", r)
}
