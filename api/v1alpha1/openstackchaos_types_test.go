package v1alpha1

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("OpenStackChaos", func() {
	var (
		key              types.NamespacedName
		created, fetched *OpenStackChaos
	)

	BeforeEach(func() {
		// Setup steps before each test
	})

	AfterEach(func() {
		// Teardown steps after each test
	})

	Context("Create API", func() {
		It("should create an object successfully", func() {
			testInstance := "testInstance"
			testSecretName := "testSecretName"
			key = types.NamespacedName{
				Name:      "foo",
				Namespace: "default",
			}

			created = &OpenStackChaos{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
				},
				Spec: OpenStackChaosSpec{
					Action: VmRestart,
					OpenStackSelector: OpenStackSelector{
						VMID: testInstance,
					},
					SecretName: &testSecretName,
				},
			}

			By("creating an API obj")
			Expect(k8sClient.Create(context.TODO(), created)).To(Succeed())

			fetched = &OpenStackChaos{}
			Expect(k8sClient.Get(context.TODO(), key, fetched)).To(Succeed())
			Expect(fetched).To(Equal(created))

			By("deleting the created object")
			Expect(k8sClient.Delete(context.TODO(), created)).To(Succeed())
			Expect(k8sClient.Get(context.TODO(), key, created)).ToNot(Succeed())
		})
	})
})
