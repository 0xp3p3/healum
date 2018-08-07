package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	restful "github.com/emicklei/go-restful"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/micro/go-micro/errors"
	log "github.com/sirupsen/logrus"
)

type Error struct {
	Domain string `json:"domain"`
	Reason string `json:"reason"`
}

type ErrResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Errors  []*Error `json:"errors"`
}

// WriteErrorResponseWithCode responses with status code
func WriteErrorResponseWithCode(rsp *restful.Response, err error, domain, reason string) {
	rsp.AddHeader("Content-Type", "application/json")
	rsp_status := int(errors.Parse(err.Error()).Code)
	rsp.WriteHeaderAndEntity(rsp_status, &ErrResponse{
		Code:    rsp_status,
		Message: reason,
		Errors: []*Error{{
			Domain: domain,
			Reason: err.Error(),
		}},
	})
}

func WriteErrorResponse(rsp *restful.Response, err error, domain, reason string) {
	rsp.AddHeader("Content-Type", "application/json")
	//FIXME:fix response status code here
	//rsp_status := int(errors.Parse(err.Error()).Code)
	//rsp.WriteHeaderAndEntity(rsp_status, &ErrResponse{
	rsp.WriteHeaderAndEntity(http.StatusBadRequest, &ErrResponse{
		Code:    http.StatusBadRequest,
		Message: reason,
		Errors: []*Error{
			{
				Domain: domain,
				Reason: err.Error(),
			},
		},
	})
}

func NoAuthorizedResponse(rsp *restful.Response, err error, domain, reason string) {
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusUnauthorized, &ErrResponse{
		Code:    http.StatusUnauthorized,
		Message: reason,
		Errors: []*Error{
			{
				Domain: domain,
				Reason: err.Error(),
			},
		},
	})
}

func UnmarshalAny(req *restful.Request, rsp *restful.Response, obj proto.Message) error {
	// getting json object
	req_obj := new(map[string]interface{})
	log.Debug("req: ", req)
	err := req.ReadEntity(req_obj)
	if err != nil {
		log.Debug("ReadEntity: ", err)
		return err
	}
	// getting json string from json object
	b, err := json.Marshal(req_obj)
	if err != nil {
		log.Debug("Marshal: ", err)
		return err
	}
	// getting Any object from json string
	if err := jsonpb.Unmarshal(strings.NewReader(string(b)), obj); err != nil {
		log.Debug("Unmarshal: ", err)
		return err
	}

	return nil
}

func MarshalAny(rsp *restful.Response, res proto.Message) interface{} {
	// getting json string from Any object
	marshaler := jsonpb.Marshaler{}
	marshaler.OrigName = true
	js, err := marshaler.MarshalToString(res)
	if err != nil {
		WriteErrorResponse(rsp, err, "go.micro.srv.response.Bind", "MarshalError")
	}

	// getting json object from json string
	var data map[string]interface{}
	decoder := json.NewDecoder(bytes.NewReader([]byte(js)))
	err = decoder.Decode(&data)
	if err != nil {
		WriteErrorResponse(rsp, err, "go.micro.srv.response.Bind", "MarshalError")
	}
	data["code"], err = strconv.Atoi(data["code"].(string))
	if err != nil {
		WriteErrorResponse(rsp, err, "go.micro.srv.response.Bind", "MarshalError")
	}
	return data
}
