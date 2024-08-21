package server

import (
	"duck/kernel/extractors"
	"duck/kernel/log"
	"duck/kernel/manage"
	"duck/kernel/server/core"
	"strconv"
)

var GetRoutes map[string]core.Handler
var PostRoutes map[string]core.Handler

func InitRoutes(manager *manage.Manager) {
	GetRoutes = make(map[string]core.Handler)
	PostRoutes = make(map[string]core.Handler)
	
	GetRoutes["/"] = HomeHandler

	GetRoutes["/tasks"] = func(context *core.Context) {
		tasks := manager.GetAllTasks()
		context.Json(core.J{"tasks": tasks})
	}
	
	// GetRoutes["/taskSets"] = func(context *core.Context) {
	// 	taskSets := schedule.GetAllTaskSets()
	// 	context.Json(core.J{"taskSets": taskSets})
	// }

	GetRoutes["/task/:taskNo"] = func(context *core.Context) {
		taskNo, err := strconv.Atoi(context.RouteParams["taskNo"])

		if err != nil {
			log.Error(err)
		}
		
		task := manager.GetTaskByNo(int64(taskNo))
		context.Json(core.J{"task": task})
	}

	// GetRoutes["/taskSet/:taskSetNo"] = func(context *core.Context) {
	// 	taskSetNo, err := strconv.Atoi(context.RouteParams["taskSetNo"])

	// 	if err != nil {
	// 		log.Error(err)
	// 	}
		
	// 	taskSet := schedule.GetTaskByNo(int64(taskSetNo))
	// 	context.Json(core.J{"taskSet": taskSet})
	// }

	GetRoutes["/extract"] = func(context *core.Context) {
		url := context.GetParam("url")
		extractor := context.GetParam("extractor")
		options, optionTypes, err := extractors.Extract(url, extractor)
		if err != nil {
			context.Json(core.J{
				"errMessage": err.Error(),
			})
		} else {			
			context.Json(core.J{
				"options": options,
				"optionTypes": optionTypes,
			})
		}
	}

	PostRoutes["/download"] = func(context *core.Context) {
		options := context.BodyParams()
		url := options["Url"]
		extractor := options["Extractor"]

		task, err := extractors.NewTask(url.(string), extractor.(string), options)
		if err != nil {
			context.Json(core.J{
				"errMessage": err.Error(),
			})
			return
		}
		manager.AddTaskToQueue(task)
		context.Json(core.J{
			"status": "succeed",
		})
	}

	PostRoutes["/task/:taskNo/update"] = func(context *core.Context) {
		
	}

	PostRoutes["/taskSet/:taskSetNo/update"] = func(context *core.Context) {
		
	}

	PostRoutes["/task/:taskNo/delete"] = func(context *core.Context) {
		
	}
	
	PostRoutes["/taskSet/:taskSetNo/delete"] = func(context *core.Context) {
		
	}
}

func HomeHandler(context *core.Context) {
	context.HTML("<strong>hello!</strong>")
}
