// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package receivercreator // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/receivercreator"

import (
	"go.opentelemetry.io/collector/component"

	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer"
)

// resolvedReceiver captures a running sub-receiver along with the resolved state
// that was used to create it, enabling change detection on endpoint updates.
type resolvedReceiver struct {
	component      component.Component
	templateKey    string        // key into config.receiverTemplates ("" for hints-based)
	resolvedCfg    userConfigMap // result of expandConfig(template.config, env)
	endpointTarget string        // endpoint.Target
}

// receiverMap is a multimap for mapping one id to many receivers. It does
// not deduplicate the same value being associated with the same key.
type receiverMap map[observer.EndpointID][]resolvedReceiver

// Put rcvr into key id. If rcvr is a duplicate it will still be added.
func (rm receiverMap) Put(id observer.EndpointID, rcvr resolvedReceiver) {
	rm[id] = append(rm[id], rcvr)
}

// Get receivers by id.
func (rm receiverMap) Get(id observer.EndpointID) []resolvedReceiver {
	return rm[id]
}

// Remove all receivers by id.
func (rm receiverMap) RemoveAll(id observer.EndpointID) {
	delete(rm, id)
}

// Get all receivers in the map.
func (rm receiverMap) Values() (out []resolvedReceiver) {
	for _, m := range rm {
		out = append(out, m...)
	}
	return out
}

// Size is the number of total receivers in the map.
func (rm receiverMap) Size() (out int) {
	for _, m := range rm {
		out += len(m)
	}
	return out
}
