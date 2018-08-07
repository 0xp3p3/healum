package gce

import (
	micro_storage "server/cloudkey-srv/storage"
	mdb "server/cloudkey-srv/proto/record"
	"sync"
	"google.golang.org/api/cloudkms/v1"
	"golang.org/x/oauth2/google"
	"fmt"

	"encoding/base64"
	"golang.org/x/net/context"
)

var (
	// Default bucket
        BucketName = "go-project-test"

	// Default key for stubbing
	DefaultKey = "dGhpc2lzYXNlY3JldGtleXBsZWFzZWRvbm90c2hhcmU="

	// Project id to set by the client
	ProjectID = "server-vision"

	// Project id to set by the client
	DefaultCryptoKey = "org"
)

type gceDriver struct{}

type gceDB struct {
	sync.RWMutex
	bucket    string
}

func init() {
	// Other drives should be added the same way (aws, etc.)
	micro_storage.Drivers["gce"] = new(gceDriver)
}

func (d *gceDriver) NewStorage() (micro_storage.ST, error) {

	return &gceDB{
                bucket:      BucketName,
	}, nil
}

func (d *gceDB) Init(mdb *mdb.Storage) error {
	return nil
}

func (d *gceDB) Close() error {
	return nil
}

// Create a generic record from a storage
func (d *gceDB) CreateKey(orgid, cryptoKey string) (error) {
	d.RLock()
	defer d.RUnlock()

	cl, err := google.DefaultClient(context.Background(), cloudkms.CloudPlatformScope)
	if err != nil {
		return err
	}
	cloudkmsService, err := cloudkms.New(cl)
	parent := fmt.Sprintf("%v/locations/global/", ProjectID)


	_, err = cloudkmsService.Projects.Locations.KeyRings.Create(parent, &cloudkms.KeyRing{
		Name: orgid,
	}).Do()
	if err != nil {
		return err
	}
	if len(cryptoKey) == 0 {
		cryptoKey = DefaultCryptoKey
	}

	_, err = cloudkmsService.Projects.Locations.KeyRings.CryptoKeys.Create(parent, &cloudkms.CryptoKey{
		Name: cryptoKey,
		Purpose: "ENCRYPT_DECRYPT",

	}).Do()
	if err != nil {
		return err
	}

	return nil
}


// Encrypt a generic record from a storage
func (d *gceDB) EncryptKey(orgid, cryptoKey, dek  string) (string, error) {
	d.RLock()
	defer d.RUnlock()
	cl, err := google.DefaultClient(context.Background(), cloudkms.CloudPlatformScope)
	if err != nil {
		return "", err
	}
	cloudkmsService, err := cloudkms.New(cl)

	parentName := fmt.Sprintf("%s/locations/%s/%s/%s/",
		ProjectID, "global", orgid, cryptoKey)
	if len(cryptoKey) == 0 {
		cryptoKey = DefaultCryptoKey
	}

	resp, err := cloudkmsService.Projects.Locations.KeyRings.CryptoKeys.Encrypt(parentName, &cloudkms.EncryptRequest{
		Plaintext: base64.StdEncoding.EncodeToString([]byte(dek)),
	}).Do()
	if err != nil {
		return "", err
	}
	res, _ := base64.StdEncoding.DecodeString(resp.Ciphertext)
	return string(res), nil
}

// Decrypt a generic record from a storage
func (d *gceDB) DecryptKey(orgid, cryptoKey, encryptedDek  string) (string, error) {
	d.RLock()
	defer d.RUnlock()

	cl, err := google.DefaultClient(context.Background(), cloudkms.CloudPlatformScope)
	if err != nil {
		return "", err
	}
	cloudkmsService, err := cloudkms.New(cl)

	parentName := fmt.Sprintf("%s/locations/%s/%s/%s/",
		ProjectID, "global", orgid, cryptoKey)
	if len(cryptoKey) == 0 {
		cryptoKey = DefaultCryptoKey
	}

	resp, err := cloudkmsService.Projects.Locations.KeyRings.CryptoKeys.Decrypt(parentName, &cloudkms.DecryptRequest{
		Ciphertext: base64.StdEncoding.EncodeToString([]byte(encryptedDek)),
	}).Do()
	if err != nil {
		return "", err
	}

	res, _ := base64.StdEncoding.DecodeString(resp.Plaintext)
	return string(res), nil
}
