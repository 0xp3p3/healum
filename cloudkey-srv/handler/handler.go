package handler

import (
	"golang.org/x/net/context"

	cloudkey_proto "server/cloudkey-srv/proto/record"
	"server/cloudkey-srv/storage"
)

type CloudKeyService struct{}

func (s *CloudKeyService) CreateKey(ctx context.Context, req *cloudkey_proto.CreateKeyRequest, rsp *cloudkey_proto.CreateKeyResponse) error {
	return storage.DefaultStorage.CreateKey(nil, req.Orgid, req.CryptoKey)
}

func (s *CloudKeyService) EncryptKey(ctx context.Context, req *cloudkey_proto.EncryptKeyRequest, rsp *cloudkey_proto.EncryptKeyResponse) error {
	res, err := storage.DefaultStorage.EncryptKey(nil, req.Orgid, req.CryptoKey, req.Dek)
	if err != nil {
		return err
	}
	rsp.EncryptedDek = res

	return nil
}

func (s *CloudKeyService) DecryptKey(ctx context.Context, req *cloudkey_proto.DecryptKeyRequest, rsp *cloudkey_proto.DecryptKeyResponse) error {
	res, err := storage.DefaultStorage.DecryptKey(nil, req.Orgid, req.CryptoKey, req.Dek)
	if err != nil {
		return err
	}
	rsp.EncryptedDek = res

	return nil
}
