package schematics

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	apiv1 "github.com/johandry/gics/schematics/api/v1"
)

const (
	uploadTarWorkspaceTimeout = 50
)

// UploadTar upload a compressed (Tar) file/content into the workspace
func (w *Workspace) UploadTar(body io.Reader) error {
	// Delete Timeout
	ctx, cancelFunc := context.WithTimeout(context.Background(), uploadTarWorkspaceTimeout*time.Second)
	defer cancelFunc()

	params := &apiv1.UploadTemplateTarParams{}
	resp, err := w.service.clientWithResponses.UploadTemplateTarWithBodyWithResponse(ctx, w.ID, templateIDDefault, params, "multipart/form-data", body)
	if err != nil {
		return err
	}
	if code := resp.StatusCode(); code != 200 {
		return getAPIError("failed to upload the compressed code", resp.Body)
	}
	response := resp.JSON200 // WorkspaceDeleteResponse

	fmt.Printf("[DEBUG] code uploaded. Response: %+v\n", response)

	if !*response.HasReceivedFile {
		return fmt.Errorf("failed to upload the code")
	}

	// TODO: Complete the UploadTar func

	return nil
}

func (w *Workspace) tarMemFiles() (io.Reader, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	for name, body := range w.tfCodeFiles {
		hdr := &tar.Header{
			Name: name,
			Mode: 0600,
			Size: int64(len(body)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return nil, err
		}
		if _, err := tw.Write([]byte(body)); err != nil {
			return nil, err
		}
	}
	if err := tw.Close(); err != nil {
		return nil, err
	}

	return &buf, nil
}

func (w *Workspace) tarCode() (io.Reader, error) {
	return nil, nil
}

func (w *Workspace) tarDir(dir string) (io.Reader, error) {
	return nil, nil
}
