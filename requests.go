package clarifai

import (
	"encoding/json"
	"errors"
	"math/big"
)

// InfoResp represents the expected JSON response from /info/
type InfoResp struct {
	StatusCode    string `json:"status_code"`
	StatusMessage string `json:"status_msg"`
	Results       struct {
		MaxImageSize      int     `json:"max_image_size"`
		DefaultLanguage   string  `json:"default_language"`
		MaxVideoSize      int     `json:"max_video_size"`
		MaxImageBytes     int     `json:"max_image_bytes"`
		DefaultModel      string  `json:"default_model"`
		MaxVideoBytes     int     `json:"max_video_bytes"`
		MaxVideoDuration  int     `json:"max_video_duration"`
		MaxVideoBatchSize int     `json:"max_video_batch_size"`
		MinVideoSize      int     `json:"min_video_size"`
		MinImageSize      int     `json:"min_image_size"`
		MaxBatchSize      int     `json:"max_batch_size"`
		APIVersion        float32 `json:"api_version"`
	}
}

// TagRequest represents a JSON request for /tag/
type TagRequest struct {
	URLs     []string `json:"url"`
	LocalIDs []string `json:"local_ids,omitempty"`
	Model    string   `json:"model,omitempty"`
}

// TagResp represents the expected JSON response from /tag/
type TagResp struct {
	StatusCode    string `json:"status_code" bson:"status_code"`
	StatusMessage string `json:"status_msg" bson:"status_msg"`
	Meta          struct {
		Tag struct {
			Timestamp json.Number `json:"timestamp" bson:"timestamp"`
			Model     string      `json:"model" bson:"model"`
			Config    string      `json:"config" bson:"config"`
		} `json:"tag" bson:"tag"`
	} `json:"meta" bson:"meta"`
	Results []TagResult `json:"results" bson:"results"`
}

// TagResult represents the expected data for a single tag result
type TagResult struct {
	DocID         *big.Int `json:"docid" bson:"docid"`
	URL           string   `json:"url" bson:"url"`
	StatusCode    string   `json:"status_code" bson:"status_code"`
	StatusMessage string   `json:"status_msg" bson:"status_msg"`
	LocalID       string   `json:"local_id" bson:"local_id"`
	Result        struct {
		Tag struct {
			Classes []string  `json:"classes" bson:"classes"`
			CatIDs  []string  `json:"catids" bson:"catids"`
			Probs   []float32 `json:"probs" bson:"probs"`
		} `json:"tag" bson:"tag"`
	} `json:"result" bson:"result"`
	DocIDString string `json:"docid_str"`
}

// ColorRequest represents the JSON request to /color/
type ColorRequest struct {
	URLs     []string `json:"url"`
	LocalIDs []string `json:"local_ids,omitempty"`
}

// ColorResp is the expected response from the /color/ endpoint
type ColorResp struct {
	StatusCode    string `json:"status_code" bson:"status_code"`
	StatusMessage string `json:"status_msg" bson:"status_msg"`
	Results       []struct {
		DocID       *big.Int `json:"docid" bson:"docid"`
		URL         string   `json:"url" bson:"url"`
		DocIDString string   `json:"docid_str" bson:"docid_str"`
		Colors      []Color  `json:"colors" bson:"colors"`
	} `json:"results" bson:"results"`
}

// Color represents a single color in a given image
type Color struct {
	W3C struct {
		Hex  string `json:"hex" bson:"hex"`
		Name string `json:"name" bson:"name"`
	} `json:"w3c" bson:"w3c"`
	Hex     string  `json:"hex" bson:"hex"`
	Density float64 `json:"density" bson:"density"`
}

// FeedbackForm is used to send feedback back to Clarifai
type FeedbackForm struct {
	DocIDs           []string `json:"docids,omitempty"`
	URLs             []string `json:"url,omitempty"`
	AddTags          []string `json:"add_tags,omitempty"`
	RemoveTags       []string `json:"remove_tags,omitempty"`
	DissimilarDocIDs []string `json:"dissimilar_docids,omitempty"`
	SimilarDocIDs    []string `json:"similar_docids,omitempty"`
	SearchClick      []string `json:"search_click,omitempty"`
}

// FeedbackResp is the expected response from /feedback/
type FeedbackResp struct {
	StatusCode    string `json:"status_code"`
	StatusMessage string `json:"status_msg"`
}

// Info will return the current status info for the given client
func (client *Client) Info() (*InfoResp, error) {
	res, err := client.commonHTTPRequest(nil, "info", "GET", false)

	if err != nil {
		return nil, err
	}

	info := new(InfoResp)
	err = json.Unmarshal(res, info)

	return info, err
}

// Tag allows the client to request tag data on a single, or multiple photos
func (client *Client) Tag(req TagRequest) (*TagResp, error) {
	if len(req.URLs) < 1 {
		return nil, errors.New("Requires at least one url")
	}

	res, err := client.commonHTTPRequest(req, "tag", "POST", false)

	if err != nil {
		return nil, err
	}

	tagres := new(TagResp)
	err = json.Unmarshal(res, tagres)

	return tagres, err
}

// Color makes a request for a series of images to be color tagged
func (client *Client) Color(req ColorRequest) (*ColorResp, error) {
	if len(req.URLs) < 1 {
		return nil, errors.New("Requires at least one url")
	}

	res, err := client.commonHTTPRequest(req, "color", "POST", false)

	if err != nil {
		return nil, err
	}

	colorResponse := new(ColorResp)
	err = json.Unmarshal(res, colorResponse)

	return colorResponse, err
}

// Feedback allows the user to provide contextual feedback to Clarifai in order to improve their results
func (client *Client) Feedback(form FeedbackForm) (*FeedbackResp, error) {
	if form.DocIDs == nil && form.URLs == nil {
		return nil, errors.New("Requires at least one docid or url")
	}

	if form.DocIDs != nil && form.URLs != nil {
		return nil, errors.New("Request must provide exactly one of the following fields: {'DocIDs', 'URLs'}")
	}

	res, err := client.commonHTTPRequest(form, "feedback", "POST", false)

	feedbackres := new(FeedbackResp)
	err = json.Unmarshal(res, feedbackres)

	return feedbackres, err

}
