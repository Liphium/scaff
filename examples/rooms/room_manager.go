package rooms

type RoomManager struct {
	start *Room
}

type Room struct {
	top, down, right, left *Room
}
