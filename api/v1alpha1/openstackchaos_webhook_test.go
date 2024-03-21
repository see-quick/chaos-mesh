package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("openstackchaos_webhook", func() {
	Context("webhook.Validator of openstackchaos", func() {
		It("Validate", func() {

			type TestCase struct {
				name    string
				chaos   OpenStackChaos
				execute func(chaos *OpenStackChaos) error
				expect  string
			}
			testVMID := "testVMID"
			tcs := []TestCase{
				{
					name: "simple ValidateCreate for VmRestart",
					chaos: OpenStackChaos{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: metav1.NamespaceDefault,
							Name:      "foo1",
						},
						Spec: OpenStackChaosSpec{
							Action: VmRestart,
						},
					},
					execute: func(chaos *OpenStackChaos) error {
						_, err := chaos.ValidateCreate()
						return err
					},
					expect: "error",
				},
				{
					name: "unknown action",
					chaos: OpenStackChaos{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: metav1.NamespaceDefault,
							Name:      "foo6",
						},
						Spec: OpenStackChaosSpec{},
					},
					execute: func(chaos *OpenStackChaos) error {
						_, err := chaos.ValidateCreate()
						return err
					},
					expect: "error",
				},
				{
					name: "validate the VmRestart with VMID",
					chaos: OpenStackChaos{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: metav1.NamespaceDefault,
							Name:      "foo7",
						},
						Spec: OpenStackChaosSpec{
							Action: VmRestart,
							OpenStackSelector: OpenStackSelector{
								VMID: testVMID,
							},
						},
					},
					execute: func(chaos *OpenStackChaos) error {
						_, err := chaos.ValidateCreate()
						return err
					},
					expect: "no error",
				},
			}

			for _, tc := range tcs {
				err := tc.execute(&tc.chaos)
				if tc.expect == "error" {
					Expect(err).To(HaveOccurred())
				} else {
					Expect(err).NotTo(HaveOccurred())
				}
			}
		})
	})
})
