import { useEffect, useState } from 'react';
import { marked } from 'marked';
import { getPosts } from './getPosts';

import './Posts.css';
import UserIcon from 'assets/images/User.svg';
import UpVoteIcon from './assets/UpVote.js';
import DownVoteIcon from './assets/DownVote.js';
import CommentsIcon from './assets/Comments.svg';

function Posts() {
  const [postsList, setPostsList] = useState([]);

  useEffect(() => {
    const fetchPosts = async () => {
      try {
        const posts = await getPosts();
        if (posts) {
          setPostsList(posts);
        } else {
          console.error('Error fetching posts');
        }
      } catch (error) {
        console.error('Error fetching posts', error);
      }
    };

    fetchPosts();
  }, [])

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
              <div className="post-text-area" dangerouslySetInnerHTML={{__html: marked.parse(post.body)}}></div>
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
                <button><UpVoteIcon className="post-action-vote" /></button>
                <div className="post-action-separator"></div>
                <button><DownVoteIcon className="post-action-vote" /></button>
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
