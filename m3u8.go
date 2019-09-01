package geektimedl

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func parseM3u8(m3u8Path string) ([]string, error) {
	f, err := os.Open(m3u8Path)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var tsFilenames []string
	for _, l := range strings.Split(string(data), "\n") {
		l = strings.TrimSpace(l)
		if l == "" || strings.HasPrefix(l, "#") {
			continue
		}
		tsFilenames = append(tsFilenames, l)
	}

	return tsFilenames, nil
}

func listM3U8Paths(dir string) []string {
	var m3u8Paths []string
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".m3u8") {
			m3u8Paths = append(m3u8Paths, path)
		}
		return nil
	})
	return m3u8Paths
}

// M3U8ToMP4Converter converts the videos with the m3u8 format to the mp4
type M3U8ToMP4Converter struct {
	outputDir, inputDir string
}

// NewM3U8ToMP4Converter creates the converter to convert the m3u8-video to the mp4
func NewM3U8ToMP4Converter(inputDir, outputDir string) (*M3U8ToMP4Converter, error) {
	if inputDir == "" {
		return nil, errors.New("empty input dir")
	}

	if outputDir == "" {
		return nil, errors.New("empty output dir")
	}

	err := makeSureDirExist(outputDir)
	if err != nil {
		return nil, err
	}

	return &M3U8ToMP4Converter{outputDir, inputDir}, nil
}

// Run doing the  converting
func (c *M3U8ToMP4Converter) Run() {
	for _, m := range listM3U8Paths(c.inputDir) {
		outputPath := filepath.Join(c.outputDir, filepath.Base(filepath.Dir(m))) + ".mp4"
		fmt.Printf("convert %s, to %s\n", m, outputPath)
		ret, err := m3u8ToMP4(m, outputPath)

		msg := "success!"
		if err != nil {
			msg = err.Error()
		}
		fmt.Printf("result:%s, message:%s\n", ret, msg)
	}
}