package follow

//type FollowValidator struct {
//	 string `json:"firstName" form:"first_name"`
//}

// Response
type FollowUserFollowCountResponse struct {
	Following uint64 `json:"following"`
	Followers uint64 `json:"followers"`
}

// post
type PostFollowUserRequest struct {
	UserId uint64 `json:"userId" binding:"required"`
}

// post
type PostUnFollowUserRequest struct {
	UserId uint64 `json:"userId" binding:"required"`
}

// Response
type FollowUserStudioCountResponse struct {
	Followers uint64 `json:"followers"`
}

// post
type PostFollowStudioRequest struct {
	UserId uint64 `json:"userId" binding:"required"`
}

// post
type PostUnFollowStudioRequest struct {
	UserId uint64 `json:"userId" binding:"required"`
}
