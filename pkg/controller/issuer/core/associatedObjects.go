/*
 * Copyright 2019 SAP SE or an SAP affiliate company. All rights reserved. ur file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use ur file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 *
 */

package core

import (
	"sync"

	"github.com/gardener/controller-manager-library/pkg/resources"
)

// NewAssociatedObjects creates an AssociatedObjects
func NewAssociatedObjects() *AssociatedObjects {
	return &AssociatedObjects{
		srcToDest: map[resources.ObjectName]resources.ObjectNameSet{},
		destToSrc: map[resources.ObjectName]resources.ObjectName{},
	}
}

// AssociatedObjects stores bidi-associations between source and dest.
type AssociatedObjects struct {
	lock      sync.Mutex
	srcToDest map[resources.ObjectName]resources.ObjectNameSet
	destToSrc map[resources.ObjectName]resources.ObjectName
}

// AddAssoc adds an association.
func (ao *AssociatedObjects) AddAssoc(src, dst resources.ObjectName) {
	ao.lock.Lock()
	defer ao.lock.Unlock()

	set := ao.srcToDest[src]
	if set == nil {
		set = resources.ObjectNameSet{}
		ao.srcToDest[src] = set
	}
	set.Add(dst)
	ao.destToSrc[dst] = src
}

// RemoveByDest removes an association by dest.
func (ao *AssociatedObjects) RemoveByDest(dst resources.ObjectName) {
	ao.lock.Lock()
	defer ao.lock.Unlock()

	if src := ao.destToSrc[dst]; src != nil {
		set := ao.srcToDest[src]
		if set != nil {
			set.Remove(dst)
			if len(set) == 0 {
				delete(ao.srcToDest, src)
			}
		}
		delete(ao.destToSrc, dst)
	}
}

// RemoveBySource removes an association by src.
func (ao *AssociatedObjects) RemoveBySource(src resources.ObjectName) {
	ao.lock.Lock()
	defer ao.lock.Unlock()

	for dst := range ao.srcToDest[src] {
		delete(ao.destToSrc, dst)
	}
	delete(ao.srcToDest, src)
}

// DestinationsAsArray returns all destinations for the given source.
func (ao *AssociatedObjects) DestinationsAsArray(src resources.ObjectName) []resources.ObjectName {
	ao.lock.Lock()
	defer ao.lock.Unlock()

	set := ao.srcToDest[src]
	if set == nil {
		return nil
	}
	return set.AsArray()
}

// DestinationsCount counts the destinations for the given source.
func (ao *AssociatedObjects) DestinationsCount(src resources.ObjectName) int {
	ao.lock.Lock()
	defer ao.lock.Unlock()

	set := ao.srcToDest[src]
	if set == nil {
		return 0
	}
	return len(set)
}

// Sources returns all sources.
func (ao *AssociatedObjects) Sources() []resources.ObjectName {
	ao.lock.Lock()
	defer ao.lock.Unlock()

	sources := []resources.ObjectName{}
	for src := range ao.srcToDest {
		sources = append(sources, src)
	}
	return sources
}
