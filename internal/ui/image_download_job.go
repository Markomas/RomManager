package ui

import (
	"RomManager/internal/config"
	"RomManager/internal/util"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
)

type ImageDownloadJob struct {
	running                  atomic.Bool
	downloadingImages        util.OrderedSet[string]
	downloadingImagesM       sync.Mutex
	downloadingImagesRunning atomic.Bool
	config                   *config.Config
}

func NewImageDownloadJob(c *config.Config) *ImageDownloadJob {
	return &ImageDownloadJob{
		running:                  atomic.Bool{},
		downloadingImages:        *util.NewOrderedSet[string](),
		downloadingImagesM:       sync.Mutex{},
		downloadingImagesRunning: atomic.Bool{},
		config:                   c,
	}
}

func (j *ImageDownloadJob) Start() {
	// start only if not running
	if !j.running.CompareAndSwap(false, true) {
		return
	}

	go func() {
		defer j.running.Store(false) // allow restart after finishing
		for j.run() {

		}
	}()
}

func (j *ImageDownloadJob) run() bool {
	j.downloadingImagesM.Lock()
	defer j.downloadingImagesM.Unlock()

	url, exists := j.downloadingImages.Last()
	if !exists {
		return false
	}

	j.downloadImage(url)
	j.downloadingImages.Remove(url)
	return true

}

func (j *ImageDownloadJob) CheckIfFileIsDownloaded(url string) bool {
	localPath := j.UrlToLocalPath(url)
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		return false
	}
	return true
}

func (j *ImageDownloadJob) downloadImage(url string) {
	localPath := j.UrlToLocalPath(url)
	cacheDir := filepath.Join(j.config.System.CachePath, "rommanager", "images")

	if _, err := os.Stat(localPath); os.IsExist(err) {
		return
	}

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		log.Printf("Failed to create image cache directory %s: %v", cacheDir, err)
		return
	}

	tempPath := localPath + ".tmp"
	// Ensure the temporary file is removed on exit, especially if an error occurs
	defer os.Remove(tempPath)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to download image from %s: %v", url, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to download image from %s: status code %d", url, resp.StatusCode)
		return
	}

	outFile, err := os.Create(tempPath)
	if err != nil {
		log.Printf("Failed to create temporary file %s: %v", tempPath, err)
		return
	}

	_, err = io.Copy(outFile, resp.Body)
	// Close the file handle before attempting to rename it.
	outFile.Close()
	if err != nil {
		log.Printf("Failed to save image to %s: %v", tempPath, err)
		return
	}

	// Atomically rename the temporary file to the final url.
	if err := os.Rename(tempPath, localPath); err != nil {
		log.Printf("Failed to rename temporary file %s to %s: %v", tempPath, localPath, err)
	}
}

func (j *ImageDownloadJob) UrlToLocalPath(path string) string {
	cacheDir := filepath.Join(j.config.System.CachePath, "rommanager", "images")
	hash := sha1.New()
	hash.Write([]byte(path))
	ext := filepath.Ext(path)
	if strings.Contains(ext, ".php") {
		ext = ".png"
	}
	fileName := hex.EncodeToString(hash.Sum(nil)) + ext
	localPath := filepath.Join(cacheDir, fileName)
	return localPath
}

func (j *ImageDownloadJob) AddDownloadJob(path string) {
	j.downloadingImagesM.Lock()
	defer j.downloadingImagesM.Unlock()

	if !j.downloadingImages.Exists(path) {
		j.downloadingImages.Add(path)
	} else {
		j.downloadingImages.MoveToFront(path)
	}
}
