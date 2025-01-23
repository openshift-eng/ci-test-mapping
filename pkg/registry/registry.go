package registry

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/apiserverauth"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/awsloadbalanceroperator"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/baremetalhardwareprovisioning"
	baremetalhardwareprovisioningbaremetaloperator "github.com/openshift-eng/ci-test-mapping/pkg/components/baremetalhardwareprovisioning/baremetaloperator"
	baremetalhardwareprovisioningclusterapiprovider "github.com/openshift-eng/ci-test-mapping/pkg/components/baremetalhardwareprovisioning/clusterapiprovider"
	baremetalhardwareprovisioningclusterbaremetaloperator "github.com/openshift-eng/ci-test-mapping/pkg/components/baremetalhardwareprovisioning/clusterbaremetaloperator"
	baremetalhardwareprovisioningironic "github.com/openshift-eng/ci-test-mapping/pkg/components/baremetalhardwareprovisioning/ironic"
	baremetalhardwareprovisioningosimageprovider "github.com/openshift-eng/ci-test-mapping/pkg/components/baremetalhardwareprovisioning/osimageprovider"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/bmerevents"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/bpfman"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/build"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/certmanager"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/climanager"
	cloudcomputebaremetalprovider "github.com/openshift-eng/ci-test-mapping/pkg/components/cloudcompute/baremetalprovider"
	cloudcomputecloudcontrollermanager "github.com/openshift-eng/ci-test-mapping/pkg/components/cloudcompute/cloudcontrollermanager"
	cloudcomputeclusterapiproviders "github.com/openshift-eng/ci-test-mapping/pkg/components/cloudcompute/clusterapiproviders"
	cloudcomputeclusterautoscaler "github.com/openshift-eng/ci-test-mapping/pkg/components/cloudcompute/clusterautoscaler"
	cloudcomputecontrolplanemachineset "github.com/openshift-eng/ci-test-mapping/pkg/components/cloudcompute/controlplanemachineset"
	cloudcomputeexternalprovider "github.com/openshift-eng/ci-test-mapping/pkg/components/cloudcompute/externalprovider"
	cloudcomputeibmprovider "github.com/openshift-eng/ci-test-mapping/pkg/components/cloudcompute/ibmprovider"
	cloudcomputekubevirtprovider "github.com/openshift-eng/ci-test-mapping/pkg/components/cloudcompute/kubevirtprovider"
	cloudcomputelibvirtprovider "github.com/openshift-eng/ci-test-mapping/pkg/components/cloudcompute/libvirtprovider"
	cloudcomputemachineapiproviders "github.com/openshift-eng/ci-test-mapping/pkg/components/cloudcompute/machineapiproviders"
	cloudcomputemachinecsrapprover "github.com/openshift-eng/ci-test-mapping/pkg/components/cloudcompute/machinecsrapprover"
	cloudcomputemachinehealthcheck "github.com/openshift-eng/ci-test-mapping/pkg/components/cloudcompute/machinehealthcheck"
	cloudcomputenutanixprovider "github.com/openshift-eng/ci-test-mapping/pkg/components/cloudcompute/nutanixprovider"
	cloudcomputeopenstackprovider "github.com/openshift-eng/ci-test-mapping/pkg/components/cloudcompute/openstackprovider"
	cloudcomputeovirtprovider "github.com/openshift-eng/ci-test-mapping/pkg/components/cloudcompute/ovirtprovider"
	cloudcomputeunknown "github.com/openshift-eng/ci-test-mapping/pkg/components/cloudcompute/unknown"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/cloudcompute/vsphereprovider"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/cloudcredentialoperator"
	cloudnativeeventscloudeventproxy "github.com/openshift-eng/ci-test-mapping/pkg/components/cloudnativeevents/cloudeventproxy"
	cloudnativeeventscloudnativeevents "github.com/openshift-eng/ci-test-mapping/pkg/components/cloudnativeevents/cloudnativeevents"
	cloudnativeeventshardwareeventproxy "github.com/openshift-eng/ci-test-mapping/pkg/components/cloudnativeevents/hardwareeventproxy"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/clusterloader"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/clusterresourceoverrideadmissionoperator"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/clusterversionoperator"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/cnfcerttnf"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/complianceoperator"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/confidentialcomputeattestation"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/configoperator"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/consolemetal3plugin"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/consolestorageplugin"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/containers"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/crc"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/descheduler"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/devconsole"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/drivertoolkit"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/etcd"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/externaldnsoperator"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/fileintegrityoperator"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/gitopsztp"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/hawkular"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/helm"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/hive"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/hypershift"
	hypershiftagent "github.com/openshift-eng/ci-test-mapping/pkg/components/hypershift/agent"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/hypershift/aro"
	hypershiftocpvirtualization "github.com/openshift-eng/ci-test-mapping/pkg/components/hypershift/ocpvirtualization"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/hypershift/openstack"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/hypershift/rosa"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/ibmrokstoolkit"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/imageregistry"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/imagestreams"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/insightsoperator"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/insightsruntimeextractor"
	installeragentbasedinstallation "github.com/openshift-eng/ci-test-mapping/pkg/components/installer/agentbasedinstallation"
	installeralibabacloud "github.com/openshift-eng/ci-test-mapping/pkg/components/installer/alibabacloud"
	installerassistedinstaller "github.com/openshift-eng/ci-test-mapping/pkg/components/installer/assistedinstaller"
	installeribmcloud "github.com/openshift-eng/ci-test-mapping/pkg/components/installer/ibmcloud"
	installernutanix "github.com/openshift-eng/ci-test-mapping/pkg/components/installer/nutanix"
	installeropenshiftansible "github.com/openshift-eng/ci-test-mapping/pkg/components/installer/openshiftansible"
	installeropenshiftinstaller "github.com/openshift-eng/ci-test-mapping/pkg/components/installer/openshiftinstaller"
	installeropenshiftonbaremetalipi "github.com/openshift-eng/ci-test-mapping/pkg/components/installer/openshiftonbaremetalipi"
	installeropenshiftonopenstack "github.com/openshift-eng/ci-test-mapping/pkg/components/installer/openshiftonopenstack"
	installeropenshiftonrhv "github.com/openshift-eng/ci-test-mapping/pkg/components/installer/openshiftonrhv"
	installerpowervs "github.com/openshift-eng/ci-test-mapping/pkg/components/installer/powervs"
	installersinglenodeopenshift "github.com/openshift-eng/ci-test-mapping/pkg/components/installer/singlenodeopenshift"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/installer/vsphere"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/isvoperators"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/jenkins"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/kmm"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/ksanstorage"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/kubeapiserver"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/kubecontrollermanager"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/kubescheduler"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/kubestorageversionmigrator"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/lcaoperator"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/lightspeed"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/logging"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/logicalvolumemanagerstorage"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/lowlatencyvalidationtooling"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/machineconfigoperator"
	machineconfigoperatorplatformbaremetal "github.com/openshift-eng/ci-test-mapping/pkg/components/machineconfigoperator/platformbaremetal"
	machineconfigoperatorplatformnone "github.com/openshift-eng/ci-test-mapping/pkg/components/machineconfigoperator/platformnone"
	machineconfigoperatorplatformopenstack "github.com/openshift-eng/ci-test-mapping/pkg/components/machineconfigoperator/platformopenstack"
	machineconfigoperatorplatformovirtrhv "github.com/openshift-eng/ci-test-mapping/pkg/components/machineconfigoperator/platformovirtrhv"
	machineconfigoperatorplatformvsphere "github.com/openshift-eng/ci-test-mapping/pkg/components/machineconfigoperator/platformvsphere"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/managementconsole"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/meteringoperator"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/microshift"
	microshiftnetworking "github.com/openshift-eng/ci-test-mapping/pkg/components/microshift/networking"
	microshiftstorage "github.com/openshift-eng/ci-test-mapping/pkg/components/microshift/storage"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/monitoring"
	monitoringgrafana "github.com/openshift-eng/ci-test-mapping/pkg/components/monitoring/grafana"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/multiarch"
	multiarcharm "github.com/openshift-eng/ci-test-mapping/pkg/components/multiarch/arm"
	multiarchibmpandz "github.com/openshift-eng/ci-test-mapping/pkg/components/multiarch/ibmpandz"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/multiarchtuningoperator"
	networkingcloudnetworkconfigcontroller "github.com/openshift-eng/ci-test-mapping/pkg/components/networking/cloudnetworkconfigcontroller"
	networkingclusternetworkoperator "github.com/openshift-eng/ci-test-mapping/pkg/components/networking/clusternetworkoperator"
	networkingdns "github.com/openshift-eng/ci-test-mapping/pkg/components/networking/dns"
	networkingdpu "github.com/openshift-eng/ci-test-mapping/pkg/components/networking/dpu"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/networking/frrk8s"
	networkingingressnodefirewall "github.com/openshift-eng/ci-test-mapping/pkg/components/networking/ingressnodefirewall"
	networkingkubernetesnmstate "github.com/openshift-eng/ci-test-mapping/pkg/components/networking/kubernetesnmstate"
	networkingkubernetesnmstateoperator "github.com/openshift-eng/ci-test-mapping/pkg/components/networking/kubernetesnmstateoperator"
	networkingkuryr "github.com/openshift-eng/ci-test-mapping/pkg/components/networking/kuryr"
	networkingmetallb "github.com/openshift-eng/ci-test-mapping/pkg/components/networking/metallb"
	networkingmultus "github.com/openshift-eng/ci-test-mapping/pkg/components/networking/multus"
	networkingnetobs "github.com/openshift-eng/ci-test-mapping/pkg/components/networking/netobs"
	networkingnetworkingconsoleplugin "github.com/openshift-eng/ci-test-mapping/pkg/components/networking/networkingconsoleplugin"
	networkingnetworktools "github.com/openshift-eng/ci-test-mapping/pkg/components/networking/networktools"
	networkingnmstateconsoleplugin "github.com/openshift-eng/ci-test-mapping/pkg/components/networking/nmstateconsoleplugin"
	networkingonpremdns "github.com/openshift-eng/ci-test-mapping/pkg/components/networking/onpremdns"
	networkingonpremhostnetworking "github.com/openshift-eng/ci-test-mapping/pkg/components/networking/onpremhostnetworking"
	networkingonpremloadbalancer "github.com/openshift-eng/ci-test-mapping/pkg/components/networking/onpremloadbalancer"
	networkingopenshiftsdn "github.com/openshift-eng/ci-test-mapping/pkg/components/networking/openshiftsdn"
	networkingovnkubernetes "github.com/openshift-eng/ci-test-mapping/pkg/components/networking/ovnkubernetes"
	networkingptp "github.com/openshift-eng/ci-test-mapping/pkg/components/networking/ptp"
	networkingrouter "github.com/openshift-eng/ci-test-mapping/pkg/components/networking/router"
	networkingruntimecfg "github.com/openshift-eng/ci-test-mapping/pkg/components/networking/runtimecfg"
	networkingsriov "github.com/openshift-eng/ci-test-mapping/pkg/components/networking/sriov"
	nodecpumanager "github.com/openshift-eng/ci-test-mapping/pkg/components/node/cpumanager"
	nodecrio "github.com/openshift-eng/ci-test-mapping/pkg/components/node/crio"
	nodedevicemanager "github.com/openshift-eng/ci-test-mapping/pkg/components/node/devicemanager"
	nodeinstasliceoperator "github.com/openshift-eng/ci-test-mapping/pkg/components/node/instasliceoperator"
	nodekubelet "github.com/openshift-eng/ci-test-mapping/pkg/components/node/kubelet"
	nodememorymanager "github.com/openshift-eng/ci-test-mapping/pkg/components/node/memorymanager"
	nodenodeproblemdetector "github.com/openshift-eng/ci-test-mapping/pkg/components/node/nodeproblemdetector"
	nodenumaawarescheduling "github.com/openshift-eng/ci-test-mapping/pkg/components/node/numaawarescheduling"
	nodepodresourceapi "github.com/openshift-eng/ci-test-mapping/pkg/components/node/podresourceapi"
	nodetopologymanager "github.com/openshift-eng/ci-test-mapping/pkg/components/node/topologymanager"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/nodefeaturediscoveryoperator"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/nodemaintenanceoperator"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/nodeobservabilityoperator"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/nodetuningoperator"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/none"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/nvidia"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/oauthapiserver"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/oauthproxy"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/observabilityui"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/oc"
	occlustercompare "github.com/openshift-eng/ci-test-mapping/pkg/components/oc/clustercompare"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/oc/nodeimage"
	ococmirror "github.com/openshift-eng/ci-test-mapping/pkg/components/oc/ocmirror"
	ocupdate "github.com/openshift-eng/ci-test-mapping/pkg/components/oc/update"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/occompliance"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/ocmirror"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/olm"
	olmoperatorhub "github.com/openshift-eng/ci-test-mapping/pkg/components/olm/operatorhub"
	olmregistry "github.com/openshift-eng/ci-test-mapping/pkg/components/olm/registry"
	opctcli "github.com/openshift-eng/ci-test-mapping/pkg/components/opct/cli"
	opctother "github.com/openshift-eng/ci-test-mapping/pkg/components/opct/other"
	opctresultsgeneral "github.com/openshift-eng/ci-test-mapping/pkg/components/opct/results/general"
	opctresultsvcsp "github.com/openshift-eng/ci-test-mapping/pkg/components/opct/results/vcsp"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/openshiftapiserver"
	openshiftcontrollermanagerapps "github.com/openshift-eng/ci-test-mapping/pkg/components/openshiftcontrollermanager/apps"
	openshiftcontrollermanagerbuild "github.com/openshift-eng/ci-test-mapping/pkg/components/openshiftcontrollermanager/build"
	openshiftcontrollermanagercontrollermanager "github.com/openshift-eng/ci-test-mapping/pkg/components/openshiftcontrollermanager/controllermanager"
	openshiftupdateserviceoperand "github.com/openshift-eng/ci-test-mapping/pkg/components/openshiftupdateservice/operand"
	openshiftupdateserviceoperator "github.com/openshift-eng/ci-test-mapping/pkg/components/openshiftupdateservice/operator"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/operatorsdk"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/performanceaddonoperator"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/podautoscaler"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/poisonpilloperator"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/registryconsole"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/release"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/rhcos"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/rhmimonitoring"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/routecontrollermanager"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/runoncedurationoverride"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/samplesoperator"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/sandboxedcontainers"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/secondaryscheduleroperator"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/security"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/securityprofilesoperator"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/servicebinding"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/servicebroker"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/serviceca"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/servicecatalog"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/specialresourceoperator"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/storage"
	storagekubernetes "github.com/openshift-eng/ci-test-mapping/pkg/components/storage/kubernetes"
	storagekubernetesexternalcomponents "github.com/openshift-eng/ci-test-mapping/pkg/components/storage/kubernetesexternalcomponents"
	storagelocalstorageoperator "github.com/openshift-eng/ci-test-mapping/pkg/components/storage/localstorageoperator"
	storageopenstackcsidrivers "github.com/openshift-eng/ci-test-mapping/pkg/components/storage/openstackcsidrivers"
	storageoperators "github.com/openshift-eng/ci-test-mapping/pkg/components/storage/operators"
	storageovirtcsidriver "github.com/openshift-eng/ci-test-mapping/pkg/components/storage/ovirtcsidriver"
	storagesharedresourcecsidriver "github.com/openshift-eng/ci-test-mapping/pkg/components/storage/sharedresourcecsidriver"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/talmoperator"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/telcoperformance"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/telemeter"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/templates"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/testframework"
	testframeworkopenstack "github.com/openshift-eng/ci-test-mapping/pkg/components/testframework/openstack"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/testinfrastructure"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/twonode/arbiter"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/unknown"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/windowscontainers"
)

type Registry struct {
	Components map[string]v1.Component
}

func NewComponentRegistry() *Registry {
	var r Registry

	r.Register("AWS Load Balancer Operator", &awsloadbalanceroperator.AWSLoadBalancerOperatorComponent)
	r.Register("BMER Events", &bmerevents.BMEREventsComponent)
	r.Register("Bare Metal Hardware Provisioning / OS Image Provider", &baremetalhardwareprovisioningosimageprovider.OSImageProviderComponent)
	r.Register("Bare Metal Hardware Provisioning / baremetal-operator", &baremetalhardwareprovisioningbaremetaloperator.BaremetalOperatorComponent)
	r.Register("Bare Metal Hardware Provisioning / cluster-api-provider", &baremetalhardwareprovisioningclusterapiprovider.ClusterAPIProviderComponent)
	r.Register("Bare Metal Hardware Provisioning / cluster-baremetal-operator", &baremetalhardwareprovisioningclusterbaremetaloperator.ClusterBaremetalOperatorComponent)
	r.Register("Bare Metal Hardware Provisioning / ironic", &baremetalhardwareprovisioningironic.IronicComponent)
	r.Register("Bare Metal Hardware Provisioning", &baremetalhardwareprovisioning.BareMetalHardwareProvisioningComponent)
	r.Register("Build", &build.BuildComponent)
	r.Register("CNF-Cert-TNF", &cnfcerttnf.CNFCertTNFComponent)
	r.Register("Cloud Compute / BareMetal Provider", &cloudcomputebaremetalprovider.BareMetalProviderComponent)
	r.Register("Cloud Compute / Cloud Controller Manager", &cloudcomputecloudcontrollermanager.CloudControllerManagerComponent)
	r.Register("Cloud Compute / Cluster API Providers", &cloudcomputeclusterapiproviders.ClusterAPIProvidersComponent)
	r.Register("Cloud Compute / Cluster Autoscaler", &cloudcomputeclusterautoscaler.ClusterAutoscalerComponent)
	r.Register("Cloud Compute / ControlPlaneMachineSet", &cloudcomputecontrolplanemachineset.ControlPlaneMachineSetComponent)
	r.Register("Cloud Compute / External Provider", &cloudcomputeexternalprovider.ExternalProviderComponent)
	r.Register("Cloud Compute / IBM Provider", &cloudcomputeibmprovider.IBMProviderComponent)
	r.Register("Cloud Compute / KubeVirt Provider", &cloudcomputekubevirtprovider.KubeVirtProviderComponent)
	r.Register("Cloud Compute / Libvirt Provider", &cloudcomputelibvirtprovider.LibvirtProviderComponent)
	r.Register("Cloud Compute / Machine API Providers", &cloudcomputemachineapiproviders.MachineAPIProvidersComponent)
	r.Register("Cloud Compute / Machine CSR Approver", &cloudcomputemachinecsrapprover.MachineCSRApproverComponent)
	r.Register("Cloud Compute / MachineHealthCheck", &cloudcomputemachinehealthcheck.MachineHealthCheckComponent)
	r.Register("Cloud Compute / Nutanix Provider", &cloudcomputenutanixprovider.NutanixProviderComponent)
	r.Register("Cloud Compute / OpenStack Provider", &cloudcomputeopenstackprovider.OpenStackProviderComponent)
	r.Register("Cloud Compute / Unknown", &cloudcomputeunknown.UnknownComponent)
	r.Register("Cloud Compute / oVirt Provider", &cloudcomputeovirtprovider.OVirtProviderComponent)
	r.Register("Cloud Credential Operator", &cloudcredentialoperator.CloudCredentialOperatorComponent)
	r.Register("Cloud Native Events / Cloud Event Proxy", &cloudnativeeventscloudeventproxy.CloudEventProxyComponent)
	r.Register("Cloud Native Events / Cloud Native Events", &cloudnativeeventscloudnativeevents.CloudNativeEventsComponent)
	r.Register("Cloud Native Events / Hardware Event Proxy", &cloudnativeeventshardwareeventproxy.HardwareEventProxyComponent)
	r.Register("Cluster Loader", &clusterloader.ClusterLoaderComponent)
	r.Register("Cluster Resource Override Admission Operator", &clusterresourceoverrideadmissionoperator.ClusterResourceOverrideAdmissionOperatorComponent)
	r.Register("Cluster Version Operator", &clusterversionoperator.ClusterVersionOperatorComponent)
	r.Register("Compliance Operator", &complianceoperator.ComplianceOperatorComponent)
	r.Register("Console Metal3 Plugin", &consolemetal3plugin.ConsoleMetal3PluginComponent)
	r.Register("Console Storage Plugin", &consolestorageplugin.ConsoleStoragePluginComponent)
	r.Register("Containers", &containers.ContainersComponent)
	r.Register("Dev Console", &devconsole.DevConsoleComponent)
	r.Register("Driver Toolkit", &drivertoolkit.DriverToolkitComponent)
	r.Register("Etcd", &etcd.EtcdComponent)
	r.Register("ExternalDNS Operator", &externaldnsoperator.ExternalDNSOperatorComponent)
	r.Register("File Integrity Operator", &fileintegrityoperator.FileIntegrityOperatorComponent)
	r.Register("GitOps ZTP", &gitopsztp.GitOpsZTPComponent)
	r.Register("Hawkular", &hawkular.HawkularComponent)
	r.Register("Helm", &helm.HelmComponent)
	r.Register("Hive", &hive.HiveComponent)
	r.Register("HyperShift / Agent", &hypershiftagent.AgentComponent)
	r.Register("HyperShift / OCP Virtualization", &hypershiftocpvirtualization.OCPVirtualizationComponent)
	r.Register("HyperShift", &hypershift.HyperShiftComponent)
	r.Register("ISV Operators", &isvoperators.ISVOperatorsComponent)
	r.Register("Image Registry", &imageregistry.ImageRegistryComponent)
	r.Register("ImageStreams", &imagestreams.ImageStreamsComponent)
	r.Register("Insights Operator", &insightsoperator.InsightsOperatorComponent)
	r.Register("Installer / Agent based installation", &installeragentbasedinstallation.AgentBasedInstallationComponent)
	r.Register("Installer / Alibaba Cloud", &installeralibabacloud.AlibabaCloudComponent)
	r.Register("Installer / Assisted installer", &installerassistedinstaller.AssistedInstallerComponent)
	r.Register("Installer / IBM Cloud", &installeribmcloud.IBMCloudComponent)
	r.Register("Installer / Nutanix", &installernutanix.NutanixComponent)
	r.Register("Installer / OpenShift on Bare Metal IPI", &installeropenshiftonbaremetalipi.OpenShiftOnBareMetalIPIComponent)
	r.Register("Installer / OpenShift on OpenStack", &installeropenshiftonopenstack.OpenShiftOnOpenStackComponent)
	r.Register("Installer / OpenShift on RHV", &installeropenshiftonrhv.OpenShiftOnRHVComponent)
	r.Register("Installer / PowerVS", &installerpowervs.PowerVSComponent)
	r.Register("Installer / Single Node OpenShift", &installersinglenodeopenshift.SingleNodeOpenShiftComponent)
	r.Register("Installer / openshift-ansible", &installeropenshiftansible.OpenshiftAnsibleComponent)
	r.Register("Installer / openshift-installer", &installeropenshiftinstaller.OpenshiftInstallerComponent)
	r.Register("Jenkins", &jenkins.JenkinsComponent)
	r.Register("LCA operator", &lcaoperator.LCAOperatorComponent)
	r.Register("Logging", &logging.LoggingComponent)
	r.Register("Logical Volume Manager Storage", &logicalvolumemanagerstorage.LogicalVolumeManagerStorageComponent)
	r.Register("Low latency validation tooling", &lowlatencyvalidationtooling.LowLatencyValidationToolingComponent)
	r.Register("Machine Config Operator / platform-baremetal", &machineconfigoperatorplatformbaremetal.PlatformBaremetalComponent)
	r.Register("Machine Config Operator / platform-none", &machineconfigoperatorplatformnone.PlatformNoneComponent)
	r.Register("Machine Config Operator / platform-openstack", &machineconfigoperatorplatformopenstack.PlatformOpenstackComponent)
	r.Register("Machine Config Operator / platform-ovirt-rhv", &machineconfigoperatorplatformovirtrhv.PlatformOvirtRhvComponent)
	r.Register("Machine Config Operator / platform-vsphere", &machineconfigoperatorplatformvsphere.PlatformVsphereComponent)
	r.Register("Machine Config Operator", &machineconfigoperator.MachineConfigOperatorComponent)
	r.Register("Management Console", &managementconsole.ManagementConsoleComponent)
	r.Register("Metering Operator", &meteringoperator.MeteringOperatorComponent)
	r.Register("MicroShift / Networking", &microshiftnetworking.NetworkingComponent)
	r.Register("MicroShift / Storage", &microshiftstorage.StorageComponent)
	r.Register("MicroShift", &microshift.MicroShiftComponent)
	r.Register("Monitoring / Grafana", &monitoringgrafana.GrafanaComponent)
	r.Register("Monitoring", &monitoring.MonitoringComponent)
	r.Register("Multi-Arch / ARM", &multiarcharm.ARMComponent)
	r.Register("Multi-Arch / IBM P and Z", &multiarchibmpandz.IBMPAndZComponent)
	r.Register("Multi-Arch", &multiarch.MultiArchComponent)
	r.Register("Multiarch Tuning Operator", &multiarchtuningoperator.MultiarchTuningOperatorComponent)
	r.Register("NVIDIA", &nvidia.NVIDIAComponent)
	r.Register("Networking / DNS", &networkingdns.DNSComponent)
	r.Register("Networking / DPU", &networkingdpu.DPUComponent)
	r.Register("Networking / Metal LB", &networkingmetallb.MetalLBComponent)
	r.Register("Networking / NetObs", &networkingnetobs.NetObsComponent)
	r.Register("Networking / On-Prem DNS", &networkingonpremdns.OnPremDNSComponent)
	r.Register("Networking / On-Prem Host Networking", &networkingonpremhostnetworking.OnPremHostNetworkingComponent)
	r.Register("Networking / On-Prem Load Balancer", &networkingonpremloadbalancer.OnPremLoadBalancerComponent)
	r.Register("Networking / SR-IOV", &networkingsriov.SRIOVComponent)
	r.Register("Networking / cloud-network-config-controller", &networkingcloudnetworkconfigcontroller.CloudNetworkConfigControllerComponent)
	r.Register("Networking / cluster-network-operator", &networkingclusternetworkoperator.ClusterNetworkOperatorComponent)
	r.Register("Networking / ingress-node-firewall", &networkingingressnodefirewall.IngressNodeFirewallComponent)
	r.Register("Networking / kubernetes-nmstate", &networkingkubernetesnmstate.KubernetesNmstateComponent)
	r.Register("Networking / kubernetes-nmstate-operator", &networkingkubernetesnmstateoperator.KubernetesNmstateOperatorComponent)
	r.Register("Networking / kuryr", &networkingkuryr.KuryrComponent)
	r.Register("Networking / multus", &networkingmultus.MultusComponent)
	r.Register("Networking / network-tools", &networkingnetworktools.NetworkToolsComponent)
	r.Register("Networking / networking-console-plugin", &networkingnetworkingconsoleplugin.NetworkingConsolePluginComponent)
	r.Register("Networking / nmstate-console-plugin", &networkingnmstateconsoleplugin.NmstateConsolePluginComponent)
	r.Register("Networking / openshift-sdn", &networkingopenshiftsdn.OpenshiftSdnComponent)
	r.Register("Networking / ovn-kubernetes", &networkingovnkubernetes.OvnKubernetesComponent)
	r.Register("Networking / ptp", &networkingptp.PtpComponent)
	r.Register("Networking / router", &networkingrouter.RouterComponent)
	r.Register("Networking / runtime-cfg", &networkingruntimecfg.RuntimeCfgComponent)
	r.Register("Node / CPU manager", &nodecpumanager.CPUManagerComponent)
	r.Register("Node / CRI-O", &nodecrio.CRIOComponent)
	r.Register("Node / Device Manager", &nodedevicemanager.DeviceManagerComponent)
	r.Register("Node / Instaslice-operator", &nodeinstasliceoperator.InstasliceOperatorComponent)
	r.Register("Node / Kubelet", &nodekubelet.KubeletComponent)
	r.Register("Node / Memory manager", &nodememorymanager.MemoryManagerComponent)
	r.Register("Node / Node Problem Detector", &nodenodeproblemdetector.NodeProblemDetectorComponent)
	r.Register("Node / Numa aware Scheduling", &nodenumaawarescheduling.NumaAwareSchedulingComponent)
	r.Register("Node / Pod resource API", &nodepodresourceapi.PodResourceAPIComponent)
	r.Register("Node / Topology manager", &nodetopologymanager.TopologyManagerComponent)
	r.Register("Node Feature Discovery Operator", &nodefeaturediscoveryoperator.NodeFeatureDiscoveryOperatorComponent)
	r.Register("Node Maintenance Operator", &nodemaintenanceoperator.NodeMaintenanceOperatorComponent)
	r.Register("Node Tuning Operator", &nodetuningoperator.NodeTuningOperatorComponent)
	r.Register("Node-observability-operator", &nodeobservabilityoperator.NodeObservabilityOperatorComponent)
	r.Register("OLM / OperatorHub", &olmoperatorhub.OperatorHubComponent)
	r.Register("OLM / Registry", &olmregistry.RegistryComponent)
	r.Register("OLM", &olm.OLMComponent)
	r.Register("OPCT / CLI", &opctcli.CLIComponent)
	r.Register("OPCT / Other", &opctother.OtherComponent)
	r.Register("OPCT / Results / General", &opctresultsgeneral.GeneralComponent)
	r.Register("OPCT / Results / VCSP", &opctresultsvcsp.VCSPComponent)
	r.Register("Observability UI", &observabilityui.ObservabilityUIComponent)
	r.Register("OpenShift Update Service / operand", &openshiftupdateserviceoperand.OperandComponent)
	r.Register("OpenShift Update Service / operator", &openshiftupdateserviceoperator.OperatorComponent)
	r.Register("Operator SDK", &operatorsdk.OperatorSDKComponent)
	r.Register("Performance Addon Operator", &performanceaddonoperator.PerformanceAddonOperatorComponent)
	r.Register("Pod Autoscaler", &podautoscaler.PodAutoscalerComponent)
	r.Register("Poison Pill Operator", &poisonpilloperator.PoisonPillOperatorComponent)
	r.Register("RHCOS", &rhcos.RHCOSComponent)
	r.Register("RHMI Monitoring", &rhmimonitoring.RHMIMonitoringComponent)
	r.Register("Registry Console", &registryconsole.RegistryConsoleComponent)
	r.Register("Release", &release.ReleaseComponent)
	r.Register("Samples Operator", &samplesoperator.SamplesOperatorComponent)
	r.Register("Security Profiles Operator", &securityprofilesoperator.SecurityProfilesOperatorComponent)
	r.Register("Security", &security.SecurityComponent)
	r.Register("Service Binding", &servicebinding.ServiceBindingComponent)
	r.Register("Service Broker", &servicebroker.ServiceBrokerComponent)
	r.Register("Service Catalog", &servicecatalog.ServiceCatalogComponent)
	r.Register("Special Resource Operator", &specialresourceoperator.SpecialResourceOperatorComponent)
	r.Register("Storage / Kubernetes External Components", &storagekubernetesexternalcomponents.KubernetesExternalComponentsComponent)
	r.Register("Storage / Kubernetes", &storagekubernetes.KubernetesComponent)
	r.Register("Storage / Local Storage Operator", &storagelocalstorageoperator.LocalStorageOperatorComponent)
	r.Register("Storage / OpenStack CSI Drivers", &storageopenstackcsidrivers.OpenStackCSIDriversComponent)
	r.Register("Storage / Operators", &storageoperators.OperatorsComponent)
	r.Register("Storage / Shared Resource CSI Driver", &storagesharedresourcecsidriver.SharedResourceCSIDriverComponent)
	r.Register("Storage / oVirt CSI Driver", &storageovirtcsidriver.OVirtCSIDriverComponent)
	r.Register("Storage", &storage.StorageComponent)
	r.Register("TALM Operator", &talmoperator.TALMOperatorComponent)
	r.Register("Telco Performance", &telcoperformance.TelcoPerformanceComponent)
	r.Register("Telemeter", &telemeter.TelemeterComponent)
	r.Register("Templates", &templates.TemplatesComponent)
	r.Register("Test Framework / OpenStack", &testframeworkopenstack.OpenStackComponent)
	r.Register("Test Framework", &testframework.TestFrameworkComponent)
	r.Register("Test Infrastructure", &testinfrastructure.TestInfrastructureComponent)
	r.Register("Unknown", &unknown.UnknownComponent)
	r.Register("Windows Containers", &windowscontainers.WindowsContainersComponent)
	r.Register("apiserver-auth", &apiserverauth.ApiserverAuthComponent)
	r.Register("cert-manager", &certmanager.CertManagerComponent)
	r.Register("cli-manager", &climanager.CliManagerComponent)
	r.Register("confidential-compute-attestation", &confidentialcomputeattestation.ConfidentialComputeAttestationComponent)
	r.Register("config-operator", &configoperator.ConfigOperatorComponent)
	r.Register("crc", &crc.CrcComponent)
	r.Register("descheduler", &descheduler.DeschedulerComponent)
	r.Register("ibm-roks-toolkit", &ibmrokstoolkit.IbmRoksToolkitComponent)
	r.Register("kmm", &kmm.KmmComponent)
	r.Register("kube-apiserver", &kubeapiserver.KubeApiserverComponent)
	r.Register("kube-controller-manager", &kubecontrollermanager.KubeControllerManagerComponent)
	r.Register("kube-scheduler", &kubescheduler.KubeSchedulerComponent)
	r.Register("kube-storage-version-migrator", &kubestorageversionmigrator.KubeStorageVersionMigratorComponent)
	r.Register("none", &none.NoneComponent)
	r.Register("oauth-apiserver", &oauthapiserver.OauthApiserverComponent)
	r.Register("oauth-proxy", &oauthproxy.OauthProxyComponent)
	r.Register("oc / cluster-compare", &occlustercompare.ClusterCompareComponent)
	r.Register("oc / oc-mirror", &ococmirror.OcMirrorComponent)
	r.Register("oc / update", &ocupdate.UpdateComponent)
	r.Register("oc", &oc.OcComponent)
	r.Register("oc-compliance", &occompliance.OcComplianceComponent)
	r.Register("openshift-apiserver", &openshiftapiserver.OpenshiftApiserverComponent)
	r.Register("openshift-controller-manager / apps", &openshiftcontrollermanagerapps.AppsComponent)
	r.Register("openshift-controller-manager / build", &openshiftcontrollermanagerbuild.BuildComponent)
	r.Register("openshift-controller-manager / controller-manager", &openshiftcontrollermanagercontrollermanager.ControllerManagerComponent)
	r.Register("route-controller-manager", &routecontrollermanager.RouteControllerManagerComponent)
	r.Register("run-once-duration-override", &runoncedurationoverride.RunOnceDurationOverrideComponent)
	r.Register("sandboxed-containers", &sandboxedcontainers.SandboxedContainersComponent)
	r.Register("secondary-scheduler-operator", &secondaryscheduleroperator.SecondarySchedulerOperatorComponent)
	r.Register("service-ca", &serviceca.ServiceCaComponent)
	r.Register("bpfman", &bpfman.BpfmanComponent)
	r.Register("Cloud Compute / vSphere Provider", &cloudcomputevsphereprovider.VSphereProviderComponent)
	r.Register("HyperShift / ARO", &hypershiftaro.AROComponent)
	r.Register("HyperShift / OpenStack", &hypershiftopenstack.OpenStackComponent)
	r.Register("HyperShift / ROSA", &hypershiftrosa.ROSAComponent)
	r.Register("insights-runtime-extractor", &insightsruntimeextractor.InsightsRuntimeExtractorComponent)
	r.Register("Installer / vSphere", &installervsphere.VSphereComponent)
	r.Register("kSAN Storage", &ksanstorage.KSANStorageComponent)
	r.Register("Lightspeed", &lightspeed.LightspeedComponent)
	r.Register("Networking / FRR-K8s", &networkingfrrk8s.FRRK8sComponent)
	r.Register("oc / node-image", &ocnodeimage.NodeImageComponent)
	r.Register("oc-mirror", &ocmirror.OcMirrorComponent)
	r.Register("TwoNode / Arbiter", &twonodearbiter.ArbiterComponent)
	// New components go here

	return &r
}

func (r *Registry) Register(name string, component v1.Component) {
	if r.Components == nil {
		r.Components = make(map[string]v1.Component)
	}

	r.Components[name] = component
}

func (r *Registry) GetForJiraComponent(name string) v1.Component {
	for _, c := range r.Components {
		for _, j := range c.JiraComponents() {
			if j == name {
				return c
			}
		}
	}

	return nil
}

func (r *Registry) Deregister(name string) {
	delete(r.Components, name)
}
