package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	account_proto "server/account-srv/proto/account"
	"server/common"
	db_proto "server/db-srv/proto/db"
	user_proto "server/user-srv/proto/user"
	"strings"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/micro/go-micro/client"
)

type clientWrapper struct {
	Db_client db_proto.DBClient
}

var (
	ClientWrapper *clientWrapper
	ErrNotFound   = errors.New("not found")
)

// Storage for a db microservice client
func NewClientWrapper(serviceClient client.Client) *clientWrapper {
	cl := db_proto.NewDBClient("", serviceClient)

	return &clientWrapper{
		Db_client: cl,
	}
}

// Init initializes healum databases
func Init(serviceClient client.Client) error {
	ClientWrapper = NewClientWrapper(serviceClient)
	// if _, err := ClientWrapper.Db_client.Init(context.TODO(), &db_proto.InitRequest{}); err != nil {
	// 	log.Fatal(err)
	// 	return err
	// }
	return nil
}

// RemoveDb removes healum database (for testing)
func RemoveDb(ctx context.Context, serviceClient client.Client) error {
	ClientWrapper = NewClientWrapper(serviceClient)
	if _, err := ClientWrapper.Db_client.RemoveDb(ctx, &db_proto.RemoveDbRequest{}); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func runQuery(ctx context.Context, q string, table string) (*db_proto.RunQueryResponse, error) {
	return ClientWrapper.Db_client.RunQuery(ctx, &db_proto.RunQueryRequest{
		Database: &db_proto.Database{
			Name:     common.DbHealumName,
			Table:    table,
			Driver:   common.DbHealumDriver,
			Metadata: common.SearchableMetaMap,
		},
		Query: q,
	})
}

func accountToRecord(account *account_proto.Account, salt string, options ...interface{}) (string, error) {
	data, err := common.MarhalToObject(account, options)
	if err != nil {
		return "", err
	}
	// FIXME: when creating a account it should either have email/password OR phone/passcode not both(related to EmitDefault)
	if len(account.Email) > 0 {
		delete(data, "phone")
		delete(data, "passcode")
	} else if len(account.Phone) > 0 {
		delete(data, "email")
		delete(data, "password")
	}

	d := map[string]interface{}{
		"_key":       account.Id,
		"id":         account.Id,
		"created":    account.Created,
		"updated":    account.Updated,
		"name":       account.Email,
		"parameter1": account.Phone,
		"parameter2": salt,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToAccount(r *db_proto.Record) (*account_proto.Account, error) {
	var p account_proto.Account
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToAccountStatusInfo(r *db_proto.Record) (*account_proto.AccountStatusInfo, error) {
	var p account_proto.AccountStatusInfo
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func DedupAccount(ctx context.Context, account *account_proto.Account) (*account_proto.Account, error) {
	query := `FILTER`
	if len(account.Email) == 0 && len(account.Phone) == 0 {
		return nil, errors.New("Account info is invalid")
	}
	if len(account.Email) > 0 {
		query += fmt.Sprintf(` && doc.data.email == "%v"`, account.Email)
	}
	if len(account.Phone) > 0 {
		query += fmt.Sprintf(` && doc.data.phone == "%v"`, account.Phone)
	}
	query = strings.Replace(query, `FILTER && `, `FILTER `, -1)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		RETURN doc`, common.DbAccountTable, query)
	resp, err := runQuery(ctx, q, common.DbAccountTable)
	if err != nil {
		common.ErrorLog(common.AccountSrv, DedupAccount, nil, "Query running is failed")
		return nil, err
	}
	if err != nil || len(resp.Records) == 0 {
		return nil, ErrNotFound
	}

	data, err := recordToAccount(resp.Records[0])
	return data, err
}

func Create(ctx context.Context, account *account_proto.Account) error {
	if len(account.Email) > 0 {
		account.Email = strings.ToLower(account.Email)
	}
	if account.Created == 0 {
		account.Created = time.Now().Unix()
	}
	if account.Status == 0 {
		account.Status = account_proto.AccountStatus_INACTIVE
	}
	account.Updated = time.Now().Unix()

	record, err := accountToRecord(account, "")
	if err != nil {
		common.ErrorLog(common.AccountSrv, Create, nil, "Account marshaling is invalid")
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { _key: "%v" } 
		INSERT %v 
		UPDATE %v 
		IN %v`, account.Id, record, record, common.DbAccountTable)
	_, err = runQuery(ctx, q, common.DbAccountTable)
	return err
}

func Update(ctx context.Context, account *account_proto.Account) (*account_proto.Account, error) {
	if len(account.Email) > 0 {
		account.Email = strings.ToLower(account.Email)
	}
	account.Updated = time.Now().Unix()
	record, err := accountToRecord(account, "", true)
	if err != nil {
		common.ErrorLog(common.AccountSrv, Create, nil, "Account marshaling is invalid")
		return nil, err
	}
	if len(record) == 0 {
		return nil, errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { _key: "%v" } 
		INSERT %v 
		UPDATE %v 
		IN %v
		RETURN NEW`, account.Id, record, record, common.DbAccountTable)
	resp, err := runQuery(ctx, q, common.DbAccountTable)
	if err != nil || len(resp.Records) == 0 {
		common.ErrorLog(common.AccountSrv, Read, err, "Read query is failed")
		return nil, ErrNotFound
	}
	data, err := recordToAccount(resp.Records[0])
	return data, err
}

func ConfirmAccountAndUpdatePass(ctx context.Context, account *account_proto.Account, salt string) error {
	if len(account.Email) > 0 {
		account.Email = strings.ToLower(account.Email)
	}
	account.Updated = time.Now().Unix()
	if account.Status == 0 {
		account.Status = account_proto.AccountStatus_INACTIVE
	}

	record, err := accountToRecord(account, salt)
	if err != nil {
		common.ErrorLog(common.AccountSrv, ConfirmAccountAndUpdatePass, nil, "Account marshaling is invalid")
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { _key: "%v" } 
		INSERT %v 
		UPDATE %v 
		IN %v`, account.Id, record, record, common.DbAccountTable)
	_, err = runQuery(ctx, q, common.DbAccountTable)
	return err
}

func Read(ctx context.Context, account *account_proto.Account) (*account_proto.Account, error) {
	query := `FILTER`
	if len(account.Id) > 0 {
		query += fmt.Sprintf(` && doc._key == "%v"`, account.Id)
	}
	if len(account.Email) > 0 {
		query += fmt.Sprintf(` && doc.data.email == "%v"`, account.Email)
	} else if len(account.Phone) > 0 {
		query += fmt.Sprintf(` && doc.data.phone == "%v"`, account.Phone)
	}

	query = common.QueryClean(query)
	if query == `FILTER` {
		common.ErrorLog(common.AccountSrv, Read, nil, "Account information is invalid")
		return nil, errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		RETURN doc`, common.DbAccountTable, query)
	resp, err := runQuery(ctx, q, common.DbAccountTable)
	if err != nil || len(resp.Records) == 0 {
		common.ErrorLog(common.AccountSrv, Read, err, "Read query is failed")
		return nil, ErrNotFound
	}

	data, err := recordToAccount(resp.Records[0])
	return data, err
}

func UpdatePassword(ctx context.Context, id, password, salt string) error {
	q := fmt.Sprintf(`
		FOR doc IN %v
		FILTER doc._key == "%v"
		UPDATE doc WITH {updated:%v, parameter2:"%v", data:{password:"%v",updated:%v}} IN %v`, common.DbAccountTable, id, time.Now().Unix(), salt, password, time.Now().Unix(), common.DbAccountTable)
	_, err := runQuery(ctx, q, common.DbAccountTable)
	return err
}

func UpdatePasscode(ctx context.Context, id, passcode, salt string) error {
	q := fmt.Sprintf(`
		FOR doc IN %v
		FILTER doc._key == "%v"
		UPDATE doc WITH {updated:%v, parameter2:"%v", data:{passcode:"%v",updated:%v}} IN %v`, common.DbAccountTable, id, time.Now().Unix(), salt, passcode, time.Now().Unix(), common.DbAccountTable)
	_, err := runQuery(ctx, q, common.DbAccountTable)
	return err
}

func UpdateStatus(ctx context.Context, id string, status account_proto.AccountStatus, confirmed bool) error {
	q := fmt.Sprintf(`
		FOR doc IN %v
		FILTER doc._key == "%v"	
		UPDATE doc WITH {updated:%v, data:{status:"%v",confirmed:%v,updated:%v}} IN %v`, common.DbAccountTable, id, time.Now().Unix(), status, confirmed, time.Now().Unix(), common.DbAccountTable)
	_, err := runQuery(ctx, q, common.DbAccountTable)
	return err
}

func UpdateLockStatus(ctx context.Context, id string) error {
	q := fmt.Sprintf(`
		FOR doc IN %v
		FILTER doc._key == "%v"	
		UPDATE doc WITH {updated:%v, data:{status:"%v",updated:%v}} IN %v`, common.DbAccountTable, id, time.Now().Unix(), account_proto.AccountStatus_LOCKED, time.Now().Unix(), common.DbAccountTable)
	_, err := runQuery(ctx, q, common.DbAccountTable)
	return err
}

func SaltAndPassword(ctx context.Context, id string, mode user_proto.ContactDetailType) (string, string, error) {
	q := fmt.Sprintf(`
		FOR doc IN %v
		FILTER doc._key == "%v"
		RETURN doc`, common.DbAccountTable, id)
	resp, err := runQuery(ctx, q, common.DbAccountTable)
	if err != nil || len(resp.Records) == 0 {
		common.ErrorLog(common.AccountSrv, SaltAndPassword, err, "Run query is failed")
		return "", "", ErrNotFound
	}

	record := resp.Records[0]
	salt := record.Parameter2
	data, err := recordToAccount(resp.Records[0])
	if err != nil {
		common.ErrorLog(common.AccountSrv, SaltAndPassword, err, "Unmarshaling is failed")
		return "", "", err
	}

	var pass string
	if mode == user_proto.ContactDetailType_EMAIL {
		pass = data.Password
	} else if mode == user_proto.ContactDetailType_PHONE {
		pass = data.Passcode
	}
	return salt, pass, nil
}

func SetAccountStatus(ctx context.Context, userId string, status account_proto.AccountStatus) error {
	q := fmt.Sprintf(`
		FOR a IN OUTBOUND "%v/%v" %v
		UPDATE a WITH {updated:%v, data:{status:"%v",updated:%v}} IN %v`,
		common.DbUserTable, userId, common.DbUserAccountEdgeTable,
		time.Now().Unix(), status, time.Now().Unix(), common.DbAccountTable)

	if _, err := runQuery(ctx, q, common.DbAccountTable); err != nil {
		common.ErrorLog(common.AccountSrv, SetAccountStatus, err, "SetAccountStatus query is failed")
		return err
	}
	return nil
}

func GetAccountStatus(ctx context.Context, userId string) (*account_proto.AccountStatusInfo, error) {
	q := fmt.Sprintf(`
		FOR acc IN OUTBOUND "%v/%v" %v
		return {data:{status:acc.data.status}}`,
		common.DbUserTable, userId, common.DbUserAccountEdgeTable)

	resp, err := runQuery(ctx, q, common.DbAccountTable)

	if err != nil || len(resp.Records) == 0 {
		return nil, ErrNotFound
	}
	data, err := recordToAccountStatusInfo(resp.Records[0])
	return data, err
}

func ResetUserPassword(ctx context.Context, accountId, salt, pass string, mode user_proto.ContactDetailType) error {
	var passquery string
	if mode == user_proto.ContactDetailType_EMAIL {
		passquery = fmt.Sprintf(`password:"%v"`, pass)
	} else if mode == user_proto.ContactDetailType_PHONE {
		passquery = fmt.Sprintf(`passcode:"%v"`, pass)
	}

	q := fmt.Sprintf(`
		FOR a IN %v
		FILTER a._key == "%v"
		UPDATE a WITH {updated:%v, parameter2:"%v", data:{%v,updated:%v}} IN %v`,
		common.DbAccountTable, accountId,
		time.Now().Unix(), salt, passquery, time.Now().Unix(), common.DbAccountTable)
	if _, err := runQuery(ctx, q, common.DbAccountTable); err != nil {
		common.ErrorLog(common.AccountSrv, ResetUserPassword, err, "ResetUserPassword query is failed")
		return err
	}
	return nil
}

// GetAccountByUser returns account from user_id in user_account_edge collection
func GetAccountByUser(ctx context.Context, userId string) (*account_proto.Account, error) {
	q := fmt.Sprintf(`
		FOR doc IN OUTBOUND "%v/%v" %v
		RETURN doc`,
		common.DbUserTable, userId, common.DbUserAccountEdgeTable)

	resp, err := runQuery(ctx, q, common.DbUserAccountEdgeTable)
	if err != nil || len(resp.Records) == 0 {
		common.ErrorLog(common.AccountSrv, GetAccountByUser, err, "GetAccountByUser query is failed")
		return nil, ErrNotFound
	}

	return recordToAccount(resp.Records[0])
}

// IsReadEmployee checks valid employee
func IsReadEmployee(ctx context.Context, userid, orgid string) bool {
	q := fmt.Sprintf(`
		FOR doc IN %v
		FILTER doc._from == "%v/%v" && doc._to == "%v/%v"
		RETURN doc`, common.DbEmployeeTable, common.DbUserTable, userid, common.DbOrganisationTable, orgid)
	resp, err := runQuery(ctx, q, common.DbUserAccountEdgeTable)
	if err != nil || len(resp.Records) == 0 {
		return false
	}
	return true
}
