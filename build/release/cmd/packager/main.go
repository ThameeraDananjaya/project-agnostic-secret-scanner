package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type entry struct {
	absolute string
	relative string
	mode     os.FileMode
	size     int64
}

func main() {
	root := flag.String("root", "", "regular-file-only input directory")
	output := flag.String("output", "", "archive output path")
	format := flag.String("format", "", "tar.gz or zip")
	epoch := flag.Int64("epoch", 0, "fixed Unix timestamp")
	flag.Parse()
	if flag.NArg() != 0 || *root == "" || *output == "" || (*format != "tar.gz" && *format != "zip") || *epoch < 315532800 {
		fatal("complete deterministic archive arguments are required")
	}
	entries, err := collect(*root)
	if err != nil {
		fatal(err.Error())
	}
	when := time.Unix(*epoch, 0).UTC()
	if *format == "tar.gz" {
		err = writeTarGzip(*output, entries, when)
	} else {
		err = writeZip(*output, entries, when)
	}
	if err != nil {
		_ = os.Remove(*output)
		fatal(err.Error())
	}
}

func collect(root string) ([]entry, error) {
	absoluteRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	var entries []entry
	err = filepath.WalkDir(absoluteRoot, func(path string, value os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == absoluteRoot || value.IsDir() {
			return nil
		}
		info, err := value.Info()
		if err != nil || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("non-regular archive input rejected: %s", path)
		}
		relative, err := filepath.Rel(absoluteRoot, path)
		if err != nil {
			return err
		}
		relative = filepath.ToSlash(relative)
		if relative == "" || strings.HasPrefix(relative, "/") || strings.Contains(relative, "../") || strings.Contains(relative, "//") {
			return fmt.Errorf("unsafe archive path rejected: %s", relative)
		}
		mode := os.FileMode(0o644)
		if strings.HasPrefix(relative, "bin/") {
			mode = 0o755
		}
		entries = append(entries, entry{absolute: path, relative: relative, mode: mode, size: info.Size()})
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("archive input is empty")
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].relative < entries[j].relative })
	return entries, nil
}

func writeTarGzip(output string, entries []entry, when time.Time) error {
	file, err := os.OpenFile(output, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	gzipWriter, err := gzip.NewWriterLevel(file, gzip.BestCompression)
	if err != nil {
		_ = file.Close()
		return err
	}
	gzipWriter.Header.ModTime = when
	gzipWriter.Header.OS = 255
	tarWriter := tar.NewWriter(gzipWriter)
	for _, value := range entries {
		header := &tar.Header{Name: value.relative, Mode: int64(value.mode.Perm()), Size: value.size, ModTime: when, AccessTime: time.Time{}, ChangeTime: time.Time{}, Uid: 0, Gid: 0, Uname: "", Gname: "", Format: tar.FormatPAX}
		if err := tarWriter.WriteHeader(header); err != nil {
			return closeTar(file, tarWriter, gzipWriter, err)
		}
		if err := copyFile(tarWriter, value.absolute); err != nil {
			return closeTar(file, tarWriter, gzipWriter, err)
		}
	}
	return closeTar(file, tarWriter, gzipWriter, nil)
}

func closeTar(file *os.File, tarWriter *tar.Writer, gzipWriter *gzip.Writer, prior error) error {
	if err := tarWriter.Close(); prior == nil {
		prior = err
	}
	if err := gzipWriter.Close(); prior == nil {
		prior = err
	}
	if err := file.Close(); prior == nil {
		prior = err
	}
	return prior
}

func writeZip(output string, entries []entry, when time.Time) error {
	file, err := os.OpenFile(output, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	writer := zip.NewWriter(file)
	for _, value := range entries {
		header := &zip.FileHeader{Name: value.relative, Method: zip.Deflate}
		header.SetMode(value.mode)
		header.SetModTime(when)
		destination, err := writer.CreateHeader(header)
		if err != nil {
			return closeZip(file, writer, err)
		}
		if err := copyFile(destination, value.absolute); err != nil {
			return closeZip(file, writer, err)
		}
	}
	return closeZip(file, writer, nil)
}

func closeZip(file *os.File, writer *zip.Writer, prior error) error {
	if err := writer.Close(); prior == nil {
		prior = err
	}
	if err := file.Close(); prior == nil {
		prior = err
	}
	return prior
}

func copyFile(destination io.Writer, path string) error {
	source, err := os.Open(path)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(destination, source)
	closeErr := source.Close()
	if copyErr != nil {
		return copyErr
	}
	return closeErr
}

func fatal(message string) {
	fmt.Fprintln(os.Stderr, "packager:", message)
	os.Exit(1)
}
