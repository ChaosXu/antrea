// Copyright 2021 Antrea Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package networkpolicy

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"

	"antrea.io/antrea/pkg/apis/controlplane"
	crdv1alpha1 "antrea.io/antrea/pkg/apis/crd/v1alpha1"
	crdv1alpha3 "antrea.io/antrea/pkg/apis/crd/v1alpha3"
	"antrea.io/antrea/pkg/controller/networkpolicy/store"
	antreatypes "antrea.io/antrea/pkg/controller/types"
)

// addClusterGroup is responsible for processing the ADD event of a ClusterGroup resource.
func (c *NetworkPolicyController) addClusterGroup(curObj interface{}) {
	cg := curObj.(*crdv1alpha3.ClusterGroup)
	key := internalGroupKeyFunc(cg)
	klog.V(2).Infof("Processing ADD event for ClusterGroup %s", cg.Name)
	newGroup := c.processClusterGroup(cg)
	klog.V(2).Infof("Creating new internal Group %s", newGroup.UID)
	c.internalGroupStore.Create(newGroup)
	c.enqueueInternalGroup(key)
}

// updateClusterGroup is responsible for processing the UPDATE event of a ClusterGroup resource.
func (c *NetworkPolicyController) updateClusterGroup(oldObj, curObj interface{}) {
	cg := curObj.(*crdv1alpha3.ClusterGroup)
	og := oldObj.(*crdv1alpha3.ClusterGroup)
	key := internalGroupKeyFunc(cg)
	klog.V(2).Infof("Processing UPDATE event for ClusterGroup %s", cg.Name)
	newGroup := c.processClusterGroup(cg)
	oldGroup := c.processClusterGroup(og)

	selectorUpdated := func() bool {
		return getNormalizedNameForSelector(newGroup.Selector) != getNormalizedNameForSelector(oldGroup.Selector)
	}
	svcRefUpdated := func() bool {
		oldSvc, newSvc := oldGroup.ServiceReference, newGroup.ServiceReference
		if oldSvc != nil && newSvc != nil && oldSvc.Name == newSvc.Name && oldSvc.Namespace == newSvc.Namespace {
			return false
		} else if oldSvc == nil && newSvc == nil {
			return false
		}
		return true
	}
	ipBlocksUpdated := func() bool {
		oldIPBs, newIPBs := sets.String{}, sets.String{}
		for _, ipb := range oldGroup.IPBlocks {
			oldIPBs.Insert(ipb.CIDR.String())
		}
		for _, ipb := range newGroup.IPBlocks {
			newIPBs.Insert(ipb.CIDR.String())
		}
		return oldIPBs.Equal(newIPBs)
	}
	childGroupsUpdated := func() bool {
		oldChildGroups, newChildGroups := sets.String{}, sets.String{}
		for _, c := range oldGroup.ChildGroups {
			oldChildGroups.Insert(c)
		}
		for _, c := range newGroup.ChildGroups {
			newChildGroups.Insert(c)
		}
		return !oldChildGroups.Equal(newChildGroups)
	}
	if !ipBlocksUpdated() && !svcRefUpdated() && !selectorUpdated() && !childGroupsUpdated() {
		// No change in the contents of the ClusterGroup. No need to enqueue for further sync.
		return
	}
	c.internalGroupStore.Update(newGroup)
	c.enqueueInternalGroup(key)
}

// deleteClusterGroup is responsible for processing the DELETE event of a ClusterGroup resource.
func (c *NetworkPolicyController) deleteClusterGroup(oldObj interface{}) {
	og, ok := oldObj.(*crdv1alpha3.ClusterGroup)
	klog.V(2).Infof("Processing DELETE event for ClusterGroup %s", og.Name)
	if !ok {
		tombstone, ok := oldObj.(cache.DeletedFinalStateUnknown)
		if !ok {
			klog.Errorf("Error decoding object when deleting ClusterGroup, invalid type: %v", oldObj)
			return
		}
		og, ok = tombstone.Obj.(*crdv1alpha3.ClusterGroup)
		if !ok {
			klog.Errorf("Error decoding object tombstone when deleting ClusterGroup, invalid type: %v", tombstone.Obj)
			return
		}
	}
	key := internalGroupKeyFunc(og)
	klog.V(2).Infof("Deleting internal Group %s", key)
	err := c.internalGroupStore.Delete(key)
	if err != nil {
		klog.Errorf("Unable to delete internal Group %s from store: %v", key, err)
	}
	c.enqueueInternalGroup(key)
}

func (c *NetworkPolicyController) processClusterGroup(cg *crdv1alpha3.ClusterGroup) *antreatypes.Group {
	internalGroup := antreatypes.Group{
		Name: cg.Name,
		UID:  cg.UID,
	}
	if len(cg.Spec.ChildGroups) > 0 {
		for _, childCGName := range cg.Spec.ChildGroups {
			internalGroup.ChildGroups = append(internalGroup.ChildGroups, string(childCGName))
		}
		return &internalGroup
	}
	if len(cg.Spec.IPBlocks) > 0 {
		for i := range cg.Spec.IPBlocks {
			ipb, _ := toAntreaIPBlockForCRD(&cg.Spec.IPBlocks[i])
			internalGroup.IPBlocks = append(internalGroup.IPBlocks, *ipb)
		}
		return &internalGroup
	}
	svcSelector := cg.Spec.ServiceReference
	if svcSelector != nil {
		// ServiceReference will be converted to groupSelector once the internalGroup is synced.
		internalGroup.ServiceReference = &controlplane.ServiceReference{
			Namespace: svcSelector.Namespace,
			Name:      svcSelector.Name,
		}
	} else {
		groupSelector := toGroupSelector("", cg.Spec.PodSelector, cg.Spec.NamespaceSelector, cg.Spec.ExternalEntitySelector)
		internalGroup.Selector = groupSelector
	}
	return &internalGroup
}

// filterInternalGroupsForService computes a list of internal Group keys which references the Service.
func (c *NetworkPolicyController) filterInternalGroupsForService(obj metav1.Object) sets.String {
	matchingKeySet := sets.String{}
	indexKey, _ := cache.MetaNamespaceKeyFunc(obj)
	matchedSvcGroups, _ := c.internalGroupStore.GetByIndex(store.ServiceIndex, indexKey)
	for i := range matchedSvcGroups {
		key, _ := store.GroupKeyFunc(matchedSvcGroups[i])
		matchingKeySet.Insert(key)
	}
	return matchingKeySet
}

func (c *NetworkPolicyController) enqueueInternalGroup(key string) {
	klog.V(4).Infof("Adding new key %s to internal Group queue", key)
	c.internalGroupQueue.Add(key)
}

func (c *NetworkPolicyController) internalGroupWorker() {
	for c.processNextInternalGroupWorkItem() {
	}
}

// Processes an item in the "internalGroup" work queue, by calling
// syncInternalGroup after casting the item to a string (Group key).
// If syncInternalGroup returns an error, this function handles it by re-queueing
// the item so that it can be processed again later. If syncInternalGroup is
// successful, the ClusterGroup is removed from the queue until we get notify
// of a new change. This function return false if and only if the work queue
// was shutdown (no more items will be processed).
func (c *NetworkPolicyController) processNextInternalGroupWorkItem() bool {
	key, quit := c.internalGroupQueue.Get()
	if quit {
		return false
	}
	defer c.internalGroupQueue.Done(key)

	err := c.syncInternalGroup(key.(string))
	if err != nil {
		// Put the item back in the workqueue to handle any transient errors.
		c.internalGroupQueue.AddRateLimited(key)
		klog.Errorf("Failed to sync internal Group %s: %v", key, err)
		return true
	}
	// If no error occurs we Forget this item so it does not get queued again until
	// another change happens.
	c.internalGroupQueue.Forget(key)
	return true
}

func (c *NetworkPolicyController) syncInternalGroup(key string) error {
	defer c.triggerCNPUpdates(key)
	defer c.triggerParentGroupSync(key)
	// Retrieve the internal Group corresponding to this key.
	grpObj, found, _ := c.internalGroupStore.Get(key)
	if !found {
		klog.V(2).Infof("Internal group %s not found.", key)
		c.groupingInterface.DeleteGroup(clusterGroupType, key)
		return nil
	}
	grp := grpObj.(*antreatypes.Group)
	originalMembersComputedStatus := grp.MembersComputed
	// Retrieve the ClusterGroup corresponding to this key.
	cg, err := c.cgLister.Get(grp.Name)
	if err != nil {
		klog.Infof("Didn't find the ClusterGroup %s, skip processing of internal group", grp.Name)
		return nil
	}
	selectorUpdated := c.processServiceReference(grp)
	if grp.Selector != nil {
		c.groupingInterface.AddGroup(clusterGroupType, grp.Name, grp.Selector)
	} else {
		c.groupingInterface.DeleteGroup(clusterGroupType, grp.Name)
	}

	membersComputed, membersComputedStatus := true, v1.ConditionFalse
	// Update the ClusterGroup status to Realized as Antrea has recognized the Group and
	// processed its group members. The ClusterGroup is considered realized if:
	//   1. It does not have child groups. The group members are immediately considered
	//      computed during syncInternalGroup, as the group selector is finalized.
	//   2. All its child groups are created and realized.
	if len(grp.ChildGroups) > 0 {
		for _, cgName := range grp.ChildGroups {
			cg, found, _ := c.internalGroupStore.Get(cgName)
			if !found || cg.(*antreatypes.Group).MembersComputed != v1.ConditionTrue {
				membersComputed = false
				break
			}
		}
	}
	if membersComputed {
		klog.V(4).Infof("Updating GroupMembersComputed Status for group %s", cg.Name)
		err = c.updateGroupStatus(cg, v1.ConditionTrue)
		if err != nil {
			klog.Errorf("Failed to update ClusterGroup %s GroupMembersComputed condition to %s: %v", cg.Name, v1.ConditionTrue, err)
		} else {
			membersComputedStatus = v1.ConditionTrue
		}
	}
	if selectorUpdated || membersComputedStatus != originalMembersComputedStatus {
		// Update the internal Group object in the store with the new selector and status.
		updatedGrp := &antreatypes.Group{
			UID:              grp.UID,
			Name:             grp.Name,
			MembersComputed:  membersComputedStatus,
			Selector:         grp.Selector,
			IPBlocks:         grp.IPBlocks,
			ServiceReference: grp.ServiceReference,
			ChildGroups:      grp.ChildGroups,
		}
		klog.V(2).Infof("Updating existing internal Group %s", key)
		c.internalGroupStore.Update(updatedGrp)
	}
	return err
}

func (c *NetworkPolicyController) triggerParentGroupSync(grp string) {
	// TODO: if the max supported group nesting level increases, a Group having children
	//  will no longer be a valid indication that it cannot have parents.
	parentGroupObjs, err := c.internalGroupStore.GetByIndex(store.ChildGroupIndex, grp)
	if err != nil {
		klog.Errorf("Error retrieving parents of ClusterGroup %s: %v", grp, err)
		return
	}
	for _, p := range parentGroupObjs {
		parentGrp := p.(*antreatypes.Group)
		c.enqueueInternalGroup(parentGrp.Name)
	}
}

// triggerCNPUpdates triggers processing of ClusterNetworkPolicies associated with the input ClusterGroup.
func (c *NetworkPolicyController) triggerCNPUpdates(cg string) {
	// If a ClusterGroup is added/updated, it might have a reference in ClusterNetworkPolicy.
	cnps, err := c.cnpInformer.Informer().GetIndexer().ByIndex(ClusterGroupIndex, cg)
	if err != nil {
		klog.Errorf("Error retrieving ClusterNetworkPolicies corresponding to ClusterGroup %s", cg)
		return
	}
	for _, obj := range cnps {
		// ClusterGroup may be used by AppliedToGroup, enqueuing them after reprocessing CNP.
		c.reprocessCNP(obj.(*crdv1alpha1.ClusterNetworkPolicy), true)
	}
}

// updateGroupStatus updates the Status subresource for a ClusterGroup.
func (c *NetworkPolicyController) updateGroupStatus(cg *crdv1alpha3.ClusterGroup, cStatus v1.ConditionStatus) error {
	condStatus := crdv1alpha3.GroupCondition{
		Status: cStatus,
		Type:   crdv1alpha3.GroupMembersComputed,
	}
	if groupMembersComputedConditionEqual(cg.Status.Conditions, condStatus) {
		// There is no change in conditions.
		return nil
	}
	condStatus.LastTransitionTime = metav1.Now()
	status := crdv1alpha3.GroupStatus{
		Conditions: []crdv1alpha3.GroupCondition{condStatus},
	}
	klog.V(4).Infof("Updating ClusterGroup %s status to %#v", cg.Name, condStatus)
	toUpdate := cg.DeepCopy()
	toUpdate.Status = status
	_, err := c.crdClient.CrdV1alpha3().ClusterGroups().UpdateStatus(context.TODO(), toUpdate, metav1.UpdateOptions{})
	return err
}

// groupMembersComputedConditionEqual checks whether the condition status for GroupMembersComputed condition
// is same. Returns true if equal, otherwise returns false. It disregards the lastTransitionTime field.
func groupMembersComputedConditionEqual(conds []crdv1alpha3.GroupCondition, condition crdv1alpha3.GroupCondition) bool {
	for _, c := range conds {
		if c.Type == crdv1alpha3.GroupMembersComputed {
			if c.Status == condition.Status {
				return true
			}
		}
	}
	return false
}

// processServiceReference knows how to process the serviceReference in the group, and set the group
// selector based on the Service referenced. It returns true if the group's selector needs to be
// updated after serviceReference processing, and false otherwise.
func (c *NetworkPolicyController) processServiceReference(group *antreatypes.Group) bool {
	svcRef := group.ServiceReference
	if svcRef == nil {
		return false
	}
	originalSelectorName := getNormalizedNameForSelector(group.Selector)
	svc, err := c.serviceLister.Services(svcRef.Namespace).Get(svcRef.Name)
	if err != nil {
		klog.V(2).Infof("Error getting Service object %s/%s: %v, setting empty selector for Group %s", svcRef.Namespace, svcRef.Name, err, group.Name)
		group.Selector = nil
		return originalSelectorName == getNormalizedNameForSelector(nil)
	}
	newSelector := c.serviceToGroupSelector(svc)
	group.Selector = newSelector
	return originalSelectorName == getNormalizedNameForSelector(newSelector)
}

// serviceToGroupSelector knows how to generate GroupSelector for a Service.
func (c *NetworkPolicyController) serviceToGroupSelector(service *v1.Service) *antreatypes.GroupSelector {
	if len(service.Spec.Selector) == 0 {
		klog.Infof("Service %s/%s is without selectors and not supported by serviceReference in ClusterGroup", service.Namespace, service.Name)
		return nil
	}
	svcPodSelector := metav1.LabelSelector{
		MatchLabels: service.Spec.Selector,
	}
	// Convert Service.spec.selector to GroupSelector by setting the Namespace to the Service's Namespace
	// and podSelector to Service's selector.
	groupSelector := toGroupSelector(service.Namespace, &svcPodSelector, nil, nil)
	return groupSelector
}

// GetAssociatedGroups retrieves the internal Groups associated with the entity being
// queried (Pod or ExternalEntity identified by name and namespace).
func (c *NetworkPolicyController) GetAssociatedGroups(name, namespace string) ([]antreatypes.Group, error) {
	// Try Pod first, then ExternalEntity.
	groups, exists := c.groupingInterface.GetGroupsForPod(namespace, name)
	if !exists {
		groups, exists = c.groupingInterface.GetGroupsForExternalEntity(namespace, name)
		if !exists {
			return nil, nil
		}
	}
	clusterGroups, exists := groups[clusterGroupType]
	if !exists {
		return nil, nil
	}
	var groupObjs []antreatypes.Group
	for _, g := range clusterGroups {
		groupObjs = append(groupObjs, c.getAssociatedGroupsByName(g)...)
	}
	// Remove duplicates in the groupObj slice.
	groupKeys, j := make(map[string]bool), 0
	for _, g := range groupObjs {
		if _, exists := groupKeys[g.Name]; !exists {
			groupKeys[g.Name] = true
			groupObjs[j] = g
			j++
		}
	}
	return groupObjs[:j], nil
}

// getAssociatedGroupsByName retrieves the internal Group and all it's parent Group objects
// (if any) by Group name.
func (c *NetworkPolicyController) getAssociatedGroupsByName(grpName string) []antreatypes.Group {
	var groups []antreatypes.Group
	groupObj, found, _ := c.internalGroupStore.Get(grpName)
	if !found {
		return groups
	}
	grp := groupObj.(*antreatypes.Group)
	groups = append(groups, *grp)
	parentGroupObjs, err := c.internalGroupStore.GetByIndex(store.ChildGroupIndex, grp.Name)
	if err != nil {
		klog.Errorf("Error retrieving parents of ClusterGroup %s: %v", grp.Name, err)
	}
	for _, p := range parentGroupObjs {
		parentGrp := p.(*antreatypes.Group)
		groups = append(groups, *parentGrp)
	}
	return groups
}

// GetGroupMembers returns the current members of a ClusterGroup.
// If the ClusterGroup is defined with IPBlocks, the returned members will be []controlplane.IPBlock.
// Otherwise, the returned members will be of type controlplane.GroupMemberSet.
func (c *NetworkPolicyController) GetGroupMembers(cgName string) (controlplane.GroupMemberSet, []controlplane.IPBlock, error) {
	groupObj, found, _ := c.internalGroupStore.Get(cgName)
	if found {
		group := groupObj.(*antreatypes.Group)
		if len(group.IPBlocks) > 0 {
			return nil, group.IPBlocks, nil
		}
		return c.getClusterGroupMemberSet(group), nil, nil
	}
	return nil, nil, fmt.Errorf("no internal Group with name %s is found", cgName)
}
