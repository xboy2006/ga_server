package cluster

import (
	"GAServer/log"
	"GAServer/service"
	_ "encoding/json"
	"gameproto/msgs"
)

//注册到center
func RegServerToCenter(s *service.ServiceData, values []*msgs.ServiceValue) bool {
	log.Info("%v reg to center...", s.Name)

	msg := msgs.AddService{
		ServiceName: s.Name,
		ServiceType: s.TypeName,
		Pid:         s.GetPID(),
		Values:      values}
	_, err := GetServicePID("center").Ask(&msg)
	if err != nil {
		log.Error("%v reg to center fail,%v", s.Name, err)
		//if err.Error() == "timeout" {
		//DisconnectService("center")
		//}
		//重连
		return false
	}
	log.Info("%v reg to center OK!", s.Name)
	return true
}

func RegServerWork(s *service.ServiceData, values []*msgs.ServiceValue) {
	go func() {
		for {
			if RegServerToCenter(s, values) {
				break
			}
		}
	}()

}

func UpdateServiceLoad(serviceName string, load uint32, state msgs.ServiceState) {
	log.Debug("%v UpdateServiceLoad %v-%v", serviceName, load, state)

	msg := msgs.UploadService{
		ServiceName: serviceName,
		Load:        load,
		State:       state}
	GetServicePID("center").Tell(&msg)
}
