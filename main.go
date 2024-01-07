package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func generateName(Url string) string {
	u, err := url.Parse(Url)
	if err != nil {
		panic(err)
	}
	outputFilePath := filepath.Base(Url)
	if outputFilePath == "" {
		outputFilePath = u.Host
	}
	return foundName(outputFilePath, outputFilePath, 0)
}

func foundName(fileName, original string, count int) string {
	for {
		if _, err := os.Stat(fileName); err == nil {
			count++
			name, extention := splitName(original)
			return foundName(fmt.Sprintf("%s (%d).%s", name, count, extention), original, count)
		} else if os.IsNotExist(err) {
			return fileName
		} else {
			fmt.Println("Error:", err)
			return fileName
		}
	}
}

func splitName(na string) (string, string) {
	if !strings.Contains(na, ".") {
		return na, ""
	}
	strs := strings.Split(na, ".")
	result := ""
	extention := ""
	for n, s := range strs {
		if n == len(strs)-1 {
			extention = s
			break
		}
		result += fmt.Sprintf("%s.", s)
	}
	return result[:len(result)-1], extention
}

func main() {
	Url := "https://ash-speed.hetzner.com/10GB.bin"
	fmt.Print("Enter Url:")
	fmt.Scanln(&Url)

	outputFilePath := generateName(Url)
	err := downloadFile(Url, outputFilePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Downloaded successfully!!!")
}

func downloadFile(url, outputPath string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,ar;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", req.Host)
	req.Header.Set("Sec-Ch-Ua", `"Not/A)Brand";v="99", "Google Chrome";v="115", "Chromium";v="115"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP GET request failed with status: %s", response.Status)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	startTime := time.Now()
	var downloadedBytes int64
	buffer := make([]byte, 2*1024)

	for {
		n, err := response.Body.Read(buffer)
		if n > 0 {
			_, writeErr := file.Write(buffer[:n])
			if writeErr != nil {
				return writeErr
			}
			downloadedBytes += int64(n)

			elapsedTime := time.Since(startTime).Seconds()
			downloadSpeedMBps := float64(downloadedBytes) / (elapsedTime * 1024 * 1024)

			fmt.Printf("\r%.2f%% downloaded | %.2f MB/s  ", float64(downloadedBytes)/float64(response.ContentLength)*100, downloadSpeedMBps)
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}

	return nil
}
