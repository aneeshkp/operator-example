 
operator-sdk build aneeshkp/operator-example:v0.0.6
docker push aneeshkp/operator-example:v0.0.6


 kubectl create -f deploy/service_account.yaml 
 kubectl create -f deploy/role
 kubectl create -f deploy/role.yaml 
 kubectl create -f deploy/role_binding.yaml 
 kubectl create -f deploy/operator.yaml 
 kubectl create -f deploy/crds/example_v1alpha1_aputturkind_crd.yaml 
 kubectl create -f deploy/operator.yaml 
 kubectl apply -f deploy/crds/example_v1alpha1_aputturkind_cr.yaml 


 kubectl describe Aputturkind aputtur-example
 



### debug
export OPERATOR_NAME=operator-example
operator-sdk up local --namespace=default
