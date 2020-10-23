package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"net"
	"sort"
	"sync"
)

type BackendConfig struct {
	Address  string
	Timeout  int
	Weight   int
	Standby  bool
}

type LinkConfig struct {
	Iface      string
	Port       int
	Timeout    int
	Mode       string
	Backend  []BackendConfig
}

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

type BackendItem struct {
	Index        int
	Address      string
	Timeout      int
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

func (n *BackendModel) Input(cfg []BackendConfig)  {
	n.RLock()
	defer n.RUnlock()

	var items []*BackendItem
	for i, v := range cfg {
		items = append(items, &BackendItem{
			Index: i, Address: v.Address,
			Timeout: v.Timeout, Weight: v.Weight,
			Standby: v.Standby,
		})
	}
	n.items = items
}

func (n *BackendModel) Output() []BackendConfig {
	n.RLock()
	defer n.RUnlock()

	var output []BackendConfig
	for _, v := range n.items {
		output = append(output, BackendConfig{
			Address: v.Address, Weight: v.Weight, Standby: v.Standby,
		})
	}
	return output
}

func (n *BackendModel)Del()  {
	n.Lock()
	defer n.Unlock()

	var idx int
	var items []*BackendItem
	for _, v := range n.items {
		if v.checked {
			continue
		}
		v.Index = idx
		items = append(items, v)
		idx++
	}

	n.items = items
	n.PublishRowsReset()
	n.Sort(n.sortColumn, n.sortOrder)
}

func (n *BackendModel)Add(addr string, timeout int, weight int, standby bool)  {
	n.Lock()
	defer n.Unlock()

	n.items = append(n.items, &BackendItem{
		Index: len(n.items),
		Timeout: timeout,
		Address: addr,
		Weight: weight,
		Standby: standby,
	})

	n.PublishRowsReset()
	n.Sort(n.sortColumn, n.sortOrder)
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
		if item.Timeout == 0 {
			return "-"
		}
		return fmt.Sprintf("%ds", item.Timeout)
	case 3:
		return item.Weight
	case 4:
		if item.Standby == true {
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
			return c(a.Timeout < b.Timeout)
		case 3:
			return c(a.Weight < b.Weight)
		case 4:
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

	var consoleIface   *walk.ComboBox
	var consoleMode    *walk.ComboBox
	var consolePort    *walk.NumberEdit
	var consoleTimeout *walk.NumberEdit

	var BackendAddr    *walk.LineEdit
	var BackendWeight  *walk.NumberEdit
	var BackendTimeout *walk.NumberEdit
	var backendStandby *walk.RadioButton
	var backendMain    *walk.RadioButton

	backendTable := new(BackendModel)
	backendTable.items = make([]*BackendItem, 0)

	var addLink LinkConfig

	addLink.Iface = "0.0.0.0"
	addLink.Port = 8080
	addLink.Timeout = 60
	addLink.Mode = LoadBalanceModeOptions()[0]

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
							addLink.Iface = consoleIface.Text()
						},
					},
					Label{
						Text: "Bind Port:",
					},
					NumberEdit{
						AssignTo: &consolePort,
						Value:    float64(addLink.Port),
						ToolTipText: "1~65535",
						MaxValue: 65535,
						MinValue: 1,
						OnValueChanged: func() {
							addLink.Port = int(consolePort.Value())
						},
					},
					Label{
						Text: "Bind Timeout:",
					},
					NumberEdit{
						AssignTo: &consoleTimeout,
						Value:    float64(addLink.Timeout),
						ToolTipText: "0~3600",
						MaxValue: 3600,
						MinValue: 0,
						Suffix: " Second",
						OnValueChanged: func() {
							addLink.Timeout = int(consoleTimeout.Value())
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
							addLink.Mode = consoleMode.Text()
						},
					},

					Label{
						Text: "Backend Address:",
					},
					LineEdit{
						AssignTo: &BackendAddr,
						ToolTipText: "192.168.1.100:8080",
						CueBanner: "192.168.1.100:8080",
						Text: "",
						OnEditingFinished: func() {
							addr := BackendAddr.Text()
							if AddressValid(addr) == false {
								BackendAddr.SetTextColor(walk.RGB(255,50,50))
								return
							} else {
								BackendAddr.SetTextColor(walk.RGB(0,0,0))
							}
						},
					},
					Label{
						Text: "Backend Timeout:",
					},
					NumberEdit{
						AssignTo: &BackendTimeout,
						Value:    float64(60),
						ToolTipText: "0~3600",
						MaxValue: 3600,
						MinValue: 0,
						Suffix: " Second",
					},
					Label{
						Text: "Weight Value:",
					},
					NumberEdit{
						AssignTo: &BackendWeight,
						Value:    float64(50),
						ToolTipText: "1~100",
						MaxValue: 100,
						MinValue: 1,
					},
					Label{
						Text: "Main or Standby:",
					},
					Composite{
						Layout: HBox{MarginsZero: true},
						Children: []Widget{
							RadioButton{
								AssignTo: &backendMain,
								Text: "Main",
								OnBoundsChanged: func() {
									backendMain.SetChecked(true)
								},
								OnClicked: func() {
									backendStandby.SetChecked(false)
								},
							},
							RadioButton{
								AssignTo: &backendStandby,
								Text: "Standby",
								OnClicked: func() {
									backendMain.SetChecked(false)
								},
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
								OnClicked: func() {
									addr := BackendAddr.Text()
									if AddressValid(addr) == false {
										BackendAddr.SetFocus()
										return
									}
									backendTable.Add(addr,
										int(BackendTimeout.Value()),
										int(BackendWeight.Value()),
										backendMain.Checked() == false )
								},
							},
							PushButton{
								Text: "Del Backend",
								OnClicked: func() {
									backendTable.Del()
								},
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
							{Title: "#", Width: 20},
							{Title: "Address", Width: 110},
							{Title: "Timeout", Width: 50},
							{Title: "Weight", Width: 50},
							{Title: "Main/Standby", Width: 80},
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
								defer func() {
									acceptPB.SetEnabled(true)
									cancelPB.SetEnabled(true)
								}()

								if ListenCheck(addLink.Iface, addLink.Port) == false {
									ErrorBoxAction(dlg,
										fmt.Sprintf("Address %s:%d binding failed!",
											addLink.Iface, addLink.Port))
									return
								}

								output := backendTable.Output()
								if len(output) == 0 {
									ErrorBoxAction(dlg, "Please add backend instance.")
									return
								}

								addLink.Backend = output
								err := LinkAdd(&addLink)
								if err != nil {
									ErrorBoxAction(dlg, err.Error())
									return
								}

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
