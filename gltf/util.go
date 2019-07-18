package gltf

type GltfMaterial struct {
	Name                 string `json:"name"`
	PbrMetallicRoughness `json:"pbrMetallicRoughness"`
}

type PbrMetallicRoughness struct {
	RoughnessFactor float32    `json:"roughnessFactor"`
	MetallicFactor  float32    `json:"metallicFactor"`
	BaseColorFactor [4]float32 `json:"baseColorFactor"`
}
