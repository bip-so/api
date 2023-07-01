package reactions

import "errors"

func (c reactionController) Create(obj CreateMentionPost, studioID uint64, userID uint64) error {
	var err error
	switch obj.Scope {
	case "block":
		err = App.Service.CreateBlockReaction(obj, studioID, userID)
		if err != nil {
			return err
		}
		break
	case "block_thread":
		err = App.Service.CreateBlockThreadReaction(obj, studioID, userID)
		if err != nil {
			return err
		}
		break
	case "block_comment":
		err = App.Service.CreateBlockThreadCommentReaction(obj, studioID, userID)
		if err != nil {
			return err
		}
		break
	case "reel":
		err = App.Service.CreateReelReaction(obj, studioID, userID)
		if err != nil {
			return err
		}
		break
	case "reel_comment":
		err = App.Service.CreateReelCommentReaction(obj, studioID, userID)
		if err != nil {
			return err
		}
		break
	default:
		return errors.New("Incorrect Input")
	}

	return nil
	//var reelComment *models.ReelComment
	//cb := models.CommentBase{
	//	0,
	//	reqData.Data,
	//	reqData.IsEdited,
	//	reqData.IsReply,
	//}
	//var parentCommentID *uint64
	//if *reqData.ParentID != 0 {
	//	parentCommentID = reqData.ParentID
	//} else {
	//	parentCommentID = nil
	//}
	//
	//reelCommetObject := reelComment.NewReelComment(cb, reelID, parentCommentID, userID, userID, userID)
	//reelComment, err := App.Repo.CreateReelComment(*reelCommetObject)
	//if err != nil {
	//	return nil, err
	//}
	//return reelComment, nil
}

func (c reactionController) Remove(obj CreateMentionPost, studioID uint64, userID uint64) error {
	var err error
	switch obj.Scope {
	case "block":
		err = App.Service.RemoveBlockReaction(obj, studioID, userID)
		if err != nil {
			return err
		}
		break
	case "block_thread":
		err = App.Service.RemoveBlockThreadReaction(obj, studioID, userID)
		if err != nil {
			return err
		}
		break
	case "block_comment":
		err = App.Service.RemoveBlockThreadCommentReaction(obj, studioID, userID)
		if err != nil {
			return err
		}
		break
	case "reel":
		err = App.Service.RemoveReelReaction(obj, studioID, userID)
		if err != nil {
			return err
		}
		break
	case "reel_comment":
		err = App.Service.RemoveReelCommentReaction(obj, studioID, userID)
		if err != nil {
			return err
		}
		break
	default:
		return errors.New("Incorrect Input")
	}

	return nil
	//var reelComment *models.ReelComment
	//cb := models.CommentBase{
	//	0,
	//	reqData.Data,
	//	reqData.IsEdited,
	//	reqData.IsReply,
	//}
	//var parentCommentID *uint64
	//if *reqData.ParentID != 0 {
	//	parentCommentID = reqData.ParentID
	//} else {
	//	parentCommentID = nil
	//}
	//
	//reelCommetObject := reelComment.NewReelComment(cb, reelID, parentCommentID, userID, userID, userID)
	//reelComment, err := App.Repo.CreateReelComment(*reelCommetObject)
	//if err != nil {
	//	return nil, err
	//}
	//return reelComment, nil
}
