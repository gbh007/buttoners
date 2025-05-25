package components

func RenderError(err error) string {
	if err == nil {
		return ""
	}

	return ErrorStyle.Render("Ошибка: "+err.Error()) + "\n\n"
}
