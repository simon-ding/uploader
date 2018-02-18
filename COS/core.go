package COS

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Client struct {
	client *http.Client
	bucket string
}

func NewClient(secretID, secretKey, bucketURL string) *Client {
	t := transport{secretID: secretID, screctKey: secretKey}
	client := &http.Client{Transport: t}
	u, err := url.Parse(bucketURL)
	u.Scheme = "https"
	if err != nil {
		log.Fatalf("url not valid: %v", err)
	}
	return &Client{client: client, bucket: u.String()}
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
func (c *Client) PutObject(object io.Reader) error {
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

func (c *Client) GetBucket() error {
	resp, err := c.client.Get(c.bucket)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Print("%s", string(b))
	return nil
}
