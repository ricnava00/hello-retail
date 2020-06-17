package accesscontrol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type wrappedRequestBody struct {
	Header wrappedRequestBodyHeader `json:"header"`
	Data   []byte                   `json:"data"`
}

type wrappedRequestBodyHeader struct {
	Permissions []permission `json:"permissions"`
}

func (r *wrappedRequestBody) toBody() (io.ReadCloser, int64, error) {
	body, err := json.Marshal(r)
	if err != nil {
		return nil, 0, err
	}

	return ioutil.NopCloser(bytes.NewBuffer(body)), int64(len(body)), nil
}

func (r *wrappedRequestBody) unwrapBody() (io.ReadCloser, int64) {
	return ioutil.NopCloser(bytes.NewBuffer(r.Data)), int64(len(r.Data))
}

func wrapRequestBodyWithHeader(header wrappedRequestBodyHeader, r *http.Request) error {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	wrapped := wrappedRequestBody{
		Header: header,
		Data:   data,
	}

	r.Body, r.ContentLength, err = wrapped.toBody()

	return err
}

func wrapRequestBody(config accessControlConfig, policy string, r *http.Request) error {
	grantedPermissions, err := config.getPermissionsForPolicy(policy)
	if err != nil {
		return err
	}

	header := wrappedRequestBodyHeader{
		Permissions: grantedPermissions,
	}

	return wrapRequestBodyWithHeader(header, r)
}

func unwrapRequestBody(r *http.Request) (wrappedRequestBodyHeader, error) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return wrappedRequestBodyHeader{}, err
	}

	wrapped := wrappedRequestBody{}

	err = json.Unmarshal(data, &wrapped)
	if err != nil {
		return wrappedRequestBodyHeader{}, err
	}

	r.Body, r.ContentLength = wrapped.unwrapBody()

	return wrapped.Header, nil
}

func checkPermissionsWithHeader(config accessControlConfig, function string, r *http.Request) error {
	requiredPermissions, err := config.getImmediatePermissionsForFunction(function)
	if err != nil {
		return err
	}

	header, err := getBodyHeaderFromRequest(r)
	if err != nil {
		return err
	}

	err = checkValidPermissions(header.Permissions, requiredPermissions)
	if err != nil {
		return err
	}

	return nil
}

func checkPermissionsWithGraphs(config accessControlConfig, function string, policy string, r *http.Request) error {
	absolutePermissions, err := config.getAbsolutePermissionsForFunction(function)
	if err != nil {
		return err
	}

	grantedPermissions, err := config.getPermissionsForPolicy(policy)
	if err != nil {
		return err
	}

	err = checkValidPermissions(grantedPermissions, absolutePermissions)
	if err != nil {
		return err
	}

	return nil
}

func checkValidPermissions(grantedPermissions []permission, requiredPermissions []permission) error {
	for _, requiredPermission := range requiredPermissions {
		valid := false

		for _, grantedPermission := range grantedPermissions {
			if requiredPermission == grantedPermission {
				valid = true
				break
			}
		}

		if !valid {
			return fmt.Errorf("missing permission: %v", requiredPermission)
		}
	}
	return nil
}

func getBodyHeaderFromRequest(r *http.Request) (wrappedRequestBodyHeader, error) {
	wrapped := wrappedRequestBody{}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return wrapped.Header, err
	}

	err = json.Unmarshal(data, &wrapped)
	if err != nil {
		return wrapped.Header, err
	}

	r.Body, r.ContentLength, err = wrapped.toBody()

	return wrapped.Header, nil
}
