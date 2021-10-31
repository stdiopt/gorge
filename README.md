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

## Packages

- `gorge` - contains mostly core components as data (light, camera, renderable,
  material, texture, font, transform,...)
- `systems/resource` - knows how to load gorge data (textures, mesh, material, fonts, ... )
  and might eventually have custom importers
- `systems/render` - knows how to render gorge components
- `systems/audio` -
- `systems/input` -

## Copyrights

- m32 contains code inspired or copied from go-gl authors https://github.com/go-gl/mathgl
- Shaders based on https://github.com/KhronosGroup/glTF-Sample-Viewer

BSD-3 Copyright ©2013 The go-gl Authors. All rights reserved.
