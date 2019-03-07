package myecharts

func getDefaltOption() EcOption {
	return EcOption{
		XAxis:  EcxAxis{Type: "category"},
		YAxis:  EcxAxis{Type: "value"},
		Legend: EcLegend{Show: true},
	}
}

// 设置显示x轴还是y轴，且设置数据内容
//data 为轴坐标内容
func (ep *EcOption) OnSetAxis(data []string, axisLabel AxisLabel, isy bool) {
	V := EcxAxis{Type: "value", AxisLabel: axisLabel}
	k := EcxAxis{Type: "category", Data: data}
	if isy {
		ep.YAxis = k
		ep.XAxis = V
	} else {
		ep.XAxis = k
		ep.YAxis = V
	}
}

//图例组件。(顶部分类器)
func (ep *EcOption) OnSetLegend(data []string) {
	ep.Legend.Data = data
	ep.Legend.Show = true
}

//添加一个数据
func (ep *EcOption) OnAddOneSeries(name string, data interface{}) {
	var tmp = EcSeries{Name: name, Type: ep._type, Data: data}
	ep.Series = append(ep.Series, tmp)
}

//添加一个数据
func (ep *EcOption) OnAddOneSeriesInt(name string, data []int) {
	var tmp []interface{}
	for _, v := range data {
		tmp = append(tmp, float32(v))
	}
	ep.OnAddOneSeries(name, tmp)
}

//
func (ep *EcOption) SetEcToolTip(ecToolTip EcToolTip) {
	//EcToolTip{Show: true, Trigger: "axis", AxisPointer: EcAxisPointer{Type: "shadow"}}
	ep.Tooltip = ecToolTip
}
