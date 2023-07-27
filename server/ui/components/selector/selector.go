package selector

import (
	"sync"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/soft-serve/server/ui/common"
)

// Selector is a list of items that can be selected.
type Selector struct {
	KeyMap list.KeyMap

	model       list.Model
	common      common.Common
	active      int
	filterState list.FilterState
	mtx         sync.RWMutex
}

// IdentifiableItem is an item that can be identified by a string. Implements
// list.DefaultItem.
type IdentifiableItem interface {
	list.DefaultItem
	ID() string
}

// ItemDelegate is a wrapper around list.ItemDelegate.
type ItemDelegate interface {
	list.ItemDelegate
}

// SelectMsg is a message that is sent when an item is selected.
type SelectMsg struct{ IdentifiableItem }

// ActiveMsg is a message that is sent when an item is active but not selected.
type ActiveMsg struct{ IdentifiableItem }

// New creates a new selector.
func New(common common.Common, items []IdentifiableItem, delegate ItemDelegate) *Selector {
	itms := make([]list.Item, len(items))
	for i, item := range items {
		itms[i] = item
	}
	l := list.New(itms, delegate, common.Width, common.Height)
	l.Styles.NoItems = common.Styles.NoItems
	s := &Selector{
		model:  l,
		common: common,
		KeyMap: list.DefaultKeyMap(),
	}
	s.SetSize(common.Width, common.Height)
	return s
}

// FilterState returns the filter state.
func (s *Selector) FilterState() list.FilterState {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.model.FilterState()
}

// SelectedItem returns the selected item.
func (s *Selector) SelectedItem() list.Item {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.model.SelectedItem()
}

// Items returns the items.
func (s *Selector) Items() []list.Item {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.model.Items()
}

// VisibleItems returns the visible items.
func (s *Selector) VisibleItems() []list.Item {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.model.VisibleItems()
}

// PerPage returns the number of items per page.
func (s *Selector) PerPage() int {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.model.Paginator.PerPage
}

// SetPage sets the current page.
func (s *Selector) SetPage(page int) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.model.Paginator.Page = page
}

// Page returns the current page.
func (s *Selector) Page() int {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.model.Paginator.Page
}

// TotalPages returns the total number of pages.
func (s *Selector) TotalPages() int {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.model.Paginator.TotalPages
}

// Select selects the item at the given index.
func (s *Selector) Select(index int) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.model.Select(index)
}

// SetShowTitle sets the show title flag.
func (s *Selector) SetShowTitle(show bool) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.model.SetShowTitle(show)
}

// SetShowHelp sets the show help flag.
func (s *Selector) SetShowHelp(show bool) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.model.SetShowHelp(show)
}

// SetShowStatusBar sets the show status bar flag.
func (s *Selector) SetShowStatusBar(show bool) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.model.SetShowStatusBar(show)
}

// DisableQuitKeybindings disables the quit keybindings.
func (s *Selector) DisableQuitKeybindings() {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.model.DisableQuitKeybindings()
}

// SetShowFilter sets the show filter flag.
func (s *Selector) SetShowFilter(show bool) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.model.SetShowFilter(show)
}

// SetShowPagination sets the show pagination flag.
func (s *Selector) SetShowPagination(show bool) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.model.SetShowPagination(show)
}

// SetFilteringEnabled sets the filtering enabled flag.
func (s *Selector) SetFilteringEnabled(enabled bool) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.model.SetFilteringEnabled(enabled)
}

// SetSize implements common.Component.
func (s *Selector) SetSize(width, height int) {
	s.common.SetSize(width, height)
	s.mtx.Lock()
	s.model.SetSize(width, height)
	s.mtx.Unlock()
}

// SetItems sets the items in the selector.
func (s *Selector) SetItems(items []IdentifiableItem) tea.Cmd {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	its := make([]list.Item, len(items))
	for i, item := range items {
		its[i] = item
	}
	return s.model.SetItems(its)
}

// Index returns the index of the selected item.
func (s *Selector) Index() int {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.model.Index()
}

// CursorUp moves the cursor up.
func (s *Selector) CursorUp() {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.model.CursorUp()
}

// CursorDown moves the cursor down.
func (s *Selector) CursorDown() {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.model.CursorDown()
}

// Init implements tea.Model.
func (s *Selector) Init() tea.Cmd {
	s.mtx.Lock()
	s.model.KeyMap = s.KeyMap
	s.mtx.Unlock()
	return s.activeCmd
}

// Update implements tea.Model.
func (s *Selector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseWheelUp:
			s.CursorUp()
		case tea.MouseWheelDown:
			s.CursorDown()
		case tea.MouseLeft:
			curIdx := s.Index()
			for i, item := range s.Items() {
				item, _ := item.(IdentifiableItem)
				// Check each item to see if it's in bounds.
				if item != nil && s.common.Zone.Get(item.ID()).InBounds(msg) {
					if i == curIdx {
						cmds = append(cmds, s.selectCmd)
					} else {
						s.Select(i)
					}
					break
				}
			}
		}
	case tea.KeyMsg:
		filterState := s.FilterState()
		switch {
		case key.Matches(msg, s.common.KeyMap.Help):
			if filterState == list.Filtering {
				return s, tea.Batch(cmds...)
			}
		case key.Matches(msg, s.common.KeyMap.Select):
			if filterState != list.Filtering {
				cmds = append(cmds, s.selectCmd)
			}
		}
	case list.FilterMatchesMsg:
		cmds = append(cmds, s.activeFilterCmd)
	}
	s.mtx.Lock()
	m, cmd := s.model.Update(msg)
	s.model = m
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	s.mtx.Unlock()
	// Track filter state and update active item when filter state changes.
	filterState := s.FilterState()
	if s.filterState != filterState {
		cmds = append(cmds, s.activeFilterCmd)
	}
	s.filterState = filterState
	// Send ActiveMsg when index change.
	if s.active != s.Index() {
		cmds = append(cmds, s.activeCmd)
	}
	s.active = s.Index()
	return s, tea.Batch(cmds...)
}

// View implements tea.Model.
func (s *Selector) View() string {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.model.View()
}

// SelectItem is a command that selects the currently active item.
func (s *Selector) SelectItem() tea.Msg {
	return s.selectCmd()
}

func (s *Selector) selectCmd() tea.Msg {
	item := s.SelectedItem()
	i, ok := item.(IdentifiableItem)
	if !ok {
		return SelectMsg{}
	}
	return SelectMsg{i}
}

func (s *Selector) activeCmd() tea.Msg {
	item := s.SelectedItem()
	i, ok := item.(IdentifiableItem)
	if !ok {
		return ActiveMsg{}
	}
	return ActiveMsg{i}
}

func (s *Selector) activeFilterCmd() tea.Msg {
	// Here we use VisibleItems because when list.FilterMatchesMsg is sent,
	// VisibleItems is the only way to get the list of filtered items. The list
	// bubble should export something like list.FilterMatchesMsg.Items().
	items := s.VisibleItems()
	if len(items) == 0 {
		return nil
	}
	item := items[0]
	i, ok := item.(IdentifiableItem)
	if !ok {
		return nil
	}
	return ActiveMsg{i}
}
