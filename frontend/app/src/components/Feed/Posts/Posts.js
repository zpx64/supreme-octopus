import { useEffect, useState } from 'react';
import { marked } from 'marked';
import { getPosts } from './getPosts';
import DOMPurify from 'dompurify';
import './Posts.css';
import User from 'assets/images/User.svg';

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

  const getUserImage = () => {

  }

  // TODO: Fix posts text overflow
  // TODO: Use markdown sanitizer while rendering posts

  return (
    <>
      <div className="posts-container">
        {postsList.toReversed().map(post => (
          <div className="post-window" key={post.id}>
            <div className="post-header">
              <p>{post.type == 1 ? `article / ${post.nickname}` : `note / ${post.nickname}`}</p>
              <img src={post.avatar_img == "default" ? User : `${process.env.REACT_APP_BACKEND_DOMAIN}/images/${post.avatar_img}`} />
            </div>
            {post.type == 1 ?
              <div className="post-text-area" dangerouslySetInnerHTML={{__html: DOMPurify.sanitize(marked.parse(post.body))}}></div>
              :
              <>
                <div className="post-text-area">{post.body}</div>
                <div className="post-images-area">{post.attachments.map(image => (
                  <img key={image} src={`${process.env.REACT_APP_BACKEND_DOMAIN}/images/${image}.webp`} />
                ))}
                </div>
              </>
            }
          </div>
        ))}
      </div>
    </>
  )
}

export default Posts;
