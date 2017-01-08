package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func serveIndex(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./static/index.html")
	if err != nil {
		log.Fatal(err)
	}
	rover := "spirit"
	limit := 1
	var p photo
	err = db.QueryRow("SELECT id, sol, rover, camera, earthdate, s3imgsrc FROM photos WHERE rover=$1 order by sol desc, id desc limit $2", rover, limit).Scan(&p.Id, &p.Sol, &p.Rover, &p.Camera, &p.EarthDate, &p.S3ImgSrc)
	if err != nil {
		log.Fatal(err)
	}
	if err := t.Execute(w, p); err != nil {
		log.Fatal(err)
	}
}

func getRoverPhotos(w http.ResponseWriter, r *http.Request) {
	var photos []photo
	rover := mux.Vars(r)["rover"]
	page, err := strconv.Atoi(mux.Vars(r)["page"])
	limit := 10
	page = page*limit + 1
	rows, err := db.Query("SELECT id, sol, rover, camera, earthdate, s3imgsrc FROM photos WHERE rover=$1 order by sol desc, id desc limit $2 offset $3", rover, limit, page)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var p photo
		err = rows.Scan(&p.Id, &p.Sol, &p.Rover, &p.Camera, &p.EarthDate, &p.S3ImgSrc)
		if err != nil {
			log.Fatal(err)
		}
		photos = append(photos, p)
	}
	j, err := json.Marshal(photos)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(j)
}
