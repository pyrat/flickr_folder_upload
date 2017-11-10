package main

import (
	"fmt"
	"gopkg.in/masci/flickr.v2"
	"gopkg.in/masci/flickr.v2/photosets"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var image_extensions = map[string]bool{"jpg": true, "gif": true, "png": true, "jpeg": true}

func getFilePaths(base_path string) []os.FileInfo {
	var only_files []os.FileInfo
	files, _ := ioutil.ReadDir(base_path)
	for _, f := range files {
		fmt.Println(f.Name())
		if f.IsDir() == false && isImage(f) {
			only_files = append(only_files, f)
		}
	}
	return only_files
}

func isImage(fileinfo os.FileInfo) bool {
	ext := filepath.Ext(fileinfo.Name())
	lowerStr := strings.ToLower(ext)

	if len(lowerStr) > 0 && image_extensions[lowerStr[1:len(lowerStr)]] == true {
		return true
	} else {
		return false
	}
}

// There should be a "find or create photoset option" (allowing to add to existing photosets etc.)
func uploadImageAndCreateSet(base_path string, fileinfo os.FileInfo, client *flickr.FlickrClient, photoset_name string, privacy bool) (string, error) {

	params := flickr.NewUploadParams()
	params.IsPublic = true
	params.IsFamily = true
	params.IsFriend = true

	path := base_path + "/" + fileinfo.Name()
	fmt.Println("Path to upload:", path)

	resp, err := flickr.UploadFile(client, path, params)
	if err != nil {
		fmt.Println("Failed uploading:", err)
		if resp != nil {
			fmt.Println(resp.ErrorMsg)
		}
	} else {
		fmt.Println("Photo uploaded:", path, resp.ID)
		removeImageFile(base_path, fileinfo)
	}
	return findOrCreatePhotoset(client, photoset_name, resp.ID)
}

func findOrCreatePhotoset(client *flickr.FlickrClient, name string, image_id string) (string, error) {
	exists, photoset_id := photosetExists(client, name)
	if exists == true {
		_, err := photosets.AddPhoto(client, photoset_id, image_id)
		return photoset_id, err
	} else {
		resp, err := photosets.Create(client, name, "", image_id)
		photoset_id = resp.Set.Id
		return photoset_id, err
	}
}

func photosetExists(client *flickr.FlickrClient, name string) (bool, string) {
	response, err := photosets.GetList(client, true, "", 0)
	if err != nil {
		fmt.Println("Error getting photosets list: ", err)
		os.Exit(2)
	}

	for _, p := range response.Photosets.Items {
		if strings.Compare(p.Title, name) == 0 {
			fmt.Println("Found a set to add to.")
			return true, p.Id
		}
	}

	return false, "notfound"

}

func uploadImageToSet(base_path string, fileinfo os.FileInfo, client *flickr.FlickrClient, photoset_id string, privacy bool) (*flickr.BasicResponse, error) {
	params := flickr.NewUploadParams()
	params.IsPublic = !privacy
	params.IsFamily = !privacy
	params.IsFriend = !privacy

	path := base_path + "/" + fileinfo.Name()
	resp, err := flickr.UploadFile(client, path, params)
	if err != nil {
		fmt.Println("Failed uploading:", err)
		if resp != nil {
			fmt.Println(resp.ErrorMsg)
		}
	} else {
		fmt.Println("Photo uploaded:", path, resp.ID)
		removeImageFile(base_path, fileinfo)
	}

	return photosets.AddPhoto(client, photoset_id, resp.ID)
}

func removeImageFile(base_path string, fileinfo os.FileInfo) {
	os.Remove(base_path + "/" + fileinfo.Name())
}
