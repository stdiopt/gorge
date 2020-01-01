# Gorge - WIP

gorge - go r? game engine

(reduced? rudimentary?)

A personal project that I am developing for learning purposes

## Demos

- [gophers](https://stdiopt.github.io/gorge/wasm/?t=gophers)
- [simple](https://stdiopt.github.io/gorge/wasm/?t=simple)
- [scene](https://stdiopt.github.io/gorge/wasm/?t=scene)
- [boxes](https://stdiopt.github.io/gorge/wasm/?t=boxes)
- [text](https://stdiopt.github.io/gorge/wasm/?t=text)

## Platforms

It was first created with wasm in mind, others were added later

- wasm
- glfw (linux, windows?, osx?)
- mobile (golang.org/x/mobile - WIP)

## Example

```go
func main() {
	opt := platform.Options{}

	platform.Start(opt, func(g *gorge.Gorge) {
        s := g.Scene(simpleScene)
        g.StartScene(s)
	})
}

func simpleScene(s *gorge.Scene) {
	gorgeutils.TrackballCamera(s)

	light := gorgeutils.NewLight()
	light.SetPosition(0, 10, -4)

	cube := primitive.Cube()
	g.AddEntity(light)
	g.AddEntity(cube)

	g.Handle(func(dt gorge.UpdateEvent) {
		cube.Rotate(0, 1*float32(dt), 0)
	})
}
```

## Packages

- `gorge` - contains mostly core components as data (light, camera, renderable,
  material, texture, font, transform,...)
- `resource` - knows how to load gorge data (textures, mesh, material, fonts, ... )
  and might eventually have custom importers
- `renderer` - knows how to render gorge components
- input - ...
- ...

## Todos

There are a couple of improvements and new features that I would like to
implement

- Scene loading, unloading
- Render textures, filters
- shadow maps, reflection probes
- renderer entity removal
- particle system
- animation, tweening

## Notes

It contains a slightly modified version of `golang.org/x/mobile` in vendor
folder as a couple of gles3 bindings were missing
