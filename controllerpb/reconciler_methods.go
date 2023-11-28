package controllerpb

import "sigs.k8s.io/controller-runtime/pkg/reconcile"

func (r *ReconcileResult) ReconcileResult() reconcile.Result {
	return reconcile.Result{
		Requeue:      r.Requeue,
		RequeueAfter: r.RequeueAfter.AsDuration(),
	}
}
