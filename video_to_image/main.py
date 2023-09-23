# -*- coding: utf-8 -*-
'''
Created Date: 2023/09/23
Author: @1chooo (Hugo ChunHo Lin)
Version: v0.0.1
'''

import cv2
import os

# Video file path
video_path = 'your_video.mp4'

# Output folder for images
output_folder = 'output_images/'

# Create the output folder if it doesn't exist
os.makedirs(output_folder, exist_ok=True)

try:
    # Open the video file
    cap = cv2.VideoCapture(video_path)

    # Check if the video opened successfully
    if not cap.isOpened():
        raise Exception("Error: Could not open the video file.")

    frame_count = 0
    frame_rate = 100  # Capture a frame every 0.01 seconds

    while True:
        ret, frame = cap.read()

        if not ret:
            break

        # Capture a frame every 0.01 seconds
        if frame_count % frame_rate == 0:
            image_filename = f"{output_folder}frame_{frame_count // frame_rate:04d}.jpg"
            cv2.imwrite(image_filename, frame)

        frame_count += 1

    print(f"Total {frame_count // frame_rate} images saved.")

except cv2.error as e:
    print(f"OpenCV Error: {e}")
except Exception as e:
    print(f"Error: {e}")
finally:
    # Close the video file
    cap.release()
