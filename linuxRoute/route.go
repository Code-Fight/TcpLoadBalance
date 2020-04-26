package linuxRoute

import (
	"TcpLoadBalance/units"
	log "github.com/Code-Fight/golog"
	"k8s.io/kubernetes/pkg/util/ipset"
	"runtime"
	"sync"
)
import "k8s.io/kubernetes/pkg/util/iptables"
import utilexec "k8s.io/utils/exec"

type Route struct {

}

var (
	iptRunner =iptables.New(utilexec.New(),iptables.ProtocolIpv4)
	ipsRunner =ipset.New(utilexec.New())
	Ipsets = sync.Map{}
) 

func NewRoute() *Route{

	if runtime.GOOS!="linux"{
		panic("linuxRoute only suports linux")
	}

	backend:=units.GetBackend()

	
	//iptables -t nat -A POSTROUTING -p tcp  --dport 22048 -j MASQUERADE
	_,err := iptRunner.EnsureRule(
		iptables.Append,
		iptables.TableNAT,
		iptables.ChainPostrouting,
		"-p","tcp","--dport","22048","-j","MASQUERADE")

	if err !=nil{
		panic(err)
		return nil
	}
	//创建所有backend的转发配置
	for  key,_:=range backend{
		//创建ipset
		tIpset := &ipset.IPSet{Name:key,SetType:ipset.HashIPPort,	MaxElem:60000}
		ipsErr:=ipsRunner.CreateSet(tIpset,	true)
		if ipsErr!=nil{
			panic(ipsErr)
		}
		//save the ipset obj to Ipset map
		Ipsets.LoadOrStore(key,tIpset)
		//iptables -t nat -A PREROUTING -p tcp -m set  --match-set test src,dst -j DNAT --to-destination 192.168.2.230:22048
		_,err := iptRunner.EnsureRule(
			iptables.Append,
			iptables.TableNAT,
			iptables.ChainPrerouting,
			"-p","tcp","-m","set","--match-set",key,"src,dst","-j","DNAT","-to-destination",key)

		if err !=nil{
			panic(err)
			return nil
		}
	}
	return &Route{}
}
// IpsetAdd add src to ipset
func(r *Route) IpsetAdd(set string,src string) bool {
	ips ,ok:=Ipsets.Load(set)
	if ok{
		 ipsObj,ipsOK:=ips.(*ipset.IPSet)
		 if ipsOK{
			 addErr:= ipsRunner.AddEntry(src,ipsObj,true)
			 if addErr ==nil{

			 	return true
			 }else {
			 	log.Error("add ipset error:",addErr.Error())
			 	return false
			 }
		 }
	}
	return false
}
// IpsetDel del src from ipset
func(r *Route) IpsetDel(set string,src string )  bool{
	delErr:= ipsRunner.DelEntry(src,set)
	if delErr ==nil{
		return true
	}else {
		return false
	}
}
