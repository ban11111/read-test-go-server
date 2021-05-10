package exporter

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"mime"
	"read-test-server/common"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type ExcelExporter struct{}

func (e *ExcelExporter) Ext() string {
	return "xlsx"
}

func (e *ExcelExporter) Export(data interface{}) (*ReaderWrapper, error) {
	f := excelize.NewFile()

	title, widths := e.makeTitle(data)
	if err := f.SetSheetRow("Sheet1", "A1", &title); err != nil {
		common.Log.Error("ExcelExporter.SetSheetRow()", zap.Error(err))
		return nil, err
	}
	for i:=0;i<len(widths);i++ {
		col := common.ConvertNumToCols(i + 1)
		_ = f.SetColWidth("Sheet1", col, col, widths[i])
	}

	contentRows, err := e.makeContent(data)
	if err != nil {
		common.Log.Error("ExcelExporter.makeContent()", zap.Error(err))
		return nil, err
	}

	for i := 0; i < len(contentRows); i++ {
		if err := f.SetSheetRow("Sheet1", "A"+strconv.Itoa(i+2), &contentRows[i]); err != nil {
			common.Log.Error("ExcelExporter.SetSheetRow()", zap.Error(err), zap.Int("i", i))
			return nil, err
		}
	}

	buffer, err := f.WriteToBuffer()
	if err != nil {
		common.Log.Error("ExcelExporter.WriteToBuffer()", zap.Error(err))
		return nil, err
	}
	contentType := mime.TypeByExtension(".xlsx")
	if contentType == "" {
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	}
	return &ReaderWrapper{
		Reader:      buffer,
		Len:         buffer.Len(),
		ContentType: contentType,
	}, nil
}

func (e *ExcelExporter) makeTitle(data interface{}) ([]string, []float64) {
	rt := reflect.TypeOf(data)
	for rt.Kind() != reflect.Struct {
		rt = rt.Elem()
	}

	var titles, widths = make([]string, 0, rt.NumField()), make([]float64, 0, rt.NumField())
	common.EnumAnObjFieldNames(rt, func(f reflect.StructField) {
		titleName := f.Tag.Get("csv")
		if titleName == "" {
			titleName = f.Name
		} else if titleName == "-" {
			return
		}
		width := float64(20)
		widthTag := f.Tag.Get("width")
		if widthTag != "" {
			width, _ = strconv.ParseFloat(widthTag, 64)
		}
		widths = append(widths, width)
		titles = append(titles, titleName)
	})
	return titles, widths
}

func (e *ExcelExporter) makeContent(datas interface{}) (contentRows [][]interface{}, err error) {
	var sliceType string
	var length int

	rTitle := reflect.TypeOf(datas)
	for rTitle.Kind() != reflect.Struct {
		rTitle = rTitle.Elem()
	}
	rv := reflect.ValueOf(datas)
	rt := rv.Type()
	switch {
	case rt.Kind() == reflect.Slice && rt.Elem().Kind() == reflect.Struct:
		sliceType = "slice-struct"
	case rt.Kind() == reflect.Slice && rt.Elem().Kind() == reflect.Ptr && rt.Elem().Elem().Kind() == reflect.Struct:
		sliceType = "slice-pointer-struct"
	default:
		return nil, errors.New("tables 请传结构体数组 或 结构提指针数组")
	}
	length = rv.Len()

	contentRows = make([][]interface{}, 0, length)
	if !common.IsEmpty(datas) {
		// 结构体 数组
		if sliceType == "slice-struct" {
			for i := 0; i < length; i++ {
				data := make([]interface{}, 0, rTitle.NumField())
				common.EnumAnObjFieldNames(rTitle, func(f reflect.StructField) {
					if f.Tag.Get("csv") == "-" {
						return
					}
					path := f.Tag.Get("path")
					if path == "" {
						path = f.Name
					}
					if strings.Contains(path, "@") {
						// 需要对数据做转换处理
						var params []reflect.Value
						paths := strings.Split(path, "@")
						fields := strings.Split(paths[0], ";")
						for _, field := range fields {
							param := complicatedPath(rv.Index(i), field, false)
							if param.Kind() != reflect.Invalid {
								params = append(params, param)
							} else {
								params = append(params, reflect.New(param.Type()))
							}
						}
						dataStr := reflect.ValueOf(data).MethodByName(paths[1]).Call(params)[0].String()
						data = append(data, dataStr)
					} else if strings.Contains(path, "#") {
						paths := strings.Split(path, "#")
						dataStr := complicatedPath(rv.Index(i), paths[0], true).String()
						data = append(data, gjson.Get(dataStr, paths[1]).String())
					} else {
						dataInter := complicatedPath(rv.Index(i), path, true).Interface()
						data = append(data, fmt.Sprintf("%v", dataInter))
					}
				})
				contentRows = append(contentRows, data)
			}
			// 结构体 指针 数组
		} else if sliceType == "slice-pointer-struct" {
			for i := 0; i < length; i++ {
				data := make([]interface{}, 0, rTitle.NumField())
				common.EnumAnObjFieldNames(rTitle, func(f reflect.StructField) {
					if f.Tag.Get("csv") == "-" {
						return
					}
					path := f.Tag.Get("path")
					if path == "" {
						path = f.Name
					}
					if strings.Contains(path, "@") {
						// 需要对数据做转换处理
						var params []reflect.Value
						paths := strings.Split(path, "@")
						fields := strings.Split(paths[0], ";")
						for _, field := range fields {
							param := complicatedPath(rv.Index(i).Elem(), field, false)
							if param.Kind() != reflect.Invalid {
								params = append(params, param)
							} else {
								params = append(params, reflect.New(param.Type()))
							}
						}
						dataStr := reflect.ValueOf(data).MethodByName(paths[1]).Call(params)[0].String()
						data = append(data, dataStr)
					} else if strings.Contains(path, "#") {
						paths := strings.Split(path, "#")
						dataStr := complicatedPath(rv.Index(i).Elem(), paths[0], true).String()
						data = append(data, gjson.Get(dataStr, paths[1]).String())
					} else {
						dataInter := complicatedPath(rv.Index(i).Elem(), path, true).Interface()
						data = append(data, fmt.Sprintf("%v", dataInter))
					}
				})
				contentRows = append(contentRows, data)
			}
		}
	} else { // 空结构体 table 不支持 复杂 结构体 (嵌套型)
		return nil, errors.New("no data to export")
	}
	return
}

// todo, 进一步优化, 暂不支持数组, 考虑有时间加上
// 规则: 比如 type example struct {
// 		First  TheFirst
//		Second *TheFirst
// }
// 如果要取 TheFirst里的字段 如: FieldA , 则:
// "First.FieldA"  以及  "Second.*FieldA"
//            同理支持多层嵌套
func complicatedPath(v reflect.Value, path string, timeTransfer bool) reflect.Value {
	var tag reflect.StructField
	paths := strings.Split(path, ".")
	for _, p := range paths {
		var ok bool
		if p[:1] == "*" {
			tag, ok = v.Type().FieldByName(common.StrToTF1(p[1:]))
			v = v.FieldByName(common.StrToTF1(p[1:])).Elem()
			// 如果取到了空指针, 则直接返回一个空字符串
			if !ok || v.Kind() == reflect.Invalid {
				return reflect.ValueOf("")
			}
		} else {
			tag, ok = v.Type().FieldByName(common.StrToTF1(p))
			v = v.FieldByName(common.StrToTF1(p))
			if !ok || v.Kind() == reflect.Invalid {
				return reflect.ValueOf("")
			}
		}

	}
	if !v.IsValid() {
		panic("数据结构体中, 找不到" + path + "字段! 请检查")
	}
	// 对时间 time.Time 做特殊处理
	defaultFormat := "2006-01-02 15:04:05"
	if tagFormat := tag.Tag.Get("timeFormat"); tagFormat != "" {
		defaultFormat = tagFormat
	}
	if timeTransfer {
		tStr := v.Type().String()
		if tStr == "time.Time" {
			return reflect.ValueOf(v.Interface().(time.Time).Format(defaultFormat))
		} else if tStr == "*time.Time" && v.IsValid() && !v.IsNil() {
			return reflect.ValueOf(v.Interface().(*time.Time).Format(defaultFormat))
		} else if tStr == "*time.Time" && (v.IsNil() || !v.IsValid()) {
			return reflect.ValueOf("")
		}
	}
	return v
}
