package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"server/common"
	content_proto "server/content-srv/proto/content"
	db_proto "server/db-srv/proto/db"
	common_proto "server/static-srv/proto/common"
	static_proto "server/static-srv/proto/static"
	user_proto "server/user-srv/proto/user"
	"strings"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/micro/go-micro/client"
	"github.com/pborman/uuid"
	log "github.com/sirupsen/logrus"
)

type clientWrapper struct {
	Db_client            db_proto.DBClient
	ContentServiceClient content_proto.ContentServiceClient
}

var (
	ClientWrapper *clientWrapper
	ErrNotFound   = errors.New("not found")
)

// Storage for a db microservice client
func NewClientWrapper(serviceClient client.Client) *clientWrapper {
	cl := db_proto.NewDBClient("", serviceClient)
	cl1 := content_proto.NewContentServiceClient("", serviceClient)

	return &clientWrapper{
		Db_client:            cl,
		ContentServiceClient: cl1,
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

func sourceToRecord(source *content_proto.Source) (string, error) {
	data, err := common.MarhalToObject(source)
	if err != nil {
		return "", err
	}
	// remove unneccearies
	common.FilterObject(data, "type", source.Type)

	d := map[string]interface{}{
		"_key":       source.Id,
		"id":         source.Id,
		"created":    source.Created,
		"updated":    source.Updated,
		"name":       source.Name,
		"parameter1": source.OrgId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToSource(r *db_proto.Record) (*content_proto.Source, error) {
	var p content_proto.Source
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func taxonomyToRecord(taxonomy *static_proto.Taxonomy) (string, error) {
	data, err := common.MarhalToObject(taxonomy)
	if err != nil {
		return "", err
	}
	d := map[string]interface{}{
		"_key":       taxonomy.Id,
		"id":         taxonomy.Id,
		"created":    taxonomy.Created,
		"updated":    taxonomy.Updated,
		"name":       taxonomy.Name,
		"parameter1": taxonomy.OrgId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToTaxonomy(r *db_proto.Record) (*static_proto.Taxonomy, error) {
	var p static_proto.Taxonomy
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func contentCategoryItemToRecord(contentCategoryItem *static_proto.ContentCategoryItem) (string, error) {
	data, err := common.MarhalToObject(contentCategoryItem)
	if err != nil {
		return "", err
	}
	common.FilterObject(data, "taxonomy", contentCategoryItem.Taxonomy)
	common.FilterObject(data, "category", contentCategoryItem.Category)

	d := map[string]interface{}{
		"_key":       contentCategoryItem.Id,
		"id":         contentCategoryItem.Id,
		"created":    contentCategoryItem.Created,
		"updated":    contentCategoryItem.Updated,
		"name":       contentCategoryItem.Name,
		"parameter1": contentCategoryItem.OrgId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToContentCategoryItem(r *db_proto.Record) (*static_proto.ContentCategoryItem, error) {
	var p static_proto.ContentCategoryItem
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func contentToRecord(content *content_proto.Content) (string, error) {
	data, err := common.MarhalToObject(content)
	if err != nil {
		return "", err
	}
	var createdById string
	if content.CreatedBy != nil {
		createdById = content.CreatedBy.Id
	}
	common.FilterObject(data, "createdBy", content.CreatedBy)
	common.FilterObject(data, "source", content.Source)
	common.FilterObject(data, "type", content.Type)
	common.FilterObject(data, "category", content.Category)
	if len(content.Tags) > 0 {
		var arr []interface{}
		for _, item := range content.Tags {
			arr = append(arr, map[string]string{
				"id": item.Id,
			})
		}
		data["tags"] = arr
	} else {
		delete(data, "tags")
	}

	// filter shares
	if len(content.Shares) > 0 {
		var arr []interface{}
		for _, item := range content.Shares {
			arr = append(arr, map[string]string{"id": item.Id})
		}
		data["shares"] = arr
	} else {
		delete(data, "shares")
	}

	d := map[string]interface{}{
		"_key":       content.Id,
		"id":         content.Id,
		"created":    content.Created,
		"updated":    content.Updated,
		"name":       content.Title,
		"parameter1": content.OrgId,
		"parameter2": createdById,
		"data":       data,
	}

	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToContent(r *db_proto.Record) (*content_proto.Content, error) {
	var p content_proto.Content
	unmarshaler := jsonpb.Unmarshaler{}
	unmarshaler.AllowUnknownFields = true
	if err := unmarshaler.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToContentDetail(r *db_proto.Record) (*content_proto.GetContentDetailResponse_Data, error) {
	var p content_proto.GetContentDetailResponse_Data
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func contentRuleToRecord(contentRule *content_proto.ContentRule) (string, error) {
	data, err := common.MarhalToObject(contentRule)
	if err != nil {
		return "", err
	}
	common.FilterObject(data, "source", contentRule.Source)
	common.FilterObject(data, "sourceType", contentRule.SourceType)
	common.FilterObject(data, "contentType", contentRule.ContentType)
	common.FilterObject(data, "parentCategory", contentRule.ParentCategory)
	common.FilterObject(data, "category", contentRule.Category)
	if len(contentRule.CategoryItems) > 0 {
		var arr []interface{}
		for _, item := range contentRule.CategoryItems {
			arr = append(arr, map[string]string{
				"id": item.Id,
			})
		}
		data["categoryItems"] = arr
	} else {
		delete(data, "categoryItems")
	}

	d := map[string]interface{}{
		"_key":    contentRule.Id,
		"id":      contentRule.Id,
		"created": contentRule.Created,
		"updated": contentRule.Updated,
		// "name":       contentRule.Name,
		"parameter1": contentRule.OrgId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToContentRule(r *db_proto.Record) (*content_proto.ContentRule, error) {
	var p content_proto.ContentRule
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func sharedToRecord(from, to, orgId string, shared *content_proto.ShareContentUser) (string, error) {
	data, err := common.MarhalToObject(shared)
	if err != nil {
		return "", err
	}
	common.FilterObject(data, "content", shared.Content)
	common.FilterObject(data, "user", shared.User)
	common.FilterObject(data, "shared_by", shared.SharedBy)
	var sharedById string
	if shared.SharedBy != nil {
		sharedById = shared.SharedBy.Id
	}
	d := map[string]interface{}{
		"_from":      from,
		"_to":        to,
		"_key":       shared.Id,
		"id":         shared.Id,
		"created":    shared.Created,
		"updated":    shared.Updated,
		"parameter1": orgId,
		"parameter2": sharedById,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToShared(r *db_proto.Record) (*content_proto.ShareContentUser, error) {
	var p content_proto.ShareContentUser
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToCategoryResponse(r *db_proto.Record) (*content_proto.CategoryResponse, error) {
	var p content_proto.CategoryResponse
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToContentResponse(r *db_proto.Record) (*content_proto.ContentResponse, error) {
	var p content_proto.ContentResponse
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToRecommendation(r *db_proto.Record) (*content_proto.Recommendation, error) {
	var p content_proto.Recommendation
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToSharedContent(r *db_proto.Record) (*content_proto.SharedContent, error) {
	var p content_proto.SharedContent
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToShareContentUser(r *db_proto.Record) (*content_proto.ShareContentUser, error) {
	var p content_proto.ShareContentUser
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func contentRecommendationToRecord(from, to string, recommendation *content_proto.ContentRecommendation) (string, error) {
	data, err := common.MarhalToObject(recommendation)
	if err != nil {
		return "", err
	}

	delete(data, "content")

	d := map[string]interface{}{
		"_from":      from,
		"_to":        to,
		"created":    recommendation.Created,
		"updated":    recommendation.Updated,
		"parameter1": recommendation.OrgId,
		"parameter2": recommendation.UserId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToContentRecommendation(r *db_proto.Record) (*content_proto.ContentRecommendation, error) {
	var p content_proto.ContentRecommendation
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToContentCategoryItemResponse(r *db_proto.Record) (*content_proto.ContentCategoryItemResponse, error) {
	var p content_proto.ContentCategoryItemResponse
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToShareableContent(r *db_proto.Record) (*user_proto.ShareableContent, error) {
	var p user_proto.ShareableContent
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func queryFilterSharedForUser(userId, shared_collection, bookmarked_collection, org_query string) (string, string, string) {
	shared_query := ""
	bookmarked_query := ""
	filter_shared_query := ""
	if len(userId) > 0 {
		shared_query = fmt.Sprintf(`LET shared = (
			FOR e, doc IN INBOUND "%v/%v" %v
			%s
			RETURN doc._from
		)`, common.DbUserTable, userId, shared_collection, org_query)
		bookmarked_query = fmt.Sprintf(`LET bookmarked = (
			FOR e, doc IN OUTBOUND "%v/%v" %v
			RETURN doc._to
		)`, common.DbUserTable, userId, bookmarked_collection)
		filter_shared_query = "FILTER doc._id NOT IN shared && doc._id NOT IN bookmarked"
	}
	return shared_query, bookmarked_query, filter_shared_query
}

func queryMerge() string {
	q := fmt.Sprintf(`	
		LET tags = (
			FOR t IN OUTBOUND doc %v
			OPTIONS {
				bfs: true,
				uniqueVertices: "global"
			}
			FILTER NOT_NULL(t.data)
				RETURN t.data
		)
		LET u = (FOR u IN %v FILTER doc.data.createdBy.id == u._key RETURN u)
		LET t = (FOR t IN %v FILTER doc.data.type.id == t._key RETURN t)
		LET s = (FOR s IN %v FILTER doc.data.source.id == s._key RETURN s)
		LET c = (FOR c IN %v FILTER doc.data.category.id == c._key RETURN c)
		LET shares = (FILTER NOT_NULL(doc.data.shares) FOR cu IN doc.data.shares FOR p IN %v FILTER cu.id == p._key RETURN p.data)
		RETURN MERGE_RECURSIVE(doc,{data:{
			createdBy:u[0].data,
			tags:tags,
			type:t[0].data,
			source:s[0].data,
			category:c[0].data,
			shares:shares
		}})`,
		common.DbContentTagEdgeTable, common.DbUserTable,
		common.DbContentTypeTable, common.DbSourceTable, common.DbContentCategoryTable,
		common.DbUserTable,
	)
	return q
}

func queryContentRuleMerge() string {
	q := fmt.Sprintf(`
		LET source = (FOR p IN %v FILTER doc.data.source.id == p._key RETURN p)
		LET sourceType = (FOR p IN %v FILTER doc.data.sourceType.id == p._key RETURN p)
		LET contentType = (FOR p IN %v FILTER doc.data.contentType.id == p._key RETURN p)
		LET parentCategory = (FOR p IN %v FILTER doc.data.parentCategory.id == p._key RETURN p)
		LET category = (FOR p IN %v FILTER doc.data.category.id == p._key RETURN p)
		LET categoryItems = (
			FILTER NOT_NULL(doc.data.categoryItems)
			FOR item IN doc.data.categoryItems
			FOR p IN %v
			FILTER item.id == p._key RETURN p.data)
		RETURN MERGE_RECURSIVE(doc,{data:{
			source:source[0].data,
			sourceType:sourceType[0].data,
			contentType:contentType[0].data,
			parentCategory:parentCategory[0].data,
			category:category[0].data,
			categoryItems:categoryItems
		}})`,
		common.DbSourceTable, common.DbContentSourceTypeTable, common.DbContentTypeTable,
		common.DbContentParentCategoryTable, common.DbContentCategoryTable, common.DbContentCategoryItemTable,
	)
	return q
}

// AllSources get all sources
func AllSources(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*content_proto.Source, error) {
	var sources []*content_proto.Source
	query := common.QueryAuth(`FILTER`, orgId, "")
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		LET t = (FOR t IN %v FILTER doc.data.type.id == t._key RETURN t)
		RETURN MERGE_RECURSIVE(doc,{data:{type:t[0].data}})`, common.DbSourceTable, query, sort_query, limit_query, common.DbContentSourceTypeTable)

	resp, err := runQuery(ctx, q, common.DbSourceTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if source, err := recordToSource(r); err == nil {
			sources = append(sources, source)
		}
	}
	return sources, nil
}

// CreateSource creates a source
func CreateSource(ctx context.Context, source *content_proto.Source) error {
	if source.Created == 0 {
		source.Created = time.Now().Unix()
	}
	source.Updated = time.Now().Unix()

	record, err := sourceToRecord(source)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { _key: "%v" } 
		INSERT %v 
		UPDATE %v 
		IN %v`, source.Id, record, record, common.DbSourceTable)
	_, err = runQuery(ctx, q, common.DbSourceTable)
	return err
}

// ReadSource reads a source by ID
func ReadSource(ctx context.Context, id, orgId, teamId string) (*content_proto.Source, error) {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	query = common.QueryAuth(query, orgId, teamId)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		LET t = (FOR t IN %v FILTER doc.data.type.id == t._key RETURN t)
		RETURN MERGE_RECURSIVE(doc,{data:{type:t[0].data}})`, common.DbSourceTable, query, common.DbContentSourceTypeTable)

	resp, err := runQuery(ctx, q, common.DbSourceTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToSource(resp.Records[0])
	return data, err
}

// DeleteSource deletes a source by ID
func DeleteSource(ctx context.Context, id, orgId, teamId string) error {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	query = common.QueryAuth(query, orgId, teamId)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		REMOVE doc IN %v`, common.DbSourceTable, query, common.DbSourceTable)
	_, err := runQuery(ctx, q, common.DbSourceTable)
	return err
}

// AllTaxonomys get all taxonomys
func AllTaxonomys(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*static_proto.Taxonomy, error) {
	var taxonomys []*static_proto.Taxonomy
	query := common.QueryAuth(`FILTER`, orgId, "")
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		RETURN doc`, common.DbTaxonomyTable, query, sort_query, limit_query)

	resp, err := runQuery(ctx, q, common.DbTaxonomyTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if taxonomy, err := recordToTaxonomy(r); err == nil {
			taxonomys = append(taxonomys, taxonomy)
		}
	}
	return taxonomys, nil
}

// CreateTaxonomy creates a taxonomy
func CreateTaxonomy(ctx context.Context, taxonomy *static_proto.Taxonomy) error {
	if taxonomy.Created == 0 {
		taxonomy.Created = time.Now().Unix()
	}
	taxonomy.Updated = time.Now().Unix()

	record, err := taxonomyToRecord(taxonomy)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { _key: "%v" } 
		INSERT %v 
		UPDATE %v 
		IN %v`, taxonomy.Id, record, record, common.DbTaxonomyTable)
	_, err = runQuery(ctx, q, common.DbTaxonomyTable)
	return err
}

// ReadTaxonomy reads a taxonomy by ID
func ReadTaxonomy(ctx context.Context, id, orgId, teamId string) (*static_proto.Taxonomy, error) {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	query = common.QueryAuth(query, orgId, "")

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		RETURN doc`, common.DbTaxonomyTable, query)

	resp, err := runQuery(ctx, q, common.DbTaxonomyTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToTaxonomy(resp.Records[0])
	return data, err
}

// DeleteTaxonomy deletes a taxonomy by ID
func DeleteTaxonomy(ctx context.Context, id, orgId, teamId string) error {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	query = common.QueryAuth(query, orgId, teamId)

	q := fmt.Sprintf(`
		FOR doc
		IN %v
		%s
		REMOVE doc IN %v`, common.DbTaxonomyTable, query, common.DbTaxonomyTable)
	_, err := runQuery(ctx, q, common.DbTaxonomyTable)
	return err
}

// AllContentCategoryItems get all contentCategoryItems
func AllContentCategoryItems(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*static_proto.ContentCategoryItem, error) {
	var contentCategoryItems []*static_proto.ContentCategoryItem
	query := common.QueryAuth(`FILTER`, orgId, "")
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		LET t = (FOR t IN %v FILTER doc.data.taxonomy.id == t._key RETURN t.data)
		LET c = (FOR c IN %v FILTER doc.data.category.id == c._key RETURN c.data)
		RETURN MERGE_RECURSIVE(doc,{data:{taxonomy:t[0],category:c[0]}})`,
		common.DbContentCategoryItemTable, query, sort_query, limit_query,
		common.DbTaxonomyTable, common.DbContentCategoryTable,
	)

	resp, err := runQuery(ctx, q, common.DbContentCategoryItemTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if contentCategoryItem, err := recordToContentCategoryItem(r); err == nil {
			contentCategoryItems = append(contentCategoryItems, contentCategoryItem)
		} else {
			log.Error(err)
		}
	}
	return contentCategoryItems, nil
}

// CreateContentCategoryItem creates a contentCategoryItem
func CreateContentCategoryItem(ctx context.Context, contentCategoryItem *static_proto.ContentCategoryItem) error {
	if contentCategoryItem.Created == 0 {
		contentCategoryItem.Created = time.Now().Unix()
	}
	contentCategoryItem.Updated = time.Now().Unix()

	record, err := contentCategoryItemToRecord(contentCategoryItem)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { _key: "%v" } 
		INSERT %v 
		UPDATE %v 
		IN %v`, contentCategoryItem.Id, record, record, common.DbContentCategoryItemTable)
	_, err = runQuery(ctx, q, common.DbContentCategoryItemTable)
	return err
}

// ReadContentCategoryItem reads a contentCategoryItem by ID
func ReadContentCategoryItem(ctx context.Context, id, orgId, teamId string) (*static_proto.ContentCategoryItem, error) {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	query = common.QueryAuth(query, orgId, "")

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		LET t = (FOR t IN %v FILTER doc.data.taxonomy.id == t._key RETURN t.data)
		LET c = (FOR c IN %v FILTER doc.data.category.id == c._key RETURN c.data)
		RETURN MERGE_RECURSIVE(doc,{data:{taxonomy:t[0],category:c[0]}})`,
		common.DbContentCategoryItemTable, query,
		common.DbTaxonomyTable, common.DbContentCategoryTable,
	)

	resp, err := runQuery(ctx, q, common.DbContentCategoryItemTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToContentCategoryItem(resp.Records[0])
	return data, err
}

// DeleteContentCategoryItem deletes a contentCategoryItem by ID
func DeleteContentCategoryItem(ctx context.Context, id, orgId, teamId string) error {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	query = common.QueryAuth(query, orgId, teamId)

	q := fmt.Sprintf(`
		FOR doc
		IN %v
		%s
		REMOVE doc IN %v`, common.DbContentCategoryItemTable, query, common.DbContentCategoryItemTable)
	_, err := runQuery(ctx, q, common.DbContentCategoryItemTable)
	return err
}

// AllContents get all contents
func AllContents(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*content_proto.Content, error) {
	var contents []*content_proto.Content
	query := common.QueryAuth(`FILTER`, orgId, "")
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		%s`,
		common.DbContentTable, query, sort_query, limit_query,
		queryMerge(),
	)

	resp, err := runQuery(ctx, q, common.DbContentTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if content, err := recordToContent(r); err == nil {
			contents = append(contents, content)
		} else {
			log.Error(err)
		}
	}
	return contents, nil
}

// CreateContent creates a content
func CreateContent(ctx context.Context, content *content_proto.Content) error {
	if len(content.Id) == 0 {
		content.Id = uuid.NewUUID().String()
	}
	if content.Created == 0 {
		content.Created = time.Now().Unix()
	}
	content.Updated = time.Now().Unix()

	record, err := contentToRecord(content)
	if err != nil {
		log.WithField("err", err).Error("parsing failed")
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { _key: "%v" } 
		INSERT %v 
		UPDATE %v 
		IN %v`, content.Id, record, record, common.DbContentTable)
	_, err = runQuery(ctx, q, common.DbContentTable)
	if err != nil {
		return err
	}

	// remove releated edge first
	_from := fmt.Sprintf(`%v/%v`, common.DbContentTable, content.Id)
	q = fmt.Sprintf(`
		FOR doc IN %v 
		FILTER doc._from == "%v"
		REMOVE doc IN %v`, common.DbContentTagEdgeTable, _from, common.DbContentTagEdgeTable)
	if _, err := runQuery(ctx, q, common.DbContentTagEdgeTable); err != nil {
		return common.InternalServerError(common.ContentSrv, CreateContent, err, "Remove ContentTagEdge failed")
	}
	//FIXME:for unknown tags, this will throw error.
	for _, tag := range content.Tags {
		field := fmt.Sprintf(`{_from:"%v",_to:"%v/%v"} `, _from, common.DbContentCategoryItemTable, tag.Id)
		q = fmt.Sprintf(`INSERT %v INTO %v`, field, common.DbContentTagEdgeTable)
		if _, err := runQuery(ctx, q, common.DbContentTagEdgeTable); err != nil {
			return err
		}
	}

	return nil
}

// ReadContent reads a content by ID
func ReadContent(ctx context.Context, id, orgId, teamId string) (*content_proto.Content, error) {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	query = common.QueryAuth(query, orgId, teamId)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s`, common.DbContentTable, query,
		queryMerge(),
	)

	resp, err := runQuery(ctx, q, common.DbContentTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToContent(resp.Records[0])
	return data, err
}

// DeleteContent deletes a content by ID
func DeleteContent(ctx context.Context, id, orgId, teamId string) error {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	query = common.QueryAuth(query, orgId, teamId)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		REMOVE doc IN %v`, common.DbContentTable, query, common.DbContentTable)
	_, err := runQuery(ctx, q, common.DbContentTable)
	if err != nil {
		return err
	}

	query = fmt.Sprintf(`FILTER doc._from == "%v/%v"`, common.DbContentTable, id)
	q = fmt.Sprintf(`
		FOR doc IN %v
		%s
		REMOVE doc IN %v`, common.DbContentTagEdgeTable, query, common.DbContentTagEdgeTable)
	_, err = runQuery(ctx, q, common.DbContentTagEdgeTable)
	return err
}

// AllContentRules get all contentRules
func AllContentRules(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*content_proto.ContentRule, error) {
	var contentRules []*content_proto.ContentRule
	query := common.QueryAuth(`FILTER`, orgId, "")
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		%s`, common.DbContentRuleTable, query, sort_query, limit_query,
		queryContentRuleMerge(),
	)

	resp, err := runQuery(ctx, q, common.DbContentRuleTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if contentRule, err := recordToContentRule(r); err == nil {
			contentRules = append(contentRules, contentRule)
		} else {
			log.Error(err)
		}
	}
	return contentRules, nil
}

// CreateContentRule creates a contentRule
func CreateContentRule(ctx context.Context, contentRule *content_proto.ContentRule) error {
	if contentRule.Created == 0 {
		contentRule.Created = time.Now().Unix()
	}
	contentRule.Updated = time.Now().Unix()

	record, err := contentRuleToRecord(contentRule)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { _key: "%v" } 
		INSERT %v 
		UPDATE %v 
		IN %v`, contentRule.Id, record, record, common.DbContentRuleTable)
	_, err = runQuery(ctx, q, common.DbContentRuleTable)
	return err
}

// ReadContentRule reads a contentRule by ID
func ReadContentRule(ctx context.Context, id, orgId, teamId string) (*content_proto.ContentRule, error) {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	query = common.QueryAuth(query, orgId, teamId)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s`, common.DbContentRuleTable, query,
		queryContentRuleMerge(),
	)
	resp, err := runQuery(ctx, q, common.DbContentRuleTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToContentRule(resp.Records[0])
	return data, err
}

// DeleteContentRule deletes a contentRule by ID
func DeleteContentRule(ctx context.Context, id, orgId, teamId string) error {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	query = common.QueryAuth(query, orgId, "")

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		REMOVE doc IN %v`, common.DbContentRuleTable, query, common.DbContentRuleTable)
	_, err := runQuery(ctx, q, common.DbContentRuleTable)
	return err
}

func FilterContent(ctx context.Context, req *content_proto.FilterContentRequest) ([]*content_proto.Content, error) {
	var contents []*content_proto.Content

	query := `FILTER`

	if len(req.Sources) > 0 {
		sources := common.QueryStringFromArray(req.Sources)
		query += fmt.Sprintf(" && doc.data.source.id IN [%v]", sources)
	}

	// if len(req.SourceTypes) > 0 {
	// 	sourceTypes := common.QueryStringFromArray(req.SourceTypes)
	// 	query += fmt.Sprintf(" && doc.data.sourceTypeId IN [%v]", sourceTypes)
	// }

	if len(req.Type) > 0 {
		types := common.QueryStringFromArray(req.Type)
		query += fmt.Sprintf(` && doc.data.item["@type"] IN [%v]`, types)
	}

	if len(req.ContentTypes) > 0 {
		contentTypes := common.QueryStringFromArray(req.ContentTypes)
		query += fmt.Sprintf(" && doc.data.type.id IN [%v]", contentTypes)
	}

	if len(req.CreatedBy) > 0 {
		createdBy := common.QueryStringFromArray(req.CreatedBy)
		query += fmt.Sprintf(" && doc.data.createdBy.id IN [%v]", createdBy)
	}

	if len(req.ContentCategories) > 0 {
		contentCategories := common.QueryStringFromArray(req.ContentCategories)
		query += fmt.Sprintf(" && c[0].data.name IN [%v]", contentCategories)
	}

	if len(req.ContentCategoryItems) > 0 {
		contentCategoryItems := common.QueryStringFromArray(req.ContentCategoryItems)
		query += fmt.Sprintf(" && tags[*].id ANY IN [%v]", contentCategoryItems)
	}
	query = common.QueryAuth(query, req.OrgId, "")
	limit_query := common.QueryPaginate(req.Offset, req.Limit)
	sort_query := common.QuerySort(req.SortParameter, req.SortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		LET c = (FOR c IN %v FILTER doc.data.category.id == c._key RETURN c)
		LET tags = (
			FOR t IN OUTBOUND doc %v
			OPTIONS {
				bfs: true,
				uniqueVertices: "global"
			}
			FILTER NOT_NULL(t.data)
				RETURN t.data
		)
		%s
		%s
		%s
		LET u = (FOR u IN %v FILTER doc.data.createdBy.id == u._key RETURN u)
		LET t = (FOR t IN %v FILTER doc.data.type.id == t._key RETURN t)
		LET s = (FOR s IN %v FILTER doc.data.source.id == s._key RETURN s)
		RETURN MERGE_RECURSIVE(doc,{data:{
			createdBy:u[0].data,
			tags:tags,
			type:t[0].data,
			source:s[0].data,
			category:c[0].data
		}})`, common.DbContentTable, common.DbContentCategoryTable, common.DbContentTagEdgeTable,
		query, sort_query, limit_query,
		common.DbUserTable, common.DbContentTypeTable, common.DbSourceTable,
	)

	resp, err := runQuery(ctx, q, common.DbContentTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if content, err := recordToContent(r); err == nil {
			contents = append(contents, content)
		}
	}

	return contents, nil
}

func ShareContent(ctx context.Context, contents []*content_proto.Content, users []*user_proto.User, sharedBy *user_proto.User, orgId string) ([]string, error) {
	userids := []string{}
	for _, content := range contents {
		for _, user := range users {
			shared := &content_proto.ShareContentUser{
				Id:       uuid.NewUUID().String(),
				Status:   static_proto.ShareStatus_SHARED,
				Updated:  time.Now().Unix(),
				Created:  time.Now().Unix(),
				Content:  content,
				User:     user,
				SharedBy: sharedBy,
			}

			_from := fmt.Sprintf(`%v/%v`, common.DbContentTable, content.Id)
			_to := fmt.Sprintf(`%v/%v`, common.DbUserTable, user.Id)
			record, err := sharedToRecord(_from, _to, orgId, shared)
			if err != nil {
				return nil, err
			}
			if len(record) == 0 {
				return nil, errors.New("server serialization")
			}

			field := fmt.Sprintf(`{_from:"%v",_to:"%v"} `, _from, _to)
			q := fmt.Sprintf(`
				UPSERT %v
				INSERT %v
				UPDATE %v
				INTO %v
				RETURN {data:{user_id: OLD ? "" : NEW.data.user.id}}`, field, record, record, common.DbShareContentUserEdgeTable)

			resp, err := runQuery(ctx, q, common.DbShareContentUserEdgeTable)
			if err != nil {
				return nil, err
			}

			// parsing to check whether this was an update (returns nothing) or insert (returns inserted user_id)
			b, err := common.RecordToInsertedUserId(resp.Records[0])
			if err != nil {
				return nil, err
			}
			if len(b) > 0 {
				userids = append(userids, b)
			}

			// save pending
			any, err := common.FilteredAnyFromObject(common.CONTENT_TYPE, content.Id)
			if err != nil {
				return nil, err
			}
			// save pending
			pending := &common_proto.Pending{
				Id:         uuid.NewUUID().String(),
				Created:    shared.Created,
				Updated:    shared.Updated,
				SharedBy:   sharedBy,
				SharedWith: user,
				Item:       any,
				OrgId:      orgId,
			}

			q1, err1 := common.SavePending(pending, content.Id)
			if err1 != nil {
				return nil, err1
			}

			if _, err := runQuery(ctx, q1, common.DbPendingTable); err != nil {
				return nil, err
			}
		}
	}
	return userids, nil
}

func GetContentCategorys(ctx context.Context) ([]*content_proto.CategoryResponse, error) {
	var categories []*content_proto.CategoryResponse

	q := fmt.Sprintf(`
		FOR c IN %v
		FOR doc IN %v 
		FILTER c.data.category == doc._key
		RETURN DISTINCT {data:{
			category_id:doc.data.id,
			name:doc.data.name,
			icon_slug:doc.data.icon_slug}}`,
		common.DbContentTable, common.DbContentCategoryTable,
	)

	resp, err := runQuery(ctx, q, common.DbContentTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if category, err := recordToCategoryResponse(r); err == nil {
			categories = append(categories, category)
		}
	}
	return categories, nil
}

func GetContentDetail(ctx context.Context, contentId string) (*content_proto.GetContentDetailResponse_Data, error) {
	content, err := ReadContent(ctx, contentId, "", "")
	if err != nil {
		common.NotFound(common.ContentSrv, GetContentDetail, err, "ReadContent query is failed")
		return nil, err
	}

	query := fmt.Sprintf(`FILTER doc._key == "%v"`, contentId)
	_to := fmt.Sprintf("%v/%v", common.DbContentTable, contentId)
	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		LET bookmark  = (
			FOR u, edge IN ANY "%v" %v
			RETURN COUNT(edge)
		)
		LET rating = (
			FOR u, edge IN ANY "%v" %v
			RETURN edge.data.rating
		)
		RETURN {"data":{"bookmarked":bookmark[0]>0, "rating":rating[0]}}`,
		common.DbContentTable, query,
		_to, common.DbBookmarkEdgeTable,
		_to, common.DbContentRatingEdgeTable)

	resp, err := runQuery(ctx, q, common.DbContentTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}
	detail, err := recordToContentDetail(resp.Records[0])
	detail.Content = content
	return detail, err
}

func GetContentByCategory(ctx context.Context, categoryId string, offset, limit int64) ([]*content_proto.ContentResponse, error) {
	var contents []*content_proto.ContentResponse

	limit_query := common.QueryPaginate(offset, limit)
	q := fmt.Sprintf(`
		FOR doc IN %v
		FILTER doc.data.category.id == "%v"
		%v
		LET category = (FOR category IN %v FILTER category._key == doc.data.category.id RETURN category)
		LET s = (FOR s IN %v FILTER doc.data.source.id == s._key RETURN s.data)
		RETURN {data:{
			content_id:doc.id,
			image:doc.data.image,
			title:doc.name,
			author:doc.data.author,
			source:s[0],
			category_id:category[0].data.id,
			category_icon_slug:category[0].data.data.icon_slug,
			category_name:category[0].data.data.name}}`,
		common.DbContentTable, categoryId, limit_query,
		common.DbContentCategoryTable, common.DbSourceTable,
	)

	resp, err := runQuery(ctx, q, common.DbContentTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if content, err := recordToContentResponse(r); err == nil {
			contents = append(contents, content)
		}
	}
	return contents, nil
}

func GetContentCategoryItemsByCategory(ctx context.Context, categoryId string) ([]*static_proto.ContentCategoryItem, error) {
	var contentCategoryItems []*static_proto.ContentCategoryItem
	query := fmt.Sprintf(`FILTER doc.data.category.id == "%v"`, categoryId)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		LET t = (FOR t IN %v FILTER doc.data.taxonomy.id == t._key RETURN t.data)
		LET c = (FOR c IN %v FILTER doc.data.category.id == c._key RETURN c.data)
		RETURN MERGE_RECURSIVE(doc,{data:{taxonomy:t[0],category:c[0]}})`,
		common.DbContentCategoryItemTable, query,
		common.DbTaxonomyTable, common.DbContentCategoryTable,
	)

	resp, err := runQuery(ctx, q, common.DbContentCategoryItemTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if contentCategoryItem, err := recordToContentCategoryItem(r); err == nil {
			contentCategoryItems = append(contentCategoryItems, contentCategoryItem)
		}
	}
	return contentCategoryItems, nil
}

func FilterCategoryAutocomplete(ctx context.Context, categoryId, name string) ([]*static_proto.ContentCategoryItem, error) {
	var items []*static_proto.ContentCategoryItem
	query := fmt.Sprintf(`FILTER CONTAINS(doc.name, "%v")`, name)
	if len(categoryId) > 0 {
		query += fmt.Sprintf(` && doc.data.category.id == "%v"`, categoryId)
	}
	q := fmt.Sprintf(`
		FOR doc IN %v
		%v
		LET t = (FOR t IN %v FILTER doc.data.taxonomy.id == t._key RETURN t.data)
		LET c = (FOR c IN %v FILTER doc.data.category.id == c._key RETURN c.data)
		RETURN MERGE_RECURSIVE(doc,{data:{taxonomy:t[0],category:c[0]}})`,
		common.DbContentCategoryItemTable, query,
		common.DbTaxonomyTable, common.DbContentCategoryTable,
	)

	resp, err := runQuery(ctx, q, common.DbContentCategoryItemTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if item, err := recordToContentCategoryItem(r); err == nil {
			items = append(items, item)
		}
	}
	return items, nil
}

func FilterContentInParticularCategory(ctx context.Context, categoryId string, contentCategoryItems []string) ([]*content_proto.ContentResponse, error) {
	var contents []*content_proto.ContentResponse

	query := fmt.Sprintf(`FILTER doc.data.category.id == "%v"`, categoryId)
	if len(contentCategoryItems) > 0 {
		tags := common.QueryStringFromArray(contentCategoryItems)
		query += fmt.Sprintf(" && tags[*]._key ANY IN [%v]", tags)
	}

	q := fmt.Sprintf(`
		FOR doc IN %v
		LET tags = (
			FOR t IN OUTBOUND doc %v
			OPTIONS {
				bfs: true,
				uniqueVertices: "global"
			}
			RETURN t
		)
		%v
		LET s = (FOR s IN %v FILTER doc.data.source.id == s._key RETURN s.data)
		LET c = (FOR c IN %v FILTER doc.data.category.id == c._key RETURN c.data)
		RETURN {data:{
			content_id:doc.id,
			image:doc.data.image,
			title:doc.name,
			author:doc.data.author,
			source:s[0],
			category_id:c[0].id,
			category_icon_slug:c[0].icon_slug,
			category_name:c[0].name}}`,
		common.DbContentTable, common.DbContentTagEdgeTable,
		query, common.DbSourceTable, common.DbContentCategoryTable,
	)

	resp, err := runQuery(ctx, q, common.DbContentTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if content, err := recordToContentResponse(r); err == nil {
			contents = append(contents, content)
		}
	}
	return contents, nil
}

func CreateContentRecommendation(ctx context.Context, recommendation *content_proto.ContentRecommendation) error {
	if len(recommendation.Id) == 0 {
		recommendation.Id = uuid.NewUUID().String()
	}
	if recommendation.Created == 0 {
		recommendation.Created = time.Now().Unix()
	}
	recommendation.Updated = time.Now().Unix()

	//this should be the other way around
	_from := fmt.Sprintf(`%v/%v`, common.DbUserTable, recommendation.UserId)
	_to := fmt.Sprintf(`%v/%v`, common.DbContentTable, recommendation.Content.Id)
	record, err := contentRecommendationToRecord(_from, _to, recommendation)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	field := fmt.Sprintf(`{_from:"%v",_to:"%v"} `, _from, _to)
	q := fmt.Sprintf(`
		UPSERT %v
		INSERT %v 
		UPDATE %v 
		IN %v`, field, record, record, common.DbContentRecommendationEdgeTable)
	_, err = runQuery(ctx, q, common.DbContentRecommendationEdgeTable)
	return err
}

func GetContentRecommendationByUser(ctx context.Context, userId, orgId string, offset, limit int64, sortParameter, sortDirection string) ([]*content_proto.Recommendation, error) {
	recommendations := []*content_proto.Recommendation{}
	_from := fmt.Sprintf(`%v/%v`, common.DbUserTable, userId)

	query := common.QueryAuth(`FILTER`, orgId, "")
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR content, doc IN OUTBOUND "%v" %v
		%v
		%v
		%v
		LET tags = (FOR tag IN OUTBOUND content %v FILTER NOT_NULL(tag) RETURN tag.data)
		LET s = (FOR s IN %v FILTER content.data.source.id == s._key RETURN s)
		LET c = (FOR c IN %v FILTER content.data.category.id == c._key RETURN c)
		RETURN {data:{
			content_id:content.data.id,
			content_image:content.data.image,
			content_title:content.data.title,
			content_author:content.data.author,
			content_source:s[0].data,
			category_id:c[0].data.id,
			category_icon_slug:c[0].data.icon_slug,
			category_name:c[0].data.name,
			user_id:doc.data.user_id,
			tags:tags}}`,
		_from, common.DbContentRecommendationEdgeTable,
		query,
		sort_query,
		limit_query,
		common.DbContentTagEdgeTable,
		common.DbSourceTable, common.DbContentCategoryTable)

	resp, err := runQuery(ctx, q, common.DbContentRecommendationEdgeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if recommendation, err := recordToRecommendation(r); err == nil {
			recommendations = append(recommendations, recommendation)
		}
	}
	return recommendations, nil
}

func GetContentRecommendationByCategory(ctx context.Context, userId, orgId, categoryId string, offset, limit int64, sortParameter, sortDirection string) ([]*content_proto.Recommendation, error) {
	recommendations := []*content_proto.Recommendation{}

	query := common.QueryAuth(`FILTER`, orgId, userId)
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%v
		FOR content IN OUTBOUND doc._from %v
		FILTER content.data.category.id == "%v"
		%v
		%v
		LET tags = (FOR tag IN OUTBOUND content %v 
			OPTIONS {
				bfs: true,
				uniqueVertices: "global"
			}
			FILTER NOT_NULL(tag) 
			RETURN tag.data
		)
		LET s = (FOR s IN %v FILTER content.data.source.id == s._key RETURN s)
		LET c = (FOR c IN %v FILTER content.data.category.id == c._key RETURN c)
		RETURN {data:{
			content_id:content.data.id,
			content_image:content.data.image,
			content_title:content.data.title,
			content_author:content.data.author,
			content_source:s[0].data,
			category_id:c[0].data.id,
			category_icon_slug:c[0].data.icon_slug,
			category_name:c[0].data.name,
			user_id:doc.data.user_id,
			tags:tags}}`,
		common.DbContentRecommendationEdgeTable,
		query,
		common.DbContentRecommendationEdgeTable,
		categoryId,
		sort_query,
		limit_query,
		common.DbContentTagEdgeTable,
		common.DbSourceTable, common.DbContentCategoryTable)

	resp, err := runQuery(ctx, q, common.DbContentRecommendationEdgeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if recommendation, err := recordToRecommendation(r); err == nil {
			recommendations = append(recommendations, recommendation)
		}
	}
	return recommendations, nil
}

func GetRandomItems(ctx context.Context, count int64) ([]*content_proto.Content, error) {
	contents := []*content_proto.Content{}

	q := fmt.Sprintf(`
		FOR doc IN %v
		SORT RAND() 
		COLLECT category_id = doc.data.category.id INTO data
		LIMIT %v 
		RETURN data[0].doc`, common.DbContentTable, count)
	resp, err := runQuery(ctx, q, common.DbContentTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if content, err := recordToContent(r); err == nil {
			contents = append(contents, content)
		}
	}
	return contents, nil
}

func GetAllSharedContents(ctx context.Context, userId, orgId string, offset, limit int64) ([]*content_proto.ShareContentUser, error) {
	shares := []*content_proto.ShareContentUser{}

	_to := fmt.Sprintf("%v/%v", common.DbUserTable, userId)
	q := fmt.Sprintf(`
		FOR content,doc IN INBOUND "%v" %v
		FOR u IN %v FILTER u._key == "%v"
		FOR shared_by IN %v FILTER shared_by._key == doc.data.shared_by.id
		RETURN MERGE_RECURSIVE(doc, {data:{
			content:content.data,
			user:u.data,
			shared_by:shared_by.data
		}})`, _to, common.DbShareContentUserEdgeTable, common.DbUserTable, userId, common.DbUserTable)

	resp, err := runQuery(ctx, q, common.DbShareContentUserEdgeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if s, err := recordToShareContentUser(r); err == nil {
			shares = append(shares, s)
		}
	}

	return shares, nil
}

func GetContentRecommendations(ctx context.Context, userId, orgId string) ([]*content_proto.ContentRecommendation, error) {
	recommendations := []*content_proto.ContentRecommendation{}

	query := common.QueryAuth(`FILTER`, orgId, "")
	q := fmt.Sprintf(`
		FOR content,doc IN OUTBOUND "%v/%v" %v
		%v
		LET tags = (
			FOR t IN OUTBOUND content %v
			OPTIONS {
				bfs: true,
				uniqueVertices: "global"
			}
			RETURN t.data
		)
		RETURN MERGE_RECURSIVE(doc, {data:{
			content:content.data,
			tags:tags
		}})`, common.DbUserTable, userId,
		common.DbContentRecommendationEdgeTable, query,
		common.DbContentTagEdgeTable,
	)

	resp, err := runQuery(ctx, q, common.DbContentRecommendationEdgeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if recommendation, err := recordToContentRecommendation(r); err == nil {
			recommendations = append(recommendations, recommendation)
		}
	}
	return recommendations, nil
}

func GetContentFiltersByPreference(ctx context.Context, userId, orgId string) ([]*static_proto.ContentCategoryItem, error) {
	contentCategoryItems := []*static_proto.ContentCategoryItem{}

	q := fmt.Sprintf(`
		LET ret = (
			FOR doc IN %v
			FILTER doc.parameter2 == "%v"
			RETURN APPEND(doc.data.allergies, APPEND(doc.data.conditions, APPEND(doc.data.food, APPEND(doc.data.cuisines, doc.data.ethinicties))))
		)[**]
		FOR r IN ret
		RETURN {data:r}`, common.DbPreferenceTable, userId)

	resp, err := runQuery(ctx, q, common.DbPreferenceTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if item, err := recordToContentCategoryItem(r); err == nil {
			contentCategoryItems = append(contentCategoryItems, item)
		}
	}
	return contentCategoryItems, nil
}

func FilterContentRecommendations(ctx context.Context, userId, orgId string, contentCategoryItems []*static_proto.ContentCategoryItem) ([]*content_proto.ContentResponse, error) {
	response := []*content_proto.ContentResponse{}

	query := common.QueryAuth(`FILTER`, orgId, "")
	var tags_query string
	if len(contentCategoryItems) > 0 {
		items := []string{}
		for _, item := range contentCategoryItems {
			items = append(items, `"`+item.Id+`"`)
		}
		tags := strings.Join(items[:], ",")
		tags_query = fmt.Sprintf("FILTER tags[*].id ANY IN [%v]", tags)
	}

	q := fmt.Sprintf(`
		FOR content,doc IN OUTBOUND "%v/%v" %v
		%v
		LET tags = (
			FOR t IN OUTBOUND content %v
			OPTIONS {
				bfs: true,
				uniqueVertices: "global"
			}
			RETURN t.data
		)
		%v
		LET s = (FOR s IN %v FILTER content.data.source.id == s._key RETURN s.data)
		LET c = (FOR c IN %v FILTER content.data.category.id == c._key RETURN c.data)
		RETURN {data:{
			content_id: content.data.id,
			image: content.data.image,
			title: content.data.title,
			author: content.data.author,
			source: s[0],
			category_id: c[0].id,
			category_icon_slug: c[0].icon_slug,
			category_name: c[0].name
		}}`, common.DbUserTable, userId, common.DbContentRecommendationEdgeTable, query,
		common.DbContentTagEdgeTable, tags_query, common.DbSourceTable, common.DbContentCategoryTable)

	resp, err := runQuery(ctx, q, common.DbContentRecommendationEdgeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if resp, err := recordToContentResponse(r); err == nil {
			response = append(response, resp)
		}
	}
	return response, nil
}

func AutocompleteTags(ctx context.Context, orgId, name string) ([]string, error) {
	var tags []string

	var query string
	if len(orgId) > 0 {
		query = fmt.Sprintf(`FILTER doc.parameter1 == "%v"`, orgId)
	}

	q := fmt.Sprintf(`
		LET tags = (
			FOR doc IN %v
			%v
			FOR tag IN OUTBOUND doc._id %v
			RETURN tag.name
		)
		FOR t IN tags
		FILTER LIKE(t,"%v",true)
		LET ret = {parameter1:t}
		RETURN DISTINCT ret`,
		common.DbContentTable, query,
		common.DbContentTagEdgeTable, name,
	)
	resp, err := runQuery(ctx, q, common.DbContentTable)

	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		tags = append(tags, r.Parameter1)
	}
	return tags, nil
}

func AutocompleteContentCategoryItem(ctx context.Context, categoryId, name string) ([]*content_proto.ContentCategoryItemResponse, error) {
	var response []*content_proto.ContentCategoryItemResponse
	// search categoryitem
	q := fmt.Sprintf(`
		FOR doc IN %v
		FILTER doc.data.category.id == "%v" && LIKE(doc.name, "%v",true)
		FOR cat in content_category
		FILTER cat.id == doc.data.category.id
		RETURN {data:{
			category_id: cat.id,
			category_nameslug: cat.data.name_slug,
			categoryitem_id: doc.data.id,
			categoryitem_name: doc.data.name
		}}`, common.DbContentCategoryItemTable, categoryId, `%`+name+`%`)
	resp, err := runQuery(ctx, q, common.DbContentCategoryItemTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if contentCategoryItem, err := recordToContentCategoryItemResponse(r); err == nil {
			response = append(response, contentCategoryItem)
		}
	}
	return response, nil
}

func AllContentCategoryItemByNameslug(ctx context.Context, categoryId string) ([]*content_proto.ContentCategoryItemResponse, error) {
	var response []*content_proto.ContentCategoryItemResponse
	// search categoryitem
	q := fmt.Sprintf(`
		FOR doc IN %v
		FILTER doc.data.category.id == "%v"
		RETURN {data:{
			category_id: doc.data.category.id,
			category_nameslug: doc.data.category.name_slug,
			categoryitem_id: doc.data.id,
			categoryitem_name: doc.data.name
		}}`, common.DbContentCategoryItemTable, categoryId)
	resp, err := runQuery(ctx, q, common.DbContentCategoryItemTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if contentCategoryItem, err := recordToContentCategoryItemResponse(r); err == nil {
			response = append(response, contentCategoryItem)
		}
	}
	return response, nil
}

func SearchContent(ctx context.Context, title, summary, description string, orgId, teamId string, offset, limit int64) (*content_proto.SearchContentResponse_Data, error) {
	data := &content_proto.SearchContentResponse_Data{}
	var contents []*content_proto.Content

	query := `FILTER`
	if len(title) > 0 {
		query += fmt.Sprintf(` && LIKE(doc.data.title, "%s",true)`, `%`+title+`%`)
	}
	if len(description) > 0 {
		query += fmt.Sprintf(` && LIKE(doc.data.description, "%v",true)`, `%`+description+`%`)
	}
	if len(summary) > 0 {
		query += fmt.Sprintf(` && LIKE(doc.data.summary, "%v",true)`, `%`+summary+`%`)
	}
	query = common.QueryAuth(query, orgId, "")
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort("", "")

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		RETURN doc`, common.DbContentTable, query, sort_query, limit_query)
	resp, err := runQuery(ctx, q, common.DbGoalTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if content, err := recordToContent(r); err == nil {
			contents = append(contents, content)
		}
	}
	data.Contents = contents
	return data, nil
}

//FIXME:Add type filter and shared_by filter
func GetShareableContents(ctx context.Context, createdBy, typez []string, search_term, userId, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*user_proto.ShareableContent, error) {
	var response []*user_proto.ShareableContent
	query := common.QueryAuth(`FILTER`, orgId, "")
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	shared_query, bookmarked_query, filter_shared_query := queryFilterSharedForUser(userId, common.DbShareContentUserEdgeTable, common.DbBookmarkEdgeTable, query)

	filter_query := `FILTER`
	//filter by createdBy
	if len(createdBy) > 0 {
		createdBys := common.QueryStringFromArray(createdBy)
		filter_query += fmt.Sprintf(" && doc.data.createdBy.id IN [%v]", createdBys)
	}

	//filter by search term
	filter_query += common.QuerySharedResourceSearch(filter_query, search_term, "doc")

	q := fmt.Sprintf(`
		%s
		%s
		FOR doc IN %v
		%s
		%s
		%s
		%s
		%s
		LET createdBy = (FOR p IN %v FILTER doc.data.createdBy.id == p._key RETURN p.data)
		RETURN {data:{
			id:doc.id,
			title:doc.name,
			org_id: doc.data.org_id,
			summary: doc.data.summary,
			shared_by: {"id": createdBy[0].id, "firstname": createdBy[0].firstname, "lastname": createdBy[0].lastname, "avatar_url": createdBy[0].avatar_url},
			image:doc.data.image,
			item: doc.data.item
		}}`,
		shared_query,
		bookmarked_query,
		common.DbContentTable,
		filter_shared_query,
		query,
		common.QueryClean(filter_query),
		sort_query, limit_query,
		common.DbUserTable,
	)

	resp, err := runQuery(ctx, q, common.DbGoalTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if res, err := recordToShareableContent(r); err == nil {
			response = append(response, res)
		}
	}
	return response, nil
}

func ReadTaxonomyByNameslug(ctx context.Context, name_slug string) (*static_proto.Taxonomy, error) {
	query := fmt.Sprintf(`FILTER doc.data.name_slug == "%v"`, name_slug)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		RETURN doc`, common.DbTaxonomyTable, query)

	resp, err := runQuery(ctx, q, common.DbTaxonomyTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToTaxonomy(resp.Records[0])
	return data, err
}
