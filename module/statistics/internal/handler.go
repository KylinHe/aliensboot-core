package internal

import (
	"github.com/KylinHe/aliensboot-core/log"
	"github.com/KylinHe/aliensboot-core/module/statistics/elastics"
	"github.com/KylinHe/aliensboot-core/module/statistics/model"
	"github.com/sirupsen/logrus"
	//"github.com/KylinHe/aliensboot-core/cluster/center"
	"github.com/KylinHe/aliensboot-core/module/statistics/conf"
	"github.com/KylinHe/aliensboot-core/module/statistics/constant"
	"github.com/KylinHe/aliensboot-core/task"
)

var esHandler = elastics.NewESHandler(conf.Game)

//服务调用统计信息、一分钟一次
var serviceStatistics = make(map[string]map[int32]*model.CallInfo) //服务名 - 服务编号 - 调用信息

var serviceFields = logrus.Fields{}

var onlineFields = logrus.Fields{}

func init() {
	skeleton.RegisterChanRPC(constant.INTERNAL_STATISTICS_SERVICE_CALL, handleServiceStatic)
	skeleton.RegisterChanRPC(constant.INTERNAL_STATISTICS_ONLINE, handleOnlineStatic)
	cron, err := task.NewCronExpr("*/1 * * * *")
	if err != nil {
		log.Error("init service statistics timer error : %v", err)
	}

	//每天凌晨12点执行一次
	dayCron, err := task.NewCronExpr("0 0 * * *")
	if err != nil {
		log.Error("init dump timer error : %v", err)
	}
	skeleton.CronFunc(dayCron, esHandler.UpdateDayPrefix)
	skeleton.CronFunc(cron, handleTimer)
}

func handleOnlineStatic(args []interface{}) {
	userCount := args[0].(int)     //用户数量
	visitorCount := args[1].(int)  //空连接数量
	onlineFields["node"] = "node1" //center.ClusterCenter.GetNodeID()
	onlineFields["u_count"] = userCount
	onlineFields["v_count"] = visitorCount
	esHandler.HandleDayESLog("online", "", onlineFields)
}

//处理服务信息统计
func handleServiceStatic(args []interface{}) {
	service := args[0].(string)   //服务名称
	serviceNo := args[1].(int32)  //服务处理编号
	interval := args[2].(float64) //服务处理时间间隔
	callInfos := serviceStatistics[service]
	if callInfos == nil {
		callInfos = make(map[int32]*model.CallInfo)
		serviceStatistics[service] = callInfos
	}

	callInfo := callInfos[serviceNo]
	if callInfo == nil {
		callInfo = &model.CallInfo{}
		callInfos[serviceNo] = callInfo
	}
	callInfo.AddCall(interval)
}

func handleTimer() {
	for service, callInfos := range serviceStatistics {
		for serviceNo, callInfo := range callInfos {
			serviceFields["service"] = service
			serviceFields["node"] = "node1" //center.ClusterCenter.GetNodeID()
			serviceFields["no"] = serviceNo
			result, count, avg := callInfo.DumpData()
			if !result {
				continue
			}
			serviceFields["count"] = count
			serviceFields["avg"] = avg
			esHandler.HandleDayESLog("service", "", serviceFields)
		}
	}
}
