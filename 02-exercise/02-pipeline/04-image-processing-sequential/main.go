package main

import (
	"fmt"
	"image"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/disintegration/imaging"
)

type result struct {
	srcImagePath   string
	thumbnailImage *image.NRGBA
	err            error
}

// Image processing - sequential
// Input - directory with images.
// output - thumbnail images
func main() {
	if len(os.Args) < 2 {
		log.Fatal("need to send directory path of images")
	}
	start := time.Now()

	err := setupPipeline(os.Args[1])

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Time taken: %s\n", time.Since(start))
}

func setupPipeline(path string) error {
	done := make(chan struct{})
	defer close(done)

	//first stage
	paths, errc := walkFiles(done, path)

	//second stage
	resized := processImage(done, paths)

	//third stage

	for r := range resized {
		if r.err != nil {
			return r.err
		}

		saveThumbnail(r.srcImagePath, r.thumbnailImage)
	}
	if err := <-errc; err != nil {
		return err
	}
	return nil
}

// walfiles - take diretory path as input
// does the file walk
// generates thumbnail images
// saves the image to thumbnail directory.
func walkFiles(done <-chan struct{}, root string) (<-chan string, <-chan error) {
	paths := make(chan string)
	errc := make(chan error, 1)

	go func() {
		defer close(paths)
		errc <- filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

			// filter out error
			if err != nil {
				return err
			}

			// check if it is file
			if !info.Mode().IsRegular() {
				return nil
			}

			// check if it is image/jpeg
			contentType, _ := getFileContentType(path)
			if contentType != "image/jpeg" {
				return nil
			}
			select {
			case paths <- path:
			case <-done:
				return fmt.Errorf("walk was cancelled")
			}

			return nil
		})
	}()

	return paths, errc
}

// processImage - takes image file as input
// return pointer to thumbnail image in memory.
func processImage(done <-chan struct{}, paths <-chan string) <-chan *result {

	results := make(chan *result)

	thumbnailer := func(done <-chan struct{}, paths <-chan string) {
		for path := range paths {
			srcImage, err := imaging.Open(path)
			if err != nil {
				select {
				case results <- &result{path, nil, err}:
				case <-done:
					return
				}
			}

			// scale the image to 100px * 100px
			thumbnailImage := imaging.Thumbnail(srcImage, 100, 100, imaging.Lanczos)

			select {
			case results <- &result{path, thumbnailImage, nil}:
			case <-done:
				return
			}
		}
	}

	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			thumbnailer(done, paths)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

// saveThumbnail - save the thumnail image to folder
func saveThumbnail(srcImagePath string, thumbnailImage *image.NRGBA) error {
	filename := filepath.Base(srcImagePath)
	dstImagePath := "thumbnail/" + filename

	// save the image in the thumbnail folder.
	err := imaging.Save(thumbnailImage, dstImagePath)
	if err != nil {
		return err
	}
	fmt.Printf("%s -> %s\n", srcImagePath, dstImagePath)
	return nil
}

// getFileContentType - return content type and error status
func getFileContentType(file string) (string, error) {

	out, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err = out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
