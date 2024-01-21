import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { marked } from 'marked';
import insane from 'insane';
import { getPosts } from './getPosts';
import { sendRemoveVote, sendVoteData } from './processVote';
import { UpVoteButton, DownVoteButton } from './postsInteractions';
import SanitizeRules from 'utils/SanitizeRules/Sanitize.json';
import { refreshTokens } from 'utils/TokensManagment/TokensManagment';
import notificationStore from 'utils/Notifications/notificationsStore';
import { NOTIFICATIONS } from 'utils/Notifications/notificationConstants';

import './Posts.css';
import UserIcon from 'assets/images/User.svg';
import CommentsIcon from './assets/Comments.svg';

function Posts() {
  const [postsList, setPostsList] = useState([]);
  const navigate = useNavigate();

  const postAction = {
    Increase: 1,
    Decrease: 2,
  }


  useEffect(() => {
    const fetchPosts = async () => {
      let posts;
      
      try {
        const posts = await getPosts();

        if (posts) {
          const postsData = posts.map(post => ({
            ...post,
            vote_lock: false
          }));

          setPostsList(postsData);
        } else {
          console.error("Error fetching posts, trying to refresh tokens...");
          const refreshResult = refreshTokens();

          if (refreshResult !== "success") {
            console.warn("Failed to refresh tokens, redirecting to login page...");
            navigate("/login");
          }

          const posts = await getPosts();

          if (posts) {
            const postsData = posts.map(post => ({
              ...post,
              vote_lock: false
            }));

            setPostsList(postsData);
          } else {
            console.error("Failed to load feed");
            notificationStore.addNotification("Failed to load feed", "err");
          }
        }
      } catch (error) {
        console.error("Error fetching posts", error);
      }
    };

    fetchPosts();
  }, [navigate])

  const processBody = (body) => {
    const parser = new DOMParser();
    const doc = parser.parseFromString(body, 'text/html');

    const images = doc.querySelectorAll('img');

    images.forEach(function(image) {
        image.classList.add('markdown-sandboxed-image');
    });

    return insane(doc.body.innerHTML, SanitizeRules)
  }

  const handleVoteUpdate = (postId, action, size) => {
    if (!size) {
      size = 1;
    }

    const updatedPosts = postsList.map(post => {
      if (post.id === postId) {
        if (action === 1) {
          console.log(`[PostID: ${postId}] Votes amount increased on the client side`);
          return { ...post, votes_amount: post.votes_amount + size };
        } else if (action === 2) {
          console.log(`[PostID: ${postId}] Votes amount decreased on the client side`);
          return { ...post, votes_amount: post.votes_amount - size };
        }
      }
      return post;
    });

    setPostsList(updatedPosts);
  }

  const getVoteStatus = (postId, action) => {
    let returnableValue = false;
    let prevVoteStatus = 0;

    const updatedPosts = postsList.map(post => {
      if (post.id === postId) {
        prevVoteStatus = post.vote_action;
        console.log(`[PostID: ${postId}] Preparing Post Action: ${post.vote_action}...`);
        if (post.vote_action === action) {
          console.log(`[PostID: ${postId}] Vote action changed to 0`);
          post.vote_action = 0;
          returnableValue = true;
        } else {
          console.log(`[PostID: ${postId}] Vote action changed to ${action}`);
          post.vote_action = action;
          returnableValue = false
        }
      }

      return post;
    });

    setPostsList(updatedPosts);
    return { returnableValue, prevVoteStatus };
  }

  const getLockStatus = (postId) => {
    const post = postsList.find(post => post.id === postId);
    return post ? post.vote_lock : undefined;
  }

  const setLockStatus = (postId, status) => {
    setPostsList(prevPosts => {
      const updatedPosts = prevPosts.map(post => {
        if (post.id === postId) {
          return {
            ...post,
            vote_lock: status
          };
        }
        return post;
      });
      return updatedPosts;
    });
  };

  const handleUpVote = async (postId) => {
    const lockStatus = await getLockStatus(postId);

    if (lockStatus === false) {
      const voteStatus = getVoteStatus(postId, postAction.Increase);

      if (voteStatus.returnableValue) {
        console.log(`[PostID: ${postId}] Removing Upvote...`);

        await setLockStatus(postId, true)
        const result = await sendRemoveVote(postId);

        if (result) {
          await handleVoteUpdate(postId, postAction.Decrease);
        } else {
          console.error(`[PostID: ${postId}] Failed to process action`);
          notificationStore.addNotification(NOTIFICATIONS.FAILED_VOTE.message, NOTIFICATIONS.FAILED_VOTE.type);
        }

        await setLockStatus(postId, false);
      } else {
        console.log(`[PostID: ${postId}] Upvoting post...`);

        await setLockStatus(postId, true)
        const result = await sendVoteData(postId, postAction.Increase);

        if (result) {
          if (voteStatus.prevVoteStatus === postAction.Decrease) {
            await handleVoteUpdate(postId, postAction.Increase, 2);
          } else {
            await handleVoteUpdate(postId, postAction.Increase, 1);
          }
        } else {
          console.error(`[PostID: ${postId}] Failed to process action`);
          notificationStore.addNotification(NOTIFICATIONS.FAILED_VOTE.message, NOTIFICATIONS.FAILED_VOTE.type);
        }

      }

      await setLockStatus(postId, false);
    } else {
      console.warn(`[PostID: ${postId}] Votes is locked`);
    }
  }

  const handleDownVote = async (postId) => {
    const lockStatus = await getLockStatus(postId);

    if (lockStatus === false) {
      const voteStatus = getVoteStatus(postId, postAction.Decrease);

      if (voteStatus.returnableValue) {
        console.log(`[PostID: ${postId}] Removing Downvote...`);

        await setLockStatus(postId, true);
        const result = await sendRemoveVote(postId);

        if (result) {
          handleVoteUpdate(postId, postAction.Increase);
        } else {
          console.error(`[PostID: ${postId}] Failed to process action`);
          notificationStore.addNotification(NOTIFICATIONS.FAILED_VOTE.message, NOTIFICATIONS.FAILED_VOTE.type);
        }
      } else {
        console.log(`[PostID: ${postId}] Downvoting post...`);

        await setLockStatus(postId, true);
        const result = await sendVoteData(postId, postAction.Decrease);

        if (result) {
          if (voteStatus.prevVoteStatus === postAction.Increase) {
            await handleVoteUpdate(postId, postAction.Decrease, 2);
          } else {
            await handleVoteUpdate(postId, postAction.Decrease, 1);
          }
        } else {
          console.error(`[PostID: ${postId}] Failed to process action`);
          notificationStore.addNotification(NOTIFICATIONS.FAILED_VOTE.message, NOTIFICATIONS.FAILED_VOTE.type);
        }
      }

      await setLockStatus(postId, false);
    } else {
      console.warn(`[PostID: ${postId}] Post is locked`);
    }
  }

  // TODO: Fix posts text overflow

  return (
    <>
      <div className="posts-container">
        {postsList.toReversed().map(post => (
          <div className="post-window" key={post.id}>
            <div className="post-header">
              <p>{post.type === 1 ? `article / ${post.nickname}` : `note / ${post.nickname}`}</p>
              <img src={post.avatar_img === "default" ? UserIcon : `${process.env.REACT_APP_BACKEND_DOMAIN}/images/${post.avatar_img}`} alt="" />
            </div>
            {post.type === 1 ?
              <div className="post-text-area" dangerouslySetInnerHTML={{__html: marked.parse(processBody(post.body))}}></div>
              :
              <>
                <div className="post-text-area">{post.body}</div>
                <div className="post-images-area">{post.attachments.map(image => (
                  <img key={image} src={`${process.env.REACT_APP_BACKEND_DOMAIN}/images/${image}.webp`} alt="" />
                ))}
                </div>
              </>
            }
            <div className="post-actions">
              <div>
                <UpVoteButton 
                  handleUpVote={handleUpVote} 
                  postId={post.id} 
                  buttonStatus={post.vote_action} 
                  reqAction={postAction.Increase} 
                />
                <div className="post-action-separator"></div>
                <DownVoteButton 
                  handleDownVote={handleDownVote} 
                  postId={post.id} 
                  buttonStatus={post.vote_action} 
                  reqAction={postAction.Decrease} 
                />
                <p>{post.votes_amount}</p>
              </div>
              <div>
                <img src={CommentsIcon} alt="" />
                <p>{post.comments_amount}</p>
              </div>
            </div>
          </div>
        ))}
      </div>
    </>
  )
}

export default Posts;
