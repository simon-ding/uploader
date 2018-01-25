package COS

import (
	"fmt"
	"io"
	"net/http"
)

type COS struct {
	client *http.Client
	bucket string
}

func NewCOS(secretID, secretKey, bucketURL string) *COS {
	t := transport{secretID: secretID, screctKey: secretKey}
	client := &http.Client{Transport: t}
	return &COS{client: client, bucket: bucketURL}
}

type transport struct {
	secretID  string
	screctKey string
}

func (t transport) RoundTrip(req *http.Request) (*http.Response, error) {
	sign(req, t.screctKey, t.secretID)
	return http.DefaultTransport.RoundTrip(req)
}

//put a local object into bucket
func (c *COS) PutObject(object io.Reader) error {
	req, err := http.NewRequest("PUT", "http://"+c.bucket, object)
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response status error: %s", resp.Status)
	}
	return nil
}
