package engine

type NewRoomFunc func(RoomID, CRHandler) Room

type Creator struct {
	new NewRoomFunc
	max int
}

type CRHandler struct {
	creators map[RoomName]Creator
}

func NewCRHandler() *CRHandler {
	return &CRHandler{
		creators: make(map[RoomName]Creator),
	}
}

func (rh *CRHandler) RegisterCreator(name RoomName, new NewRoomFunc, max int) {
	rh.creators[name] = Creator{new: new, max: max}
}
