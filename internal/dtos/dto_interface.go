package dtos

type IEntityTransformer interface {
	ToEntity() (interface{}, error)
	IsValid() bool
}
