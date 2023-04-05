package provider

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type Environment struct {
	Host string
}

type Client struct {
	M2MToken    string
	Environment Environment
	HTTPClient  *http.Client
}

func (c *Client) GetOASDocument(title, version string) (interface{}, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/documents/%s?version=%s", c.Environment.Host, title, version),
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.M2MToken))
	req.Header.Set("Accept", "application/json")

	// Send the HTTP request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(string(respBody))
	}

	return respBody, nil
}

func (c *Client) UploadOASDocument(oasFile *os.File, contentType string) ([]byte, error) {
	// Get stats about the OAS file
	oasFileInfo, err := oasFile.Stat()
	if err != nil {
		return nil, err
	}

	// Create a new multipart form
	multiBuf := &bytes.Buffer{}
	multiWriter := multipart.NewWriter(multiBuf)

	// Add the file to the multipart form
	part, err := multiWriter.CreatePart(map[string][]string{
		"Content-Disposition": {"form-data; name=\"document\"; filename=\"" + oasFileInfo.Name() + "\""},
		"Content-Type":        {contentType},
	})
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, oasFile)
	if err != nil {
		return nil, err
	}

	// Close the multipart writer
	err = multiWriter.Close()
	if err != nil {
		return nil, err
	}

	// Create a new HTTP request with the multipart form as the body
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/documents", c.Environment.Host), multiBuf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.M2MToken))
	req.Header.Set("Content-Type", multiWriter.FormDataContentType())
	req.Header.Set("Accept", "application/json")

	// Send the HTTP request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf(string(respBody))
	}

	return respBody, nil
}

func (c *Client) DeleteOASDocument(title, version string) error {
	// Create a new HTTP request
	req, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("%s/documents/%s?version=%s", c.Environment.Host, title, version),
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.M2MToken))
	req.Header.Set("Accept", "application/json")

	// Send the HTTP request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		return err
	}

	return nil
}
