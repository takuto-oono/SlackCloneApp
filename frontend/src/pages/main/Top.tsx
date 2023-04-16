import { Link } from 'react-router-dom';
import { useCookies } from "react-cookie";
import router from 'next/router';


export default function Top() {
  const [cookies, setCookie, removeCookie] = useCookies(['token', 'user_id']);
  console.log(cookies.token)
  if (cookies.token) {
    
    return (
      <>
        <nav>
          <ul>
            <li>
              <p>main</p>
            </li>
          </ul>
        </nav>
      </>)
  } else {
    return (<>
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
