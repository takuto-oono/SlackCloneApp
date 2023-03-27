import { Link } from 'react-router-dom';

export const Home = () => (
  <>
    <h2>Home</h2>
    <nav>
      <ul>
        <li>
          <Link to="login_form">ログイン</Link>
        </li>
        <li>
          <Link to="workspace">Workspace一覧</Link>
        </li>
      </ul>
    </nav>
  </>
);
