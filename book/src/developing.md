# Developing

Here we describe good practices for developing components for our framework.

Rules for components are simple:
- Component should declare some `Input` and an optional `Output` (we use that so we can skip or cache any component results)
- Components should be isolated, they should not return anything except basic types like `int`, `string`, `maps` or `structs`
- Component **must** have documentation under [Components](./framework/components/overview.md), here is an [example](./framework/components/chainlink/node.md)