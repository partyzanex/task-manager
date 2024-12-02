package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	cmd := &cobra.Command{
		Use:   "bla-bla",
		Short: "123",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			host, err := cmd.Flags().GetString("host")
			if err != nil {
				log.Printf("failed to get host: %v", err)
				return
			}

			port, err := cmd.Flags().GetUint16("port")
			if err != nil {
				log.Printf("failed to get port: %v", err)
				return
			}

			err = createTask(cmd.Context(), host, port)
			if err != nil {
				log.Printf("failed to create task: %v", err)
				return
			}
		},
	}

	cmd.PersistentFlags().String("host", "localhost", "hostname of task server")
	cmd.PersistentFlags().Uint16("port", 8080, "port of task server")

	if err := cmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

type Task struct {
	ID          int64      `json:"id"`
	Type        string     `json:"type"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Status      uint8      `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type CreateTaskResponse struct {
	ID     int64   `json:"id"`
	Status uint8   `json:"status"`
	Error  *string `json:"error,omitempty"`
}

func createTask(ctx context.Context, host string, port uint16) error {
	uri, err := url.Parse(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}

	task := &Task{
		ID:          uuid.New(),
		Type:        "test-task",
		Name:        "Test Name",
		Description: "Test task description",
		Status:      0,
		CreatedAt:   time.Now(),
		UpdatedAt:   nil,
	}

	payload, err := json.Marshal(&task)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), bytes.NewReader(payload))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	jsonResp := CreateTaskResponse{}

	err = json.NewDecoder(resp.Body).Decode(&jsonResp)
	if err != nil {
		return err
	}

	return nil
}
