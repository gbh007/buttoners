package button

import (
	"html/template"
	"io"
)

var badgeTemplate = template.Must(template.New("").Parse(`<svg xmlns="http://www.w3.org/2000/svg" width="{{ .totalWidth }}" height="20">
    <mask id="fillmask">
        <rect width="{{ .totalWidth }}" height="20" rx="3" fill="#fff" />
    </mask>

    <g mask="url(#fillmask)">
        <path fill="#fdb759" d="M0 0 h 50 v 20 H 0 z" />
        <path fill="#cdcdcd" d="M 50 0 h {{ .width }} v 20 H 50 z" />
    </g>

    <g fill="#000000" text-anchor="left" font-family="sans-serif" font-size="11">
        <text x="5" y="14">За🥖ил</text>
        <text x="55" y="14">{{ .text }}</text>
    </g>
</svg>`))

func renderBadgeTemplate(w io.Writer, count int) error {
	d := 0
	s := ""

	if count == 0 {
		s = "0"
		d = 1
	}

	for count > 0 {
		if d%3 == 0 && d > 0 {
			s = " " + s
		}
		s = string('0'+count%10) + s //nolint // так и задумано
		d++
		count /= 10
	}

	emoji := ""

	switch d {
	case 2:
		emoji = "💪"
	case 3:
		emoji = "🔥"
	case 4:
		emoji = "🙏"
	case 5:
		emoji = "😱"
	case 6:
		emoji = "🍌"
	case 7:
		emoji = "🍆"
	}

	if d > 7 {
		emoji = "🙈"
	}

	size := d*7 + 10
	if emoji != "" {
		size += 14
		s += " " + emoji
	}

	return badgeTemplate.Execute(w, map[string]any{
		"text":       s,
		"width":      size,
		"totalWidth": 50 + size,
	})
}
