package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"math/rand"
	"os"
	"strings"
	"time"
)

func SaveProfileImage(base64Data string) (string, error) {
	// Decode Base64
	// Data URI format: data:image/png;base64,.....
	idx := strings.IndexByte(base64Data, ',')
	if idx != -1 {
		base64Data = base64Data[idx+1:]
	}

	dec, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", err
	}

	// Decode Image (Automatically detects format if imported)
	img, _, err := image.Decode(bytes.NewReader(dec))
	if err != nil {
		return "", err
	}

	// Ensure directory exists
	if err := os.MkdirAll("uploads/images", 0755); err != nil {
		return "", err
	}

	// create file
	// Save as JPEG for better compression
	filename := fmt.Sprintf("%d-%d.jpeg", time.Now().Unix(), rand.Intn(10000))
	filepath := fmt.Sprintf("uploads/images/%s", filename)

	out, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Compress and save image
	if err := jpeg.Encode(out, img, &jpeg.Options{Quality: 75}); err != nil {
		return "", err
	}

	return "/" + filepath, nil
}
