package accesscontrol

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
)

const accessControlIdHeader = "X-AccessControlId"

var inflightRequests = map[int]wrappedRequestBodyHeader{}
var mutex sync.Mutex

func trackRequest(header wrappedRequestBodyHeader, r *http.Request) error {
	mutex.Lock()
	defer mutex.Unlock()

	identifier := rand.Int()

	r.Header.Add(accessControlIdHeader, strconv.Itoa(identifier))

	inflightRequests[identifier] = header

	return nil
}

func untrackRequest(r *http.Request) (wrappedRequestBodyHeader, error) {
	mutex.Lock()
	defer mutex.Unlock()

	identifier, err := strconv.Atoi(r.Header.Get(accessControlIdHeader))
	if err != nil {
		return wrappedRequestBodyHeader{}, err
	}

	r.Header.Del(accessControlIdHeader)

	header, ok := inflightRequests[identifier]
	if !ok {
		return wrappedRequestBodyHeader{}, fmt.Errorf("invalid identifier: %d", identifier)
	}

	return header, nil
}
