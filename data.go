package main

import (
	"bufio"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type manifestResponse struct {
	Manifest Manifest `json:"photo_manifest"`
}

type Manifest struct {
	Name        string
	LandingDate string `json:"landing_date"`
	LaunchDate  string `json:"launch_date"`
	Status      string
	MaxSol      int    `json:"max_sol"`
	MaxDate     string `json:"max_date"`
	TotalPhotos int    `json:"total_photos"`
	Photos      []Sol
}

type Sol struct {
	Sol         int
	TotalPhotos int `json:"total_photos"`
}

type Photo struct {
	Id         int
	Sol        int
	Rover      string `json:"rover.name"`
	Camera     Camera
	EarthDate  string `json:"earth_date"`
	NasaImgSrc string `json:"img_src"`
	S3ImgSrc   string
}

type Camera struct {
	Name string
}

func (p *Photo) save() (err error) {
	statement := "INSERT INTO photos (id, sol, rover, camera, earthdate, nasaimgsrc, s3imgsrc) VALUES($1, $2, $3, $4, $5, $6, $7) returning id"
	stmt, err := db.Prepare(statement)
	if err != nil {
		return fmt.Errorf("saving image %d to db: %v", p.Id, err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(p.Id, p.Sol, p.Rover, p.Camera.Name, p.EarthDate, p.NasaImgSrc, p.S3ImgSrc).Scan(&p.Id)
	if err != nil {
		return fmt.Errorf("saving image %d to db: %v", p.Id, err)
	} else {
		fmt.Printf("successfully saved data for image %d \n", p.Id)
	}
	return
}

func (p *Photo) copyToS3(region string, bucket string) (err error) {
	res, err := http.Get(p.NasaImgSrc)
	if err != nil {
		return fmt.Errorf("retrieving image %d from nasa: %v", p.Id, err)
	} else {
		defer res.Body.Close()
		reader := bufio.NewReader(res.Body)
		uploader := s3manager.NewUploader(session.New(&aws.Config{Region: aws.String(region)}))
		result, err := uploader.Upload(&s3manager.UploadInput{
			Body:   reader,
			Bucket: aws.String(bucket),
			Key:    aws.String(fmt.Sprintf("%s/%d.jpg", p.Rover, p.Id)),
		})
		if err != nil {
			return fmt.Errorf("uploading image %d to s3: %v", p.Id, err)
		}
		fmt.Printf("completed upload to s3 url: %s \n", result.Location)
		p.S3ImgSrc = result.Location
	}
	return
}
