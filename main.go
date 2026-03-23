package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/corona10/goimagehash"
)

type MediaHash struct {
	Size  uint64
	VHash string // For videos
	PHash uint64 // For photos
}

func main() {
	// 1. Setup flags
	dirPtr := flag.String("dir", ".", "Path to folder for scanning")
	deletePtr := flag.Bool("delete", false, "Delete found duplicates (default false)")
	flag.Parse()

	absPath, err := filepath.Abs(*dirPtr)
	if err != nil {
		fmt.Printf("Error with path: %v\n", err)
		return
	}

	fmt.Printf("Scanning: %s\n", absPath)
	if *deletePtr {
		fmt.Println("Delete mode enabled (use with caution!)")
	} else {
		fmt.Println("Read-only mode (use -delete to enable deletion)")
	}

	db := make(map[MediaHash]string)
	duplicatesCount := 0

	err = filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		var mHash MediaHash
		mHash.Size = uint64(info.Size())

		switch ext {
		case ".jpg", ".jpeg", ".png", ".heic":
			phash, err := getPhotoHash(path, ext)
			if err != nil {
				return nil
			}
			mHash.PHash = phash
		case ".mp4", ".mov", ".m4v":
			vhash, err := getVideoQuickHash(path, info.Size())
			if err != nil {
				return nil
			}
			mHash.VHash = vhash
		default:
			return nil
		}

		if original, exists := db[mHash]; exists {
			fmt.Printf("🗑️ Duplicate: %s\n   (Original: %s)\n", filepath.Base(path), filepath.Base(original))
			duplicatesCount++

			if *deletePtr {
				if err := os.Remove(path); err != nil {
					fmt.Printf("Error removing %s: %v\n", path, err)
				} else {
					fmt.Println("   Deleted")
				}
			}
		} else {
			db[mHash] = path
		}
		return nil
	})

	fmt.Printf("\nFound duplicates: %d\n", duplicatesCount)
}

func getPhotoHash(path string, ext string) (uint64, error) {
	var img image.Image
	var decodeErr error

	if ext == ".heic" {
		// use macOS built-in sips to convert HEIC to JPEG in-memory
		tmp := path + ".tmp.jpg"
		defer os.Remove(tmp)
		cmd := exec.Command("sips", "-s", "format", "jpeg", path, "--out", tmp)
		if err := cmd.Run(); err != nil {
			return 0, err
		}
		f, _ := os.Open(tmp)
		defer f.Close()
		img, _, decodeErr = image.Decode(f)
	} else {
		f, err := os.Open(path)
		if err != nil {
			return 0, err
		}
		defer f.Close()
		img, _, decodeErr = image.Decode(f)
	}

	if decodeErr != nil {
		return 0, decodeErr
	}

	hash, _ := goimagehash.PerceptionHash(img)
	return hash.GetHash(), nil
}

func getVideoQuickHash(path string, size int64) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if size < 200000 {
		io.Copy(h, f)
	} else {
		io.CopyN(h, f, 100000)
		f.Seek(-100000, io.SeekEnd)
		io.CopyN(h, f, 100000)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
