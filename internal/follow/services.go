package follow

type FollowUserService interface {
	FollowerCountUse(user uint64) FollowUserFollowCountResponse
	StudioFollowCount(studioid uint64) FollowUserStudioCountResponse
}

func (fs followService) FollowerCountUse(userID uint64) FollowUserFollowCountResponse {

	followersCount, _ := App.Repo.UserCountFollowing(userID)
	followingCount, _ := App.Repo.UserCountFollower(userID)
	fc := FollowUserFollowCountResponse{
		Followers: followersCount,
		Following: followingCount,
	}
	// update redis
	updateCahceFollowerCount(fc, userID)
	return fc
}

func (fs followService) StudioFollowCount(studioid uint64) FollowUserStudioCountResponse {
	followersCount, _ := App.Repo.StudioCountFollowing(studioid)
	fc := FollowUserStudioCountResponse{
		Followers: followersCount,
	}
	updateCahceStudioFollowerCount(fc, studioid)
	return fc
}
