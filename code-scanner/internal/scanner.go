package internal

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"

	pb "github.com/majorperfect/guardrails-test/code-analyzer/proto/code-analyzer"
	"github.com/majorperfect/guardrails-test/code-scanner/connection"
)

type ScannerFunc = func(ctx context.Context, path string) error

func Scanner(scanFiles scanFilesFunc) ScannerFunc {
	return func(ctx context.Context, path string) error {
		var reCurr func(p string) error
		reCurr = func(p string) error {

			dirs, err := scanFiles(ctx, p)
			if err != nil {
				return err
			}
			for _, d := range dirs {
				reCurr(d)
			}
			return nil
		}

		return reCurr(path)

	}

}

type scanFilesFunc = func(ctx context.Context, path string) ([]string, error)

func ScanFiles(goScanFile GoScanFilesFunc) scanFilesFunc {
	return func(ctx context.Context, path string) ([]string, error) {

		screenFiles, err := ioutil.ReadDir(path)
		if err != nil {
			log.Fatal(err)
		}

		files := []string{}
		dirs := []string{}
		for _, f := range screenFiles {
			fmt.Println(f.Name(), f.IsDir())
			if f.IsDir() {
				dirs = append(dirs[:], path+"/"+f.Name())
			} else { // a file
				files = append(files[:], path+"/"+f.Name())
			}
		}
		fmt.Println(runtime.NumCPU())

		var wg sync.WaitGroup
		batch := 20 // no. of records to process in one go-routine

		N := 20 // this is the channel size, can be used to limit the number of the goroutines

		sem := make(chan struct{}, N)

		for i := 0; i < len(files); i += batch {
			j := i + batch
			if j > len(files) {
				j = len(files)
			}
			sem <- struct{}{}
			wg.Add(1)

			go goScanFile(ctx, files[i:j], sem, &wg)

		}
		wg.Wait()
		close(sem)

		return dirs, nil
	}

}

type GoScanFilesFunc = func(ctx context.Context, files []string, ch chan struct{}, wg *sync.WaitGroup)

func GoScanFiles(codeAnlyzer connection.AnalyzerMicroserviceFunc) GoScanFilesFunc {
	return func(ctx context.Context, files []string, ch chan struct{}, wg *sync.WaitGroup) {
		defer wg.Done()
		ch <- struct{}{}
		defer func() {
			// Reading from the channel decrements the semaphore
			// (frees up buffer slot).
			<-ch
		}()
		for _, file := range files {

			conn, client, err := codeAnlyzer(ctx)
			if err != nil {
				log.Print("connection error codeAnlyzer : ", err)
				return
			}
			if err := StreamFile(ctx, file, client); err != nil {
				log.Print("StreamFile error : ", err)
			}
			// if err := SendFile(ctx, file, client); err != nil {
			// 	log.Print("StreamFile error : ", err)
			// }
			defer conn.Close()
		}

	}

}

func StreamFile(ctx context.Context, file string, client pb.CodeAnalyzerServiceClient) error {
	openFile, err := os.Open(file)
	if err != nil {
		return err
	}
	fmt.Println("StreamFile:", file)
	fileName := strings.Trim(file, "./")

	//TODO: handle error
	// Maximum 1KB size per stream.
	buf := make([]byte, 1024)
	stream, err := client.AnalyzeUploader(ctx)
	if err != nil {
		return err
	}
	for {
		num, err := openFile.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("134 StreamFile:", err)
			return stream.CloseSend()

		}

		if err := stream.Send(&pb.UploadRequest{Chunk: buf[:num], Name: fileName}); err != nil {
			fmt.Println("134 StreamFile:", err)

			err = stream.CloseSend()
			if err != nil {
				return err
			}
			break

		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		return err
	}
	return nil
}

func SendFile(ctx context.Context, file string, client pb.CodeAnalyzerServiceClient) error {
	openFile, err := os.Open(file)
	if err != nil {
		return err
	}
	fmt.Println("StreamFile:", file)
	reader := bufio.NewReader(openFile)
	content, _ := ioutil.ReadAll(reader)

	// Encode as base64.
	encoded := base64.StdEncoding.EncodeToString(content)

	req := &pb.FileUploadRequest{
		Base64: encoded,
		Name:   file,
	}
	_, err = client.FileUpload(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
