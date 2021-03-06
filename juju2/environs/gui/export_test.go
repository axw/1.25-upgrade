// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package gui

import (
	"github.com/juju/1.25-upgrade/juju2/environs/simplestreams"
)

var (
	StreamsVersion = streamsVersion
	DownloadType   = downloadType
)

func NewConstraint(stream string, majorVersion int) *constraint {
	return &constraint{
		LookupParams: simplestreams.LookupParams{Stream: stream},
		majorVersion: majorVersion,
	}
}
