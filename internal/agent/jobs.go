package agent

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

func registerCronJobs(config *Config) (*cron.Cron, error) {
	c := cron.New(cron.WithLogger(cron.DefaultLogger))
	for _, file := range config.Files {
		log.Default().Printf("register cron for file '%s' with interval '%s'\n", file.Path, file.Interval)
		_, err := c.AddFunc(file.Interval, func() {
			log.Default().Printf("running job for file '%s'", file.Path)
			err := sendFileToArchive1Server(config.ArchiveServer, config.AgentName, config.AgentKey, &file)
			if err != nil {
				log.Default().Printf("job fails. file: %s, error: [%s]", file.String(), err.Error())
			}
		})

		if err != nil {
			return nil, err
		}
	}
	c.Start()
	return c, nil
}

func sendFileToArchive1Server(server string, name string, key string, file *File) error {
	client := &http.Client{}
	correlationId := uuid.New().String()

	log.Default().Printf("correlation-id:'%s', target-file: '%s'\n", correlationId, file.Path)

	// read file
	f, err := os.Open(file.Path)
	if err != nil {
		return fmt.Errorf("correlation-id:'%s', error: %s", correlationId, err.Error())
	}
	defer f.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if file.Filename != "" {
		writer.WriteField("filename", file.Filename)
	}
	writer.WriteField("rotate", strconv.FormatInt(file.Rotate, 10))
	part, err := writer.CreateFormFile("file", filepath.Base(f.Name()))
	if err != nil {
		return fmt.Errorf("correlation-id:'%s', error: %s", correlationId, err.Error())
	}
	io.Copy(part, f)
	writer.Close()

	requestUrl := fmt.Sprintf("%s%s", server, "/api/v1/servers/store/file")

	req, err := http.NewRequest(http.MethodPost, requestUrl, body)
	if err != nil {
		return fmt.Errorf("correlation-id:'%s', error: %s", correlationId, err.Error())
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", key)
	req.Header.Set("X-Agent1-Name", name)
	req.Header.Set("X-Correlation-Id", correlationId)

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("correlation-id:'%s', error: %s", correlationId, err.Error())
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("correlation-id:'%s', error: %s", correlationId, err.Error())
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("correlation-id:'%s', error: %s",
			correlationId,
			fmt.Sprintf(
				"non 200 status code received. response: %s",
				resBody,
			),
		)
	}
	log.Default().Printf("correlation-id:'%s', response: %s\n", correlationId, resBody)

	return nil
}
