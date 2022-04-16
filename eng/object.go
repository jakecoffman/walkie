package eng

var objectId int

func GetObjectId() int {
	objectId++
	return objectId
}

type Object struct {
	lastPosition *Vector
}
