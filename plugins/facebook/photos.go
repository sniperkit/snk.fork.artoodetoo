package facebook

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/Tympanix/artoodetoo/data"
	"github.com/Tympanix/artoodetoo/event"
)

// Photos listen for new photo uploads or tags on facebook
type Photos struct {
	event.Base
	facebookStyle
	LastSeen string
	Token    Token   `io:"input"`
	Interval float64 `io:"input"`

	URL   string      `io:"output"`
	Image data.Stream `io:"output"`
}

func init() {
	event.Register(new(Photos))
}

// Describe describes what the facebook photo event does
func (p *Photos) Describe() string {
	return "Triggers whenever you are tagged in a photo on facebook"
}

// Listen starts listening for new photo uploads and tags on facebook
func (p *Photos) Listen(stop <-chan struct{}) error {
	if p.Interval < 1000 {
		return errors.New("Interval must be above 1 second")
	}

	err := p.updateLast()

	if err != nil {
		return err
	}

	ticker := time.NewTicker(time.Duration(p.Interval) * time.Millisecond)

	for {
		select {
		case <-ticker.C:
			if err := p.checkNew(); err != nil {
				return err
			}
		case <-stop:
			return nil
		}
	}
}

func (p *Photos) checkNew() error {
	photos, err := p.getPhotos()

	if err != nil {
		return err
	}

	for _, photo := range photos {
		if photo.ID != p.LastSeen {
			details, err := p.getPhotoDeails(photo.ID)
			if err != nil {
				return err
			}
			if len(details.Images) > 0 {
				p.URL = details.Images[0].Source
				stream, err := data.FromURL(p.URL)
				if err != nil {
					return err
				}
				p.Image = stream
				p.Trigger()
			} else {
				return errors.New("No images found")
			}
		} else {
			break
		}
	}
	if len(photos) > 0 {
		p.LastSeen = photos[0].ID
	}
	return nil
}

type photoInfo struct {
	CreatedTime string    `json:"created_time"`
	ID          string    `json:"id"`
	Created     time.Time `json:"-"`
}

func (p *photoInfo) UnmarshalJSON(data []byte) error {
	type alias photoInfo
	var photo alias
	if err := json.Unmarshal(data, &photo); err != nil {
		return err
	}
	*p = photoInfo(photo)
	time, err := time.Parse(TIME, p.CreatedTime)
	if err != nil {
		return err
	}
	p.Created = time
	return nil
}

func (p *Photos) updateLast() error {
	photos, err := p.getPhotos()
	if err != nil {
		return err
	}
	if len(photos) > 0 {
		p.LastSeen = photos[0].ID
	}
	return nil
}

func (p *Photos) getPhotos() (photos []*photoInfo, err error) {
	if _, err = getData("/me/photos", p.Token, Options{}, &photos); err != nil {
		return
	}

	return
}

type photoDetails struct {
	ID     string   `json:"id"`
	Images []*image `json:"images"`
}

type image struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	Source string `json:"source"`
}

func (p *Photos) getPhotoDeails(id string) (details photoDetails, err error) {
	options := Options{
		Fields: "images",
	}
	if err = getNode(id, p.Token, options, &details); err != nil {
		return
	}
	return
}
