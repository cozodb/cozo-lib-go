/*
 * Copyright 2022, The Cozo Project Authors.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
 * If a copy of the MPL was not distributed with this file,
 * You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cozo

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/objx"
	"unsafe"
)

/*
#include <string.h>
#include "cozo_c.h"
#cgo LDFLAGS: -lcozo_c -lstdc++ -lm
#cgo windows LDFLAGS: -lbcrypt -lwsock32 -lws2_32 -lshlwapi -lrpcrt4
#cgo darwin LDFLAGS: -framework Security
*/
import "C"

type CozoDB struct {
	Id C.int32_t
}

var emptyMap = C.CString("{}")

type Map = objx.Map

type QueryError struct {
	Data Map
}

type NamedRows struct {
	Headers []string `json:"headers"`
	Rows    [][]any  `json:"rows"`
	Took    float64  `json:"took"`
	Ok      bool     `json:"ok"`
}

func (m QueryError) Error() string {
	msg := ""
	msg = m.Data.Get("display").String()
	if len(msg) > 0 {
		return msg
	}

	msg = m.Data.Get("message").String()
	if len(msg) > 0 {
		return msg
	}
	return "Unknown error"
}

func New(engine string, path string, options Map) (CozoDB, error) {
	var ret CozoDB
	cEngine := C.CString(engine)
	defer C.free(unsafe.Pointer(cEngine))

	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	var cOptions *C.char
	if options == nil {
		cOptions = emptyMap
	} else {
		b, err := json.Marshal(options)
		cOptions = C.CString(string(b))
		defer C.free(unsafe.Pointer(cOptions))
		if err != nil {
			return ret, err
		}
	}

	err := C.cozo_open_db(cEngine, cPath, cOptions, &ret.Id)
	if err != nil {
		goErr := C.GoString(err)
		C.cozo_free_str(err)
		return ret, errors.New(goErr)
	}
	return ret, nil
}

func (db *CozoDB) Close() {
	C.cozo_close_db(db.Id)
}

func (db *CozoDB) Run(query string, params Map) (NamedRows, error) {
	var paramsStr *C.char
	result := NamedRows{
		Ok:   false,
		Took: -1.0,
	}

	if params == nil {
		paramsStr = emptyMap
	} else {
		b, err := json.Marshal(params)
		paramsStr = C.CString(string(b))
		defer C.free(unsafe.Pointer(paramsStr))

		if err != nil {
			return result, err
		}
	}

	queryStr := C.CString(query)
	defer C.free(unsafe.Pointer(queryStr))

	res := C.cozo_run_query(db.Id, queryStr, paramsStr)
	defer C.cozo_free_str(res)

	cLen := C.int(C.strlen(res))
	gBytes := C.GoBytes(unsafe.Pointer(res), cLen)

	jsonErr := json.Unmarshal(gBytes, &result)
	if jsonErr != nil {
		return result, jsonErr
	}

	if result.Ok {
		return result, nil
	} else {
		var raw Map

		jsonErr := json.Unmarshal(gBytes, &raw)
		if jsonErr != nil {
			return result, jsonErr
		}

		return result, QueryError{
			Data: raw,
		}
	}
}

func (db *CozoDB) ImportRelations(payload Map) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	cPayload := C.CString(string(b))
	defer C.free(unsafe.Pointer(cPayload))

	var res Map
	cRes := C.cozo_import_relations(db.Id, cPayload)
	defer C.cozo_free_str(cRes)
	cLen := C.int(C.strlen(cRes))
	gBytes := C.GoBytes(unsafe.Pointer(cRes), cLen)

	jsonErr := json.Unmarshal(gBytes, &res)
	if jsonErr != nil {
		return jsonErr
	}
	if !res.Get("ok").Bool(false) {
		return QueryError{
			Data: res,
		}
	} else {
		return nil
	}
}

func (db *CozoDB) ExportRelations(relations []string) (Map, error) {
	payloadMap := map[string][]string{
		"relations": relations,
	}

	b, err := json.Marshal(payloadMap)
	if err != nil {
		return nil, err
	}

	payload := C.CString(string(b))
	defer C.free(unsafe.Pointer(payload))

	var res Map
	cRes := C.cozo_export_relations(db.Id, payload)
	defer C.cozo_free_str(cRes)
	cLen := C.int(C.strlen(cRes))
	gBytes := C.GoBytes(unsafe.Pointer(cRes), cLen)

	jsonErr := json.Unmarshal(gBytes, &res)
	if jsonErr != nil {
		return nil, jsonErr
	}
	if !res.Get("ok").Bool(false) {
		return nil, QueryError{
			Data: res,
		}
	} else {
		return res.Get("data").MustObjxMap(), nil
	}
}

func (db *CozoDB) Backup(path string) error {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	cRes := C.cozo_backup(db.Id, cPath)
	defer C.cozo_free_str(cRes)

	cLen := C.int(C.strlen(cRes))
	gBytes := C.GoBytes(unsafe.Pointer(cRes), cLen)

	var res Map
	jsonErr := json.Unmarshal(gBytes, &res)
	if jsonErr != nil {
		return jsonErr
	}
	if !res.Get("ok").Bool(false) {
		return QueryError{
			Data: res,
		}
	} else {
		return nil
	}
}

func (db *CozoDB) Restore(path string) error {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	cRes := C.cozo_restore(db.Id, cPath)
	defer C.cozo_free_str(cRes)

	cLen := C.int(C.strlen(cRes))
	gBytes := C.GoBytes(unsafe.Pointer(cRes), cLen)

	var res Map
	jsonErr := json.Unmarshal(gBytes, &res)
	if jsonErr != nil {
		return jsonErr
	}
	if !res.Get("ok").Bool(false) {
		return QueryError{
			Data: res,
		}
	} else {
		return nil
	}
}

func (db *CozoDB) ImportRelationsFromBackup(path string, relations []string) error {
	payloadMap := map[string]any{
		"relations": relations,
		"path":      path,
	}

	b, err := json.Marshal(payloadMap)
	if err != nil {
		return err
	}

	payload := C.CString(string(b))
	defer C.free(unsafe.Pointer(payload))

	var res Map
	cRes := C.cozo_import_from_backup(db.Id, payload)
	defer C.cozo_free_str(cRes)
	cLen := C.int(C.strlen(cRes))
	gBytes := C.GoBytes(unsafe.Pointer(cRes), cLen)

	jsonErr := json.Unmarshal(gBytes, &res)
	if jsonErr != nil {
		return jsonErr
	}
	if !res.Get("ok").Bool(false) {
		return QueryError{
			Data: res,
		}
	} else {
		return nil
	}
}
