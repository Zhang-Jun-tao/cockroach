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
	"math/rand"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/context"

	"github.com/pkg/errors"

	"github.com/cockroachdb/cockroach/pkg/acceptance/cluster"
	"github.com/cockroachdb/cockroach/pkg/gossip"
	"github.com/cockroachdb/cockroach/pkg/internal/client"
	"github.com/cockroachdb/cockroach/pkg/util"
	"github.com/cockroachdb/cockroach/pkg/util/httputil"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/timeutil"
)

const longWaitTime = 2 * time.Minute
const shortWaitTime = 20 * time.Second

type checkGossipFunc func(map[string]gossip.Info) error

// checkGossip fetches the gossip infoStore from each node and invokes the given
// function. The test passes if the function returns 0 for every node,
// retrying for up to the given duration.
func checkGossip(
	ctx context.Context, t *testing.T, c cluster.Cluster, d time.Duration, f checkGossipFunc,
) {
	err := util.RetryForDuration(d, func() error {
		select {
		case <-stopper.ShouldStop():
			t.Fatalf("interrupted")
			return nil
		case <-time.After(1 * time.Second):
		}

		var infoStatus gossip.InfoStatus
		for i := 0; i < c.NumNodes(); i++ {
			if err := httputil.GetJSON(cluster.HTTPClient, c.URL(ctx, i)+"/_status/gossip/local", &infoStatus); err != nil {
				return errors.Wrapf(err, "failed to get gossip status from node %d", i)
			}
			if err := f(infoStatus.Infos); err != nil {
				return errors.Errorf("node %d: %s", i, err)
			}
		}

		return nil
	})
	if err != nil {
		t.Fatal(errors.Errorf("condition failed to evaluate within %s: %s", d, err))
	}
}

// hasPeers returns a checkGossipFunc that passes when the given
// number of peers are connected via gossip.
func hasPeers(expected int) checkGossipFunc {
	return func(infos map[string]gossip.Info) error {
		count := 0
		for k := range infos {
			if strings.HasPrefix(k, "node:") {
				count++
			}
		}
		if count != expected {
			return errors.Errorf("expected %d peers, found %d", expected, count)
		}
		return nil
	}
}

// hasSentinel is a checkGossipFunc that passes when the sentinel gossip is present.
func hasSentinel(infos map[string]gossip.Info) error {
	if _, ok := infos[gossip.KeySentinel]; !ok {
		return errors.Errorf("sentinel not found")
	}
	return nil
}

// hasClusterID is a checkGossipFunc that passes when the cluster ID gossip is present.
func hasClusterID(infos map[string]gossip.Info) error {
	if _, ok := infos[gossip.KeyClusterID]; !ok {
		return errors.Errorf("cluster ID not found")
	}
	return nil
}

func TestGossipPeerings(t *testing.T) {
	runTestOnConfigs(t, testGossipPeeringsInner)
}

func testGossipPeeringsInner(
	ctx context.Context, t *testing.T, c cluster.Cluster, cfg cluster.TestConfig,
) {
	num := c.NumNodes()

	deadline := timeutil.Now().Add(cfg.Duration)

	waitTime := longWaitTime
	if cfg.Duration < waitTime {
		waitTime = shortWaitTime
	}

	for timeutil.Now().Before(deadline) {
		checkGossip(ctx, t, c, waitTime, hasPeers(num))

		// Restart the first node.
		log.Infof(ctx, "restarting node 0")
		if err := c.Restart(ctx, 0); err != nil {
			t.Fatal(err)
		}
		checkGossip(ctx, t, c, waitTime, hasPeers(num))

		// Restart another node (if there is one).
		var pickedNode int
		if num > 1 {
			pickedNode = rand.Intn(num-1) + 1
		}
		log.Infof(ctx, "restarting node %d", pickedNode)
		if err := c.Restart(ctx, pickedNode); err != nil {
			t.Fatal(err)
		}
		checkGossip(ctx, t, c, waitTime, hasPeers(num))
	}
}

// TestGossipRestart verifies that the gossip network can be
// re-bootstrapped after a time when all nodes were down
// simultaneously.
func TestGossipRestart(t *testing.T) {
	// TODO(bram): #4559 Limit this test to only the relevant cases. No chaos
	// agents should be required.
	runTestOnConfigs(t, testGossipRestartInner)
}

func testGossipRestartInner(
	ctx context.Context, t *testing.T, c cluster.Cluster, cfg cluster.TestConfig,
) {
	// This already replicates the first range (in the local setup).
	// The replication of the first range is important: as long as the
	// first range only exists on one node, that node can trivially
	// acquire the range lease. Once the range is replicated, however,
	// nodes must be able to discover each other over gossip before the
	// lease can be acquired.
	num := c.NumNodes()

	deadline := timeutil.Now().Add(cfg.Duration)

	waitTime := longWaitTime
	if cfg.Duration < waitTime {
		waitTime = shortWaitTime
	}

	for timeutil.Now().Before(deadline) {
		log.Infof(ctx, "waiting for initial gossip connections")
		checkGossip(ctx, t, c, waitTime, hasPeers(num))
		checkGossip(ctx, t, c, waitTime, hasClusterID)
		checkGossip(ctx, t, c, waitTime, hasSentinel)

		log.Infof(ctx, "killing all nodes")
		for i := 0; i < num; i++ {
			if err := c.Kill(ctx, i); err != nil {
				t.Fatal(err)
			}
		}

		log.Infof(ctx, "restarting all nodes")
		for i := 0; i < num; i++ {
			if err := c.Restart(ctx, i); err != nil {
				t.Fatal(err)
			}
		}

		log.Infof(ctx, "waiting for gossip to be connected")
		checkGossip(ctx, t, c, waitTime, hasPeers(num))
		checkGossip(ctx, t, c, waitTime, hasClusterID)
		checkGossip(ctx, t, c, waitTime, hasSentinel)

		for i := 0; i < num; i++ {
			db := c.NewClient(ctx, t, i)
			if i == 0 {
				if err := db.Del(ctx, "count"); err != nil {
					t.Fatal(err)
				}
			}
			var kv client.KeyValue
			if err := db.Txn(ctx, func(txn *client.Txn) error {
				var err error
				kv, err = txn.Inc("count", 1)
				return err
			}); err != nil {
				t.Fatal(err)
			} else if v := kv.ValueInt(); v != int64(i+1) {
				t.Fatalf("unexpected value %d for write #%d (expected %d)", v, i, i+1)
			}
		}
	}
}
