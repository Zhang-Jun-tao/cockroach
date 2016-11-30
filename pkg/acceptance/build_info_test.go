// Copyright 2015 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.
//
// Author: Peter Mattis (peter@cockroachlabs.com)

package acceptance

import (
	"testing"
	"time"

	"golang.org/x/net/context"

	"github.com/cockroachdb/cockroach/pkg/acceptance/cluster"
	"github.com/cockroachdb/cockroach/pkg/server/serverpb"
	"github.com/cockroachdb/cockroach/pkg/util"
	"github.com/cockroachdb/cockroach/pkg/util/httputil"
)

func TestBuildInfo(t *testing.T) {
	runTestOnConfigs(t, testBuildInfoInner)
}

func testBuildInfoInner(
	ctx context.Context, t *testing.T, c cluster.Cluster, cfg cluster.TestConfig,
) {
	checkGossip(ctx, t, c, 20*time.Second, hasPeers(c.NumNodes()))

	var details serverpb.DetailsResponse
	util.SucceedsSoon(t, func() error {
		select {
		case <-stopper.ShouldStop():
			t.Fatalf("interrupted")
		default:
		}
		return httputil.GetJSON(cluster.HTTPClient, c.URL(ctx, 0)+"/_status/details/local", &details)
	})

	bi := details.BuildInfo
	testData := map[string]string{
		"go_version":   bi.GoVersion,
		"tag":          bi.Tag,
		"time":         bi.Time,
		"dependencies": bi.Dependencies,
	}
	for key, val := range testData {
		if val == "" {
			t.Errorf("build info not set for \"%s\"", key)
		}
	}
}
