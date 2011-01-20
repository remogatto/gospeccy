package spectrum

type Renderer interface {
	Resize(app *Application, scale2x, fullscreen bool)
}
