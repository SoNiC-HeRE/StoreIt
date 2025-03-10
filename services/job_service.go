package services

import (
    "math/rand"
    "net/http"
    "time"
    "image"
    _ "image/jpeg"
    _ "image/png"
    "retail_pulse/repository"
)

func ProcessImage(imgURL string) (int, error) {
    response, err := http.Get(imgURL)
    if err != nil {
        return 0, err
    }
    defer response.Body.Close()

    imgConfig, _, err := image.DecodeConfig(response.Body)
    if err != nil {
        return 0, err
    }

    perimeter := 2 * (imgConfig.Height + imgConfig.Width)
    time.Sleep(time.Duration(rand.Intn(300)+100) * time.Millisecond)
    
    return perimeter, nil
}
