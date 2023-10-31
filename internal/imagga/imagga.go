package imagga

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Imagga struct {
	Config     Config
	httpClient *http.Client
}

func New(cfg Config) *Imagga {
	return &Imagga{
		Config:     cfg,
		httpClient: &http.Client{},
	}
}

func (i *Imagga) DetectFace(imgURL string) (faceID string, err error) {
	const imaggaURL = "https://api.imagga.com/v2/faces/detections?return_face_id=1&image_url="

	req, _ := http.NewRequest(http.MethodGet, imaggaURL+imgURL, nil)
	req.SetBasicAuth(i.Config.APIKey, i.Config.APISecret)

	resp, err := i.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	log.Println("imagga response:", resp.Status)

	var body struct {
		Result struct {
			Faces []struct {
				FaceID string `json:"face_id"`
			} `json:"faces"`
		} `json:"result"`
	}

	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return "", err
	}

	if len(body.Result.Faces) == 0 {
		return "", nil
	}

	return body.Result.Faces[0].FaceID, nil
}

func (i *Imagga) FaceSimilarity(faceID1, faceID2 string) (similarity float64, err error) {
	const imaggaURL = "https://api.imagga.com/v2/faces/similarity?face_id=%s&second_face_id=%s"

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf(imaggaURL, faceID1, faceID2), nil)
	req.SetBasicAuth(i.Config.APIKey, i.Config.APISecret)

	resp, err := i.httpClient.Do(req)
	if err != nil {
		return 0, err
	}

	log.Println("imagga response:", resp.Status)

	var body struct {
		Result struct {
			Score float64 `json:"score"`
		} `json:"result"`
	}

	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return 0, err
	}

	log.Println("imagga response similarity score:", body.Result.Score)

	return body.Result.Score, nil
}
