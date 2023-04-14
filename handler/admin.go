package handler

import (
	R "Project/MyProject/response"
	"Project/MyProject/services"
	"Project/MyProject/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

/*
	管理员界面需要有的功能
1. 路线添加、查询所有、删除
2. 列车添加、查询所有、删除
3. 售票信息添加、条件查询、查询所有、删除

*/

type AdminHandler struct {
	RouteService  services.RouteServiceImplement
	TrainService  services.TrainServiceImplement
	TicketService services.TicketServiceImplement
}

// 管理员界面所有服务接口
type AdminImplement interface {
	RouteImplement
	TrainImplement
	TicketImplement
}

// 路线服务接口
type RouteImplement interface {
	AddRoute(ctx *gin.Context)
	GetRoute(ctx *gin.Context)
	GetRouteByID(ctx *gin.Context)
	ListRoutes(*gin.Context)
	DeleteRoute(ctx *gin.Context)
}

// 列车服务接口
type TrainImplement interface {
	AddTrain(ctx *gin.Context)
	GetTrain(ctx *gin.Context)
	ListTrains(ctx *gin.Context)
	DeleteTrain(ctx *gin.Context)
}

// 售票服务接口
type TicketImplement interface {
	AddTicket(ctx *gin.Context)
	GetTicket(ctx *gin.Context)
	ListTicketsOnSale(ctx *gin.Context)
	ListTickets(ctx *gin.Context)
	ListAllTicketsOnSale(ctx *gin.Context)
	ListAllTickets(ctx *gin.Context)
	StopSell(ctx *gin.Context)
}

func NewAdminHandler(routeService services.RouteServiceImplement,
	trainService services.TrainServiceImplement,
	ticketService services.TicketServiceImplement) AdminImplement {
	return &AdminHandler{
		RouteService:  routeService,
		TrainService:  trainService,
		TicketService: ticketService,
	}
}

func (handler *AdminHandler) AddRoute(ctx *gin.Context) {
	var addRouteService services.AddRouteService
	if err := ctx.ShouldBind(&addRouteService); err == nil {
		if err := handler.RouteService.AddRoute(&addRouteService); err == nil {
			R.Ok(ctx, "添加路线成功", nil)
		} else {
			R.Error(ctx, err.Error(), nil)
		}
	} else {
		R.Response(ctx, http.StatusBadRequest, http.StatusBadRequest, "参数错误", err.Error())
	}
}

func (handler *AdminHandler) GetRoute(ctx *gin.Context) {
	var getRouteService services.GetRouteService
	if err := ctx.ShouldBind(&getRouteService); err == nil {
		if route, err := handler.RouteService.GetRoute(&getRouteService); err == nil {
			R.Ok(ctx, "查询路线成功", gin.H{
				"route": route,
			})
		} else {
			R.Error(ctx, err.Error(), nil)
		}
	} else {
		R.Response(ctx, http.StatusBadRequest, http.StatusBadRequest, "参数错误", err.Error())
	}
}

func (handler *AdminHandler) GetRouteByID(ctx *gin.Context) {
	if routeID, err := strconv.Atoi(ctx.Param("route_id")); err == nil {
		if route, err := handler.RouteService.GetRouteByID(routeID); err == nil {
			R.Ok(ctx, "查询路线成功", gin.H{
				"route": route,
			})
		} else {
			R.Error(ctx, err.Error(), nil)
		}
	} else {
		R.Response(ctx, http.StatusBadRequest, http.StatusBadRequest, "参数错误", err.Error())
	}
}

func (handler *AdminHandler) ListRoutes(ctx *gin.Context) {
	if routes, err := handler.RouteService.ListRoutes(); err == nil {
		R.Ok(ctx, "查询路线成功", gin.H{
			"routes": routes,
		})
	} else {
		R.Error(ctx, err.Error(), nil)
	}
}

func (handler *AdminHandler) DeleteRoute(ctx *gin.Context) {
	if routeID, err := strconv.Atoi(ctx.Param("route_id")); err == nil {
		_, err = handler.RouteService.GetRouteByID(routeID)
		if err != nil {
			R.Error(ctx, err.Error(), nil)
			return
		}
		if err = handler.RouteService.DeleteRoute(routeID); err == nil {
			R.Ok(ctx, "删除路线成功", nil)
		} else {
			R.Error(ctx, err.Error(), nil)
		}
	} else {
		R.Response(ctx, http.StatusBadRequest, http.StatusBadRequest, "参数错误", err.Error())
	}
}

func (handler *AdminHandler) AddTrain(ctx *gin.Context) {
	var addTrainService services.AddTrainService
	if err := ctx.ShouldBind(&addTrainService); err == nil {
		if err := handler.TrainService.AddTrain(&addTrainService); err == nil {
			R.Ok(ctx, "添加列车成功", nil)
		} else {
			R.Error(ctx, err.Error(), nil)
		}
	} else {
		R.Response(ctx, http.StatusBadRequest, http.StatusBadRequest, "参数错误", err.Error())
	}
}

func (handler *AdminHandler) GetTrain(ctx *gin.Context) {
	trainID := ctx.Param("train_id")
	if train, err := handler.TrainService.GetTrain(trainID); err == nil {
		R.Ok(ctx, "查询列车信息成功", gin.H{
			"train": train,
		})
	} else {
		R.Error(ctx, err.Error(), nil)
	}
}

func (handler *AdminHandler) ListTrains(ctx *gin.Context) {
	if trains, err := handler.TrainService.ListTrains(); err == nil {
		R.Ok(ctx, "查询列车信息成功", gin.H{
			"trains": trains,
		})
	} else {
		R.Error(ctx, err.Error(), nil)
	}
}

func (handler *AdminHandler) DeleteTrain(ctx *gin.Context) {
	trainID := ctx.Param("train_id")
	_, err := handler.TrainService.GetTrain(trainID)
	if err != nil {
		R.Error(ctx, err.Error(), nil)
		return
	}
	if err = handler.TrainService.DeleteTrain(trainID); err == nil {
		R.Ok(ctx, "删除列车信息成功", nil)
	} else {
		R.Error(ctx, err.Error(), nil)
	}
}

func (handler *AdminHandler) AddTicket(ctx *gin.Context) {
	var addTicketService services.AddTicketService
	if err := ctx.ShouldBind(&addTicketService); err == nil {
		// 先通过trainID和routeID获取列车和路线信息，再计算车票需要的其他信息
		train, err := handler.TrainService.GetTrain(addTicketService.TrainID)
		if err != nil {
			R.Error(ctx, err.Error(), nil)
			return
		}
		route, err := handler.RouteService.GetRouteByID(addTicketService.RouteID)
		if err != nil {
			R.Error(ctx, err.Error(), nil)
			return
		}
		// 需要通过 departureTime 计算 arrivalTime
		departureTime, err := utils.ParseStringToTime(addTicketService.DepartureTime)
		if err != nil {
			R.Response(ctx, http.StatusBadRequest, http.StatusBadRequest, "参数错误", err.Error())
			return
		}
		addTicketService.Stock = train.Seats
		addTicketService.Start = route.Start
		addTicketService.End = route.End
		expectedHour := time.Duration(route.Length / train.Speed)        // 预计需要 expectdHour 小时
		route.Length = route.Length % train.Speed                        // 经过 expectedHour 后剩余的路程
		expectedMinute := time.Duration(route.Length * 60 / train.Speed) // 预计需要 expectedMinute 分钟
		addTicketService.Duration = fmt.Sprintf("%d时%d分", expectedHour, expectedMinute)

		addTicketService.ArrivalTime = departureTime.Add(expectedHour * time.Hour).Add(expectedMinute * time.Minute)
		if err := handler.TicketService.AddTicket(&addTicketService); err == nil {
			R.Ok(ctx, "车票发售成功", nil)
		} else if err == R.ErrInvalidParam {
			R.Response(ctx, http.StatusBadRequest, http.StatusBadRequest, "参数错误", err.Error())
		} else {
			R.Error(ctx, err.Error(), nil)
		}

	} else {
		R.Response(ctx, http.StatusBadRequest, http.StatusBadRequest, "参数错误", err.Error())
	}
}

func (handler *AdminHandler) GetTicket(ctx *gin.Context) {
	if ticketID, err := strconv.Atoi(ctx.Param("ticket_id")); err == nil {
		if ticket, err := handler.TicketService.GetTicket(ticketID); err == nil {
			R.Ok(ctx, "查询售票信息成功", gin.H{
				"ticket": ticket,
			})
		} else {
			R.Error(ctx, err.Error(), nil)
		}
	} else {
		R.Response(ctx, http.StatusBadRequest, http.StatusBadRequest, "参数错误", err.Error())
	}
}

func (handler *AdminHandler) ListTicketsOnSale(ctx *gin.Context) {
	var listTicketOnSaleService services.ListTicketsOnSaleService
	if err := ctx.ShouldBind(&listTicketOnSaleService); err == nil {
		if tickets, err := handler.TicketService.ListTicketsOnSale(&listTicketOnSaleService); err == nil {
			R.Ok(ctx, "查询该路线当前在售车票信息成功", gin.H{
				"tickets": tickets,
			})
		} else {
			R.Error(ctx, err.Error(), nil)
		}
	} else {
		R.Response(ctx, http.StatusBadRequest, http.StatusBadRequest, "参数错误", err.Error())
	}
}

func (handler *AdminHandler) ListTickets(ctx *gin.Context) {
	var listTicketService services.ListTicketsService
	if err := ctx.ShouldBind(&listTicketService); err == nil {
		if tickets, err := handler.TicketService.ListTickets(&listTicketService); err == nil {
			R.Ok(ctx, "查询该路线售票信息成功", gin.H{
				"tickets": tickets,
			})
		} else {
			R.Error(ctx, err.Error(), nil)
		}
	} else {
		R.Response(ctx, http.StatusBadRequest, http.StatusBadRequest, "参数错误", err.Error())
	}
}

func (handler *AdminHandler) ListAllTicketsOnSale(ctx *gin.Context) {
	tickets, err := handler.TicketService.GetAllTicketsOnSale()
	if err == nil {
		R.Ok(ctx, "查询当前在售车票信息成功", gin.H{
			"tickets": tickets,
		})
	} else {
		R.Error(ctx, err.Error(), nil)
	}
}

func (handler *AdminHandler) ListAllTickets(ctx *gin.Context) {
	tickets, err := handler.TicketService.GetAllTickets()
	if err == nil {
		R.Ok(ctx, "查询所有车票信息成功", gin.H{
			"tickets": tickets,
		})
	} else {
		R.Error(ctx, err.Error(), nil)
	}
}

func (handler *AdminHandler) StopSell(ctx *gin.Context) {
	var stopSellTicketService services.StopSellTicketService
	if err := ctx.ShouldBind(&stopSellTicketService); err == nil {
		if err = handler.TicketService.StopSellTicket(&stopSellTicketService); err == nil {
			R.Ok(ctx, "停止售票成功", nil)
		} else {
			R.Error(ctx, err.Error(), nil)
		}
	} else {
		R.Response(ctx, http.StatusBadRequest, http.StatusBadRequest, "参数错误", err.Error())
	}
}
