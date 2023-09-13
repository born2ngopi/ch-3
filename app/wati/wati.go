package wati

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type (
	WatiApp struct {
		Host        string
		AccessToken string
		Operator    string
	}

	watiReqPayload struct {
		TemplateName  string           `json:"template_name"`
		BroadcastName string           `json:"broadcast_name"`
		Parameters    []WatiParameters `json:"parameters"`
	}

	WatiParameters struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	WatiInteractiveListPayload struct {
		Header     string                    `json:"header"`
		Body       string                    `json:"body"`
		Footer     string                    `json:"footer"`
		ButtonText string                    `json:"buttonText"`
		Sections   []WatiInteractiveSections `json:"sections"`
	}

	WatiInteractiveRows struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	WatiInteractiveSections struct {
		Title string                `json:"title"`
		Rows  []WatiInteractiveRows `json:"rows"`
	}
)

func NewClient() *WatiApp {
	return &WatiApp{
		Host:        os.Getenv("WATI_HOST"),
		AccessToken: os.Getenv("WATI_ACCESS_TOKEN"),
		Operator:    os.Getenv("WATI_OPERATOR"),
	}
}

func (w *WatiApp) SendTemplate(ctx context.Context, waid, templateName, broadcastName string, parameters []WatiParameters) error {

	payload := watiReqPayload{
		TemplateName:  templateName,
		BroadcastName: broadcastName,
		Parameters:    parameters,
	}

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(payload)
	if err != nil {
		return err
	}

	return w.send(ctx, "/api/v1/sendTemplateMessage?whatsappNumber="+waid, buf)
}

func (w *WatiApp) SendInteractiveList(ctx context.Context, waid string, payload WatiInteractiveListPayload) error {

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(payload)
	if err != nil {
		return err
	}

	return w.send(ctx, "/api/v1/sendInteractiveListMessage?whatsappNumber="+waid, buf)
}

func (w *WatiApp) AssignOperator(ctx context.Context, waid string) error {
	buf := new(bytes.Buffer)
	return w.send(ctx, "/api/v1/assignOperator?email="+w.Operator+"&whatsappNumber="+waid, buf)
}

func (w *WatiApp) send(ctx context.Context, path string, buf *bytes.Buffer) error {

	req, err := http.NewRequest(http.MethodPost, w.Host+path, buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+w.AccessToken)
	req = req.WithContext(ctx)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
		return fmt.Errorf("error send message")
	}
	return nil
}
