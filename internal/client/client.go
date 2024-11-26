package client

import (
	"OnlineMusic/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type APIClient struct {
	URL        string
	HttpClient *http.Client
}

type SuccessResponse struct {
	StatusCode int         `json:"status"`
	Data       interface{} `json:"data"`
}

func NewAPIClient(url string) *APIClient {
	return &APIClient{
		URL: url,
		HttpClient: &http.Client{
			Timeout: time.Second * 15,
		},
	}
}

func (api *APIClient) FetchInfoMusic(ctx context.Context, group, song string) (model.SongInfoDetail, error) {
	url := fmt.Sprintf("%s?group=%s&song=%s", api.URL, group, song)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return model.SongInfoDetail{}, err
	}

	req = req.WithContext(ctx)
	res := model.SongInfoDetail{}
	if err := api.sendRequest(req, &res); err != nil {
		return model.SongInfoDetail{}, err
	}
	return res, nil
}

func (api *APIClient) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")

	res, err := api.HttpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes error
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Error())
		}

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	successResponse := SuccessResponse{
		Data: v,
	}

	if err = json.NewDecoder(res.Body).Decode(&successResponse); err != nil {
		return err
	}

	return nil
}
