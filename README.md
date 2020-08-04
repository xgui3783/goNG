# goNG
Yet another package that converts popular 3D formats (gii, obj, off, stl, gltf) and [neuroglancer precomputed mesh format](https://github.com/google/neuroglancer/tree/5bfa8c3/src/neuroglancer/datasource/precomputed#legacy-single-resolution-mesh-format)

# Requirements

- Ubuntu >= 18.10*

Most of the functionality should work in all OS, except for ng precomputed mesh format. Windows and MacOS seems to have a hard time parsing `:` in the filename.

# Usage

TODO code examples

## Input formats

- [neuroglancer precomputed mesh format](https://github.com/google/neuroglancer/tree/5bfa8c3/src/neuroglancer/datasource/precomputed#legacy-single-resolution-mesh-format)
- STL ascii
- STL binary
- gifti
- obj
- vtk
- OFF ascii

## Output formats

- [neuroglancer precomputed mesh format](https://github.com/google/neuroglancer/tree/5bfa8c3/src/neuroglancer/datasource/precomputed#legacy-single-resolution-mesh-format)
- STL ascii
- STL binary
- gifti
- obj
- vtk
- OFF ascii

# Development

requirements: 
```
go >= 1.11
git
git lfs # for downloading testing meshes
```
# License

MIT