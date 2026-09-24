package panels

const (
	PanelSizeHalf = iota
	PanelSizeQuarter
	PanelSizeQuarterSlim
	PanelSizeQuarterHigh
	PanelSizeThird
	PanelSizeFull
)

type PanelSize byte

func WithPanelSize[T interface {
	Height(h uint32) T
	Span(w uint32) T
}](data T, size PanelSize) T {
	var h, w uint32

	switch size {
	case PanelSizeQuarterSlim:
		h = 3
		w = 6
	case PanelSizeQuarter:
		h = 6
		w = 6
	case PanelSizeQuarterHigh:
		h = 9
		w = 6
	case PanelSizeThird:
		h = 9
		w = 8
	case PanelSizeHalf:
		h = 9
		w = 12
	case PanelSizeFull:
		h = 12
		w = 24
	}

	return data.Height(h).Span(w)
}
