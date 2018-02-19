package COS

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type Client struct {
	client *http.Client
	//bucket url
	bucket string
}

func NewClient(secretID, secretKey, bucketURL string) (*Client, error) {
	t := transport{secretID: secretID, secretKey: secretKey}
	client := &http.Client{Transport: t}
	u, err := url.Parse(bucketURL)
	u.Scheme = "https"
	if err != nil {
		return nil, fmt.Errorf("url not valid: %v", err)
	}
	return &Client{client: client, bucket: u.String()}, nil
}

type transport struct {
	secretID  string
	secretKey string
}

func (t transport) RoundTrip(req *http.Request) (*http.Response, error) {
	sign(req, t.secretKey, t.secretID)
	//req.Write(os.Stdout)
	return http.DefaultTransport.RoundTrip(req)
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

//Put Object 接口请求可以将本地的文件(Object)上传至指定 Bucket 中。该操作需要请求者对 Bucket 有 WRITE 权限。
func (c *Client) PutObject(file string, dir string) error {
	info, err := os.Stat(file)
	if err != nil {
		return err
	}
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	u, err := url.Parse(c.bucket)
	if err != nil {
		return fmt.Errorf("API URL is not valid: %v", err)
	}
	u, err = u.Parse(dir)
	if err != nil {
		return fmt.Errorf("dir is not valid: %v", err)
	}
	u, err = u.Parse(info.Name())
	if err != nil {
		return fmt.Errorf("file name not valid: %v", err)
	}
	//fmt.Println(dir, u.String())
	req, err := http.NewRequest("PUT", u.String(), f)
	if err != nil {
		return fmt.Errorf("form request error: %v", err)
	}
	//set Content-Length headers, required
	size := info.Size()
	req.ContentLength = size
	//req.Header.Add("Content-Length", strconv.Itoa(int(size)))

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("request error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(b))
		return fmt.Errorf("Status Code: %s", resp.Status)
	}
	return nil
}
