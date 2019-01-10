package aputturkind

import (
	"context"
	"log"
	"reflect"

	examplev1alpha1 "github.com/aneeshkp/operator-example/pkg/apis/example/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Aputturkind Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileAputturkind{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("aputturkind-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Aputturkind
	err = c.Watch(&source.Kind{Type: &examplev1alpha1.Aputturkind{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Aputturkind
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &examplev1alpha1.Aputturkind{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileAputturkind{}

// ReconcileAputturkind reconciles a Aputturkind object
type ReconcileAputturkind struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Examplekind object and makes changes based on the state read
// and what is in the Examplekind.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileAputturkind) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log.Printf("Reconciling Aputturkind %s/%s\n", request.Namespace, request.Name)

	// Fetch the Examplekind instance
	instance := &examplev1alpha1.Aputturkind{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Print("Request object not found, could have been deleted after reconcile request.")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Print("Error reading the object - requeue the request.")
		return reconcile.Result{}, err
	}

	// Check if the deployment already exists, if not create a new one
	found := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep := r.newDeploymentForCR(instance)
		log.Printf("Creating a new Deployment %s/%s\n", dep.Namespace, dep.Name)
		err = r.client.Create(context.TODO(), dep)
		if err != nil {
			log.Printf("Failed to create new Deployment: %v\n", err)
			return reconcile.Result{}, err
		}
		// Deployment created successfully - return and requeue
		log.Print("Deployment created successfully - return and requeue")
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		log.Printf("Failed to get Deployment: %v\n", err)
		return reconcile.Result{}, err
	}

	// Ensure the deployment Count is the same as the spec
	count := instance.Spec.Count
	if *found.Spec.Replicas != count {
		found.Spec.Replicas = &count
		err = r.client.Update(context.TODO(), found)
		if err != nil {
			log.Printf("Failed to update Deployment: %v\n", err)
			return reconcile.Result{}, err
		}
		// Spec updated - return and requeue
		log.Print(" Spec updated - return and requeue")
		return reconcile.Result{Requeue: true}, nil
	}

	// List the pods for this deployment
	log.Print("Listing pods")
	podList := &corev1.PodList{}

	log.Printf("Listing pods size %d", podList.Size())
	labelSelector := labels.SelectorFromSet(labelsForAputturKind(instance.Name))
	listOps := &client.ListOptions{Namespace: instance.Namespace, LabelSelector: labelSelector}
	err = r.client.List(context.TODO(), listOps, podList)
	if err != nil {
		log.Printf("Failed to list pods: %v", err)
		return reconcile.Result{}, err
	}
	podNames := getPodNames(podList.Items)

	// Update status.PodNames if needed
	log.Print("Listing pod names")
	if !reflect.DeepEqual(podNames, instance.Status.PodNames) {
		instance.Status.PodNames = podNames
		err := r.client.Update(context.TODO(), instance)
		if err != nil {
			log.Printf("failed to update node status: %v", err)
			return reconcile.Result{}, err
		}
	}

	// Update AppGroup status
	log.Print("Update  ApGroup and Status")
	if instance.Spec.Group != instance.Status.AppGroup {
		instance.Status.AppGroup = instance.Spec.Group
		err := r.client.Update(context.TODO(), instance)
		if err != nil {
			log.Printf("failed to update group status: %v", err)
			return reconcile.Result{}, err
		}
		log.Print("Update  ApGroup and Status for the instance run kubectl describe Aputturkind aputtur-example")
	}

	return reconcile.Result{}, nil
}

// Reconcile reads that state of the cluster for a Aputturkind object and makes changes based on the state read
// and what is in the Aputturkind.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileAputturkind) Reconcileold(request reconcile.Request) (reconcile.Result, error) {
	log.Printf("Reconciling Aputturkind %s/%s\n", request.Namespace, request.Name)

	// Fetch the Aputturkind instance
	instance := &examplev1alpha1.Aputturkind{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Printf("ERRORRR")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Printf("ERRORRR --again")
		return reconcile.Result{}, err
	}

	// Define a new Pod object
	//pod := newPodForCR(instance)

	// Set Aputturkind instance as the owner and controller
	//if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
	//	return reconcile.Result{}, err
	//}

	// Check if this Pod already exists
	//found := &corev1.Pod{}

	// Check if the deployment already exists, if not create a new one
	found := &appsv1.Deployment{}

	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		dep := r.newDeploymentForCR(instance)
		log.Printf("Creating a new deployment %s dep.Name %s\n", dep.Namespace, dep.Name)
		err = r.client.Create(context.TODO(), dep)
		if err != nil {
			log.Printf("Failed to create new Deployment: %v\n", err)
			return reconcile.Result{}, err
		}

		// deployment created successfully -  requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		log.Printf("Failed to get deployment")
		return reconcile.Result{}, err
	}

	//ensure the deployment count is same
	count := instance.Spec.Count
	if *found.Spec.Replicas != count {
		found.Spec.Replicas = &count
		err = r.client.Update(context.TODO(), found)
		if err != nil {
			log.Printf("Failed to update deployment")
			return reconcile.Result{}, err
		}

		return reconcile.Result{Requeue: true}, nil
	}

	// List the pods for this deployment
	podList := &corev1.PodList{}
	labelSelector := labels.SelectorFromSet(labelsForAputturKind(instance.Name))
	listOps := &client.ListOptions{Namespace: instance.Namespace, LabelSelector: labelSelector}
	err = r.client.List(context.TODO(), listOps, podList)
	if err != nil {
		log.Printf("Failed to list pods: %v", err)
		return reconcile.Result{}, err
	}
	podNames := getPodNames(podList.Items)

	// Update status.PodNames if needed
	if !reflect.DeepEqual(podNames, instance.Status.PodNames) {
		instance.Status.PodNames = podNames
		err := r.client.Update(context.TODO(), instance)
		if err != nil {
			log.Printf("failed to update node status: %v", err)
			return reconcile.Result{}, err
		}
	}

	// Update AppGroup status
	if instance.Spec.Group != instance.Status.AppGroup {
		instance.Status.AppGroup = instance.Spec.Group
		err := r.client.Update(context.TODO(), instance)
		if err != nil {
			log.Printf("failed to update group status: %v", err)
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *examplev1alpha1.Aputturkind) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}
}

//getPodNames
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

//Set labes in a map
func labelsForAputturKind(name string) map[string]string {
	return map[string]string{"app": "Operator-Example", "operatorexample_cr": name}

}

// Create newDeploymentForCR method to create a deployment.
func (r *ReconcileAputturkind) newDeploymentForCR(m *examplev1alpha1.Aputturkind) *appsv1.Deployment {
	labels := labelsForAputturKind(m.Name)
	replicas := m.Spec.Count
	dep := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: m.Spec.Image,
						Name:  m.Name,
						Ports: []corev1.ContainerPort{{
							ContainerPort: m.Spec.Port,
							Name:          m.Name,
						}},
					}},
				},
			},
		},
	}
	// Set Examplekind instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep

}
