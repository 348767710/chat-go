package args

type ContactArg struct {
	PageArg
	Userid   int64  `json:"userid" form:"userid"`
	FriendId int64  `json:"friend_id" form:"friend_id"`
	Type     string `json:"type" form:"type"`
	PetName  string `json:"pet_name" form:"pet_name"`
}

type GroupArg struct {
	GroupId int64 `json:"groupId" form:"groupId"`
	UserId  int64 `json:"userId" form:"userId"`
}

type GroupChatDetailIds struct {
	PageArg
	Userid     int64  `json:"userid" form:"userid"`
	SendStatus string `json:"sendStatus" form:"sendStatus"`
	CreateTime int64  `json:"create_time" form:"create_time"`
}
