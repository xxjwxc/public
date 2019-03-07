package myecharts

//echarts option 根数据
type EcOption struct {
	Color []string `json:"color,omitempty"` //调色盘颜色列表。如果系列没有设置颜色，则会依次循环从该列表中取颜色作为系列颜色。
	//标题组件，包含主标题和副标题。
	Title   EcTitle    `json:"title,omitempty"`
	Tooltip EcToolTip  `json:"tooltip,omitempty"`
	Grid    EcGrid     `json:"grid,omitempty"`
	XAxis   EcxAxis    `json:"xAxis,omitempty"`   //直角坐标系 grid 中的 x 轴
	YAxis   EcxAxis    `json:"yAxis,omitempty"`   //直角坐标系 grid 中的 y 轴
	Legend  EcLegend   `json:"legend,omitempty"`  //图例组件。
	Series  []EcSeries `json:"series,omitempty"`  //数据内容
	Toolbox Toolbox    `json:"toolbox,omitempty"` //
	_type   string
}

//系列列表。每个系列通过 type 决定自己的图表类型
type EcSeries struct {
	Name string      `json:"name,omitempty"` //系列名称
	Type string      `json:"type,omitempty"` //line:线 bar:柱状图
	Data interface{} `json:"data,omitempty"` //系列中的数据内容数组
}

//图例组件。
type EcLegend struct {
	Show bool     `json:"show,omitempty"`
	Data []string `json:"data,omitempty"` //图例的数据数组
}

//直角坐标系内绘图网格
type EcxAxis struct {
	Type      string    `json:"type,omitempty"`
	Data      []string  `json:"data,omitempty"`
	AxisLabel AxisLabel `json:"axisLabel,omitempty"`
}

//
type AxisLabel struct {
	Formatter string `json:"formatter,omitempty"`
}

//直角坐标系内绘图网格
type EcGrid struct {
	Left         string `json:"left,omitempty"` //组件离容器左侧的距离。
	Right        string `json:"right,omitempty"`
	Bottom       string `json:"bottom,omitempty"`
	ContainLabel bool   `json:"containLabel,omitempty"` //grid 区域是否包含坐标轴的刻度标签
}

//提示框组件。
type EcToolTip struct {
	Show bool `json:"show,omitempty"`
	/*
		触发类型。
		'item' 数据项图形触发，主要在散点图，饼图等无类目轴的图表中使用。
		'axis' 坐标轴触发，主要在柱状图，折线图等会使用类目轴的图表中使用。
		'none' 什么都不触发。
	*/
	Trigger string `json:"trigger,omitempty"`
	//坐标轴指示器配置项。
	AxisPointer EcAxisPointer `json:"axisPointer,omitempty"`
}

//坐标轴指示器配置项。
type EcAxisPointer struct {
	/*
		'line' 直线指示器
		'shadow' 阴影指示器
		'none' 无指示器
		'cross' 十字准星指示器
	*/
	Type string `json:"type,omitempty"`
}

//标题组件，包含主标题和副标题。
type EcTitle struct {
	Text    string `json:"text,omitempty"`    //主标题文本，支持使用 \n 换行。
	SubText string `json:"subtext,omitempty"` //副标题文本，支持使用 \n 换行。
}

//
type Toolbox struct {
	Show    bool    `json:"show,omitempty"`    //
	Feature Feature `json:"feature,omitempty"` //
}

//
type Feature struct {
}
