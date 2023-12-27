package pod

/*
import (
	"context"
	"log"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
)

const kubeConfig = "<path to kubeconfig file>"
const testNS = "test-ns"

// var k8s kubernetes.Interface
var k8sClient client.Client
var ctx context.Context
var testEnv *envtest.Environment
var scheme *runtime.Scheme

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pod Test")
}

var _ = BeforeSuite(func() {
	log.Println("test")

	var err error

	// create context
	ctx = context.TODO()

	// scheme
	scheme = runtime.NewScheme()
	err = clientgoscheme.AddToScheme(scheme)
	Expect(err).NotTo(HaveOccurred())

	// testEnv
	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{ }, //filepath.Join("..", "..", "config", "crd", "bases")
		ErrorIfCRDPathMissing: true,
	}
	cfg, err := testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	// create k8s client
	//config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	//Expect(err).NotTo(HaveOccurred())
	//clientset, err := kubernetes.NewForConfig(cfg)
	//Expect(err).NotTo(HaveOccurred())
	//k8s = clientset
	k8sClient, err = client.New(cfg, client.Options{})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	// create test namespace
	err = k8sClient.Create(ctx, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: testNS,
		},
	})
	Expect(err).NotTo(HaveOccurred())
})

var _ = Describe("Test", func() {
	BeforeEach(func() {
		var err error
		err = k8sClient.DeleteAllOf(ctx, &corev1.ConfigMap{}, client.InNamespace("test"))
		Expect(err).NotTo(HaveOccurred())
		err = k8sClient.DeleteAllOf(ctx, &appsv1.Deployment{}, client.InNamespace("test"))
		Expect(err).NotTo(HaveOccurred())
		svcs := &corev1.ServiceList{}
		err = k8sClient.List(ctx, svcs, client.InNamespace("test"))
		Expect(err).NotTo(HaveOccurred())
		for _, svc := range svcs.Items {
			err := k8sClient.Delete(ctx, &svc)
			Expect(err).NotTo(HaveOccurred())
		}
		time.Sleep(100 * time.Millisecond)
	})

	Describe("The container names in the same pod can be obtained", func() {
		Context("when the number of container is only one", func() {
			It("can obtain the pod name", func() {
				By("create test pod", func() {
					Expect(
						k8sClient.Create(ctx, &corev1.Pod{
							ObjectMeta: metav1.ObjectMeta{
								Name:      "test-pod",
								Namespace: testNS,
							},
							Spec: corev1.PodSpec{
								Containers: []corev1.Container{
									{
										Image: "test",
										Name:  "test-container-1",
									},
								},
							},
						}),
					).NotTo(HaveOccurred())

					Eventually(func() string {
						p := &corev1.Pod{}
						k8sClient.Get(ctx, client.ObjectKey{Namespace: testNS, Name: "test-pod"}, p)

						return p.Name
					}).WithTimeout(time.Second * 5).Should(Equal("test-pod"))
				})

				By("Get Container Names")
				log.Println("test")
				//Expect("error").Should(Equal("fail"))
			})
		})
	})
})
*/
