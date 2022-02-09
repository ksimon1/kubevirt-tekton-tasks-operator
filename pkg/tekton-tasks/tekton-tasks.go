package tekton_tasks

import (
	"fmt"

	"github.com/kubevirt/kubevirt-tekton-tasks-operator/pkg/common"
	"github.com/kubevirt/kubevirt-tekton-tasks-operator/pkg/operands"
	tektonbundle "github.com/kubevirt/kubevirt-tekton-tasks-operator/pkg/tekton-bundle"
	pipeline "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	v1 "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// +kubebuilder:rbac:groups=tekton.dev,resources=clustertasks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterroles;rolebindings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kubevirt.io,resources=virtualmachineinstances;virtualmachines,verbs=create;update;get;list;watch;delete
// +kubebuilder:rbac:groups=subresources.kubevirt.io,resources=virtualmachines/restart;virtualmachines/start;virtualmachines/stop,verbs=update
// +kubebuilder:rbac:groups=template.openshift.io,resources=templates,verbs=get;list;watch;create;patch
// +kubebuilder:rbac:groups=template.openshift.io,resources=processedtemplates,verbs=create
// +kubebuilder:rbac:groups=cdi.kubevirt.io,resources=datavolumes,verbs=*
// +kubebuilder:rbac:groups=kubevirt.io,resources=virtualmachines/finalizers,verbs=*
// +kubebuilder:rbac:groups=*,resources=persistentvolumeclaims,verbs=*

const (
	operandName      = "tekton-tasks"
	operandComponent = common.AppComponentTemplating
)

var allowedTasks = map[string]bool{
	"cleanup-vm":                      true,
	"copy-template":                   true,
	"create-datavolume-from-manifest": true,
	"create-vm-from-manifest":         true,
	"disk-virt-customize":             true,
	"disk-virt-sysprep":               true,
	"modify-vm-template":              true,
	"wait-for-vmi-status":             true,
}

func init() {
	utilruntime.Must(pipeline.AddToScheme(common.Scheme))
}

type tektonTasks struct {
	clusterTasks    []pipeline.ClusterTask
	serviceAccounts []v1.ServiceAccount
	roleBinding     []rbac.RoleBinding
	clusterRole     []rbac.ClusterRole
}

var _ operands.Operand = &tektonTasks{}

func New(bundle *tektonbundle.Bundle) *tektonTasks {
	tt := &tektonTasks{
		clusterTasks:    bundle.ClusterTasks,
		serviceAccounts: bundle.ServiceAccounts,
		roleBinding:     bundle.RoleBindings,
		clusterRole:     bundle.ClusterRoles,
	}

	tt.filterUnusedObjects()

	return tt
}

func (t *tektonTasks) Name() string {
	return operandName
}

func (t *tektonTasks) WatchClusterTypes() []client.Object {
	return []client.Object{
		&rbac.ClusterRole{},
		&pipeline.ClusterTask{},
	}
}

func (t *tektonTasks) WatchTypes() []client.Object {
	return []client.Object{
		&rbac.RoleBinding{},
		&v1.ServiceAccount{},
	}
}

func (t *tektonTasks) RequiredCrds() []string {
	return []string{
		"tasks.tekton.dev",
	}
}
func (t *tektonTasks) filterUnusedObjects() {
	newTasks := []pipeline.ClusterTask{}
	for _, task := range t.clusterTasks {
		if _, ok := allowedTasks[task.Name]; ok {
			newTasks = append(newTasks, task)
		}
	}
	t.clusterTasks = newTasks

	newSA := []v1.ServiceAccount{}
	for _, sa := range t.serviceAccounts {
		if _, ok := allowedTasks[sa.Name+"-task"]; ok {
			newSA = append(newSA, sa)
		}
	}
	t.serviceAccounts = newSA

	newRB := []rbac.RoleBinding{}
	for _, rb := range t.roleBinding {
		if _, ok := allowedTasks[rb.Name+"-task"]; ok {
			newRB = append(newRB, rb)
		}
	}
	t.roleBinding = newRB

	newCR := []rbac.ClusterRole{}
	for _, cr := range t.clusterRole {
		if _, ok := allowedTasks[cr.Name+"-task"]; ok {
			newCR = append(newCR, cr)
		}
	}
	t.clusterRole = newCR

}
func (t *tektonTasks) Reconcile(request *common.Request) ([]common.ReconcileResult, error) {
	var results []common.ReconcileResult
	var reconcileFunc []common.ReconcileFunc
	reconcileFunc = append(reconcileFunc, reconcileTektonTasksFuncs(t.clusterTasks)...)
	reconcileFunc = append(reconcileFunc, reconcileClusterRoleFuncs(t.clusterRole)...)
	reconcileFunc = append(reconcileFunc, reconcileServiceAccountsFuncs(t.serviceAccounts)...)
	reconcileFunc = append(reconcileFunc, reconcileRoleBindingFuncs(t.roleBinding)...)

	reconcileTektonBundleResults, err := common.CollectResourceStatus(request, reconcileFunc...)
	if err != nil {
		return nil, err
	}

	upgradingNow := isUpgradingNow(request)
	for _, r := range reconcileTektonBundleResults {
		if !upgradingNow && (r.OperationResult == common.OperationResultUpdated) {
			request.Logger.Info(fmt.Sprintf("Changes reverted in tekton tasks: %s", r.Resource.GetName()))
		}
	}
	return append(results, reconcileTektonBundleResults...), nil
}

func (t *tektonTasks) Cleanup(request *common.Request) ([]common.CleanupResult, error) {
	var objects []client.Object

	return common.DeleteAll(request, objects...)
}

func isUpgradingNow(request *common.Request) bool {
	return request.Instance.Status.ObservedVersion != common.GetOperatorVersion()
}

func reconcileTektonTasksFuncs(tasks []pipeline.ClusterTask) []common.ReconcileFunc {
	funcs := make([]common.ReconcileFunc, 0, len(tasks))
	for i := range tasks {
		task := &tasks[i]
		funcs = append(funcs, func(request *common.Request) (common.ReconcileResult, error) {
			return common.CreateOrUpdate(request).
				ClusterResource(task).
				WithAppLabels(operandName, operandComponent).
				UpdateFunc(func(newRes, foundRes client.Object) {
					newTask := newRes.(*pipeline.ClusterTask)
					foundTask := foundRes.(*pipeline.ClusterTask)
					foundTask.Spec = newTask.Spec
				}).
				Reconcile()
		})
	}
	return funcs
}

func reconcileClusterRoleFuncs(crs []rbac.ClusterRole) []common.ReconcileFunc {
	funcs := make([]common.ReconcileFunc, 0, len(crs))
	for i := range crs {
		cr := &crs[i]
		funcs = append(funcs, func(request *common.Request) (common.ReconcileResult, error) {
			return common.CreateOrUpdate(request).
				ClusterResource(cr).
				WithAppLabels(operandName, operandComponent).
				UpdateFunc(func(newRes, foundRes client.Object) {
					newTask := newRes.(*rbac.ClusterRole)
					foundTask := foundRes.(*rbac.ClusterRole)
					foundTask.Rules = newTask.Rules
				}).
				Reconcile()
		})
	}
	return funcs
}

func reconcileServiceAccountsFuncs(sas []v1.ServiceAccount) []common.ReconcileFunc {
	funcs := make([]common.ReconcileFunc, 0, len(sas))
	for i := range sas {
		sa := &sas[i]
		funcs = append(funcs, func(request *common.Request) (common.ReconcileResult, error) {
			return common.CreateOrUpdate(request).
				ClusterResource(sa).
				WithAppLabels(operandName, operandComponent).
				Reconcile()
		})
	}
	return funcs
}

func reconcileRoleBindingFuncs(rbs []rbac.RoleBinding) []common.ReconcileFunc {
	funcs := make([]common.ReconcileFunc, 0, len(rbs))
	for i := range rbs {
		rb := &rbs[i]
		funcs = append(funcs, func(request *common.Request) (common.ReconcileResult, error) {
			return common.CreateOrUpdate(request).
				ClusterResource(rb).
				WithAppLabels(operandName, operandComponent).
				UpdateFunc(func(newRes, foundRes client.Object) {
					newTask := newRes.(*rbac.RoleBinding)
					foundTask := foundRes.(*rbac.RoleBinding)
					foundTask.RoleRef = newTask.RoleRef
					foundTask.Subjects = newTask.Subjects
				}).
				Reconcile()
		})
	}
	return funcs
}
