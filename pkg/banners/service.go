package banners

import (
	"strconv"
	"os"
	"net/http"
	"errors"
	"context"
	"sync"
	"path/filepath"
	 "io"
)

type Service struct {
	mu sync.RWMutex
	items []*Banner
}

func NewService() *Service {
	return &Service{items: make([]*Banner,0)}
}

var BannerID int64 = 0

type Banner struct {
	ID int64
	Title string
	Content string
	Button string
	Link string
	Image string
}

func (s *Service) All(ctx context.Context) ([]*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.items, nil
}

func (s *Service) Save(request *http.Request, item *Banner) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if item.ID == 0 {
		BannerID += 1
		fileName, err := UploadImage(request, BannerID)
		if err != nil {
			return nil, err
		}
		item.ID = BannerID
		item.Image = fileName
		s.items = append(s.items, item)
		return item, nil
	}
	for i, banner := range s.items {
		if  banner.ID == item.ID {
			fileName, err := UploadImage(request, item.ID)
			if err != nil {
				return nil, err
			}
			item.Image = fileName
			s.items[i] = item
			return item, nil
		}
	}
	return nil, errors.New("Item not found")
}

func (s *Service) RemoveByID(ctx context.Context, id int64) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i, banner := range s.items {
		if banner.ID == id {
			s.items = append(s.items[:i], s.items[i+1:]...)
			return banner, nil
		}
	}
	return nil, errors.New("Item not found")
}

func (s *Service) ByID(ctx context.Context, id int64) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, banner := range s.items {
		if banner.ID == id {
			return banner, nil
		}
	}
	return nil, errors.New("item not found")
}

func UploadImage(request *http.Request, bannerID int64) (string, error) {
	if err := request.ParseMultipartForm(10 * 1024 * 1024); err != nil {
		return "", err
	}
	file, handler, err := request.FormFile("image")
	if err != nil {
		return "", err
	}
	defer file.Close()
	fileName := strconv.FormatInt(bannerID,10) + filepath.Ext(handler.Filename)
	absPath, err := filepath.Abs("web/banners/" + fileName)
	if err != nil {
		return "", err
	}
	dst, err := os.Create(absPath)
	defer dst.Close()
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}
	return fileName, nil
}