/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package deprovisioning

import (
	"context"
	"time"

	v1 "k8s.io/api/core/v1"

	deprovisioningevents "github.com/aws/karpenter-core/pkg/controllers/deprovisioning/events"
	"github.com/aws/karpenter-core/pkg/events"
	"github.com/aws/karpenter-core/pkg/utils/pretty"
)

// Reporter is used to periodically report node statuses regarding deprovisioning. This gives observers awareness of why
// deprovisioning of a particular node isn't occurring.
type Reporter struct {
	cm       *pretty.ChangeMonitor
	recorder events.Recorder
}

func NewReporter(recorder events.Recorder) *Reporter {
	// This change monitor is used by the deprovisioning reporter to report why nodes can't be deprovisioned
	// periodically.  The reporter can be called as often as is convenient and it will prevent these notifications from
	// flooding events.
	cm := pretty.NewChangeMonitor()
	cm.Reconfigure(15 * time.Minute)

	return &Reporter{
		recorder: recorder,
		cm:       cm,
	}
}

// RecordUnconsolidatableReason is used to periodically report why a node is unconsolidatable to it can be logged
func (r *Reporter) RecordUnconsolidatableReason(ctx context.Context, node *v1.Node, reason string) {
	if r.cm.HasChanged(string(node.UID), "consolidation") {
		r.recorder.Publish(deprovisioningevents.UnconsolidatableReason(node, reason))
	}
}
