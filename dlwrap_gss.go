// Copyright 2013 Apcera Inc. All rights reserved.

// +build darwin linux

package gssapi

// This file provides the wrapper functions for bouncing through GSSAPI

/*
#include <gssapi.h>

OM_uint32
wrap_goGss_one_buffer(void *wrapped_func, OM_uint32 *min, gss_buffer_t buf) {
	return ((OM_uint32(*)(OM_uint32*, gss_buffer_t))wrapped_func)(min, buf);
}

*/
import "C"

import (
	"fmt"
	"unsafe"
)

type PopulateError struct {
	Symbol string
	DlErr  error
}

func (pe *PopulateError) Error() string {
	return fmt.Sprintf("missing symbol %q: %s", pe.Symbol, pe.DlErr)
}

func (lib *GssapiLib) symbolResolveOne(
	name string,
) (
	value unsafe.Pointer,
	okay bool,
) {
	var v unsafe.Pointer
	v, lib.populateErr = lib.DlSym(name)
	if lib.populateErr != nil {
		return nil, false
	}
	return v, true
}

func (lib *GssapiLib) Populate() error {
	lib.populate.Do(func() {
		var ok bool
		if lib.fp_gss_release_buffer, ok = lib.symbolResolveOne("gss_release_buffer"); !ok {
			return
		}
		if !lib.populateNameFunctions() {
			return
		}
	})
	if lib.populateErr != nil {
		return lib.populateErr
	}
	return nil
}

func (lib *GssapiLib) gss_release_buffer(
	min *C.OM_uint32,
	buf C.gss_buffer_t,
) C.OM_uint32 {
	return C.wrap_goGss_one_buffer(lib.fp_gss_release_buffer, min, buf)
}
