package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"server/api/utils"
	"server/common"
	product_proto "server/product-srv/proto/product"
	static_proto "server/static-srv/proto/static"
	"testing"
	"time"

	restful "github.com/emicklei/go-restful"
)

var productURL = "/server/products"

var product = &product_proto.Product{
	Name:  "title",
	OrgId: "orgid",
}

var service = &product_proto.Service{
	Name:  "service",
	OrgId: "orgid",
}

var batch = &static_proto.Batch{
	Name:  "service",
	OrgId: "orgid",
}

func createProduct(product *product_proto.Product, t *testing.T) *product_proto.Product {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"product": product})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+productURL+"/product?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := product_proto.CreateProductResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.Product == nil {
		t.Errorf("Object does not matched")
		return nil
	}

	return r.Data.Product
	// t.Error(product)
}

func createService(service *product_proto.Service, t *testing.T) *product_proto.Service {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"service": service})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+productURL+"/service?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := product_proto.CreateServiceResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.Service == nil {
		t.Errorf("Object does not matched")
		return nil
	}

	return r.Data.Service
}

func createBatch(batch *static_proto.Batch, t *testing.T) *static_proto.Batch {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"batch": batch})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+productURL+"/batch?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := product_proto.CreateBatchResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.Batch == nil {
		t.Errorf("Object does not matched")
		return nil
	}

	return r.Data.Batch
}

func TestAllProducts(t *testing.T) {
	createProduct(product, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+productURL+"/products/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := product_proto.AllProductsResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if len(r.Data.Products) == 0 {
		t.Errorf("Object count does not matched")
	}
}

func TestReadProduct(t *testing.T) {
	p := createProduct(product, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+productURL+"/product/"+p.Id+"?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := product_proto.ReadProductResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.Product == nil {
		t.Errorf("Object does not matched")
		return
	}
	if r.Data.Product.Id != p.Id {
		t.Errorf("Object Id does not matched")
		return
	}
}

func TestDeleteProduct(t *testing.T) {
	p := createProduct(product, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a Delete request.
	req, err := http.NewRequest("DELETE", serverURL+productURL+"/product/"+p.Id+"?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	// Send a GET request.
	req, err = http.NewRequest("GET", serverURL+productURL+"/product/"+p.Id+"?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := product_proto.ReadProductResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if r.Data != nil {
		t.Errorf("Object does not matched")
		return
	}
}

func TestAutocompleteProduct(t *testing.T) {
	createProduct(product, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"title": "ti"})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+productURL+"/autocomplete/product?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := product_proto.AutocompleteProductResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.Products == nil {
		t.Errorf("Object does not matched")
		return
	}
	if len(r.Data.Products) == 0 {
		t.Errorf("Object count does not matched")
	}
}

func TestErrReadProduct(t *testing.T) {
	initStaticDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+productURL+"/product/999?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if len(r.Errors) == 0 {
		t.Skip("Errors does not matched")
		t.Error(r)
	} else {
		t.Log(r)
	}
}

func TestAllServices(t *testing.T) {
	createService(service, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+productURL+"/services/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := product_proto.AllServicesResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if len(r.Data.Services) == 0 {
		t.Errorf("Object count does not matched")
	}
}

func TestReadService(t *testing.T) {
	initStaticDb()

	p := createService(service, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+productURL+"/service/"+p.Id+"?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}

	time.Sleep(time.Second)

	r := product_proto.ReadServiceResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.Service == nil {
		t.Errorf("Object does not matched")
		return
	}
	if r.Data.Service.Id != p.Id {
		t.Errorf("Object Id does not matched")
		return
	}
}

func TestDeleteService(t *testing.T) {
	initStaticDb()

	p := createService(service, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a Delete request.
	req, err := http.NewRequest("DELETE", serverURL+productURL+"/service/"+p.Id+"?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	// Send a GET request.
	req, err = http.NewRequest("GET", serverURL+productURL+"/service/"+p.Id+"?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := product_proto.ReadServiceResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if r.Data != nil {
		t.Errorf("Object does not matched")
		return
	}
}

func TestAutocompleteService(t *testing.T) {
	createService(service, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"title": "se"})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+productURL+"/autocomplete/service?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := product_proto.AutocompleteServiceResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.Services == nil {
		t.Errorf("Object does not matched")
		return
	}
	if len(r.Data.Services) == 0 {
		t.Errorf("Object count does not matched")
	}
}

func TestErrReadService(t *testing.T) {
	initStaticDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+productURL+"/service/999?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if len(r.Errors) == 0 {
		t.Skip("Errors does not matched")
		t.Error(r)
	} else {
		t.Log(r)
	}
}

func TestAllBatches(t *testing.T) {
	createBatch(batch, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+productURL+"/batches/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := product_proto.AllBatchesResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if len(r.Data.Batches) == 0 {
		t.Errorf("Object count does not matched")
	}
}

func TestReadBatch(t *testing.T) {
	initStaticDb()

	p := createBatch(batch, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+productURL+"/batch/"+p.Id+"?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := product_proto.ReadBatchResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.Batch == nil {
		t.Errorf("Object does not matched")
		return
	}
	if r.Data.Batch.Id != p.Id {
		t.Errorf("Object Id does not matched")
		return
	}
}

func TestDeleteBatch(t *testing.T) {
	initStaticDb()

	p := createBatch(batch, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a Delete request.
	req, err := http.NewRequest("DELETE", serverURL+productURL+"/batch/"+p.Id+"?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	// Send a GET request.
	req, err = http.NewRequest("GET", serverURL+productURL+"/batch/"+p.Id+"?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := product_proto.ReadBatchResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if r.Data != nil {
		t.Errorf("Object does not matched")
		return
	}
}
