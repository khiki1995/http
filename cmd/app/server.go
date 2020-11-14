package app

import (
	"io"
	"os"
	"path/filepath"
	"encoding/json"
	"strconv"
	"log"
	"github.com/khiki1995/http/pkg/banners"
	"net/http"
)

type Server struct {
	mux *http.ServeMux
	bannersSvc *banners.Service

}
func NewServer(mux *http.ServeMux, bannersSvc *banners.Service) *Server {
	return &Server{mux: mux, bannersSvc: bannersSvc}
}
func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}
func (s *Server) Init(){
	s.mux.HandleFunc("/banners.getAll", s.handleGetAllBanners)
	s.mux.HandleFunc("/banners.getById", s.handleGetBannerByID)
	s.mux.HandleFunc("/banners.save", s.handleSaveBanner)
	s.mux.HandleFunc("/banners.removeById", s.handleRemoveByID)
}
func (s *Server) handleGetBannerByID(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	item, err := s.bannersSvc.ByID(request.Context(), id)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type","application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}
func (s *Server) handleSaveBanner(writer http.ResponseWriter, request *http.Request) {
	idParam := request.PostFormValue("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	fileName, err := UploadImage(writer, request, idParam)
	if err != nil {
		log.Print(err)
		return
	}
	var banner = &banners.Banner{
		ID: 		 id,
		Title: 	  	 request.PostFormValue("title"),
		Content: 	 request.PostFormValue("content"),
		Button:		 request.PostFormValue("button"),
		Link:		 request.PostFormValue("link"),
		Image:		 fileName,
	}
	item, err := s.bannersSvc.Save(request.Context(),banner)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type","application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}
func (s *Server) handleGetAllBanners(writer http.ResponseWriter, request *http.Request) {
	items, _ := s.bannersSvc.All(request.Context())
	data, err := json.Marshal(items)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type","application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}
func (s *Server) handleRemoveByID(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")
	id, err := strconv.ParseInt(idParam,10,64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	item, err := s.bannersSvc.RemoveByID(request.Context(), id)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type","application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}
func UploadImage(writer http.ResponseWriter, request *http.Request, fileName string) (string, error) {
	if err := request.ParseMultipartForm(10 * 1024 * 1024); err != nil {
		log.Print(err)
		return "", err
	}
	file, handler, err := request.FormFile("image")
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer file.Close()
	if fileName == "0" {
		fileName = "1"
	}
	fileName = fileName + filepath.Ext(handler.Filename)
	absPath, err := filepath.Abs("web/banners/"+fileName)
	if err != nil {
		log.Print(err)
		return "", err
	}
	dst, err := os.Create(absPath)
	defer dst.Close()
	if err != nil {
		http.Error(writer, err.Error(),http.StatusInternalServerError)
		return "", err
	}

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return "", err
	}
	return fileName, nil
}