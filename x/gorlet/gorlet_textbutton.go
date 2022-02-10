package gorlet

type WTextButton struct {
	Widget[WTextButton]
	btn *WButton
	lbl *WLabel
}

func TextButton(t string) *WTextButton {
	return Build(&WTextButton{}).SetText(t)
}

func (w *WTextButton) Build(b *B) {
	w.SetAnchor(0, 0, 1, 0)
	w.SetSize(0, 3)

	w.btn = b.BeginButton().
		SetAnchor(0, 0, 1, 1).
		SetSize(0)
	{
		w.lbl = b.Label("").
			SetColor(0).
			FillParent()
	}
	b.EndButton()
}

func (w *WTextButton) SetText(t string) *WTextButton {
	w.lbl.SetText(t)
	return w
}

func (w *WTextButton) SetFontScale(s float32) *WTextButton {
	w.lbl.SetFontScale(s)
	return w
}

func (w *WTextButton) Color(vs ...float32) *WTextButton {
	w.btn.Color(vs...)
	return w
}

func (w *WTextButton) Highlight(vs ...float32) *WTextButton {
	w.btn.Highlight(vs...)
	return w
}

func (w *WTextButton) Pressed(vs ...float32) *WTextButton {
	w.btn.Pressed(vs...)
	return w
}

func (w *WTextButton) FadeFactor(f float32) *WTextButton {
	w.btn.FadeFactor(f)
	return w
}

func (w *WTextButton) OnClick(fn func()) *WTextButton {
	w.btn.OnClick(fn)
	return w
}

func (b *B) TextButton(t string) *WTextButton {
	w := TextButton(t)
	b.Add(w)
	return w
}
