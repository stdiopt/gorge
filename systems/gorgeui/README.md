# gorge ui system

gorgeui system handles entities with specific components such as
`RectComponent` and `ElementComponent`

## handler

Gorge ui system handles the following events

- input.EventPointer
- gorge.EventPreUpdate
- gorge.EventPostUpdate
- gorge.EventAddEntity
- gorge.EventRemoveEntity

## components

- RectComponent  
  derives `gorge.TransformComponent` and adds extra data for 2D Rects
- ElementComponent  
  Contains a `event.Bus` and `gorge.Container`
