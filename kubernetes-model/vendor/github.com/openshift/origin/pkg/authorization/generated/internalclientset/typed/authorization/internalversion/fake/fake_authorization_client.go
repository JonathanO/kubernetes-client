/**
 * Copyright (C) 2015 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package fake

import (
	internalversion "github.com/openshift/origin/pkg/authorization/generated/internalclientset/typed/authorization/internalversion"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeAuthorization struct {
	*testing.Fake
}

func (c *FakeAuthorization) ClusterPolicies() internalversion.ClusterPolicyInterface {
	return &FakeClusterPolicies{c}
}

func (c *FakeAuthorization) ClusterPolicyBindings() internalversion.ClusterPolicyBindingInterface {
	return &FakeClusterPolicyBindings{c}
}

func (c *FakeAuthorization) ClusterRoles() internalversion.ClusterRoleInterface {
	return &FakeClusterRoles{c}
}

func (c *FakeAuthorization) ClusterRoleBindings() internalversion.ClusterRoleBindingInterface {
	return &FakeClusterRoleBindings{c}
}

func (c *FakeAuthorization) LocalResourceAccessReviews(namespace string) internalversion.LocalResourceAccessReviewInterface {
	return &FakeLocalResourceAccessReviews{c, namespace}
}

func (c *FakeAuthorization) LocalSubjectAccessReviews(namespace string) internalversion.LocalSubjectAccessReviewInterface {
	return &FakeLocalSubjectAccessReviews{c, namespace}
}

func (c *FakeAuthorization) Policies(namespace string) internalversion.PolicyInterface {
	return &FakePolicies{c, namespace}
}

func (c *FakeAuthorization) PolicyBindings(namespace string) internalversion.PolicyBindingInterface {
	return &FakePolicyBindings{c, namespace}
}

func (c *FakeAuthorization) ResourceAccessReviews() internalversion.ResourceAccessReviewInterface {
	return &FakeResourceAccessReviews{c}
}

func (c *FakeAuthorization) Roles(namespace string) internalversion.RoleInterface {
	return &FakeRoles{c, namespace}
}

func (c *FakeAuthorization) RoleBindings(namespace string) internalversion.RoleBindingInterface {
	return &FakeRoleBindings{c, namespace}
}

func (c *FakeAuthorization) RoleBindingRestrictions(namespace string) internalversion.RoleBindingRestrictionInterface {
	return &FakeRoleBindingRestrictions{c, namespace}
}

func (c *FakeAuthorization) SelfSubjectRulesReviews(namespace string) internalversion.SelfSubjectRulesReviewInterface {
	return &FakeSelfSubjectRulesReviews{c, namespace}
}

func (c *FakeAuthorization) SubjectAccessReviews() internalversion.SubjectAccessReviewInterface {
	return &FakeSubjectAccessReviews{c}
}

func (c *FakeAuthorization) SubjectRulesReviews(namespace string) internalversion.SubjectRulesReviewInterface {
	return &FakeSubjectRulesReviews{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeAuthorization) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
