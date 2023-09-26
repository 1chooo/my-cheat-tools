// Created Date: 2023/09/26
// Author: @1chooo (Hugo ChunHo Lin)
// Version: v0.0.1

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gocv.io/x/gocv"
)

type VideoFrameExtractor struct {
	videoPath     string
	outputFolder  string
	frameRate     int
	video         *gocv.VideoCapture
	frameCount    int
}

func NewVideoFrameExtractor(videoPath, outputFolder string, frameRate int) *VideoFrameExtractor {
	return &VideoFrameExtractor{
		videoPath:    videoPath,
		outputFolder: outputFolder,
		frameRate:    frameRate,
	}
}

func (v *VideoFrameExtractor) Initialize() error {
	// Create the output folder if it doesn't exist
	if err := os.MkdirAll(v.outputFolder, os.ModePerm); err != nil {
		return err
	}

	// Open the video file
	video, err := gocv.VideoCaptureFile(v.videoPath)
	if err != nil {
		return err
	}
	v.video = video

	return nil
}

func (v *VideoFrameExtractor) ProcessVideo() {
	for {
		frame := gocv.NewMat()
		if ok := v.video.Read(&frame); !ok {
			break
		}
		defer frame.Close()

		// Capture a frame every v.frameRate frames
		if v.frameCount%v.frameRate == 0 {
			imageFilename := filepath.Join(v.outputFolder, fmt.Sprintf("frame_%04d.jpg", v.frameCount/v.frameRate))
			if ok := gocv.IMWrite(imageFilename, frame); !ok {
				fmt.Printf("Error writing image: %v\n", imageFilename)
			}
		}

		v.frameCount++
	}
}

func (v *VideoFrameExtractor) Close() {
	v.video.Close()
}

func main() {
	videoPath := "your_video.mp4"
	outputFolder := "output_images"
	frameRate := 100

	extractor := NewVideoFrameExtractor(videoPath, outputFolder, frameRate)
	if err := extractor.Initialize(); err != nil {
		fmt.Printf("Error initializing video frame extractor: %v\n", err)
		return
	}
	defer extractor.Close()

	extractor.ProcessVideo()

	fmt.Printf("Total %d images saved.\n", extractor.frameCount/extractor.frameRate)
}
