# Split Mesh

From time to time, there may be a need to split a single mesh into multiple meshes. 

## Strategies

Currently, the only supported strategy is via labelled vertices.

### By labelled vertices

Supply the `-splitByVertexPath` argument to the path of a text file, with the following format:

```
# this is a comment
# {vertex_index} {label}
# vertex_index needs to be able to be parsed into utin32, and label needs to be an non empty string
0 label_a
1 label_a
2 label_b
```

#### Ambiguity

There will inevitably result in ambiguities. When one or more of the vertices disagree with the labelling of the face. To resolve the ambiguity, supply `-splitMeshAmbiguousStrategy` with the desired strategy

- EMPTY_LABEL *(default)*

All ambiguous faces will be labelled as EMPTY_LABEL

- MAJORITY_OR_FIRST_INDEX

Majority will take precedence. If there is no clear majority, the first index will be used. 