package myecharts

//获取线图基础数据
/*
	t_str
		line:线
		bar:柱
		pie:饼
		scatter:散点（气泡）	effectScatter:带有涟漪特效动画的散点	radar:雷达
		tree:树	treemap:面积	sunburst:旭日图	boxplot:箱形图	candlestick:K线图	sankey:桑基图
		heatmap:热力图	map:地图	parallel: 平行坐标系的系列	lines:线图	graph:关系图
		funnel:漏斗图	gauge:仪表盘	pictorialBar:象形柱图	themeRiver:主题河流	custom:自定义系列
*/
func OnGetBaseInfo(t_str string) (ec EcOption) {
	ec = getDefaltOption()
	ec._type = t_str
	return ec
}

/*
	获取基础柱状图
*/
func OnGetBarInfo(ecToolTip EcToolTip) (ec EcOption) {
	ec = OnGetBaseInfo("bar")
	return ec
}
