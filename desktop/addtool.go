package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"net"
	"sort"
	"sync"
)

type AddLink struct {
	Iface     string
	Port      int
	Mode      string
	BackEnd []string
}

var consoleIface   *walk.ComboBox
var consoleBackend *walk.ComboBox
var consoleMode    *walk.ComboBox
var consolePort    *walk.NumberEdit

func IfaceOptions() []string {
	ifaces, err := net.Interfaces()
	if err != nil {
		logs.Error(err.Error())
	}
	output := []string{"0.0.0.0"}
	for _, v := range ifaces {
		if v.Flags & net.FlagUp == 0 {
			continue
		}
		address, err := InterfaceLocalIP(&v)
		if err != nil {
			continue
		}
		if len(address) == 0 {
			continue
		}
		output = append(output, address[0].String())
	}
	return output
}

func LoadBalanceModeOptions() []string {
	return []string{
		"Random","RoundRobin","WeightRoundRobin","AddressHash","MainStandby",
	}
}

func MainStandbyOptions() []string {
	return []string{
		"Random","RoundRobin","WeightRoundRobin","AddressHash","MainStandby",
	}
}

type BackendItem struct {
	Index        int
	Address      string
	Weight       int
	Standby      bool

	checked      bool
}

type BackendModel struct {
	sync.RWMutex

	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder

	items      []*BackendItem
}

func (n *BackendModel)RowCount() int {
	return len(n.items)
}

func (n *BackendModel)Value(row, col int) interface{} {
	item := n.items[row]
	switch col {
	case 0:
		return item.Index
	case 1:
		return item.Address
	case 2:
		return item.Weight
	case 3:
		if item.Standby {
			return "standby"
		}
		return "main"
	}
	panic("unexpected col")
}

func (n *BackendModel) Checked(row int) bool {
	return n.items[row].checked
}

func (n *BackendModel) SetChecked(row int, checked bool) error {
	n.items[row].checked = checked
	return nil
}

func (m *BackendModel) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order
	sort.SliceStable(m.items, func(i, j int) bool {
		a, b := m.items[i], m.items[j]
		c := func(ls bool) bool {
			if m.sortOrder == walk.SortAscending {
				return ls
			}
			return !ls
		}
		switch m.sortColumn {
		case 0:
			return c(a.Index < b.Index)
		case 1:
			return c(a.Address < b.Address)
		case 2:
			return c(a.Weight < b.Weight)
		case 3:
			return c(a.Standby)
		}
		panic("unreachable")
	})
	return m.SorterBase.Sort(col, order)
}


func AddToolBar()  {
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton
	var backendView *walk.TableView

	backendTable := new(BackendModel)
	backendTable.items = make([]*BackendItem, 0)

	cnt, err := Dialog{
		AssignTo: &dlg,
		Title: "Add Link",
		Icon: ICON_TOOL_ADD,
		DefaultButton: &acceptPB,
		CancelButton: &cancelPB,
		Size: Size{350, 500},
		MinSize: Size{350, 500},
		Layout:  VBox{ Margins: Margins{Top: 10, Bottom: 10, Left: 10, Right: 10}},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "Bind Ethernet:",
					},
					ComboBox{
						AssignTo: &consoleIface,
						CurrentIndex:  0,
						Model:         IfaceOptions(),
						OnCurrentIndexChanged: func() {
							//LocalIfaceOptionsSet(consoleIface.Text())
						},
					},
					Label{
						Text: "Bind Port:",
					},
					NumberEdit{
						AssignTo: &consolePort,
						Value:    float64(8080),
						ToolTipText: "1~65535",
						MaxValue: 65535,
						MinValue: 1,
						OnValueChanged: func() {
							//PortOptionSet(int(consolePort.Value()))
						},
					},
					Label{
						Text: "Load Balance Mode:",
					},
					ComboBox{
						AssignTo: &consoleMode,
						CurrentIndex:  0,
						Model:         LoadBalanceModeOptions(),
						OnCurrentIndexChanged: func() {

						},
					},

					Label{
						Text: "Backend Address:",
					},
					LineEdit{
						Text: "",
					},
					Label{
						Text: "Weight Value:",
					},
					NumberEdit{
						Value:    float64(50),
						ToolTipText: "1~100",
						MaxValue: 100,
						MinValue: 1,
						OnValueChanged: func() {

						},
					},
					Label{
						Text: "Main or Standby:",
					},
					Composite{
						Layout: HBox{MarginsZero: true},
						Children: []Widget{
							RadioButton{
								Text: "Main",
							},
							RadioButton{
								Text: "Standby",
							},
						},
					},
					Label{
						Text: "Backend List Edit:",
					},
					Composite{
						Layout: HBox{MarginsZero: true},
						Children: []Widget{
							PushButton{
								Text: "Add Backend",
							},
							PushButton{
								Text: "Del Backend",
							},
						},
					},
				},
			},
			Composite{
				Layout: VBox{},
				Children: []Widget{

					TableView{
						AssignTo: &backendView,
						AlternatingRowBG: true,
						ColumnsOrderable: true,
						CheckBoxes: true,
						Columns: []TableViewColumn{
							{Title: "#", Width: 30},
							{Title: "Address", Width: 120},
							{Title: "Weight", Width: 60},
							{Title: "Main&Standby", Width: 80},
						},
						StyleCell: func(style *walk.CellStyle) {
							if style.Row()%2 == 0 {
								style.BackgroundColor = walk.RGB(248, 248, 255)
							} else {
								style.BackgroundColor = walk.RGB(220, 220, 220)
							}
						},
						Model:backendTable,
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						AssignTo: &acceptPB,
						Text: "Add",
						OnClicked: func() {
							acceptPB.SetEnabled(false)
							cancelPB.SetEnabled(false)

							go func() {

								acceptPB.SetEnabled(true)
								cancelPB.SetEnabled(true)

								//if err != nil {
								//	ErrorBoxAction(dlg, err.Error())
								//	return
								//}
								dlg.Accept()
							}()
						},
					},
					PushButton{
						AssignTo: &cancelPB,
						Text: "Cancel",
						OnClicked: func() {
							dlg.Cancel()
						},
					},
				},
			},
		},
	}.Run(MainWindowsCtrl())
	if err != nil {
		logs.Error(err.Error())
	} else {
		logs.Info("add link dialog return %d", cnt)
	}
}
