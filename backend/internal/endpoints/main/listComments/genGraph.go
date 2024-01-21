package listComments

import (
	"github.com/zpx64/supreme-octopus/internal/model"
)

func getCommentFromCommentWithUser(comment model.CommentWithUser) Comment {
	return Comment{
		CommentId:    comment.CommentId,
		Nickname:     comment.Nickname,
		AvatarImg:    comment.AvatarImg,
		Body:         comment.Body,
		Attachments:  comment.Attachments,
		CreationDate: comment.CreationDate,
		VotesAmount:  comment.VotesAmount,
	}
}

// NEED SORTED BY comment_id ARRAY
// TODO: optimize. maybe rewrite in rust @electro
func ConvertArrayOfCommentsToGraphOfComments(
	comments []model.CommentWithUser,
) []Comment {
	var (
		commentsMap    = make(map[int]int, len(comments))
		commentThreads = make([]Comment, 0, len(comments))
		commentsUsed   = make(map[int]struct{}, len(comments))
	)

	// fill map
	for i, e := range comments {
		commentsMap[e.CommentId] = i
	}

	var genGraph func([]model.CommentWithUser, *Comment)
	genGraph = func(comments []model.CommentWithUser, comment *Comment) {
		for i, commentReply := range comments {
			_, commentAlreadyUsed := commentsUsed[i]
			if commentAlreadyUsed ||
				commentReply.ReplyId == nil ||
				comment.CommentId == commentReply.CommentId ||
				comment.CommentId != *commentReply.ReplyId {
				continue
			}

			commentsUsed[i] = struct{}{}

			replyedComment := getCommentFromCommentWithUser(comments[i])
			genGraph(comments, &replyedComment)
			comment.Reply = append(comment.Reply, replyedComment)
		}
	}

	for i, comment := range comments {
		_, commentAlreadyUsed := commentsUsed[i]
		if commentAlreadyUsed {
			continue
		}

		commentThread := getCommentFromCommentWithUser(comment)
		genGraph(comments, &commentThread)
		commentThreads = append(commentThreads, commentThread)
	}

	return commentThreads
}
