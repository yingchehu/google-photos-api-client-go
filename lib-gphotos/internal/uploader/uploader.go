package uploader

import (
	"log"
	"net/http"
	"os"
)

const (
	// Prefix to write at beginning of each log line.
	logPrefix = "gphotos-uploader: "

	// Define which text to prefix to each log entry generated by the Logger.
	// For example, flags Ldate | Ltime (or LstdFlags) produce,
	//	2009/01/23 01:23:23 message
	// while flags Ldate | Ltime | Lmicroseconds | Llongfile produce,
	//	2009/01/23 01:23:23.123123 /a/b/c/d.go:23: message
	//
	// See standard `log` package.
	logFlags = 0

	// API endpoint URL for upload media
	uploadEndpoint = "https://photoslibrary.googleapis.com/v1/uploads"
)

// Uploader is a client for uploading media to Google Photos.
// Original photos library does not provide `/v1/uploads` API.
type Uploader struct {
	// HTTP Client
	c *http.Client
	// URL of the endpoint to upload to
	url string
	// If Resume is true the UploadSessionStore is required.
	resume bool
	// store keeps upload session information.
	store UploadSessionStore

	log *log.Logger
}

// NewUploader returns an Uploader using the specified client or error in case
// of non valid configuration.
// The client must have the proper permissions to upload files.
//
// Use OptionResumableUploads(...), OptionLog(...) and OptionEndpoint(...) to
// customize configuration.
func NewUploader(client *http.Client, options ...Option) (*Uploader, error) {
	u := &Uploader{
		c:      client,
		url:    uploadEndpoint,
		resume: false,
		store:  nil,
		log:    log.New(os.Stderr, logPrefix, logFlags),
	}

	for _, opt := range options {
		opt(u)
	}

	if err := u.Validate(); err != nil {
		return nil, err
	}

	return u, nil

}

// OptionResumableUploads enables resumable uploads.
// Resumable uploads needs an UploadSessionStore to keep upload session information.
func OptionResumableUploads(store UploadSessionStore) Option {
	return func(c *Uploader) {
		c.resume = true
		c.store = store
	}
}

// OptionLog sets the logger to log messages.
func OptionLog(l *log.Logger) Option {
	return func(c *Uploader) {
		c.log = l
	}
}

// OptionEndpoint sets the URL of the endpoint to upload to.
func OptionEndpoint(url string) Option {
	return func(c *Uploader) {
		c.url = url
	}
}

// Validate validates the configuration of the Client.
func (u *Uploader) Validate() error {
	if u.resume && u.store == nil {
		return ErrNilStore
	}

	return nil
}

// Option defines an option for a Client
type Option func(*Uploader)