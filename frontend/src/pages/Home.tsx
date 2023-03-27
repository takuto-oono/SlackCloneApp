import { Link } from 'react-router-dom';
import { useCookies } from "react-cookie";
import router from 'next/router';


export default function Home() {
  const [cookies, setCookie, removeCookie] = useCookies(['token', 'user_id']);
  if (cookies.token) {
    console.log(cookies.token)
    return (
      <>
        <h2>Home</h2>
        <nav>
          <ul>
            <li>
              <Link to="workspace">Workspace一覧</Link>
            </li>
          </ul>
        </nav>
      </>)
  } else {
    return (<>
      <h2>Home</h2>
      <nav>
        <ul>
          <li>
            <Link to="login_form">ログイン</Link>
          </li>
        </ul>
      </nav>
    </>)
  }
};
