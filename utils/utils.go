package utils

import (
	"errors"
	"mime/multipart"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func IncludesObjectId(arr []primitive.ObjectID, id primitive.ObjectID) bool {
	for _, val := range arr {
		if val.Hex() == id.Hex() {
			return true
		}
	}

	return false
}

func IsImage(fileHeader *multipart.FileHeader) bool {
	fileType := fileHeader.Header.Get("Content-Type")
	return strings.HasPrefix(fileType, "image/")
}

func HandleImage(req *http.Request, formName string) (string, int ,error) {
	file, fileHeader, err := req.FormFile(formName)
	if err != nil {
		return "", http.StatusBadRequest, err
	}
	defer file.Close()

	isImage := IsImage(fileHeader)
	if !isImage {
		return "", http.StatusBadRequest, errors.New(formName + " can only be an image")
	}
	
	filePath, err := UploadFileToS3(file, fileHeader)
	if err != nil {
		return "", http.StatusInternalServerError, errors.New(formName + " can only be an image")
	}
	
	return filePath, http.StatusOK, nil
}