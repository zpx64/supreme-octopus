import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { marked } from 'marked';
import insane from 'insane';
import { getPosts } from './getPosts';
import { sendVoteData } from './processVote';
import SanitizeRules from 'utils/SanitizeRules/Sanitize.json';
import { refreshTokens } from 'utils/TokensManagment/TokensManagment';
import notificationStore from 'utils/Notifications/NotificationsStore';

import './Posts.css';
import UserIcon from 'assets/images/User.svg';
import UpVoteIcon from './assets/UpVote.js';
import DownVoteIcon from './assets/DownVote.js';
import CommentsIcon from './assets/Comments.svg';

function Posts() {
  const [postsList, setPostsList] = useState([]);
  const navigate = useNavigate();

  const postAction = {
    UpVote: 1,
    DownVote: 2
  }

  useEffect(() => {
    const fetchPosts = async () => {
      try {
        const posts = await getPosts();
        if (posts) {
          setPostsList(posts);
        } else {
          console.error("Error fetching posts, trying to refresh tokens...");
          const refreshResult = refreshTokens();

          if (refreshResult !== "success") {
            console.warn("Failed to refresh tokens, redirecting to login page...");
            navigate("/login");
          }

          const posts = await getPosts();

          if (posts) {
            setPostsList(posts);
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
  }, [])

  const processBody = (body) => {
    const parser = new DOMParser();
    const doc = parser.parseFromString(body, 'text/html');

    const images = doc.querySelectorAll('img');

    images.forEach(function(image) {
        image.classList.add('markdown-sandboxed-image');
    });

    return insane(doc.body.innerHTML, SanitizeRules)
  }

  const handleVoteUpdate = (action, postId) => {
    const updatedPosts = postsList.map(post => {
      if (post.id === postId) {
        if (action == 1) {
          return { ...post, votes_amount: post.votes_amount + 1 };
        } else if (action == 2) {
          return { ...post, votes_amount: post.votes_amount - 1 };
        }
      }
      return post;
    });

    setPostsList(updatedPosts);
  }

  const handleUpVote = async (postId) => {
    const result = await sendVoteData(postId, postAction.UpVote);
    if (result) {
      handleVoteUpdate(postAction.UpVote, postId);
    } else {
      notificationStore.addNotification("Failed to process your voice.", "err");
    }
  }

  const handleDownVote = async (postId) => {
    const result = await sendVoteData(postId, postAction.DownVote);
    if (result) {
      handleVoteUpdate(postAction.DownVote, postId);
    } else {
      notificationStore.addNotification("Failed to process your voice.", "err");
    }
  }

  // TODO: Fix posts text overflow
  // TODO: Use markdown sanitizer while rendering posts

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
                <button onClick={() => handleUpVote(post.id)}>
                  <UpVoteIcon className="post-action-vote" />
                </button>
                <div className="post-action-separator"></div>
                <button onClick={() => handleDownVote(post.id)}>
                  <DownVoteIcon className="post-action-vote" />
                </button>
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
