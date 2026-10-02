import { useState } from "react";

function Blog() {
  const data = {
    blogTitle: "AWS Blog",
    blogContent: "Đây là nội dung blog",
    blogAuthor: "admin",
    blogDate: "2022-01-01",
  };

  const [title, setTitle] = useState(data.blogTitle);
  const [content, setContent] = useState(data.blogContent);
  const [author, setAuthor] = useState(data.blogAuthor);
  const [date, setDate] = useState(data.blogDate);
  const [listComment, setListComment] = useState([
    {
      id: 1,
      name: "name",
      content: "test",
    },
  ]);
  const [comment, setComment] = useState("");

  const handleComment = () => {
    console.log(comment);
    setComment("");
  };

  return (
    <div>
      <h1>{title}</h1>
      <p>{content}</p>
      <p>{author}</p>
      <p>{date}</p>

      <div>
        <label htmlFor="comment">Bình luận</label>
        <input
          type="text"
          id="comment"
          value={comment}
          onChange={(e) => setComment(e.target.value)}
        />
        <button onClick={handleComment}>Gửi</button>
      </div>

      <div>
        <h2>Bình luận</h2>
        {listComment.map((item) => (
          <div key={item.id}>
            <p>{item.name}</p>
            <p>{item.content}</p>
          </div>
        ))}
      </div>
    </div>
  );
}

export default Blog;
