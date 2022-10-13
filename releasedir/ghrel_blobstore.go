package releasedir

import (
	gobytes "bytes"
	"context"
	"encoding/json"
	"os"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"

	ghrelclient "github.com/fejnartal/bosh-ghrelcli/client"
	ghrelconfig "github.com/fejnartal/bosh-ghrelcli/config"
)

type GHRelBlobstore struct {
	fs      boshsys.FileSystem
	uuidGen boshuuid.Generator
	options map[string]interface{}
}

func NewGHRelBlobstore(
	fs boshsys.FileSystem,
	uuidGen boshuuid.Generator,
	options map[string]interface{},
) GHRelBlobstore {
	return GHRelBlobstore{
		fs:      fs,
		uuidGen: uuidGen,
		options: options,
	}
}

func (b GHRelBlobstore) Get(blobID string) (string, error) {
	client, err := b.client()
	if err != nil {
		return "", err
	}

	file, err := b.fs.TempFile("bosh-ghrel-blob")
	if err != nil {
		return "", bosherr.WrapError(err, "Creating destination file")
	}

	defer file.Close()

	err = client.Get(blobID, file)
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}

func (b GHRelBlobstore) Create(path string) (string, error) {
	client, err := b.client()
	if err != nil {
		return "", err
	}

	blobID, err := b.uuidGen.Generate()
	if err != nil {
		return "", bosherr.WrapError(err, "Generating blobstore ID")
	}

	file, err := b.fs.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return "", bosherr.WrapError(err, "Opening source file")
	}

	defer file.Close()

	err = client.Put(file, blobID)
	if err != nil {
		return "", bosherr.WrapError(err, "Generating blobstore ID")
	}

	return blobID, nil
}

func (b GHRelBlobstore) CleanUp(path string) error {
	panic("CLEANUP Not implemented")
}

func (b GHRelBlobstore) Delete(blobID string) error {
	panic("DELETE Not implemented")
}

func (b GHRelBlobstore) Validate() error {
	_, err := b.client()
	return err
}

func (b GHRelBlobstore) client() (*ghrelclient.GHRelBlobstore, error) {
	bytes, err := json.Marshal(b.options)
	if err != nil {
		return nil, bosherr.WrapError(err, "Marshaling config")
	}

	conf, err := ghrelconfig.NewFromReader(gobytes.NewBuffer(bytes))
	if err != nil {
		return nil, bosherr.WrapError(err, "Reading config")
	}

	client, err := ghrelclient.New(context.Background(), &conf)
	if err != nil {
		return nil, bosherr.WrapError(err, "Validating config")
	}

	return client, nil
}
