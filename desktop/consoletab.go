package main

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"sort"
	"sync"
)

type LinkItem struct {
	Index        int
	Bind         string
	Mode         string
	Count        int
	Speed        int64
	Size         int64
	Status       string

	checked      bool
}

type LinkModel struct {
	sync.RWMutex

	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder

	items      []*LinkItem
}

func (n *LinkModel)RowCount() int {
	return len(n.items)
}

func (n *LinkModel)Value(row, col int) interface{} {
	item := n.items[row]
	switch col {
	case 0:
		return item.Index
	case 1:
		return item.Bind
	case 2:
		return item.Mode
	case 3:
		return fmt.Sprintf("%d", item.Count)
	case 4:
		if item.Speed == 0 {
			return "-"
		}
		return fmt.Sprintf("%s/s", ByteViewLite(item.Speed))
	case 5:
		return ByteView(item.Size)
	case 6:
		return item.Status
	}
	panic("unexpected col")
}

func (n *LinkModel) Checked(row int) bool {
	return n.items[row].checked
}

func (n *LinkModel) SetChecked(row int, checked bool) error {
	n.items[row].checked = checked
	return nil
}

func (m *LinkModel) Sort(col int, order walk.SortOrder) error {
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
			return c(a.Bind < b.Bind)
		case 2:
			return c(a.Mode < b.Mode)
		case 3:
			return c(a.Count < b.Count)
		case 4:
			return c(a.Speed < b.Speed)
		case 5:
			return c(a.Size < b.Size)
		case 6:
			return c(a.Status < b.Status)
		}
		panic("unreachable")
	})
	return m.SorterBase.Sort(col, order)
}

const (
	STATUS_UNLINK = "unlink"
	STATUS_LINK   = "link"
)

func StatusToIcon(status string) walk.Image {
	switch status {
	case STATUS_LINK:
		return ICON_STATUS_LINK
	case STATUS_UNLINK:
		return ICON_STATUS_UNLINK
	default:
		return ICON_STATUS_UNLINK
	}
	return nil
}

var consoleLinkTable *LinkModel

func init()  {
	consoleLinkTable = new(LinkModel)
	consoleLinkTable.items = make([]*LinkItem, 0)
}

func LinkTalbeUpdate(item []*LinkItem )  {
	lt := consoleLinkTable

	lt.Lock()
	defer lt.Unlock()

	oldItem := lt.items
	if len(oldItem) == len(item) {
		for i, v := range item {
			v.checked = oldItem[i].checked
		}
	}

	lt.items = item
	lt.PublishRowsReset()
	lt.Sort(lt.sortColumn, lt.sortOrder)
}

func (lk *LinkModel)LinkTableSelectClean()  {
	lk.Lock()
	defer lk.Unlock()

	for _, v := range lk.items {
		v.checked = false
	}

	lk.PublishRowsReset()
	lk.Sort(lk.sortColumn, lk.sortOrder)
}

func (lk *LinkModel)LinkTableSelectAll()  {
	lk.Lock()
	defer lk.Unlock()

	done := true
	for _, v := range lk.items {
		if !v.checked {
			done = false
		}
	}

	for _, v := range lk.items {
		v.checked = !done
	}

	lk.PublishRowsReset()
	lk.Sort(lk.sortColumn, lk.sortOrder)
}

func (lt *LinkModel)LinkTableSelectList() []string {
	lt.RLock()
	defer lt.RUnlock()

	var output []string
	for _, v := range lt.items {
		if v.checked {
			output = append(output, v.Bind)
		}
	}

	return output
}

func (lt *LinkModel)LinkTableSelectStatus(status string)  {
	lt.Lock()
	defer lt.Unlock()

	for _, v := range lt.items {
		v.checked = false
	}

	for _, v := range lt.items {
		if v.Status == status {
			v.checked = true
		}
	}

	lt.PublishRowsReset()
	lt.Sort(lt.sortColumn, lt.sortOrder)
}

var tableView *walk.TableView

func TableWight() []Widget {
	return []Widget{
		Label{
			Text: "Link List:",
		},
		TableView{
			AssignTo: &tableView,
			AlternatingRowBG: true,
			ColumnsOrderable: true,
			CheckBoxes: true,
			OnItemActivated: func() {
				InfoBoxAction(MainWindowsCtrl(),"")
			},
			Columns: []TableViewColumn{
				{Title: "#", Width: 30},
				{Title: "Bind", Width: 120},
				{Title: "Mode", Width: 80},
				{Title: "Connects", Width: 60},
				{Title: "Traffic", Width: 80},
				{Title: "Stat", Width: 60},
				{Title: "Status", Width: 80},
			},
			StyleCell: func(style *walk.CellStyle) {
				item := consoleLinkTable.items[style.Row()]
				if style.Row()%2 == 0 {
					style.BackgroundColor = walk.RGB(248, 248, 255)
				} else {
					style.BackgroundColor = walk.RGB(220, 220, 220)
				}
				switch style.Col() {
				case 6:
					style.Image = StatusToIcon(item.Status)
				}
			},
			Model:consoleLinkTable,
		},
		Composite{
			Layout: HBox{MarginsZero: true},
			Children: []Widget{
				PushButton{
					Text: "All",
					OnClicked: func() {
						go func() {
							consoleLinkTable.LinkTableSelectAll()
						}()
					},
				},
				PushButton{
					Text: "Linked",
					OnClicked: func() {
						go func() {
							consoleLinkTable.LinkTableSelectStatus(STATUS_LINK)
						}()
					},
				},
				PushButton{
					Text: "Unlinked",
					OnClicked: func() {
						go func() {
							consoleLinkTable.LinkTableSelectStatus(STATUS_UNLINK)
						}()
					},
				},
				HSpacer{

				},
			},
		},
	}
}

