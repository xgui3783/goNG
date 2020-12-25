# goNG

Yet another package that converts popular 3D formats (gii, obj, off, stl, gltf) and [neuroglancer precomputed mesh format](https://github.com/google/neuroglancer/tree/5bfa8c3/src/neuroglancer/datasource/precomputed#legacy-single-resolution-mesh-format)

## Requirements

- Ubuntu >= 18.10*

Most of the functionality should work in all OS, except for ng precomputed mesh format. Windows and MacOS seems to have a hard time parsing `:` in the filename.

## Usage

TODO code examples

### Input formats

- [neuroglancer precomputed mesh format](https://github.com/google/neuroglancer/tree/5bfa8c3/src/neuroglancer/datasource/precomputed#legacy-single-resolution-mesh-format)
- STL ascii
- STL binary
- gifti
- obj
- vtk
- OFF ascii

### Output formats

- [neuroglancer precomputed mesh format](https://github.com/google/neuroglancer/tree/5bfa8c3/src/neuroglancer/datasource/precomputed#legacy-single-resolution-mesh-format)
- STL ascii
- STL binary
- gifti
- obj
- vtk
- OFF ascii

### Splitting meshes

See <common/splitMesh.md>

## Build

from v2.0+ onwards, gong supports [go build constraints](https://golang.org/pkg/go/build/#hdr-Build_Constraints) using tags. This should allow for smaller binary sizes.

For example, if only the conversion bewteen `STL_ASCII` and `NG_MESH` is required, a smaller binary can be built with:

```bash
go build -tags "ng_mesh stl_ascii"
```

Available tags are:

- all
- gii
- gltf
- ng_mesh
- obj
- off_ascii
- stl_ascii
- stl_binary
- vtk

## Development

requirements: 
```
go >= 1.11
git
git lfs # for downloading testing meshes
```
## License

MIT